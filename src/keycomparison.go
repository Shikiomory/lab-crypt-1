package main

import "math"

func CalcNPCR(cipherText1, cipherText2 []byte) float64 {
	cipherLen := len(cipherText1)
	if cipherLen != len(cipherText2) {
		return 0.0
	}
	var numerator int = 0
	for i := 0; i < cipherLen; i++ {
		if cipherText1[i] != cipherText2[i] {
			numerator++
		}
	}
	return float64(numerator) / float64(cipherLen)
}

func CalcUACI(cipherText1, cipherText2 []byte) float64 {
	cipherLen := len(cipherText1)
	if cipherLen != len(cipherText2) || cipherLen == 0 {
		return 0.0
	}

	var sum float64 = 0.0
	for i := 0; i < cipherLen; i++ {
		diff := float64(cipherText1[i]) - float64(cipherText2[i])
		sum += math.Abs(diff) / 255.0
	}

	return sum / float64(cipherLen)
}
