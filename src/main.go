package main

import (
	"fmt"
)

func main() {
	var key uint32 = 42
	plainText := "Hello, world!"

	fmt.Println("Оригинальное сообщение: ", plainText)

	fmt.Println("Начало шифровки")
	cipherText := XorCipher([]byte(plainText), key)
	fmt.Printf("Сообщение зашифровано: %X\n", cipherText)

	fmt.Println("Начало дешифровки")
	decryptedText := XorCipher(cipherText, key)
	fmt.Println("Расшифрованное сообщение: ", string(decryptedText))

	if plainText == string(decryptedText) {
		fmt.Println("Сообщения совпадают")
	} else {
		fmt.Println("Сообщения не совпадают")
	}
}
