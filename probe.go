package main

import (
	"context"
	"log"
	"net/http"
	"time"
)

type Probe struct {
	ctx       context.Context
	tcpServer *TCPServer

	target             string
	interval           time.Duration
	healthyThreshold   int
	unhealthyThreshold int

	httpClient *http.Client
	ticker     *time.Ticker
}

func NewProbe(ctx context.Context, tcpServer *TCPServer, target string, interval, timeout time.Duration, healthyThreshold, unhealthyThreshold int) *Probe {
	return &Probe{
		ctx:       ctx,
		tcpServer: tcpServer,
		httpClient: &http.Client{
			Timeout: timeout,
		},

		target:             target,
		interval:           interval,
		healthyThreshold:   healthyThreshold,
		unhealthyThreshold: unhealthyThreshold,
	}
}

func (p *Probe) Start() {
	p.ticker = time.NewTicker(p.interval)
	defer p.Stop()

	state := "unhealthy"
	transition := time.Now()

	healthyCounter := 0
	unhealthyCounter := 0

	log.Printf("Probe is starting, health check state is %s", state)

	for {
		select {
		case <-p.ctx.Done():
			return
		case <-p.ticker.C:
			healthy := p.checkTarget()
			if healthy {
				healthyCounter++
			} else {
				unhealthyCounter++
			}

			if state != "healthy" && healthyCounter >= p.healthyThreshold {
				log.Printf("Health check became healthy after %s", time.Now().Sub(transition))

				state = "healthy"
				transition = time.Now()

				go func() {
					err := p.tcpServer.Start()
					if err != nil {
						log.Fatalf("Failed to start the probe. Error: %s", err.Error())
					}
				}()
			} else if state != "unhealthy" && unhealthyCounter >= p.unhealthyThreshold {
				log.Printf("Health check became unhealthy after %s", time.Now().Sub(transition))

				state = "unhealthy"
				transition = time.Now()

				p.tcpServer.Stop()
			}

			if healthy && state == "healthy" {
				unhealthyCounter = 0
			}
			if !healthy && state == "unhealthy" {
				healthyCounter = 0
			}
		}
	}
}

func (p *Probe) checkTarget() bool {
	response, err := p.httpClient.Get(p.target)
	if err != nil {
		return false
	}

	return response.StatusCode == 200
}

func (p *Probe) Stop() {
	p.tcpServer.Stop()
	p.ticker.Stop()
}
