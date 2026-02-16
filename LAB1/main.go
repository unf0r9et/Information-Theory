package main

const alphabetSizeRU = 33
const alphabetSizeEN = 26

var CharToIndexRU map[rune]int
var CharToIndexEN map[rune]int

var IndexToCharRU map[int]rune
var IndexToCharEN map[int]rune

func init() {
	CharToIndexRU = make(map[rune]int)
	CharToIndexEN = make(map[rune]int)

	IndexToCharEN = make(map[int]rune)
	IndexToCharRU = make(map[int]rune)

	upper := []rune("АБВГДЕЁЖЗИЙКЛМНОПРСТУФХЦЧШЩЪЫЬЭЮЯ")
	lower := []rune("абвгдеёжзийклмнопрстуфхцчшщъыьэюя")

	for i := 0; i < len(upper); i++ {
		CharToIndexRU[upper[i]] = i
		CharToIndexRU[lower[i]] = i
		IndexToCharRU[i] = upper[i]
	}

	upperEN := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	lowerEN := []rune("abcdefghijklmnopqrstuvwxyz")

	for i := 0; i < len(upperEN); i++ {
		CharToIndexEN[upperEN[i]] = i
		CharToIndexEN[lowerEN[i]] = i
		IndexToCharEN[i] = upperEN[i]
	}
}

func main() {
	StartUI()
}
