package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
)

// func generateBytes(n int) ([]byte, error) {
// 	b := make([]byte, n)
// 	_
// }

func EncryptAES_ECB(data, key []byte) ([]byte, error) {
	block, blockSize, err := initAES_ECB(data, key)
	if err != nil {
		return nil, err
	}
	cipherText := make([]byte, len(data))
	for i := 0; i < len(data); i += blockSize {
		block.Encrypt(cipherText[i:i+blockSize], data[i:i+blockSize])
	}

	return cipherText, nil
}

func DecryptAES_ECB(data, key []byte) ([]byte, error) {
	block, blockSize, err := initAES_ECB(data, key)
	if err != nil {
		return nil, err
	}
	plainText := make([]byte, len(data))
	for i := 0; i < len(data); i += blockSize {
		block.Decrypt(plainText[i:i+blockSize], data[i:i+blockSize])
	}

	return plainText, nil
}

func initAES_ECB(data, key []byte) (cipher.Block, int, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, 0, err
	}

	blockSize := block.BlockSize()
	if len(data)%blockSize != 0 {
		return nil, 0, errors.New("Размер данных не кратен 16 байтам")
	}

	return block, blockSize, nil
}

//CTR:

func GenerateIV(blockSize int) ([]byte, error) {
	iv := make([]byte, blockSize)
	_, err := rand.Read(iv)
	if err != nil {
		return nil, fmt.Errorf("ошибка генерации IV: %w", err)
	}
	return iv, nil
}

func CTRCipher(data, key, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if len(iv) != block.BlockSize() {
		return nil, errors.New("Размер блока не совпадает с iv")
	}

	result := make([]byte, len(data))
	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(result, data)

	return result, nil
}
