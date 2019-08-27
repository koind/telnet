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

	inputCh := make(chan string)
	responseCh := make(chan string)
	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, os.Interrupt)

	go func(stopCh <-chan os.Signal) {
		<-stopCh

		log.Println("Get SIGINT signal")
		cancel()
	}(stopCh)

	go func(ctx context.Context, inputCh chan<- string) {
		scanner := bufio.NewScanner(os.Stdin)

		for {
			select {
			case <-ctx.Done():
				fmt.Println("Write close")
				return
			default:
				if !scanner.Scan() {
					return
				}
				inputCh <- scanner.Text()
			}
		}
	}(ctx, inputCh)

	go func(ctx context.Context, conn net.Conn, responseCh chan<- string) {
		scanner := bufio.NewScanner(conn)

		for {
			select {
			case <-ctx.Done():
				fmt.Println("Read close")
				return
			default:
				if !scanner.Scan() {
					log.Printf("CANNOT SCAN")
					return
				}
				responseCh <- scanner.Text()
			}
		}
	}(ctx, conn, responseCh)

OUTER:
	for {
		select {
		case message := <-inputCh:
			log.Printf("To server %v\n", message)

			_, err := conn.Write([]byte(fmt.Sprintf("%s\n", message)))
			if err != nil {
				log.Printf("error writing to connection: %v", err)
			}
		case message := <-responseCh:
			log.Printf("From server: %s", message)
		case <-ctx.Done():
			break OUTER
		}
	}

	err = conn.Close()
	if err != nil {
		log.Fatal(err)
	}
}
