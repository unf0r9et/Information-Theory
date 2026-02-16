package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	_ "fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func StartUI() {

	//-------------------------------------------------------------------------------------------

	myApp := app.NewWithID("com.unf0r9et.cryptography")

	//myApp.Settings().SetTheme(theme.LightTheme())

	myWindow := myApp.NewWindow("Криптография")

	var filePathToRead string

	var filePathToWrite string

	labelPathFromFile := widget.NewLabel("Файл для чтения не выбран")
	labelPathToFile := widget.NewLabel("Файл для записи не выбран")

	entryFromFile := widget.NewMultiLineEntry()
	entryFromFile.SetMinRowsVisible(3)

	entryToFile := widget.NewMultiLineEntry()
	entryToFile.SetMinRowsVisible(3)

	keyEntry := widget.NewEntry()

	methodSelect := widget.NewRadioGroup([]string{"Метод Децимации (Язык: EN)", "Метод Виженера прямой ключ (Язык: RU)"}, nil)
	methodSelect.SetSelected("Метод Децимации (Язык: EN)")

	//-------------------------------------------------------------------------------------------

	process := func(opName string) {
		keyWord := keyEntry.Text
		textFromEntry := entryFromFile.Text
		method := methodSelect.Selected

		if textFromEntry == "" {
			dialog.ShowInformation("ОШИБКА", "Поле исходного текста пустое. Выберите файл или введите текст.", myWindow)
			return
		}

		var result string

		switch method {
		case "Метод Децимации (Язык: EN)":
			if IsTheKeyCorrect(keyWord) {
				if opName == "encrypt" {
					result = DecimationEncrypt(textFromEntry, convertStringToNumber(keyWord))
				} else {
					result = DecimationDecipher(textFromEntry, convertStringToNumber(keyWord))
				}
			} else {
				dialog.ShowInformation("ОШИБКА", "Неверный ключ. (Ключ должен быть взаимно простым с 26)", myWindow)
				return
			}
			if result == "" {
				dialog.ShowInformation("ОШИБКА", "Неверный текст в поле ввода. (Текст на английском языке)", myWindow)
				return
			}
			break
		case "Метод Виженера прямой ключ (Язык: RU)":
			if opName == "encrypt" {
				result = VigenereEncrypt(textFromEntry, keyWord)
			} else {
				result = VigenereDecipher(textFromEntry, keyWord)
			}
			if result == "" {
				dialog.ShowInformation("ОШИБКА", "Неверный ключ и/или поле ввода. (Текст на русском языке)", myWindow)
				return
			}
			break
		default:
			dialog.ShowInformation("ОШИБКА", "Выберите метод", myWindow)
			return
		}

		entryToFile.SetText(result)

		if filePathToWrite != "" {
			saveWindowsFile(filePathToWrite, entryToFile.Text)
		}

	}

	buttonEncrypt := widget.NewButton("Шифровать", func() { process("encrypt") })

	buttonDecipher := widget.NewButton("Дешифровать", func() { process("decrypt") })

	buttonClear := widget.NewButton("Очистить", func() {
		entryFromFile.SetText("")
		entryToFile.SetText("")
		keyEntry.SetText("")
		labelPathFromFile.SetText("Файл для чтения не выбран")
		labelPathToFile.SetText("Файл для записи не выбран")
	})

	//-------------------------------------------------------------------------------------------

	workWithFile := func(pathName string) {
		fileDialog := dialog.NewFileOpen(
			func(reader fyne.URIReadCloser, err error) {

				if err != nil {
					dialog.ShowError(err, myWindow)
					return
				}

				if reader == nil {
					return
				}

				switch pathName {
				case "READ":
					filePathToRead = reader.URI().Path()

					labelPathFromFile.SetText("Файл для чтения: " + filePathToRead)

					reader.Close()

					text, err := readWindowsFile(filePathToRead)

					if err == nil {
						entryFromFile.SetText(text)
					} else {
					}
					break
				case "WRITE":
					filePathToWrite = reader.URI().Path()

					labelPathToFile.SetText("Файл для записи: " + filePathToWrite)

					reader.Close()
					break
				}

			},
			myWindow,
		)

		fileDialog.Show()

	}

	btnOpenFileToRead := widget.NewButton("Выберете файл для чтения (опционально)", func() { workWithFile("READ") })

	btnOpenFileToWrite := widget.NewButton("Выберете файл для записи (опционально)", func() { workWithFile("WRITE") })

	//-------------------------------------------------------------------------------------------

	spacer := func(height float32) fyne.CanvasObject {
		rect := canvas.NewRectangle(color.Transparent)
		rect.SetMinSize(fyne.NewSize(1, height))
		return rect
	}

	rowKey := container.NewBorder(nil, nil, widget.NewLabel("Введите ключ: "), nil, keyEntry)

	twoButtons := container.NewGridWithColumns(2,
		buttonEncrypt,
		buttonDecipher,
	)

	horizontalContent := container.NewVBox(
		btnOpenFileToRead,
		btnOpenFileToWrite,
		labelPathFromFile,
		labelPathToFile,
	)

	myColor := color.RGBA{R: 82, G: 82, B: 82, A: 255}

	line := canvas.NewRectangle(myColor)

	line.SetMinSize(fyne.NewSize(1, 1))

	content := container.NewVBox(
		spacer(10),
		horizontalContent,
		line,
		widget.NewLabel("Выберете метод: "),
		methodSelect,
		spacer(10),
		rowKey,
		spacer(10),
		widget.NewLabel("Исходный текст: "),
		entryFromFile,
		widget.NewLabel("Полученный текст: "),
		entryToFile,
		spacer(10),
		twoButtons,
		spacer(10),
		container.NewBorder(nil, nil, nil, nil, buttonClear),
		spacer(10),
	)

	myWindow.SetContent(content)
	myWindow.Resize(fyne.NewSize(600, 400))
	myWindow.ShowAndRun()
}
