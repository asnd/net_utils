import argparse
import pyshark

def strip_bytes_from_packets(input_file, output_file, num_bytes):
    capture = pyshark.FileCapture(input_file)
    output_packets = []

    for packet in capture:
        if 'eth' in packet:
            stripped_packet = packet.eth.raw_mode[:num_bytes] + packet.eth.raw_mode[num_bytes+num_bytes:]
            output_packets.append(stripped_packet)

    with open(output_file, 'wb') as f:
        for packet in output_packets:
            f.write(packet)

    print(f"Stripped {num_bytes} bytes from each packet and saved to {output_file}")

def strip_until_second_ethernet(input_file, output_file):
    capture = pyshark.FileCapture(input_file)
    output_packets = []

    for packet in capture:
        if 'eth' in packet:
            layers = packet.layers
            if len(layers) >= 2 and layers[1].layer_name == 'eth':
                stripped_packet = packet.eth.raw_mode[layers[1].start:]
            else:
                stripped_packet = packet.eth.raw_mode
            output_packets.append(stripped_packet)

    with open(output_file, 'wb') as f:
        for packet in output_packets:
            f.write(packet)

    print(f"Stripped bytes until the second Ethernet header for all packets and saved to {output_file}")

def filter_packets(input_file, output_file, display_filter):
    capture = pyshark.FileCapture(input_file, display_filter=display_filter)
    output_packets = []

    for packet in capture:
        output_packets.append(packet.get_raw_packet())

    with open(output_file, 'wb') as f:
        for packet in output_packets:
            f.write(packet)

    print(f"Filtered packets based on the display filter '{display_filter}' and saved to {output_file}")

# Main program
if __name__ == "__main__":
    parser = argparse.ArgumentParser(description="PCAP Packet Stripper and Filter")
    parser.add_argument("input_file", help="Path to the input PCAP file")
    parser.add_argument("-o", "--output", default="output.pcap", help="Path to the output PCAP file")
    parser.add_argument("-f", "--filter", help="Wireshark display filter")

    print("Select an option:")
    print("1. Strip X bytes from each packet")
    print("2. Strip bytes until the second Ethernet header for all packets")
    print("3. Filter packets based on a Wireshark display filter")

    choice = input("Enter your choice (1, 2, or 3): ")

    args = parser.parse_args()
    input_file = args.input_file
    output_file = args.output
    display_filter = args.filter

    if choice == '1':
        num_bytes = int(input("Enter the number of bytes to strip from each packet: "))
        strip_bytes_from_packets(input_file, output_file, num_bytes)
    elif choice == '2':
        strip_until_second_ethernet(input_file, output_file)
    elif choice == '3':
        if display_filter:
            filter_packets(input_file, output_file, display_filter)
        else:
            print("No display filter provided. Please use the -f or --filter option to specify a Wireshark display filter.")
    else:
        print("Invalid choice. Exiting.")