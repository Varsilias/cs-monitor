package dislpay

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"sort"

	"github.com/fatih/color"
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
	color.Cyan("┌─────────────────────────────────────────────────┐")
	color.Cyan("│              CONTAINER STATS MONITOR            │")
	color.Cyan("├─────────────────────────────────────────────────┤")
	color.Green("│ Container ID: %-33s │\n", shortenID(stats.ID))
	color.Green("├─────────────────────────────────────────────────┤")
	color.Green("│ CPU Usage:    %-6.2f%%                            │\n", stats.CPUPercentage)
	color.Green("│ Memory:       %-10s / %-10s (%.2f%%)    │\n",
		memUsageFormatted, memLimitFormatted, stats.MemPercentage)
	color.Green("├─────────────────────────────────────────────────┤")
	color.Blue("│ Network I/O:  ↓ %-10s / ↑ %-10s      │\n",
		netInputFormatted, netOutputFormatted)
	color.Blue("│ Block I/O:    ↓ %-10s / ↑ %-10s      │\n",
		blockInputFormatted, blockOutputFormatted)
	color.Blue("└─────────────────────────────────────────────────┘")
}

// MultiContainerStats holds stats for multiple containers
type MultiContainerStats map[string]*collector.ContainerStats

// RenderMultiStats displays stats for multiple containers
func (d *Display) RenderMultiStats(stats MultiContainerStats) {
	clearScreen()

	color.Cyan("┌─────────────────────────────────────────────────────────────────────────────────────────────────┐")
	color.Cyan("│                                   CONTAINER STATS MONITOR                                       │")
	color.Cyan("├──────────────┬─────────────────────┬────────────┬─────────────────┬─────────────┬──────────────┤")
	color.Green("│ CONTAINER ID │ NAME                │ CPU %      │ MEMORY USAGE    │ NET I/O     │ BLOCK I/O    │")
	color.Green("├──────────────┼─────────────────────┼────────────┼─────────────────┼─────────────┼──────────────┤")

	if len(stats) == 0 {
		fmt.Println("│ No running containers found                                                                    │")
		fmt.Println("└──────────────┴─────────────────────┴────────────┴─────────────────┴─────────────┴──────────────┘")
		return
	}

	// Sort container IDs for consistent display
	var containerIDs []string
	for id := range stats {
		containerIDs = append(containerIDs, id)
	}
	sort.Strings(containerIDs)

	for _, id := range containerIDs {
		stat := stats[id]
		fmt.Printf("│ %-12s", shortenID(stat.ID))
		fmt.Printf("│ %-19s", truncateString(stat.Name, 19))
		fmt.Printf("│ %8.2f%%  ", stat.CPUPercentage)
		fmt.Printf("│ %8s/%-8s", formatBytes(stat.MemUsage), formatBytes(stat.MemLimit))
		fmt.Printf("│ %11s", formatBytes(stat.NetInput)+"/"+formatBytes(stat.NetOutput))
		fmt.Printf("│ %12s │\n", formatBytes(stat.BlockInput)+"/"+formatBytes(stat.BlockOutput))
	}
	fmt.Println("└──────────────┴─────────────────────┴────────────┴─────────────────┴─────────────┴──────────────┘")
	color.Red("Press Ctrl+C to exit")
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

// Helper function to truncate strings for display
func truncateString(s string, maxLen int) string {
	if len(s) > maxLen {
		return s[:maxLen-3] + "..."
	}
	return s
}
