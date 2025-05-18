package collector

import (
	"context"
	"encoding/json"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

type ContainerStats struct {
	ID            string
	Name          string
	CPUPercentage float64
	MemUsage      float64
	MemLimit      float64
	MemPercentage float64
	NetInput      float64
	NetOutput     float64
	BlockInput    float64
	BlockOutput   float64
}

type Collector struct {
	client *client.Client
}

func NewCollector(ctx context.Context) (*Collector, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}

	return &Collector{
		client: cli,
	}, nil
}

func (c *Collector) GetContainerStats(ctx context.Context, containerID string) (*ContainerStats, error) {
	stats, err := c.client.ContainerStats(ctx, containerID, false)
	if err != nil {
		return nil, err
	}
	defer stats.Body.Close()
	var statsResponse container.StatsResponse
	if err = json.NewDecoder(stats.Body).Decode(&statsResponse); err != nil {
		return nil, err
	}

	cpuPercent := calculateCPUPercentage(&statsResponse)

	memUsage := float64(statsResponse.MemoryStats.Usage)
	memLimit := float64(statsResponse.MemoryStats.Limit)

	memPercent := 0.0

	if memLimit > 0 {
		memPercent = (memUsage / memLimit) * 100.0
	}

	netInput, netOutput := calculateNetworkPercentage(&statsResponse)
	blockInput, blockOutput := calculateBlockIOUsage(&statsResponse)

	return &ContainerStats{
		ID:            containerID,
		Name:          stats.OSType,
		CPUPercentage: cpuPercent,
		MemUsage:      memPercent,
		MemLimit:      memLimit,
		MemPercentage: memPercent,
		NetInput:      netInput,
		NetOutput:     netOutput,
		BlockInput:    blockInput,
		BlockOutput:   blockOutput,
	}, nil
}

func calculateCPUPercentage(stats *container.StatsResponse) float64 {
	cpuDelta := float64(stats.CPUStats.CPUUsage.TotalUsage) - float64(stats.PreCPUStats.CPUUsage.TotalUsage)
	systemDelta := float64(stats.CPUStats.SystemUsage) - float64(stats.CPUStats.SystemUsage)

	if systemDelta > 0.0 && cpuDelta > 0.0 {
		return (cpuDelta / systemDelta) * float64(len(stats.CPUStats.CPUUsage.PercpuUsage)) * 100.0
	}

	return 0.0
}

func calculateNetworkPercentage(stats *container.StatsResponse) (float64, float64) {
	var rx, tx float64

	for _, network := range stats.Networks {
		rx += float64(network.RxBytes)
		tx += float64(network.TxBytes)
	}

	return rx, tx
}

func calculateBlockIOUsage(stats *container.StatsResponse) (float64, float64) {
	var read, write float64

	for _, io := range stats.BlkioStats.IoServiceBytesRecursive {
		if io.Op == "Read" {
			read += float64(io.Value)
		} else if io.Op == "Write" {
			write += float64(io.Value)
		}
	}

	return read, write
}
