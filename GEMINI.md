# GEMINI Project Profile: Net Utils

## Project Overview
- **Name**: Net Utils
- **Primary Language**: Go & Python
- **Purpose**: A collection of networking utilities for DNS querying, PCAP stripping, and speed tests.
- **Key Features**: High-concurrency DNS query tool (Go), PCAP payload stripper (Go & Python), and Ookla Speedtest data fetcher.

## Project Structure
```
net_utils/
├── cmd/                     # Go application entry points
│   ├── dns_query/           # DNS Query tool
│   └── pcap_stripper/       # PCAP Stripper tool
├── DNS_checker.py           # Python DNS checker
├── ookla.py                 # Ookla Speedtest data fetcher
├── pcap_stripper.py         # Python PCAP stripper
├── go.mod                   # Go module definition
└── README.md                # Documentation
```

## Build & Deployment
- **Build System**: `go build` for Go tools.
- **CI/CD**: `.gitlab-ci.yml` present (likely for build and test).

## Suggested Development Tools
- **VSCode Extensions**:
  - `golang.go`: Go language support.
  - `ms-python.python`: Python language support.
- **CLI Tools**:
  - `go`: Go compiler and toolchain.
  - `python3`: Python interpreter.
- **MCP Servers**:
  - `filesystem`: For file access.
