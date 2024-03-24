package main

import (
    "fmt"
    "io/ioutil"
    "log"
    "os"
    "strings"

    "github.com/google/gopacket"
    "github.com/google/gopacket/layers"
    "github.com/google/gopacket/pcapgo"
)

func stripBytesFromBeginning(inputFile, outputFile string, numBytes int) error {
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

        // Strip the specified number of bytes from the beginning of each packet
        strippedPacket := data[numBytes:]

        // Write the stripped packet to the output PCAP file
        writer.WritePacket(ci, strippedPacket)
    }

    fmt.Printf("Stripped %d bytes from the beginning of each packet and saved to %s\n", numBytes, outputFile)
    return nil
}

func stripBytesFromEnd(inputFile, outputFile string, numBytes int) error {
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

        // Strip the specified number of bytes from the end of each packet
        strippedPacket := data[:len(data)-numBytes]

        // Write the stripped packet to the output PCAP file
        writer.WritePacket(ci, strippedPacket)
    }

    fmt.Printf("Stripped %d bytes from the end of each packet and saved to %s\n", numBytes, outputFile)
    return nil
}

func stripUntilSecondEthernet(inputFile, outputFile string) error {
    // ... (code remains unchanged)
}

func filterPackets(inputFile, outputFile, displayFilter string) error {
    // ... (code remains unchanged)
}

func matchFilter(packet gopacket.Packet, filter string) bool {
    // ... (code remains unchanged)
}

func main() {
    if len(os.Args) < 2 {
        log.Fatal("Please provide the input PCAP file as an argument")
    }

    inputFile := os.Args[1]
    outputFile := "output.pcap"
    var displayFilter string

    fmt.Println("Select an option:")
    fmt.Println("1. Strip X bytes from the beginning of each packet")
    fmt.Println("2. Strip X bytes from the end of each packet")
    fmt.Println("3. Strip outer headers until the second Ethernet header for all packets")
    fmt.Println("4. Filter packets based on a Wireshark display filter")

    var choice string
    fmt.Print("Enter your choice (1, 2, 3, or 4): ")
    fmt.Scanln(&choice)

    var err error

    switch choice {
    case "1":
        var numBytes int
        fmt.Print("Enter the number of bytes to strip from the beginning of each packet: ")
        fmt.Scanln(&numBytes)
        err = stripBytesFromBeginning(inputFile, outputFile, numBytes)
    case "2":
        var numBytes int
        fmt.Print("Enter the number of bytes to strip from the end of each packet: ")
        fmt.Scanln(&numBytes)
        err = stripBytesFromEnd(inputFile, outputFile, numBytes)
    case "3":
        err = stripUntilSecondEthernet(inputFile, outputFile)
    case "4":
        fmt.Print("Enter the Wireshark display filter: ")
        fmt.Scanln(&displayFilter)
        err = filterPackets(inputFile, outputFile, displayFilter)
    default:
        log.Fatal("Invalid choice. Exiting.")
    }

    if err != nil {
        log.Fatal(err)
    }
}