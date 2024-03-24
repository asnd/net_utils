package main

import (
    "fmt"
    "io/ioutil"
    "log"
    "os"

    "github.com/google/gopacket"
    "github.com/google/gopacket/layers"
    "github.com/google/gopacket/pcapgo"
)

func stripBytesFromPackets(inputFile, outputFile string, numBytes int) error {
    // Open the input PCAP file
    input, err := os.Open(inputFile)
    if err != nil {
        return err
    }
    defer input.Close()

    // Create a new PCAP file reader
    reader, err := pcapgo.NewReader(input)
    if err != nil {
        return err
    }

    // Create the output PCAP file
    output, err := os.Create(outputFile)
    if err != nil {
        return err
    }
    defer output.Close()

    // Create a new PCAP file writer
    writer := pcapgo.NewWriter(output)
    writer.WriteFileHeader(1024, layers.LinkTypeEthernet)

    // Iterate over each packet in the input PCAP file
    for {
        data, ci, err := reader.ReadPacketData()
        if err != nil {
            if err.Error() == "EOF" {
                break
            }
            return err
        }

        // Strip the specified number of bytes from each packet
        strippedPacket := append(data[:numBytes], data[numBytes+numBytes:]...)

        // Write the stripped packet to the output PCAP file
        writer.WritePacket(ci, strippedPacket)
    }

    fmt.Printf("Stripped %d bytes from each packet and saved to %s\n", numBytes, outputFile)
    return nil
}

func stripUntilSecondEthernet(inputFile, outputFile string) error {
    // Open the input PCAP file
    input, err := os.Open(inputFile)
    if err != nil {
        return err
    }
    defer input.Close()

    // Create a new PCAP file reader
    reader, err := pcapgo.NewReader(input)
    if err != nil {
        return err
    }

    // Create the output PCAP file
    output, err := os.Create(outputFile)
    if err != nil {
        return err
    }
    defer output.Close()

    // Create a new PCAP file writer
    writer := pcapgo.NewWriter(output)
    writer.WriteFileHeader(1024, layers.LinkTypeEthernet)

    // Iterate over each packet in the input PCAP file
    for {
        data, ci, err := reader.ReadPacketData()
        if err != nil {
            if err.Error() == "EOF" {
                break
            }
            return err
        }

        packet := gopacket.NewPacket(data, layers.LayerTypeEthernet, gopacket.Default)
        ethernetLayers := packet.LayersByType(layers.LayerTypeEthernet)

        if len(ethernetLayers) >= 2 {
            // Find the offset of the second Ethernet header
            secondEthernetOffset := ethernetLayers[1].(*layers.Ethernet).Contents[0] & 0xf

            // Strip the outer headers before the second Ethernet header
            strippedPacket := data[secondEthernetOffset:]
            writer.WritePacket(ci, strippedPacket)
        } else {
            writer.WritePacket(ci, data)
        }
    }

    fmt.Printf("Stripped outer headers until the second Ethernet header for all packets and saved to %s\n", outputFile)
    return nil
}

func main() {
    if len(os.Args) < 2 {
        log.Fatal("Please provide the input PCAP f