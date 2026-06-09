package pkg


import (
	"math/rand/v2"
)

func BuildBankCode() string{

	// 定義包含大小寫英文字母的字元集
	letters := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	length := 3

	// 產生 3 位數的隨機字串
	b := make([]rune, length)
	for i := range b {
		b[i] = letters[rand.IntN(len(letters))]
	}

	return string(b)
}