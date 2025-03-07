package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Place your code here,
	// P.S. Do not rush to throw context down, think think if it is useful with blocking operation?
	timeout := flag.Duration("timeout", 10*time.Second, "connection timeout")
	flag.Parse()

	args := flag.Args()
	if len(args) != 2 {
		log.Fatal("Usage: go-telnet [--timeout=10s] host port")
	}

	address := net.JoinHostPort(args[0], args[1])
	client := NewTelnetClient(address, *timeout, os.Stdin, os.Stdout)

	if err := client.Connect(); err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	fmt.Fprintf(os.Stderr, "...Connected to %s\n", address)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT)
	errChan := make(chan error, 2)

	// Запуск горутин для отправки и получения данных
	go func() {
		errChan <- client.Send()
	}()

	go func() {
		errChan <- client.Receive()
	}()

	// Ожидание сигнала или ошибки
	select {
	case <-sigChan:
		fmt.Fprintln(os.Stderr, "\n...Connection was closed by client")
	case err := <-errChan:
		if err != nil {
			fmt.Fprintln(os.Stderr, "...Connection was closed by server")
		}
	}
}
