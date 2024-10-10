package main

import (
    "encoding/hex"
    "flag"
    "fmt"
    "os"

    "github.com/google/gopacket"
    "github.com/google/gopacket/pcap"
    "github.com/google/gopacket/pcapgo"
)

func stripBytes(packetData []byte, numBytes int, fromBeginning bool) []byte {
    if fromBeginning {
        if len(packetData) > numBytes {
            return packetData[numBytes:] // Strip from the beginning
        }
        return []byte{} // Packet too short, strip everything
    } else {
        if len(packetData) > numBytes {
            return packetData[:len(packetData)-numBytes] // Strip from the end
        }
        return []byte{} // Packet too short, strip everything
    }
}

func main() {
    // Command-line arguments
    inputFile := flag.String("input", "", "Input PCAP file path")
    outputFile := flag.String("output", "", "Output PCAP file path")
    numBytes := flag.Int("bytes", 0, "Number of bytes to strip")
    direction := flag.String("direction", "beginning", "Strip from 'beginning' or 'end'")
    flag.Parse()

    if *inputFile == "" || *outputFile == "" || *numBytes <= 0 {
        fmt.Println("Please provide valid input, output file paths, and a positive number of bytes to strip.")
        flag.Usage()
        return
    }

    fromBeginning := *direction == "beginning"

    // Open the input PCAP file
    handle, err := pcap.OpenOffline(*inputFile)
    if err != nil {
        fmt.Printf("Error opening input file: %v\n", err)
        return
    }
    defer handle.Close()

    // Create the output PCAP file
    f, err := os.Create(*outputFile)
    if err != nil {
        fmt.Printf("Error creating output file: %v\n", err)
        return
    }
    defer f.Close()

    // Create a new pcap writer
    writer := pcapgo.NewWriter(f)
    writer.WriteFileHeader(65535, handle.LinkType())

    // Process packets
    packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
    var strippedPackets [][]byte
    for packet := range packetSource.Packets() {
        strippedPacket := stripBytes(packet.Data(), *numBytes, fromBeginning)
        if len(strippedPacket) > 0 { // Avoid writing empty packets
            strippedPackets = append(strippedPackets, strippedPacket)
            writer.WritePacket(gopacket.CaptureInfo{
                Timestamp:      packet.Metadata().CaptureInfo.Timestamp,
                Length:         len(strippedPacket),
                CaptureLength:  len(strippedPacket),
            }, strippedPacket)
        }
    }

    // Output the new file size
    newFileInfo, err := os.Stat(*outputFile)
    if err != nil {
        fmt.Printf("Error getting new file info: %v\n", err)
        return
    }
    fmt.Printf("New file size: %d bytes\n", newFileInfo.Size())

    // Display the first 20 bytes of the first 5 packets
    fmt.Println("\nFirst 20 bytes of the first 5 packets:")
    for i, packet := range strippedPackets[:5] {
        if len(packet) >= 20 {
            fmt.Printf("Packet %d: %s\n", i+1, hex.EncodeToString(packet[:20]))
        } else {
            fmt.Printf("Packet %d: %s (less than 20 bytes)\n", i+1, hex.EncodeToString(packet))
        }
    }
}
