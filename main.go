package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
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

	monitorAllContainers := len(os.Args) == 1 || len(os.Args) == 2 && os.Args[1] == "--all"

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	ticker := time.NewTicker(1 * time.Second)

	defer ticker.Stop()

	if monitorAllContainers {
		runMultiContainerMonitor(ctx, c, d, ticker, sigs)
	} else {
		runSingleContainerMonitor(ctx, c, d, containerID, ticker, sigs)
	}

}

func runSingleContainerMonitor(
	ctx context.Context,
	c *collector.Collector,
	d *dislpay.Display,
	containerID string,
	ticker *time.Ticker,
	sigs chan os.Signal,
) {
	for {
		select {
		case <-ticker.C:
			stats, err := c.GetContainerStats(ctx, containerID)
			if err != nil {
				fmt.Printf("Error getting stats: %v\n", err)
				continue
			}
			d.RenderStats(stats)
		case <-sigs:
			fmt.Println("\nExiting...")
			return
		}

	}
}

func runMultiContainerMonitor(
	ctx context.Context,
	c *collector.Collector,
	d *dislpay.Display,
	ticker *time.Ticker,
	sigs chan os.Signal,
) {
	for {
		select {
		case <-ticker.C:
			containers, err := c.ListRunningContainers(ctx)
			if err != nil {
				fmt.Printf("Error listing containers: %v\n", err)
				continue
			}

			if len(containers) == 0 {
				d.RenderMultiStats(dislpay.MultiContainerStats{})
			}
			// Collect stats for all containers
			multiStats := make(dislpay.MultiContainerStats, len(containers))
			for _, container := range containers {
				stats, err := c.GetContainerStats(ctx, container.ID)
				if err != nil {
					fmt.Printf("Error getting stats for container %s: %v\n", container.ID, err)
					continue
				}
				multiStats[container.ID] = stats
			}

			d.RenderMultiStats(multiStats)
		case <-sigs:
			fmt.Println("\nExiting...")
			return

		}
	}
}
