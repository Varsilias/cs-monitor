## CS-MONITOR

A tool built with Go for Collecting Metrics of a Running Containers

## Usage

- Clone the repository `git clone git@github.com:Varsilias/cs-monitor.git`
- Navigate into the project directory `cd cs-monitor`
- Built the project `go build -o cs-monitor  main.go`
- Show Metrics for a single running container `./cs-monitor <container_id>`
- Show Metrics for all running containers `./cs-monitor --all`
- Press `Press Ctrl+C to exit` in either case of running the binary
