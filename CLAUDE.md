# Net Utils

Network utilities combining Go CLI tools and Python scripts for DNS, packet analysis, and speed testing.

## Tech Stack
- **Languages**: Go 1.24, Python 3.10+
- **Go Libraries**: gopacket, miekg/dns, cobra
- **Python**: Various utility scripts
- **CI/CD**: GitLab CI

## Development Commands

### Go Commands
```bash
# Build Go CLI tools
go build -o bin/ ./cmd/...

# Run tests
go test ./...

# Format
go fmt ./...

# Vet
go vet ./...
```

### Python Commands
```bash
# Run DNS checker
python3 DNS_checker.py

# Run speed test
python3 ookla.py

# Run PCAP stripper
python3 pcap_stripper.py
```

## Project Structure
- `cmd/` - Go CLI commands
- `*.py` - Python utility scripts
- `vmware_ai_agent/` - VMware AI agent integration

## Key Tools
- **dns_query_arm** - Pre-built ARM binary for DNS queries
- **DNS_checker.py** - DNS validation utility
- **ookla.py** - Speed test utility
- **pcap_stripper.py** - PCAP file processing

## Suggested Claude Code Plugins

### MCP Servers
- **filesystem** - For editing scripts and configs
- **docker** - Container management

### Skills
- **python-to-golang** - For converting Python utilities to Go
- **nsx-avi-reference** - For network infrastructure context

### Recommended Workflow
1. Use `go fmt` and `go vet` for Go code
2. Add type hints to Python scripts
3. Build cross-platform binaries with `GOOS` and `GOARCH`
