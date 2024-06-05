package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
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
	// Read the file
	fileData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Strip bytes
	var newData []byte
	if position == "Beginning" {
		if bytesToStrip >= len(fileData) {
			return fmt.Errorf("number of bytes to strip exceeds file size")
		}
		newData = fileData[bytesToStrip:]
	} else if position == "End" {
		if bytesToStrip >= len(fileData) {
			return fmt.Errorf("number of bytes to strip exceeds file size")
		}
		newData = fileData[:len(fileData)-bytesToStrip]
	} else {
		return fmt.Errorf("invalid position: %s", position)
	}

	// Save the file
	err = ioutil.WriteFile(filePath, newData, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}
