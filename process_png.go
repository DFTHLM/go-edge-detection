package main

import (
    "image"
    "image/color"
    "math"
    "sync"
)

func convert_color(c color.Color) color.Color {
    r, g, b, a := c.RGBA()

    gray := uint8((0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)) / 256.0)
    return color.RGBA{gray, gray, gray, uint8(a >> 8)}
}

func Grayscale(src image.Image) *image.Gray {

    bounds := src.Bounds()
    out := image.NewGray(bounds)

    for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
        for x := bounds.Min.X; x < bounds.Max.X; x++ {

            originalColor := src.At(x, y)
            grayColor := convert_color(originalColor)

            out.Set(x, y, grayColor)

        }
    }
    return out

}

func GaussianFilter(img *image.Gray, kernel []float64, horizontal bool) *image.Gray {
    bounds := img.Bounds()
    out := image.NewGray(bounds)
    kernelRadius := len(kernel) / 2
    width := bounds.Dx()
    height := bounds.Dy()

    var wg sync.WaitGroup

    for y := range height {
        wg.Add(1)
        go func(y int) {

            for x := range width {
                var sum float64
                var weightSum float64

                for k := -kernelRadius; k <= kernelRadius; k++ {
                    kx := x + k
                    ky := y + k
                    pos := k + kernelRadius

                    if horizontal {
                        if kx < 0 { continue }
                        if kx >= width { continue }
                        sum += float64(img.GrayAt(kx, y).Y) * kernel[pos]
                    } else {
                        if ky < 0 { continue }
                        if ky >= height { continue }
                        sum += float64(img.GrayAt(x, ky).Y) * kernel[pos]
                    }

                    weightSum += kernel[pos]
                }
                val := uint8(math.Min(math.Max(float64(sum/weightSum), 0), 255))
                out.SetGray(x, y, color.Gray{Y: val})
            }
            wg.Done()
        }(y)
    }

    wg.Wait()
    return out
}

func ApplyKernel(img *image.Gray, kernel [][]float64) *image.Gray {
    bounds := img.Bounds()
    width := bounds.Dx()
    height := bounds.Dy()
    kernelWidth := len(kernel[0])
    kernelHeight := len(kernel)

    out := image.NewGray(bounds)

    var wg sync.WaitGroup

    for y := range height {   
        wg.Add(1)
        go func(y int) {

            for x := range width {
                var sum float64

                for ky := range kernelHeight {
                    for kx := range kernelWidth {

                        xk, yk := x+kx-(kernelWidth/2), y+ky-(kernelHeight/2)

                        if xk >= 0 && xk < width && yk >= 0 && yk < height {
                            sum += float64(img.GrayAt(xk, yk).Y) * kernel[ky][kx]
                        }
                    }
                }

                sum = math.Abs(sum)
                sum = math.Min(255, sum)
                out.SetGray(x, y, color.Gray{Y: uint8(sum)})

            }
            wg.Done()

        }(y)
    }
    wg.Wait()
    return out
}

func GetMagnitude(Gx, Gy *image.Gray) (*image.Gray) {
    bounds := Gx.Bounds()
    width := bounds.Dx()
    height := bounds.Dy()
    out := image.NewGray(bounds)

    for y := range height {
        for x := range width {

            gx := float64(Gx.GrayAt(x, y).Y)
            gy := float64(Gy.GrayAt(x, y).Y)

            mag := uint8(math.Min(math.Sqrt(float64(gx*gx+gy*gy)), 255))
            out.SetGray(x, y, color.Gray{Y: mag})
        }
    }

    return out
}

func CalculateGradientDirection(Gx, Gy *image.Gray) [][]float64 {
    bounds := Gx.Bounds()
    width := bounds.Dx()
    height := bounds.Dy()

    out := make([][]float64, height)
    for i := range out {
        out[i] = make([]float64, width)
    }

    for y := range height {
        for x := range width {

            gx := float64(Gx.GrayAt(x, y).Y)
            gy := float64(Gy.GrayAt(x, y).Y)

            if gx == 0 && gy == 0 {
                out[y][x] = 0
                continue
            }

            dir := math.Atan2(float64(gy), float64(gx)) * 180 / math.Pi
            if dir < 0 {
                dir += 180
            }
            out[y][x] = dir
        }
    }

    return out
}

func NonMaxSuppression(G *image.Gray, direction [][]float64) *image.Gray {
    bounds := G.Bounds()
    out := image.NewGray(bounds)
    width := bounds.Dx()
    height := bounds.Dy()

    for y := 1; y < height-1; y++ {
        for x := 1; x < width-1; x++ {
            mag := int(G.GrayAt(x, y).Y)
            dir := direction[y][x]

            var neighbor1, neighbor2 int

            if dir < 22.5 || dir >= 157.5 {
                neighbor1 = int(G.GrayAt(x-1, y).Y)
                neighbor2 = int(G.GrayAt(x+1, y).Y)

            } else if dir < 67.5 {
                neighbor1 = int(G.GrayAt(x-1, y-1).Y)
                neighbor2 = int(G.GrayAt(x+1, y+1).Y)

            } else if dir < 112.5 {
                neighbor1 = int(G.GrayAt(x, y-1).Y)
                neighbor2 = int(G.GrayAt(x, y+1).Y)

            } else {
                neighbor1 = int(G.GrayAt(x-1, y+1).Y)
                neighbor2 = int(G.GrayAt(x+1, y-1).Y)
            }

            if mag >= neighbor1 && mag >= neighbor2 {
                out.SetGray(x, y, color.Gray{Y: uint8(mag)})

            } else {
                out.SetGray(x, y, color.Gray{Y: 0})
            }
        }
    }

    return out
}


