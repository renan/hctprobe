package main

import (
	"context"
	"flag"
	"log"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	address string

	interval, timeout, healthyThreshold, unhealthyThreshold int
)

func init() {
	flag.StringVar(&address, "address", ":8080", "Host to start the TCP Server")
	flag.IntVar(&interval, "interval", 1, "The approximate amount of time, in seconds, between health checks of an individual target")
	flag.IntVar(&timeout, "timeout", 6, "The amount of time, in seconds, during which no response means a failed health check")
	flag.IntVar(&healthyThreshold, "healthyThreshold", 3, "The number of consecutive health checks successes required before considering an unhealthy target healthy")
	flag.IntVar(&unhealthyThreshold, "unhealthyThreshold", 3, "The number of consecutive health check failures required before considering the target unhealthy")
}

func usage() {
	log.Print("Usage: hctprobe [options] http://example.org\n\nOptions:\n")
	flag.PrintDefaults()
}

func main() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		sig := <-signals

		log.Printf("Exitting with %s", sig)
		os.Exit(1)
	}()

	flag.Usage = usage
	flag.Parse()

	target := flag.Arg(0)
	_, err := url.ParseRequestURI(target)
	if err != nil {
		log.Printf("Error parsing target URI. Error: %s", err.Error())
		usage()
		os.Exit(1)
	}

	probe := NewProbe(
		ctx,
		NewTCPServer(ctx, address),
		target,
		time.Duration(interval)*time.Second,
		time.Duration(timeout)*time.Second,
		healthyThreshold,
		unhealthyThreshold,
	)

	probe.Start()
}
