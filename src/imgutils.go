package main

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"os"
)

func readImage(filePath string) ([]byte, image.Rectangle, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, image.Rectangle{}, fmt.Errorf("Ошибка при открытии файла", err)
	}

	defer file.Close()

	srcImg, err := png.Decode(file)
	if err != nil {
		return nil, image.Rectangle{}, fmt.Errorf("Ошибка при распаковке изображения")
	}

	bounds := srcImg.Bounds()
	nrgba := image.NewNRGBA(bounds)
	draw.Draw(nrgba, bounds, srcImg, bounds.Min, draw.Src)

	return nrgba.Pix, bounds, nil
}

func saveImage(pixels []byte, bounds image.Rectangle, filePath string) error {
	outFile, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("Ошибка при сохраении файла", err)
	}

	defer outFile.Close()

	img := &image.NRGBA{
		Pix:    pixels,
		Stride: bounds.Dx() * 4,
		Rect:   bounds,
	}

	err = png.Encode(outFile, img)
	if err != nil {
		return fmt.Errorf("Ошибка при кодировании файла", err)
	}

	return nil
}
