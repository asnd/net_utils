// trimPackets trims a specified number of bytes from each packet in the input pcap file
// and writes the trimmed packets to the output pcap file. The trimming can be done from
// the beginning or the end of each packet based on the provided parameter.
//
// Parameters:
//   - inputFile: Path to the input pcap file.
//   - outputFile: Path to the output pcap file.
//   - trimBytes: Number of bytes to trim from each packet (must be positive).
//   - trimFrom: Specifies whether to trim from the 'beginning' or 'end' of each packet.
//
// Returns:
//   - error: An error object if an error occurs during the trimming process, otherwise nil.
//
// The function performs the following steps:
//   1. Opens the input pcap file for reading.
//   2. Creates a pcap reader for the input file.
//   3. Creates the output pcap file for writing.
//   4. Writes the pcap file header to the output file.
//   5. Iterates through each packet in the input file, trims the specified number of bytes
//      from the beginning or end of the packet, and writes the trimmed packet to the output file.
//   6. Displays the new size of the output pcap file.
//   7. Displays the first 20 bytes of the first 5 packets in the output file.
//
// Note: If the number of bytes to trim is greater than or equal to the packet size, the packet
//       is skipped and not written to the output file.
package main

import (
    "encoding/hex"
    "flag"
    "fmt"
    "log"
    "os"

    "github.com/google/gopacket/pcapgo"
)

type trimOptions struct {
    trimBytes int
    inputFile string
    outputFile string
    trimFrom string
}

func (opts *trimOptions) validate() {
    if opts.trimBytes <= 0 {
        log.Fatalf("Please specify a positive number of bytes to trim using -x")
    }
    if opts.inputFile == "" || opts.outputFile == "" {
        log.Fatalf("Please specify input and output pcap files using -i and -o")
    }
    if opts.trimFrom != "beginning" && opts.trimFrom != "end" {
        log.Fatalf("Please specify 'beginning' or 'end' for the --from parameter")
    }
}

// main is the entry point for the pcap trimming utility. It processes command-line arguments
// to determine the number of bytes to trim from each packet, the input and output pcap files,
// and whether to trim from the beginning or end of each packet. It validates the arguments
// and calls the trimPackets function to perform the trimming operation.
//
// Usage:
//   -x int
//         Number of bytes to trim from each packet (must be positive)
//   -i string
//         Input pcap file (required)
//   -o string
//         Output pcap file (required)
//   --from string
//         Trim from 'beginning' or 'end' of each packet (default "end")
func main() {
    opts := &trimOptions{}
    flag.IntVar(&opts.trimBytes, "x", 0, "Number of bytes to trim from each packet")
    flag.StringVar(&opts.inputFile, "i", "", "Input pcap file")
    flag.StringVar(&opts.outputFile, "o", "", "Output pcap file")
    flag.StringVar(&opts.trimFrom, "from", "end", "Trim from 'beginning' or 'end' of each packet")
    flag.Parse()

    opts.validate()

    err := trimPackets(opts.inputFile, opts.outputFile, opts.trimBytes, opts.trimFrom)
    if err != nil {
        log.Fatalf("Error: %v", err)
    }
}

func trimPackets(inputFile, outputFile string, trimBytes int, trimFrom string) error {
    inputHandle, err := os.Open(inputFile)
    if err != nil {
        return fmt.Errorf("failed to open input file: %v", err)
    }
    defer inputHandle.Close()

    reader, err := pcapgo.NewReader(inputHandle)
    if err != nil {
        return fmt.Errorf("failed to create pcap reader: %v", err)
    }

    outputHandle, err := os.Create(outputFile)
    if err != nil {
        return fmt.Errorf("failed to create output file: %v", err)
    }
    defer outputHandle.Close()

    writer := pcapgo.NewWriter(outputHandle)
    err = writer.WriteFileHeader(reader.Snaplen(), reader.LinkType())
    if err != nil {
        return fmt.Errorf("failed to write file header: %v", err)
    }

    packetCount := 0

    for {
        data, ci, err := reader.ReadPacketData()
        if err != nil {
            if err.Error() == "EOF" {
                break
            }
            return fmt.Errorf("error reading packet: %v", err)
        }

        packetCount++
        pktLen := len(data)

        if trimBytes >= pktLen {
            // Skip packet if trim size is greater than or equal to packet size
            continue
        }

        var trimmedData []byte
        if trimFrom == "beginning" {
            trimmedData = data[trimBytes:]
        } else {
            trimmedData = data[:pktLen-trimBytes]
        }

        ci.CaptureLength = len(trimmedData)
        ci.Length = len(trimmedData)

        err = writer.WritePacket(ci, trimmedData)
        if err != nil {
            return fmt.Errorf("failed to write packet: %v", err)
        }
    }

    // Display new file size
    fi, err := os.Stat(outputFile)
    if err != nil {
        return fmt.Errorf("failed to get file info: %v", err)
    }
    fmt.Printf("New pcap file size: %d bytes\n", fi.Size())

    // Display first 20 bytes of first 5 packets
    fmt.Println("\nFirst 20 bytes of first 5 packets:")
    outputHandle.Seek(0, 0) // Reset file pointer to the beginning
    reader, err = pcapgo.NewReader(outputHandle)
    if err != nil {
        return fmt.Errorf("failed to create pcap reader for output file: %v", err)
    }

    for i := 0; i < 5; i++ {
        data, _, err := reader.ReadPacketData()
        if err != nil {
            if err.Error() == "EOF" {
                break
            }
            return fmt.Errorf("error reading packet: %v", err)
        }
        first20Bytes := data
        if len(data) > 20 {
            first20Bytes = data[:20]
        }
        fmt.Printf("Packet %d: %s\n", i+1, hex.EncodeToString(first20Bytes))
    }

    return nil
}
