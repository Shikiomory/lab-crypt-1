package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	var key uint32 = 42
	srcImgPath := "imgs/gradient.png"

	if len(os.Args) > 1 {
		srcImgPath = os.Args[1]
	}

	baseName := filepath.Base(srcImgPath)
	ext := filepath.Ext(baseName)
	name := strings.TrimSuffix(baseName, ext)

	srcDecImgPath := filepath.Join("results", name+"-dec"+ext)
	srcEncImgPath := filepath.Join("results", name+"-enc"+ext)

	pixels, bounds, err := readImage(srcImgPath)
	if err != nil {
		fmt.Println("Возникла проблема при запуске: ", err)
		return
	}

	fmt.Printf("Изображение размера %dx%d загружено: %X\n", bounds.Dx(), bounds.Dy(), pixels[:20])

	fmt.Println("Начало шифровки")
	cipherText := XorCipher(pixels, key)
	fmt.Println("Сообщение зашифровано: ", string(cipherText)[:20])
	err = saveImage(cipherText, bounds, srcEncImgPath)
	if err != nil {
		fmt.Println("Возникла проблема при сохранении зашифрованного файла", err)
		return
	}

	fmt.Println("Начало дешифровки")
	decryptedText := XorCipher(cipherText, key)
	fmt.Println("Сообщение расшифровано: ", string(decryptedText)[:20])
	err = saveImage(decryptedText, bounds, srcDecImgPath)
	if err != nil {
		fmt.Println("Возникла проблема при сохранении расшифрованного файла", err)
		return
	}

	if bytes.Equal(pixels, decryptedText) {
		fmt.Println("Сообщения совпадают")
	} else {
		fmt.Println("Сообщения не совпадают")
	}
}
