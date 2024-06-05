package main

import (
	"fmt"
	"os"
	"strconv"

	"https://github.com/fyne-io/fyne/"
	"https://github.com/fyne-io/fyne/app"
	"https://github.com/fyne-io/fyne/container"
	"https://github.com/fyne-io/fyne/dialog"
	"https://github.com/fyne-io/fyne/widget"
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"github.com/google/gopacket/pcapgo"
)

func main() {
	a := app.New()
	w := a.NewWindow("PCAP Stripper")

	// UI Elements
	filePathEntry := widget.NewEntry()
	filePathEntry.SetPlaceHolder("Select .cap or .pcap file")

	selectFileButton := widget.NewButton("Select File", func() {
		fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err == nil && reader != nil {
				filePathEntry.SetText(reader.URI().Path())
			}
		}, w)
		fd.SetFilter(dialog.NewExtensionFileFilter([]string{".cap", ".pcap"}))
		fd.Show()
	})

	bytesEntry := widget.NewEntry()
	bytesEntry.SetPlaceHolder("Enter number of bytes to strip")

	positionSelect := widget.NewSelect([]string{"Beginning", "End"}, func(value string) {})

	startButton := widget.NewButton("Start", func() {
		filePath := filePathEntry.Text
		bytesToStrip, err := strconv.Atoi(bytesEntry.Text)
		if err != nil {
			dialog.ShowError(fmt.Errorf("Invalid number of bytes"), w)
			return
		}
		position := positionSelect.Selected

		err = stripBytes(filePath, bytesToStrip, position)
		if err != nil {
			dialog.ShowError(err, w)
		} else {
			dialog.ShowInformation("Success", "File processed successfully", w)
		}
	})

	// Layout
	form := container.NewVBox(
		filePathEntry,
		selectFileButton,
		bytesEntry,
		positionSelect,
		startButton,
	)

	w.SetContent(form)
	w.ShowAndRun()
}

func stripBytes(filePath string, bytesToStrip int, position string) error {
	handle, err := pcap.OpenOffline(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer handle.Close()

	newFilePath := filePath + "_stripped.pcap"
	f, err := os.Create(newFilePath)
	if err != nil {
		return fmt.Errorf("failed to create new file: %w", err)
	}
	defer f.Close()

	w := pcapgo.NewWriter(f)
	w.WriteFileHeader(65536, handle.LinkType())

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		packetData := packet.Data()
		var newPacketData []byte

		if position == "Beginning" {
			if bytesToStrip >= len(packetData) {
				return fmt.Errorf("number of bytes to strip exceeds packet size")
			}
			newPacketData = packetData[bytesToStrip:]
		} else if position == "End" {
			if bytesToStrip >= len(packetData) {
				return fmt.Errorf("number of bytes to strip exceeds packet size")
			}
			newPacketData = packetData[:len(packetData)-bytesToStrip]
		} else {
			return fmt.Errorf("invalid position: %s", position)
		}

		w.WritePacket(packet.Metadata().CaptureInfo, newPacketData)
	}

	return nil
}
