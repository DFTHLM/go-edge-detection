package main

import (
    "image/jpeg"
    "os"
    "fmt"
)

func parseArgs() (string, string) {
    if len(os.Args) != 3 {
        fmt.Println("Usage: go run main.go <input_image_path> <output_image_path>")
        os.Exit(1)
    }
    return os.Args[1], os.Args[2]
}

func main() {
    inputPath, outputPath := parseArgs()

    file, err := os.Open(inputPath)
    if err != nil {
        fmt.Println("Error opening file:")
        panic(err)
    }
    defer file.Close()

    img, err := jpeg.Decode(file)
    if err != nil {
        fmt.Println("Error decoding image:")
        panic(err)
    }

    grayImg := Grayscale(img)
    kernel := []float64{0.0003,	0.0033,	0.0237,	0.0970,	0.2260,	0.2995,	0.2260,	0.0970,	0.0237,	0.0033,	0.0003}
    grayImg = GaussianFilter(grayImg, kernel, true)
    grayImg = GaussianFilter(grayImg, kernel, false)

    Kx := [][]float64{{-1.0, 0.0, 1.0}, {-2.0, 0.0,  2.0}, {-1.0, 0.0, 1.0}}
    Gx := ApplyKernel(grayImg, Kx)
    Ky := [][]float64{{-1.0, -2.0, -1.0}, {0.0, 0.0, 0.0}, {1.0, 2.0, 1.0}}
    Gy := ApplyKernel(grayImg, Ky)

    G := GetMagnitude(Gx, Gy)
    direction := CalculateGradientDirection(Gx, Gy)

    suppressed := NonMaxSuppression(G, direction)

    outFile, err := os.Create(outputPath)
    if err != nil {
        fmt.Println("Error creating output file:")
        panic(err)
    }
    defer outFile.Close()

    err = jpeg.Encode(outFile, suppressed, &jpeg.Options{Quality: 100})
    if err != nil {
        fmt.Println("Error encoding image:")
        panic(err)
    }
}
