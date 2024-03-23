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

def strip_until_second_ethernet(input_file, output_file, num_bytes):
    capture = pyshark.FileCapture(input_file)
    output_packets = []
    is_first_packet = True

    for packet in capture:
        if 'eth' in packet:
            if is_first_packet:
                stripped_packet = packet.eth.raw_mode[num_bytes:]
                is_first_packet = False
            else:
                stripped_packet = packet.eth.raw_mode
            output_packets.append(stripped_packet)

    with open(output_file, 'wb') as f:
        for packet in output_packets:
            f.write(packet)

    print(f"Stripped {num_bytes} bytes from the first packet until the second Ethernet header and saved to {output_file}")

# Main program
input_file = 'input.pcap'
output_file = 'output.pcap'

print("Select an option:")
print("1. Strip X bytes from each packet")
print("2. Strip Y bytes from the first packet until the second Ethernet header")

choice = input("Enter your choice (1 or 2): ")

if choice == '1':
    num_bytes = int(input("Enter the number of bytes to strip from each packet: "))
    strip_bytes_from_packets(input_file, output_file, num_bytes)
elif choice == '2':
    num_bytes = int(input("Enter the number of bytes to strip from the first packet until the second Ethernet header: "))
    strip_until_second_ethernet(input_file, output_file, num_bytes)
else:
    print("Invalid choice. Exiting.")