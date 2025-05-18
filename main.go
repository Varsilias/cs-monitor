package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/varsilias/cs-monitor/collector"
	dislpay "github.com/varsilias/cs-monitor/display"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: cs-monitor [container_id]")
		os.Exit(1)
	}

	containerID := os.Args[1]
	ctx := context.Background()

	c, err := collector.NewCollector(ctx)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	d := dislpay.NewDisplay()
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			stats, err := c.GetContainerStats(ctx, containerID)
			if err != nil {
				fmt.Printf("Error getting stats: %v\n", err)
				continue
			}
			d.RenderStats(stats)
		}

	}
}
