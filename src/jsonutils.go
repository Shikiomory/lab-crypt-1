package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
)

type MetaInfo struct {
	Algoritm string `json:"algoritm"`
	Nonce    string `json:"nonce`
	Source   string `json:"source"`
}

func Save2JSON(filePath string, meta MetaInfo) error {
	outFile, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("Ошибка при создании файла: %w", err)
	}

	defer outFile.Close()

	encoder := json.NewEncoder(outFile)
	encoder.SetIndent("", "  ")

	err = encoder.Encode(meta)
	if err != nil {
		return fmt.Errorf("Ошибка при сохранении файла: %w", err)
	}

	return nil
}

func LoadJSON(filepath string) (MetaInfo, []byte, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return MetaInfo{}, nil, fmt.Errorf("Ошибка при открытии файла: %w", err)
	}

	defer file.Close()

	var meta MetaInfo

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&meta)
	if err != nil {
		return MetaInfo{}, nil, fmt.Errorf("Ошибка при парсинге файла: %w", err)
	}

	iv, err := hex.DecodeString(meta.Nonce)
	if err != nil {
		return MetaInfo{}, nil, fmt.Errorf("Ошибка при декодировании iv: %w", err)
	}

	return meta, iv, nil
}
