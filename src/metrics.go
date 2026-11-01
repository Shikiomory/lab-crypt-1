package main

import (
	"encoding/csv"
	"fmt"
	"image"
	"math"
	"os"
	"strconv"
)

func CalcEntropy(data []byte) float64 {
	total := len(data)

	counts := CalcFrequency(data)

	var entropy float64 = 0.0
	for _, count := range counts {
		if count == 0 {
			continue
		}

		p := float64(count) / float64(total)
		entropy -= p * math.Log2(p)
	}

	return entropy
}

func pearsonCorr(x, y []float64) float64 {
	n := len(x)
	if n != len(y) {
		return 0.0
	}

	var sumX, sumY float64

	for i := 0; i < n; i++ {
		sumX += x[i]
		sumY += y[i]
	}

	meanX := sumX / float64(n)
	meanY := sumY / float64(n)

	var denominatorX, denominatorY, numerator float64
	for i := range n {
		diffX := x[i] - meanX
		diffY := y[i] - meanY

		numerator += diffY * diffX
		denominatorX += diffX * diffX
		denominatorY += diffY * diffY
	}

	denominator := math.Sqrt(denominatorX * denominatorY)
	if denominator == 0 {
		return 0.0
	}

	return numerator / denominator
}

func CalcCorrelation(pixels []byte, bounds image.Rectangle) (corrH float64, corrV float64) {
	height := bounds.Dy()
	width := bounds.Dx()

	getPixelVal := func(x, y int) float64 {
		idx := (width*y + x) * 4
		r := float64(pixels[idx])
		g := float64(pixels[idx+1])
		b := float64(pixels[idx+2])

		return 0.299*r + 0.587*g + 0.114*b
	}

	idx := 0
	htotal := height * (width - 1)
	hx := make([]float64, htotal)
	hy := make([]float64, htotal)
	for y := 0; y < height; y++ {
		for x := 0; x < width-1; x++ {
			hx[idx] = getPixelVal(x, y)
			hy[idx] = getPixelVal(x+1, y)
			idx++
		}
	}

	corrH = pearsonCorr(hx, hy)

	idx = 0
	vtotal := (height - 1) * width
	vx := make([]float64, vtotal)
	vy := make([]float64, vtotal)
	for y := 0; y < height-1; y++ {
		for x := 0; x < width; x++ {
			vx[idx] = getPixelVal(x, y)
			vy[idx] = getPixelVal(x, y+1)
			idx++
		}
	}

	corrV = pearsonCorr(vx, vy)
	return corrH, corrV
}

func CalcFrequency(data []byte) [256]int {
	var counts [256]int
	for _, el := range data {
		counts[el]++
	}
	return counts
}

func saveHistogram(filePath string, orig, xor, ecb, ctr [256]int) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("не удалось создать файл: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	header := []string{"byte", "original", "xor", "ecb", "ctr"}
	err = writer.Write(header)
	if err != nil {
		return fmt.Errorf("не удалось записать заголовок: %w", err)
	}

	for i := 0; i < 256; i++ {
		row := []string{
			strconv.Itoa(i),
			strconv.Itoa(orig[i]),
			strconv.Itoa(xor[i]),
			strconv.Itoa(ecb[i]),
			strconv.Itoa(ctr[i]),
		}
		err = writer.Write(row)
		if err != nil {
			return fmt.Errorf("не удалось записать строку: %w", err)
		}
	}
	return nil
}
