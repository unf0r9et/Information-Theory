package main

import (
	"encoding/binary"
	"fmt"
	"image/color"
	"math/big"
	"os"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Криптосистема Рабина")

	var selectedPath string
	var resultData []byte

	pEntry := widget.NewEntry()
	pEntry.PlaceHolder = "p (простое число)"
	qEntry := widget.NewEntry()
	qEntry.PlaceHolder = "q (простое число)"
	bEntry := widget.NewEntry()
	bEntry.PlaceHolder = "b"

	labelInfo := widget.NewLabel("Файл не выбран")

	entryOriginal := widget.NewMultiLineEntry()
	entryOriginal.SetMinRowsVisible(4)

	entryResult := widget.NewMultiLineEntry()
	entryResult.SetMinRowsVisible(4)

	btnSave := widget.NewButton("Сохранить результат на диск", func() {
		if len(resultData) == 0 {
			return
		}
		dialog.ShowFileSave(func(w fyne.URIWriteCloser, err error) {
			if w == nil {
				return
			}
			defer w.Close()
			w.Write(resultData)
			dialog.ShowInformation("Успех", "Файл сохранен!", myWindow)
		}, myWindow)
	})
	btnSave.Disable()

	formatToDec := func(data []byte, step int) string {
		var nums []string

		for i := 0; i < len(data) && i < 50*step; i += step {
			if step == 1 {
				nums = append(nums, fmt.Sprintf("%d", data[i]))
			} else {
				if i+4 > len(data) {
					break
				}
				val := binary.BigEndian.Uint32(data[i : i+4])
				nums = append(nums, fmt.Sprintf("%d", val))
			}
		}
		return strings.Join(nums, " ")
	}

validate := func() (*big.Int, *big.Int, *big.Int, bool) {
		p, _ := new(big.Int).SetString(pEntry.Text, 10)
		q, _ := new(big.Int).SetString(qEntry.Text, 10)
		b, _ := new(big.Int).SetString(bEntry.Text, 10)

		if p == nil || q == nil || b == nil {
			dialog.ShowError(fmt.Errorf("Введите корректные целые числа в поля p, q, b"), myWindow)
			return nil, nil, nil, false
		}

		n := new(big.Int).Mul(p, q)
		limit := big.NewInt(256)
		if n.Cmp(limit) <= 0 {
			dialog.ShowError(fmt.Errorf("Произведение p*q должно быть больше 256 для корректного шифрования байта"), myWindow)
			return nil, nil, nil, false
		}

		mod4 := big.NewInt(4)
		three := big.NewInt(3)
		if new(big.Int).Mod(p, mod4).Cmp(three) != 0 || new(big.Int).Mod(q, mod4).Cmp(three) != 0 {
			dialog.ShowError(fmt.Errorf("Параметры p и q должны удовлетворять условию x ≡ 3 mod 4 (напр. 7, 11, 19...)"), myWindow)
			return nil, nil, nil, false
		}

		if !p.ProbablyPrime(20) || !q.ProbablyPrime(20) {
			dialog.ShowError(fmt.Errorf("Числа p и q должны быть простыми"), myWindow)
			return nil, nil, nil, false
		}

		if b.Cmp(n) >= 0 {
			dialog.ShowError(fmt.Errorf("Параметр b должен быть меньше n (p*q)"), myWindow)
			return nil, nil, nil, false
		}

		return p, q, b, true
	}

	encryptFile := func() {
		p, q, b, ok := validate()
		if !ok || selectedPath == "" {
			return
		}
		n := new(big.Int).Mul(p, q)

		input, _ := os.ReadFile(selectedPath)

		entryOriginal.SetText(formatToDec(input, 1))

		output := make([]byte, len(input)*4)
		for i, mByte := range input {
			m := big.NewInt(int64(mByte))
			c := new(big.Int).Add(m, b)
			c.Mul(c, m).Mod(c, n)
			binary.BigEndian.PutUint32(output[i*4:(i+1)*4], uint32(c.Uint64()))
		}
		resultData = output

		entryResult.SetText(formatToDec(resultData, 4))
		btnSave.Enable()
	}

	decryptFile := func() {
		p, q, b, ok := validate()
		if !ok || selectedPath == "" {
			return
		}
		n := new(big.Int).Mul(p, q)

		input, _ := os.ReadFile(selectedPath)
		if len(input)%4 != 0 {
			dialog.ShowError(fmt.Errorf("Размер зашифрованного файла должен быть кратен 4"), myWindow)
			return
		}

		entryOriginal.SetText(formatToDec(input, 4))

		output := make([]byte, len(input)/4)
		for i := 0; i < len(input)/4; i++ {
			cVal := binary.BigEndian.Uint32(input[i*4 : (i+1)*4])
			output[i] = RabinDecryptByte(big.NewInt(int64(cVal)), b, n, p, q)
		}
		resultData = output

		entryResult.SetText(formatToDec(resultData, 1))
		btnSave.Enable()
	}

	btnOpen := widget.NewButton("Открыть файл", func() {
		dialog.ShowFileOpen(func(r fyne.URIReadCloser, e error) {
			if r != nil {
				selectedPath = r.URI().Path()
				labelInfo.SetText("Выбран: " + selectedPath)
			}
		}, myWindow)
	})

	spacer := func(h float32) fyne.CanvasObject {
		r := canvas.NewRectangle(color.Transparent)
		r.SetMinSize(fyne.NewSize(1, h))
		return r
	}

	content := container.NewVBox(
		btnOpen, labelInfo,
		spacer(5),
		container.NewGridWithColumns(3, pEntry, qEntry, bEntry),
		spacer(10),
		widget.NewLabel("Исходный файл:"),
		entryOriginal,
		spacer(10),
		container.NewGridWithColumns(2,
			widget.NewButton("Зашифровать", encryptFile),
			widget.NewButton("Расшифровать", decryptFile),
		),
		spacer(10),
		widget.NewLabel("Результат:"),
		entryResult,
		spacer(5),
		btnSave,
		widget.NewButton("Сброс", func() {
			selectedPath = ""
			labelInfo.SetText("Файл не выбран")
			resultData = nil
			entryOriginal.SetText("")
			entryResult.SetText("")
			btnSave.Disable()

		}),
	)

	myWindow.SetContent(container.NewPadded(content))
	myWindow.Resize(fyne.NewSize(650, 600))
	myWindow.ShowAndRun()
}
