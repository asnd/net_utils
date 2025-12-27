# Net Utils

Net Utils is a collection of networking utilities written in Go and Python. It provides a set of tools for performing various network-related tasks.

## Table of Contents

- [Installation](#installation)
- [Usage](#usage)
- [Tools](#tools)
  - [DNS Query (Go)](#dns-query-go)
  - [PCAP Stripper (Go)](#pcap-stripper-go)
  - [Python Tools](#python-tools)
- [Contributing](#contributing)
- [License](#license)

## Installation

To use the Net Utils tools, you need to have Go installed on your system. You can download and install Go from the official website: [https://golang.org/](https://golang.org/)

To install Net Utils, follow these steps:

1. Clone the repository:
   ```
   git clone https://gitlab.com/nikosa/net_utils.git
   ```

2. Change to the project directory:
   ```
   cd net_utils
   ```

3. Build the Go tools:
   ```bash
   go build -o dns_query ./cmd/dns_query
   go build -o pcap_stripper ./cmd/pcap_stripper
   ```

## Usage

Each tool in Net Utils is a separate program. You can run the tools individually by executing the corresponding binary or script.

## Tools

### DNS Query (Go)

The DNS Query tool allows you to perform DNS queries against multiple DNS servers concurrently and measure the response times.

**Usage:**
```bash
./dns_query [flags]
```

**Flags:**
- `--servers`: Comma-separated list of DNS servers to test (default: common public DNS servers)
- `--domains`: Comma-separated list of domains to query (default: common websites)
- `--repeat`: Number of times to repeat the test (default: 3)

**Example:**
```bash
./dns_query --servers="1.1.1.1,8.8.8.8" --domains="google.com,example.com" --repeat=5
```

### PCAP Stripper (Go)

PCAP Stripper is a command-line tool for removing a specified number of bytes from packet payloads in PCAP/CAP files.

**Usage:**
```bash
./pcap_stripper -i <input_file> -b <bytes_to_strip> [flags]
```

**Flags:**
- `-i, --input`: Input PCAP/CAP file (required)
- `-o, --output`: Output PCAP file (auto-generated if not specified)
- `-b, --bytes`: Number of bytes to strip from each packet (required)
- `-p, --position`: Position to strip from: 'beginning' or 'end' (default: "beginning")
- `--preview`: Preview first 5 packets of the output file
- `-v, --verbose`: Enable verbose output
- `-f, --force`: Force overwrite of existing output file

**Example:**
```bash
./pcap_stripper -i input.pcap -b 14 -p beginning --preview
```

### Python Tools

The repository also contains Python utilities. Ensure you have Python 3 installed.

#### DNS_checker.py
A Python equivalent of the DNS Query tool.

```bash
python3 DNS_checker.py --servers "1.1.1.1,8.8.8.8" --domains "google.com" --repeat 5
```

#### ookla.py
Downloads data extract files from Speedtest Intelligence.

```bash
python3 ookla.py --api-key YOUR_KEY --api-secret YOUR_SECRET --output-dir ./downloads
```
Or use environment variables `OOKLA_API_KEY` and `OOKLA_API_SECRET`.

#### pcap_stripper.py
A Python implementation of the PCAP stripper using `pyshark`.

```bash
python3 pcap_stripper.py strip input.pcap 14
```

## Contributing

Contributions to Net Utils are welcome! If you find any issues or have suggestions for improvements, please open an issue or submit a pull request on the GitLab repository.

When contributing to this repository, please follow the [code of conduct](CODE_OF_CONDUCT.md) and the [contributing guidelines](CONTRIBUTING.md).

## License

Net Utils is licensed under the [MIT License](LICENSE).