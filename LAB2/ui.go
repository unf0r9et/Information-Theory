package main

import (
	"fmt"
	"image/color"
	"os"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func toBitStr(data []byte) string {
	var sb strings.Builder
	for _, b := range data {
		sb.WriteString(fmt.Sprintf("%08b", b))
	}
	return sb.String()
}

func StartUI() {
	myApp := app.NewWithID("com.unf0r9ecryptography_lfsrt.")
	myWindow := myApp.NewWindow(fmt.Sprintf("Потоковое шифрование (LFSR %d)", Degree))

	var filePathToRead string
	var resultBytes []byte 

	labelPathFromFile := widget.NewLabel("Файл для чтения не выбран")

	entryOriginal := widget.NewMultiLineEntry()
	entryOriginal.SetMinRowsVisible(5)
	entryOriginal.Wrapping = fyne.TextWrapWord
	
	entryKey := widget.NewMultiLineEntry()
	entryKey.SetMinRowsVisible(5)
	entryKey.Wrapping = fyne.TextWrapWord
	
	entryResult := widget.NewMultiLineEntry()
	entryResult.SetMinRowsVisible(5)
	entryResult.Wrapping = fyne.TextWrapWord

	keyEntry := widget.NewEntry()
	keyEntry.PlaceHolder = fmt.Sprintf("Введите %d бит (только 0 и 1)...", Degree)

	keyEntry.OnChanged = func(s string) {
		var filtered string
		for _, char := range s {
			if char == '0' || char == '1' {
				filtered += string(char)
			}
		}
		if len(filtered) > Degree {
			filtered = filtered[:Degree]
		}
		if s != filtered {
			keyEntry.SetText(filtered)
		}
	}

	btnSaveAs := widget.NewButton("Сохранить результат (Save As...)", func() {
		if len(resultBytes) == 0 { return }
		dialog.ShowFileSave(func(writer fyne.URIWriteCloser, err error) {
			if err != nil || writer == nil { return }
			defer writer.Close()
			writer.Write(resultBytes)
			dialog.ShowInformation("Успех", "Файл успешно сохранен!", myWindow)
		}, myWindow)
	})
	btnSaveAs.Disable()

	process := func() {
		if len(keyEntry.Text) != Degree {
			dialog.ShowInformation("ОШИБКА", fmt.Sprintf("Ключ должен состоять ровно из %d бит!", Degree), myWindow)
			return
		}
		if filePathToRead == "" {
			dialog.ShowInformation("ОШИБКА", "Файл для чтения не выбран!", myWindow)
			return
		}

		inputBytes, err := os.ReadFile(filePathToRead)
		if err != nil { return }

		lfsr := NewLFSR(keyEntry.Text, Taps)

		limitBits := 108
		limitBytes := 14 

		outputBytes, firstKeyBytes, lastKeyBytes := lfsr.ProcessData(inputBytes, limitBytes)
		resultBytes = outputBytes

		var uiOrig, uiKey, uiCrypt string

		if len(inputBytes)*8 <= limitBits*2 {
			uiOrig = toBitStr(inputBytes)
			uiKey = toBitStr(firstKeyBytes)
			uiCrypt = toBitStr(outputBytes)
		} else {
			firstOrigBytes := inputBytes[:limitBytes]
			lastOrigBytes := inputBytes[len(inputBytes)-limitBytes:]
			firstCryptBytes := outputBytes[:limitBytes]
			lastCryptBytes := outputBytes[len(outputBytes)-limitBytes:]

			exFirst := func(b []byte) string {
				s := toBitStr(b)
				if len(s) > limitBits { return s[:limitBits] }
				return s
			}
			exLast := func(b []byte) string {
				s := toBitStr(b)
				if len(s) > limitBits { return s[len(s)-limitBits:] }
				return s
			}

			uiOrig = exFirst(firstOrigBytes) + "\n\n... (вырезано) ...\n\n" + exLast(lastOrigBytes)
			uiKey = exFirst(firstKeyBytes) + "\n\n... (вырезано) ...\n\n" + exLast(lastKeyBytes)
			uiCrypt = exFirst(firstCryptBytes) + "\n\n... (вырезано) ...\n\n" + exLast(lastCryptBytes)
		}

		entryOriginal.SetText(uiOrig)
		entryKey.SetText(uiKey)
		entryResult.SetText(uiCrypt)

		btnSaveAs.Enable()
	}

	buttonEncrypt := widget.NewButton("Применить алгоритм", process)

	buttonClear := widget.NewButton("Очистить", func() {
		entryOriginal.SetText("")
		entryKey.SetText("")
		entryResult.SetText("")
		keyEntry.SetText("")
		filePathToRead = ""
		resultBytes = nil 
		labelPathFromFile.SetText("Файл для чтения не выбран")
		btnSaveAs.Disable() 
	})

	btnOpenFileToRead := widget.NewButton("Выбрать файл", func() {
		dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
			if reader != nil {
				filePathToRead = reader.URI().Path()
				labelPathFromFile.SetText("Выбран файл: " + filePathToRead)
				reader.Close()
			}
		}, myWindow)
	})

	spacer := func(height float32) fyne.CanvasObject {
		rect := canvas.NewRectangle(color.Transparent)
		rect.SetMinSize(fyne.NewSize(1, height))
		return rect
	}

	rowKey := container.NewBorder(nil, nil, widget.NewLabel(fmt.Sprintf("Стартовый ключ (%d бит): ", Degree)), nil, keyEntry)
	horizontalContent := container.NewVBox(btnOpenFileToRead, labelPathFromFile)
	line := canvas.NewRectangle(color.RGBA{R: 82, G: 82, B: 82, A: 255})
	line.SetMinSize(fyne.NewSize(1, 1))

	outputSection := container.NewVBox(
		widget.NewLabel("Исходный файл (первые и последние 108 бит):"),
		entryOriginal,
		widget.NewLabel("Сгенерированный ключевой поток:"),
		entryKey,
		widget.NewLabel("Результат (зашифровано/расшифровано):"),
		entryResult,
	)

	content := container.NewVBox(spacer(10), horizontalContent, spacer(10), line, spacer(10), rowKey, spacer(10), buttonEncrypt, spacer(10), outputSection, spacer(10), btnSaveAs, spacer(10), buttonClear, spacer(10))

	myWindow.SetContent(container.NewVScroll(content))
	myWindow.Resize(fyne.NewSize(700, 750))
	myWindow.ShowAndRun()
}