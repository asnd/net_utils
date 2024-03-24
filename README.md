
# Net Utils

Net Utils is a collection of networking utilities written in Go and python . It provides a set of tools for performing various network-related tasks.

## Table of Contents

- [Installation](#installation)
- [Usage](#usage)
- [Tools](#tools)
  - [DNS Query](#dns-query)
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

3. Build the tools:
   ```
   go build ./...
   ```

## Usage

Each tool in Net Utils is a separate Go program. You can run the tools individually by executing the corresponding binary.

## Tools

### DNS Query

The DNS Query tool allows you to perform DNS queries against multiple DNS servers concurrently and measure the response times.

To run the DNS Query tool:
```
./dns_query
```

The tool will display the DNS query test results, showing the response times for each DNS server.

## Contributing

Contributions to Net Utils are welcome! If you find any issues or have suggestions for improvements, please open an issue or submit a pull request on the GitLab repository.

When contributing to this repository, please follow the [code of conduct](CODE_OF_CONDUCT.md) and the [contributing guidelines](CONTRIBUTING.md).

## License

Net Utils is licensed under the [MIT License](LICENSE).

---

Feel free to customize the README file based on your specific project structure, tools, and requirements. Include relevant sections, installation instructions, usage examples, and any additional information that would be helpful for users and contributors.