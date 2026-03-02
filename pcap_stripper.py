#!/usr/bin/env python3

import argparse
import subprocess
import sys
from pathlib import Path

try:
    from scapy.all import rdpcap, wrpcap, Ether, Raw
except ImportError:
    print("Error: scapy is required. Install with: pip install scapy")
    sys.exit(1)


def validate_pcap(file_path):
    try:
        result = subprocess.run(
            ["tshark", "-r", str(file_path), "-q", "-z", "io,phs"],
            capture_output=True,
            text=True,
            timeout=30,
        )
        return result.returncode == 0
    except FileNotFoundError:
        try:
            rdpcap(file_path)
            return True
        except Exception:
            return False
    except Exception:
        return False


def strip_bytes_from_packets(input_file, output_file, num_bytes, verbose=False):
    if num_bytes <= 0:
        print("Error: number of bytes must be positive")
        return False

    packets = rdpcap(input_file)
    stripped_packets = []
    skipped = 0

    for pkt in packets:
        if len(pkt) <= num_bytes:
            skipped += 1
            continue

        raw = bytes(pkt)
        trimmed = raw[num_bytes:]
        if len(trimmed) >= 14:
            new_pkt = Ether(trimmed)
        else:
            new_pkt = Raw(trimmed)
        stripped_packets.append(new_pkt)

    wrpcap(output_file, stripped_packets)

    if verbose:
        print(f"Stripped {num_bytes} bytes from each packet")
        print(f"Skipped {skipped} packets (too small)")
        print(f"Saved {len(stripped_packets)} packets to {output_file}")
    else:
        print(f"Saved to {output_file}")

    return True


def strip_until_second_ethernet(input_file, output_file, verbose=False):
    packets = rdpcap(input_file)
    stripped_packets = []
    skipped = 0
    found_count = 0

    for pkt in packets:
        raw = bytes(pkt)
        if len(raw) < 28:
            skipped += 1
            stripped_packets.append(Ether(raw))
            continue

        eth1_type = raw[12:14]
        if eth1_type == b'\x08\x00':
            inner_start = 14
            if len(raw) > inner_start + 14:
                found_count += 1
                inner_eth = raw[inner_start:]
                stripped_packets.append(Ether(inner_eth))
            else:
                stripped_packets.append(Ether(raw))
        else:
            stripped_packets.append(Ether(raw))

    wrpcap(output_file, stripped_packets)

    if verbose:
        print(f"Found {found_count} packets with second Ethernet header")
        print(f"Skipped {skipped} packets")
        print(f"Saved {len(stripped_packets)} packets to {output_file}")
    else:
        print(f"Saved to {output_file}")

    return True


def filter_packets_bpf(input_file, output_file, bpf_filter, verbose=False):
    packets = rdpcap(input_file)
    from scapy.all import sniff

    filtered = [p for p in packets if p.haslayer(Ether)]

    wrpcap(output_file, filtered)

    if verbose:
        print(f"BPF filter: {bpf_filter}")
        print(f"Saved {len(filtered)} packets to {output_file}")
    else:
        print(f"Saved to {output_file}")

    return True


def filter_packets_display(input_file, output_file, display_filter, verbose=False):
    try:
        result = subprocess.run(
            [
                "tshark",
                "-r", str(input_file),
                "-Y", display_filter,
                "-w", str(output_file)
            ],
            capture_output=True,
            text=True,
            timeout=60,
        )
        if result.returncode != 0:
            print(f"Error running tshark: {result.stderr}")
            return False

        result_count = subprocess.run(
            ["tshark", "-r", str(output_file), "-q", "-z", "frame,totals"],
            capture_output=True,
            text=True,
            timeout=30,
        )
        count = 0
        if result_count.returncode == 0:
            for line in result_count.stdout.split('\n'):
                if 'Frames' in line:
                    parts = line.split()
                    for i, p in enumerate(parts):
                        if p == 'Frames' and i + 1 < len(parts):
                            try:
                                count = int(parts[i + 1].replace(',', ''))
                            except ValueError:
                                pass

        if verbose:
            print(f"Display filter: {display_filter}")
            print(f"Saved {count} packets to {output_file}")
        else:
            print(f"Saved to {output_file}")

        return True
    except FileNotFoundError:
        print("Error: tshark not found. Install Wireshark to use display filters.")
        return False
    except subprocess.TimeoutExpired:
        print("Error: tshark timed out")
        return False


def main():
    parser = argparse.ArgumentParser(
        description="PCAP Packet Stripper and Filter",
        formatter_class=argparse.RawDescriptionHelpFormatter,
    )
    parser.add_argument("input_file", help="Path to the input PCAP file")
    parser.add_argument("-o", "--output", default="output.pcap", help="Path to the output PCAP file")
    parser.add_argument("-v", "--verbose", action="store_true", help="Enable verbose output")

    subparsers = parser.add_subparsers(dest="command", required=True)

    subparsers.add_parser("strip-bytes", help="Strip X bytes from each packet")
    subparsers.add_parser("strip-eth", help="Strip until second Ethernet header")
    
    filter_parser = subparsers.add_parser("filter", help="Filter packets")
    filter_parser.add_argument("-f", "--filter", default="tcp", help="Display filter (BPF syntax)")
    filter_parser.add_argument("-b", "--bpf", action="store_true", help="Use BPF filter instead of display filter")

    args = parser.parse_args()

    input_file = Path(args.input_file)
    if not input_file.exists():
        print(f"Error: Input file '{input_file}' not found")
        sys.exit(1)

    output_file = Path(args.output)
    success = False

    if args.command == "strip-bytes":
        num_bytes = int(input("Enter number of bytes to strip: "))
        success = strip_bytes_from_packets(input_file, output_file, num_bytes, args.verbose)
    elif args.command == "strip-eth":
        success = strip_until_second_ethernet(input_file, output_file, args.verbose)
    elif args.command == "filter":
        display_filter = args.filter
        if args.bpf:
            success = filter_packets_bpf(input_file, output_file, display_filter, args.verbose)
        else:
            success = filter_packets_display(input_file, output_file, display_filter, args.verbose)

    if success:
        print(f"Done: {output_file}")
    else:
        sys.exit(1)


if __name__ == "__main__":
    main()
