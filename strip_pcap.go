package main

import (
    "encoding/hex"
    "flag"
    "fmt"
    "log"
    "os"

    "github.com/google/gopacket/pcapgo"
)

func main() {
    // Command-line arguments
    trimBytes := flag.Int("x", 0, "Number of bytes to trim from each packet")
    inputFile := flag.String("i", "", "Input pcap file")
    outputFile := flag.String("o", "", "Output pcap file")
    trimFrom := flag.String("from", "end", "Trim from 'beginning' or 'end' of each packet")
    flag.Parse()

    if *trimBytes <= 0 {
        log.Fatal("Please specify a positive number of bytes to trim using -x")
    }
    if *inputFile == "" || *outputFile == "" {
        log.Fatal("Please specify input and output pcap files using -i and -o")
    }
    if *trimFrom != "beginning" && *trimFrom != "end" {
        log.Fatal("Please specify 'beginning' or 'end' for the --from parameter")
    }

    err := trimPackets(*inputFile, *outputFile, *trimBytes, *trimFrom)
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
