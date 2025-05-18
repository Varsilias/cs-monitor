package dislpay

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/varsilias/cs-monitor/collector"
)

type Display struct{}

func NewDisplay() *Display {
	return &Display{}
}

func (d *Display) RenderStats(stats *collector.ContainerStats) {

	clearScreen()

	// Format memory values
	memUsageFormatted := formatBytes(stats.MemUsage)
	memLimitFormatted := formatBytes(stats.MemLimit)

	// Format network values
	netInputFormatted := formatBytes(stats.NetInput)
	netOutputFormatted := formatBytes(stats.NetOutput)

	// Format block I/O values
	blockInputFormatted := formatBytes(stats.BlockInput)
	blockOutputFormatted := formatBytes(stats.BlockOutput)

	// Display stats
	fmt.Println("┌─────────────────────────────────────────────────┐")
	fmt.Println("│              CONTAINER STATS MONITOR            │")
	fmt.Println("├─────────────────────────────────────────────────┤")
	fmt.Printf("│ Container ID: %-33s │\n", shortenID(stats.ID))
	fmt.Println("├─────────────────────────────────────────────────┤")
	fmt.Printf("│ CPU Usage:    %-6.2f%%                            │\n", stats.CPUPercentage)
	fmt.Printf("│ Memory:       %-10s / %-10s (%.2f%%)    │\n",
		memUsageFormatted, memLimitFormatted, stats.MemPercentage)
	fmt.Println("├─────────────────────────────────────────────────┤")
	fmt.Printf("│ Network I/O:  ↓ %-10s / ↑ %-10s      │\n",
		netInputFormatted, netOutputFormatted)
	fmt.Printf("│ Block I/O:    ↓ %-10s / ↑ %-10s      │\n",
		blockInputFormatted, blockOutputFormatted)
	fmt.Println("└─────────────────────────────────────────────────┘")
}

func clearScreen() {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}

	cmd.Stdout = os.Stdout
	cmd.Run()
}

func shortenID(id string) string {
	if len(id) > 12 {
		return id[:12]
	}
	return id
}

func formatBytes(bytes float64) string {
	const unit = 1024.0
	if bytes < unit {
		return fmt.Sprintf("%.0f B", bytes)
	}
	div, exp := unit, 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", bytes/div, "KMGTPE"[exp])
}
