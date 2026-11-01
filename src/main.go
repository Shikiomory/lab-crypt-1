package main

import (
	"bytes"
	"crypto/aes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"image"
	"os"
	"path/filepath"
	"strings"
)

const jsonPathFile string = "meta.json"

func runXOR(pixels []byte, bounds image.Rectangle, key uint32, srcImgPath string) ([]byte, error) {
	fmt.Println("\n===== Работа XOR: =====\n")
	srcDecImgPath, srcEncImgPath := createOutPath(srcImgPath, "XOR")

	meta := MetaInfo{
		Algoritm: "aes-ctr",
		Nonce:    "",
		Source:   filepath.Base(srcDecImgPath),
	}

	err := Save2JSON(jsonPathFile, meta)
	if err != nil {
		return nil, err
	}

	fmt.Println("Начало шифровки")
	cipherText := XorCipher(pixels, key)
	fmt.Printf("Сообщение зашифровано: %X\n", cipherText[:20])
	err = saveImage(cipherText, bounds, srcEncImgPath)
	if err != nil {
		fmt.Println("Возникла проблема при сохранении зашифрованного файла", err)
		return nil, err
	}

	fmt.Println("Начало дешифровки")
	decryptedText := XorCipher(cipherText, key)
	fmt.Printf("Сообщение расшифровано: %X\n", decryptedText[:20])
	err = saveImage(decryptedText, bounds, srcDecImgPath)
	if err != nil {
		fmt.Println("Возникла проблема при сохранении расшифрованного файла", err)
		return nil, err
	}

	equalsCrypt(pixels, decryptedText)
	return cipherText, nil
}

func runECB(pixels []byte, bounds image.Rectangle, key []byte, srcImgPath string) ([]byte, error) {
	fmt.Println("\n===== Работа ECB: =====\n")
	srcDecImgPath, srcEncImgPath := createOutPath(srcImgPath, "ECB")

	meta := MetaInfo{
		Algoritm: "aes-ctr",
		Nonce:    "",
		Source:   filepath.Base(srcDecImgPath),
	}

	err := Save2JSON(jsonPathFile, meta)
	if err != nil {
		return nil, err
	}

	fmt.Println("Начало шифровки")
	cipherText, err := EncryptAES_ECB(pixels, key)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Сообщение зашифровано: %X\n", cipherText[:20])
	err = saveImage(cipherText, bounds, srcEncImgPath)
	if err != nil {
		fmt.Println("Возникла проблема при сохранении зашифрованного файла", err)
		return nil, err
	}

	fmt.Println("Начало дешифровки")
	decryptedText, err := DecryptAES_ECB(cipherText, key)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Сообщение расшифровано: %X\n", decryptedText[:20])
	err = saveImage(decryptedText, bounds, srcDecImgPath)
	if err != nil {
		fmt.Println("Возникла проблема при сохранении расшифрованного файла", err)
		return nil, err
	}

	equalsCrypt(pixels, decryptedText)
	return cipherText, nil
}

func runCTR(pixels []byte, bounds image.Rectangle, key []byte, srcImgPath string) ([]byte, error) {
	fmt.Println("\n===== Работа CTR: =====\n")
	srcDecImgPath, srcEncImgPath := createOutPath(srcImgPath, "CTR")

	iv, err := GenerateIV(aes.BlockSize)
	if err != nil {
		return nil, err
	}

	meta := MetaInfo{
		Algoritm: "aes-ctr",
		Nonce:    hex.EncodeToString(iv),
		Source:   filepath.Base(srcDecImgPath),
	}

	err = Save2JSON(jsonPathFile, meta)
	if err != nil {
		return nil, err
	}

	fmt.Println("Начало шифровки")
	cipherText, err := CTRCipher(pixels, key, iv)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Сообщение зашифровано: %X\n", cipherText[:20])
	err = saveImage(cipherText, bounds, srcEncImgPath)
	if err != nil {
		fmt.Println("Возникла проблема при сохранении зашифрованного файла", err)
		return nil, err
	}

	_, iv, err = LoadJSON(jsonPathFile)
	if err != nil {
		return nil, err
	}

	fmt.Println("Начало дешифровки")
	decryptedText, err := CTRCipher(cipherText, key, iv)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Сообщение расшифровано: %X\n", decryptedText[:20])
	err = saveImage(decryptedText, bounds, srcDecImgPath)
	if err != nil {
		fmt.Println("Возникла проблема при сохранении расшифрованного файла", err)
		return nil, err
	}

	equalsCrypt(pixels, decryptedText)
	return cipherText, nil
}

func createOutPath(srcImgPath string, alg string) (srcDecImgPath string, srcEncImgPath string) {
	baseName := filepath.Base(srcImgPath)
	ext := filepath.Ext(baseName)
	name := strings.TrimSuffix(baseName, ext)

	srcDecImgPath = filepath.Join("results", name+alg+"-dec"+ext)
	srcEncImgPath = filepath.Join("results", name+alg+"-enc"+ext)

	return srcDecImgPath, srcEncImgPath
}

func equalsCrypt(plainText []byte, decryptedText []byte) {
	if bytes.Equal(plainText, decryptedText) {
		fmt.Println("Сообщения совпадают")
	} else {
		fmt.Println("Сообщения не совпадают")
	}

	hash := sha256.Sum256(plainText)
	fmt.Printf("SHA-256: %x\n", hash)
	hashNew := sha256.Sum256(decryptedText)
	fmt.Printf("SHA-256: %x\n", hashNew)
	if bytes.Equal(plainText, decryptedText) {
		fmt.Println("Хэши совпадают")
	} else {
		fmt.Println("Хэши не совпадают")
	}
}

func main() {
	var xorKey uint32 = 42
	var aesKey = []byte("abcdefgzxrtyui16")
	srcImgPath := "imgs/gradient.png"

	if len(os.Args) > 1 {
		srcImgPath = os.Args[1]
	}

	pixels, bounds, err := readImage(srcImgPath)
	if err != nil {
		fmt.Println("Возникла проблема при запуске: ", err)
		return
	}

	fmt.Printf("Изображение размера %dx%d загружено: %X\n", bounds.Dx(), bounds.Dy(), pixels[:20])

	cipherXORText, err := runXOR(pixels, bounds, xorKey, srcImgPath)
	if err != nil {
		fmt.Println("Сбой работы XOR: ", err)
		return
	}

	cipherECBText, err := runECB(pixels, bounds, aesKey, srcImgPath)
	if err != nil {
		fmt.Println("Сбой работы ECB: ", err)
		return
	}

	cipherCTRText, err := runCTR(pixels, bounds, aesKey, srcImgPath)
	if err != nil {
		fmt.Println("Сбой работы ECB: ", err)
		return
	}

	fmt.Println("===== Энтропия: =====")
	fmt.Printf("Оригинальное изображение: %f\n", CalcEntropy(pixels))
	fmt.Printf("XOR шифровка: %f\n", CalcEntropy(cipherXORText))
	fmt.Printf("ECB шифровка: %f\n", CalcEntropy(cipherECBText))
	fmt.Printf("CTR шифровка: %f\n", CalcEntropy(cipherCTRText))

	fmt.Println("===== Корреляция: =====")
	hOrigCol, vOrigCol := CalcCorrelation(pixels, bounds)
	hXORCol, vXORCol := CalcCorrelation(cipherXORText, bounds)
	hECBCol, vECBCol := CalcCorrelation(cipherECBText, bounds)
	hCTRCol, vCTRCol := CalcCorrelation(cipherCTRText, bounds)
	fmt.Printf("Оригинальное изображение: гориз. = %f верт. = %f\n", hOrigCol, vOrigCol)
	fmt.Printf("XOR шифровка: гориз. = %f верт. = %f\n", hXORCol, vXORCol)
	fmt.Printf("ECB шифровка: гориз. = %f верт. = %f\n", hECBCol, vECBCol)
	fmt.Printf("CTR шифровка: гориз. = %f верт. = %f\n", hCTRCol, vCTRCol)
}
