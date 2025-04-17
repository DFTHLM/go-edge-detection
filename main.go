package main

import (
    "image/jpeg"
    "os"
    "fmt"
    "math"
)

func makeGaussianKernel(radius float64) []float64 {
    width := math.Ceil(radius)
    kernel := make([]float64, int(width) * 2 + 1)

    sigma := radius / 3.0
    sum := 0.0

    base := 1.0 / (math.Sqrt(2.0 * math.Pi) * sigma)
    coeff := 2 * sigma * sigma

    for i := -int(width); i <= int(width); i++ {
        value := base * math.Exp(-float64(i*i) / coeff)
        kernel[i + int(width)] = value
        sum += value
    }

    for i := range kernel {
        kernel[i] /= sum
    }

    return kernel
}

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
    kernel := makeGaussianKernel(5.0)
    grayImg = GaussianFilter(grayImg, kernel, true)
    grayImg = GaussianFilter(grayImg, kernel, false)

    Kx := [][]float64{{-1.0, 0.0, 1.0}, {-2.0, 0.0,  2.0}, {-1.0, 0.0, 1.0}}
    Gx := ApplyKernel(grayImg, Kx)
    Ky := [][]float64{{-1.0, -2.0, -1.0}, {0.0, 0.0, 0.0}, {1.0, 2.0, 1.0}}
    Gy := ApplyKernel(grayImg, Ky)

    direction, G := CalculateGradientValues(Gx, Gy)

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
