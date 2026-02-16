package main

import (
	"regexp"
	"strconv"
	"unicode/utf8"
)

func DecimationEncrypt(text string, key int) string {
	runeText := []rune(text)
	lenghtText := utf8.RuneCountInString(text)

	var cipherRunes []rune

	for i := 0; i < lenghtText; i++ {
		idx, doesItExist := CharToIndexEN[runeText[i]]
		if doesItExist {
			plainIndex := (idx * key) % alphabetSizeEN
			cipherRunes = append(cipherRunes, IndexToCharEN[plainIndex])
		}
	}

	return string(cipherRunes)
}

func DecimationDecipher(text string, key int) string {

	reverseKey, ok := modInverse(key, alphabetSizeEN)
	if !ok {
		return ""
	}

	runeText := []rune(text)
	lenghtText := utf8.RuneCountInString(text)

	var decryptedRunes []rune

	for i := 0; i < lenghtText; i++ {
		idx, doesItExist := CharToIndexEN[runeText[i]]
		if doesItExist {
			plainIndex := (idx * reverseKey) % alphabetSizeEN
			decryptedRunes = append(decryptedRunes, IndexToCharEN[plainIndex])
		}
	}

	return string(decryptedRunes)
}

func convertStringToNumber(key string) int {
	re := regexp.MustCompile(`\d+`)
	foundStr := re.FindString(key)

	if foundStr != "" {
		num, _ := strconv.Atoi(foundStr)
		return num
	} else {
		return -1
	}
}

func IsTheKeyCorrect(key string) bool {
	num := convertStringToNumber(key)
	if num != -1 {
		if num%2 == 0 || num%13 == 0 {
			return false
		}
		return true
	} else {
		return false
	}
}

func extendedGCD(a, b int) (gcd, x, y int) {
	if b == 0 {
		return a, 1, 0
	}

	gcd, x1, y1 := extendedGCD(b, a%b)

	x = y1
	y = x1 - (a/b)*y1

	return
}

func modInverse(key, mod int) (int, bool) {
	gcd, x, _ := extendedGCD(key, mod)

	if gcd != 1 {
		return 0, false
	}

	inv := x % mod
	if inv < 0 {
		inv += mod
	}

	return inv, true
}
