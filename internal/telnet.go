package internal

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"time"
)

// App options
type Options struct {
	Address string
	Timeout int64
}

// Starts execution
func Run(options Options) {
	dialer := &net.Dialer{}
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Duration(options.Timeout)*time.Millisecond)
	defer cancel()

	conn, err := dialer.DialContext(ctx, "tcp", options.Address)
	if err != nil {
		log.Fatalf("Cannot connect: %v", err)
	}

	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, os.Interrupt)

	go func() {
		<-stopCh

		log.Println("Get SIGINT signal")
		cancel()
	}()

	readRoutine(ctx, conn)
	writeRoutine(ctx, conn)

	err = conn.Close()
	if err != nil {
		log.Fatal(err)
	}
}

// Reads from a connection
func readRoutine(ctx context.Context, conn net.Conn) {
OUTER:
	for {
		select {
		case <-ctx.Done():
			return
		default:
			timeoutForRead := time.Millisecond * 150
			reader := bufio.NewReader(conn)

			for {
				err := conn.SetReadDeadline(time.Now().Add(timeoutForRead))
				if err != nil {
					log.Fatal(err)
				}

				bytes, err := reader.ReadBytes('\n')
				if err != nil {
					continue OUTER
				}

				log.Printf("From server: %s", string(bytes))
			}
		}
	}
}

// Writes to the connection
func writeRoutine(ctx context.Context, conn net.Conn) {
	inputCh := make(chan string, 1)
	go getInput(inputCh)

OUTER:
	for {
		select {
		case <-ctx.Done():
			break OUTER
		case data := <-inputCh:
			log.Printf("To server %v\n", data)

			_, err := conn.Write([]byte(fmt.Sprintf("%s\n", data)))
			if err != nil {
				log.Printf("error writing to connection: %v", err)
			}
		}

	}

	log.Printf("Finished writeRoutine")
}

// Reads data from stdin
func getInput(inputCh chan<- string) {
	for {
		in := bufio.NewReader(os.Stdin)
		result, err := in.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		inputCh <- result
	}
}
