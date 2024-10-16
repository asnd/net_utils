package main

import (
    "bytes"
    "encoding/hex"
    "flag"
    "fmt"
    "log"
    "os"

    "github.com/google/gopacket"
    "github.com/google/gopacket/pcap"
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
    handle, err := pcap.OpenOffline(inputFile)
    if err != nil {
        return fmt.Errorf("failed to open input file: %v", err)
    }
    defer handle.Close()

    // Get file info for output pcap file
    f, err := os.Create(outputFile)
    if err != nil {
        return fmt.Errorf("failed to create output file: %v", err)
    }
    defer f.Close()

    w := pcapgo.NewWriter(f)
    err = w.WriteFileHeader(65536, handle.LinkType())
    if err != nil {
        return fmt.Errorf("failed to write file header: %v", err)
    }

    packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
    packetCount := 0

    for packet := range packetSource.Packets() {
        packetCount++
        data := packet.Data()
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

        ci := gopacket.CaptureInfo{
            Timestamp:      packet.Metadata().Timestamp,
            CaptureLength:  len(trimmedData),
            Length:         len(trimmedData),
            InterfaceIndex: packet.Metadata().InterfaceIndex,
        }

        err = w.WritePacket(ci, trimmedData)
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
    outputHandle, err := pcap.OpenOffline(outputFile)
    if err != nil {
        return fmt.Errorf("failed to open output file: %v", err)
    }
    defer outputHandle.Close()

    outputPacketSource := gopacket.NewPacketSource(outputHandle, outputHandle.LinkType())

    i := 0
    for packet := range outputPacketSource.Packets() {
        i++
        data := packet.Data()
        first20Bytes := data
        if len(data) > 20 {
            first20Bytes = data[:20]
        }
        fmt.Printf("Packet %d: %s\n", i, hex.EncodeToString(first20Bytes))
        if i >= 5 {
            break
        }
    }

    return nil
}
