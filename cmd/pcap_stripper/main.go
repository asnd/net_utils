package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"github.com/google/gopacket/pcapgo"
	"github.com/spf13/cobra"
)

// Build information (set via ldflags during build)
var (
	Version   = "dev"
	BuildTime = "unknown"
	GitCommit = "unknown"
)

type ProcessingStats struct {
	TotalPackets    int
	ProcessedOK     int
	SkippedTooSmall int
	BytesStripped   int64
}

var (
	inputFile    string
	outputFile   string
	bytesToStrip int
	position     string
	verbose      bool
	force        bool
	showVersion  bool
)

func main() {
	var rootCmd = &cobra.Command{
		Use:     "pcap_stripper",
		Version: fmt.Sprintf("%s (commit: %s, built: %s)", Version, GitCommit, BuildTime),
		Short:   "PCAP Stripper - Strip bytes from packet payloads in PCAP files",
		Long: `PCAP Stripper is a command-line tool for removing a specified number of bytes
from packet payloads in PCAP/CAP files. This is useful for sanitizing captures,
removing headers, or preparing files for analysis.

The tool can strip bytes from either the beginning or end of each packet's payload,
preserving the PCAP file structure and metadata while modifying only the packet data.`,
		Example: `  # Strip 14 bytes from the beginning of each packet (e.g., Ethernet header)
  pcap_stripper -i input.pcap -o output.pcap -b 14 -p beginning

  # Strip 4 bytes from the end of each packet
  pcap_stripper -i capture.cap -o sanitized.pcap -b 4 -p end

  # Use auto-generated output filename with verbose output
  pcap_stripper -i input.pcap -b 14 -p beginning -v

  # Force overwrite existing output file
  pcap_stripper -i input.pcap -o output.pcap -b 14 -p beginning -f`,
		RunE: runStripper,
	}

	// Define flags
	rootCmd.Flags().StringVarP(&inputFile, "input", "i", "", "Input PCAP/CAP file (required)")
	rootCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output PCAP file (auto-generated if not specified)")
	rootCmd.Flags().IntVarP(&bytesToStrip, "bytes", "b", 0, "Number of bytes to strip from each packet (required)")
	rootCmd.Flags().StringVarP(&position, "position", "p", "beginning", "Position to strip from: 'beginning' or 'end'")
	rootCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")
	rootCmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite of existing output file")

	// Mark required flags
	rootCmd.MarkFlagRequired("input")
	rootCmd.MarkFlagRequired("bytes")

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func runStripper(cmd *cobra.Command, args []string) error {
	// Validate input file
	if err := validateInputFile(inputFile); err != nil {
		return err
	}

	// Validate bytes to strip
	if err := validateBytesToStrip(bytesToStrip); err != nil {
		return err
	}

	// Validate and normalize position
	pos, err := validatePosition(position)
	if err != nil {
		return err
	}
	position = pos

	// Determine output file path
	outPath, err := determineOutputPath(inputFile, outputFile)
	if err != nil {
		return err
	}

	// Check if input and output are the same file
	if err := checkSameFile(inputFile, outPath); err != nil {
		return err
	}

	// Check if output file exists
	if err := checkOutputExists(outPath, force); err != nil {
		return err
	}

	if verbose {
		fmt.Printf("Input file:      %s\n", inputFile)
		fmt.Printf("Output file:     %s\n", outPath)
		fmt.Printf("Bytes to strip:  %d\n", bytesToStrip)
		fmt.Printf("Strip position:  %s\n", position)
		fmt.Println("\nProcessing...")
	}

	// Process the file
	stats, err := stripBytes(inputFile, outPath, bytesToStrip, position, verbose)
	if err != nil {
		return fmt.Errorf("processing failed: %w", err)
	}

	// Display results
	if verbose {
		fmt.Println("\n" + strings.Repeat("=", 50))
	}
	fmt.Printf("✓ Processing completed successfully!\n\n")
	fmt.Printf("Output file:            %s\n", outPath)
	fmt.Printf("Total packets:          %d\n", stats.TotalPackets)
	fmt.Printf("Successfully processed: %d\n", stats.ProcessedOK)
	fmt.Printf("Skipped (too small):    %d\n", stats.SkippedTooSmall)
	fmt.Printf("Total bytes stripped:   %d\n", stats.BytesStripped)

	if stats.SkippedTooSmall > 0 {
		fmt.Printf("\n⚠ Warning: %d packet(s) were smaller than %d bytes and were skipped\n",
			stats.SkippedTooSmall, bytesToStrip)
	}

	return nil
}

// validateInputFile checks if the input file exists and is readable
func validateInputFile(path string) error {
	if path == "" {
		return fmt.Errorf("input file is required")
	}

	// Clean the path
	path = strings.TrimSpace(path)

	// Check if file exists
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return fmt.Errorf("input file does not exist: %s", path)
	}
	if err != nil {
		return fmt.Errorf("cannot access input file: %w", err)
	}

	// Check if it's a directory
	if info.IsDir() {
		return fmt.Errorf("input path is a directory, not a file: %s", path)
	}

	// Check if file is readable
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("cannot read input file (permission denied?): %w", err)
	}
	file.Close()

	// Validate file extension
	ext := strings.ToLower(filepath.Ext(path))
	if ext != ".pcap" && ext != ".cap" {
		return fmt.Errorf("input file must have .pcap or .cap extension, got: %s", ext)
	}

	return nil
}

// validateBytesToStrip checks if the bytes parameter is valid
func validateBytesToStrip(bytes int) error {
	if bytes <= 0 {
		return fmt.Errorf("bytes to strip must be greater than 0, got: %d", bytes)
	}

	if bytes > 65536 {
		return fmt.Errorf("bytes to strip is unusually large (%d). Maximum packet size is typically 65536 bytes", bytes)
	}

	return nil
}

// validatePosition validates and normalizes the position parameter
func validatePosition(pos string) (string, error) {
	normalized := strings.ToLower(strings.TrimSpace(pos))

	switch normalized {
	case "beginning", "begin", "start", "head", "front":
		return "beginning", nil
	case "end", "tail", "back":
		return "end", nil
	default:
		return "", fmt.Errorf("invalid position '%s': must be 'beginning' or 'end'", pos)
	}
}

// determineOutputPath generates output path if not specified
func determineOutputPath(input, output string) (string, error) {
	if output != "" {
		return strings.TrimSpace(output), nil
	}

	// Auto-generate output filename
	ext := filepath.Ext(input)
	base := strings.TrimSuffix(input, ext)
	return fmt.Sprintf("%s_stripped.pcap", base), nil
}

// checkSameFile verifies input and output are not the same file
func checkSameFile(input, output string) error {
	absInput, err := filepath.Abs(input)
	if err != nil {
		return fmt.Errorf("cannot resolve input path: %w", err)
	}

	absOutput, err := filepath.Abs(output)
	if err != nil {
		return fmt.Errorf("cannot resolve output path: %w", err)
	}

	if absInput == absOutput {
		return fmt.Errorf("input and output files cannot be the same: %s", absInput)
	}

	return nil
}

// checkOutputExists checks if output file exists and handles overwrite logic
func checkOutputExists(output string, force bool) error {
	if _, err := os.Stat(output); err == nil {
		if !force {
			return fmt.Errorf("output file already exists: %s (use -f/--force to overwrite)", output)
		}
	}
	return nil
}

// stripBytes performs the actual byte stripping operation
func stripBytes(inputPath, outputPath string, bytesToStrip int, position string, verbose bool) (ProcessingStats, error) {
	var stats ProcessingStats

	// Open input PCAP file
	handle, err := pcap.OpenOffline(inputPath)
	if err != nil {
		return stats, fmt.Errorf("failed to open input file: %w", err)
	}
	defer handle.Close()

	// Get file info for size estimation
	fileInfo, _ := os.Stat(inputPath)
	if verbose && fileInfo != nil {
		fmt.Printf("Input file size: %.2f MB\n", float64(fileInfo.Size())/(1024*1024))
	}

	// Create output file
	f, err := os.Create(outputPath)
	if err != nil {
		return stats, fmt.Errorf("failed to create output file: %w", err)
	}

	// Use buffered writer for better performance
	bufWriter := bufio.NewWriterSize(f, 1024*1024) // 1MB buffer

	// Flag to track if we should clean up on error
	var processingError error
	defer func() {
		bufWriter.Flush()
		f.Close()

		if processingError != nil {
			// Clean up partial output file on error
			if verbose {
				fmt.Printf("\nCleaning up incomplete output file: %s\n", outputPath)
			}
			os.Remove(outputPath)
		}
	}()

	// Write PCAP header
	w := pcapgo.NewWriter(bufWriter)
	if err := w.WriteFileHeader(65536, handle.LinkType()); err != nil {
		processingError = fmt.Errorf("failed to write PCAP header: %w", err)
		return stats, processingError
	}

	// Process packets
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	for packet := range packetSource.Packets() {
		stats.TotalPackets++

		packetData := packet.Data()

		// Handle empty packets
		if len(packetData) == 0 {
			stats.SkippedTooSmall++
			if verbose && stats.TotalPackets%1000 == 0 {
				fmt.Printf("\rProcessed: %d packets (skipped: %d)", stats.TotalPackets, stats.SkippedTooSmall)
			}
			continue
		}

		var newPacketData []byte

		if position == "beginning" {
			if bytesToStrip >= len(packetData) {
				// Skip packets that are too small
				stats.SkippedTooSmall++
				if verbose && stats.TotalPackets%1000 == 0 {
					fmt.Printf("\rProcessed: %d packets (skipped: %d)", stats.TotalPackets, stats.SkippedTooSmall)
				}
				continue
			}
			newPacketData = packetData[bytesToStrip:]
		} else { // position == "end"
			if bytesToStrip >= len(packetData) {
				// Skip packets that are too small
				stats.SkippedTooSmall++
				if verbose && stats.TotalPackets%1000 == 0 {
					fmt.Printf("\rProcessed: %d packets (skipped: %d)", stats.TotalPackets, stats.SkippedTooSmall)
				}
				continue
			}
			newPacketData = packetData[:len(packetData)-bytesToStrip]
		}

		// Write modified packet
		if err := w.WritePacket(packet.Metadata().CaptureInfo, newPacketData); err != nil {
			processingError = fmt.Errorf("failed to write packet %d: %w", stats.TotalPackets, err)
			return stats, processingError
		}

		stats.ProcessedOK++
		stats.BytesStripped += int64(bytesToStrip)

		// Update progress in verbose mode
		if verbose && stats.TotalPackets%1000 == 0 {
			fmt.Printf("\rProcessed: %d packets (skipped: %d)", stats.TotalPackets, stats.SkippedTooSmall)
		}
	}

	// Ensure all data is flushed
	if err := bufWriter.Flush(); err != nil {
		processingError = fmt.Errorf("failed to flush output: %w", err)
		return stats, processingError
	}

	return stats, nil
}
