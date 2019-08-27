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
	Address     string
	Timeout     int64
	ReadTimeout int64
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

	inputCh := make(chan string, 1)
	responseCh := make(chan string, 1)
	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, os.Interrupt)

	go func(stopCh <-chan os.Signal) {
		<-stopCh

		log.Println("Get SIGINT signal")
		cancel()
	}(stopCh)

	go func(inputCh chan<- string) {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			inputCh <- scanner.Text()
		}
	}(inputCh)

	go func(ctx context.Context, conn net.Conn, readTimeout int64, responseCh chan<- string) {
		timeoutForRead := time.Duration(readTimeout) * time.Millisecond
		err := conn.SetReadDeadline(time.Now().Add(timeoutForRead))
		if err != nil {
			log.Println(err)
		}

		scanner := bufio.NewScanner(conn)

		for scanner.Scan() {
			responseCh <- scanner.Text()
		}
	}(ctx, conn, options.ReadTimeout, responseCh)

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
