package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	var timeout time.Duration
	flag.DurationVar(&timeout, "timeout", 10*time.Second, "connection timeout")
	flag.Parse()

	args := flag.Args()
	if len(args) != 2 {
		fmt.Fprintln(os.Stderr, "Usage: go-telnet [--timeout=10s] host port")
		os.Exit(1)
	}

	host := args[0]
	port := args[1]
	address := net.JoinHostPort(host, port)

	client := NewTelnetClient(address, timeout, os.Stdin, os.Stdout)

	if err := client.Connect(); err != nil {
		fmt.Fprintf(os.Stderr, "...Failed to connect: %v\n", err)
		os.Exit(1)
	}
	defer client.Close()

	fmt.Fprintf(os.Stderr, "...Connected to %s\n", address)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	sendDone := make(chan error, 1)
	receiveDone := make(chan error, 1)

	go func() {
		sendDone <- client.Send()
	}()

	go func() {
		receiveDone <- client.Receive()
	}()

	select {
	case <-ctx.Done():
		fmt.Fprintln(os.Stderr, "...Interrupted by signal")
	case err := <-sendDone:
		if err != nil {
			if err == io.EOF {
				fmt.Fprintln(os.Stderr, "...EOF")
			} else {
				fmt.Fprintf(os.Stderr, "...Send error: %v\n", err)
			}
		}
	case err := <-receiveDone:
		if err != nil {
			if err == io.EOF {
				fmt.Fprintln(os.Stderr, "...Connection was closed by peer")
			} else {
				fmt.Fprintf(os.Stderr, "...Receive error: %v\n", err)
			}
		} else {
			fmt.Fprintln(os.Stderr, "...Connection was closed by peer")
		}
	}
}
