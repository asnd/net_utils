# PCAP Stripper - Help Examples

This document shows the help output you'll get from the refactored Cobra-based CLI.

## Main Help Output (`./pcap_stripper --help` or `./pcap_stripper -h`)

```
PCAP Stripper - Strip bytes from packet payloads in PCAP files

PCAP Stripper is a command-line tool for removing a specified number of bytes
from packet payloads in PCAP/CAP files. This is useful for sanitizing captures,
removing headers, or preparing files for analysis.

The tool can strip bytes from either the beginning or end of each packet's payload,
preserving the PCAP file structure and metadata while modifying only the packet data.

Usage:
  pcap_stripper [flags]

Examples:
  # Strip 14 bytes from the beginning of each packet (e.g., Ethernet header)
  pcap_stripper -i input.pcap -o output.pcap -b 14 -p beginning

  # Strip 4 bytes from the end of each packet
  pcap_stripper -i capture.cap -o sanitized.pcap -b 4 -p end

  # Use auto-generated output filename with verbose output
  pcap_stripper -i input.pcap -b 14 -p beginning -v

  # Force overwrite existing output file
  pcap_stripper -i input.pcap -o output.pcap -b 14 -p beginning -f

Flags:
  -b, --bytes int         Number of bytes to strip from each packet (required)
  -f, --force             Force overwrite of existing output file
  -h, --help              help for pcap_stripper
  -i, --input string      Input PCAP/CAP file (required)
  -o, --output string     Output PCAP file (auto-generated if not specified)
  -p, --position string   Position to strip from: 'beginning' or 'end' (default "beginning")
  -v, --verbose           Enable verbose output
```

## Error Messages Examples

### Missing Required Flag:
```
$ ./pcap_stripper
Error: required flag(s) "bytes", "input" not set
```

### Invalid Bytes Value:
```
$ ./pcap_stripper -i input.pcap -b 0
Error: bytes to strip must be greater than 0, got: 0
```

```
$ ./pcap_stripper -i input.pcap -b -5
Error: bytes to strip must be greater than 0, got: -5
```

```
$ ./pcap_stripper -i input.pcap -b 100000
Error: bytes to strip is unusually large (100000). Maximum packet size is typically 65536 bytes
```

### Input File Errors:
```
$ ./pcap_stripper -i nonexistent.pcap -b 14
Error: input file does not exist: nonexistent.pcap
```

```
$ ./pcap_stripper -i /some/directory -b 14
Error: input path is a directory, not a file: /some/directory
```

```
$ ./pcap_stripper -i file.txt -b 14
Error: input file must have .pcap or .cap extension, got: .txt
```

### Output File Errors:
```
$ ./pcap_stripper -i input.pcap -o output.pcap -b 14
Error: output file already exists: output.pcap (use -f/--force to overwrite)
```

```
$ ./pcap_stripper -i input.pcap -o input.pcap -b 14
Error: input and output files cannot be the same: /home/user/input.pcap
```

### Invalid Position:
```
$ ./pcap_stripper -i input.pcap -b 14 -p middle
Error: invalid position 'middle': must be 'beginning' or 'end'
```

## Successful Execution Examples

### Normal Mode:
```
$ ./pcap_stripper -i capture.pcap -b 14 -p beginning
✓ Processing completed successfully!

Output file:            capture_stripped.pcap
Total packets:          1542
Successfully processed: 1538
Skipped (too small):    4
Total bytes stripped:   21532
```

### Verbose Mode:
```
$ ./pcap_stripper -i capture.pcap -b 14 -p beginning -v
Input file:      capture.pcap
Output file:     capture_stripped.pcap
Bytes to strip:  14
Strip position:  beginning

Processing...
Input file size: 2.45 MB
Processed: 1000 packets (skipped: 2)
Processed: 2000 packets (skipped: 4)

==================================================
✓ Processing completed successfully!

Output file:            capture_stripped.pcap
Total packets:          2156
Successfully processed: 2152
Skipped (too small):    4
Total bytes stripped:   30128

⚠ Warning: 4 packet(s) were smaller than 14 bytes and were skipped
```

## Position Aliases

All of these are valid and equivalent:

### For Beginning:
- `-p beginning`
- `-p begin`
- `-p start`
- `-p head`
- `-p front`

### For End:
- `-p end`
- `-p tail`
- `-p back`

## Quick Reference Card

```
┌─────────────────────────────────────────────────────────────┐
│              PCAP STRIPPER QUICK REFERENCE                  │
├─────────────────────────────────────────────────────────────┤
│ Basic Usage:                                                │
│   pcap_stripper -i INPUT.pcap -b BYTES [-p POSITION]       │
│                                                             │
│ Required Flags:                                             │
│   -i, --input     Input PCAP/CAP file                      │
│   -b, --bytes     Number of bytes to strip (> 0)           │
│                                                             │
│ Optional Flags:                                             │
│   -o, --output    Output file (auto-generated if omitted)  │
│   -p, --position  beginning|end (default: beginning)       │
│   -v, --verbose   Show detailed progress                   │
│   -f, --force     Overwrite existing output                │
│   -h, --help      Show help message                        │
│                                                             │
│ Common Use Cases:                                           │
│   Remove Ethernet:  -b 14 -p beginning                     │
│   Remove FCS:       -b 4 -p end                            │
│   Remove VLAN tag:  -b 4 -p beginning (after Ethernet)     │
│                                                             │
│ Tips:                                                       │
│   • Use -v to see progress on large files                  │
│   • Output auto-named: input_stripped.pcap                 │
│   • Tool skips packets smaller than strip size             │
│   • Corrupted input will produce clear error messages      │
└─────────────────────────────────────────────────────────────┘
```
