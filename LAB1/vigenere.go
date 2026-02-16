package main

import (
	"unicode/utf8"
)

func VigenereEncrypt(text string, key string) string {

	runeText := []rune(text)
	lenghtText := utf8.RuneCountInString(text)

	correctKey := CorrectKeyVigenere(key)
	if correctKey == "" {
		return ""
	}

	runeKey := []rune(correctKey)
	lenghtKey := utf8.RuneCountInString(correctKey)
	indexKey := 0

	var cipherRunes []rune

	for i := 0; i < lenghtText; i++ {
		idx, doesItExist := CharToIndexRU[runeText[i]]
		if doesItExist {
			if indexKey >= lenghtKey {
				indexKey = 0
			}

			plainIndex := (idx + CharToIndexRU[runeKey[indexKey]]) % alphabetSizeRU
			cipherRunes = append(cipherRunes, IndexToCharRU[plainIndex])
			indexKey++
		}
	}

	return string(cipherRunes)
}

func VigenereDecipher(text string, key string) string {

	runeText := []rune(text)
	lenghtText := utf8.RuneCountInString(text)

	correctKey := CorrectKeyVigenere(key)
	if correctKey == "" {
		return ""
	}
	runeKey := []rune(correctKey)
	lenghtKey := utf8.RuneCountInString(correctKey)
	indexKey := 0

	var cipherRunes []rune

	for i := 0; i < lenghtText; i++ {
		idx, doesItExist := CharToIndexRU[runeText[i]]
		if doesItExist {
			if indexKey >= lenghtKey {
				indexKey = 0
			}
			plainIndex := (idx - CharToIndexRU[runeKey[indexKey]] + alphabetSizeRU) % alphabetSizeRU
			cipherRunes = append(cipherRunes, IndexToCharRU[plainIndex])
			indexKey++
		}
	}

	return string(cipherRunes)
}

func CorrectKeyVigenere(key string) string {
	runes := []rune(key)
	var correctKey []rune
	count := utf8.RuneCountInString(key)

	for i := 0; i < count; i++ {
		_, doesKeyExist := CharToIndexRU[runes[i]]
		if doesKeyExist {
			correctKey = append(correctKey, runes[i])
		}
	}
	return string(correctKey)
}
