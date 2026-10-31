package main

import "fmt"

func main() {
	rng := NewXorShift32(42)

	fmt.Println("Начало шифровки")

	for i := 0; i < 5; i++ {
		b := rng.NextByte()
		fmt.Printf("Байт %d: %3d (hex: 0x%02X)\n", i+1, b, b)
	}
}
