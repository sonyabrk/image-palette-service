package processor

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"sort"
)

type RGB struct {
	R, G, B int
}

type bucket struct {
	pixels []RGB
}

func extractPalette(img image.Image, n int) []RGB {
	pixels := collectPixels(img)

	buckets := []bucket{{pixels: pixels}}

	for len(buckets) < n {
		largest := findLargestBucket(buckets)
		left, right := splitBucket(buckets[largest])

		buckets = append(buckets[:largest], buckets[largest+1:]...)
		buckets = append(buckets, left, right)
	}

	palette := make([]RGB, len(buckets))
	for i, b := range buckets {
		palette[i] = averageColor(b.pixels)
	}

	return palette
}

func collectPixels(img image.Image) []RGB {
	bounds := img.Bounds()

	pixels := make([]RGB, 0, bounds.Dx()*bounds.Dy())

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			pixels = append(pixels, RGB{
				R: int(r / 257),
				G: int(g / 257),
				B: int(b / 257),
			})
		}
	}

	return pixels
}

func findLargestBucket(buckets []bucket) int {
	maxRange := -1
	maxIdx := 0

	for i, b := range buckets {
		r := colorRange(b.pixels)
		if r > maxRange {
			maxRange = r
			maxIdx = i
		}
	}

	return maxIdx
}

func colorRange(pixels []RGB) int {
	if len(pixels) == 0 {
		return 0
	}

	minR, maxR := 255, 0
	minG, maxG := 255, 0
	minB, maxB := 255, 0

	for _, p := range pixels {
		if p.R < minR {
			minR = p.R
		}
		if p.R > maxR {
			maxR = p.R
		}

		if p.G < minG {
			minG = p.G
		}
		if p.G > maxG {
			maxG = p.G
		}

		if p.B < minB {
			minB = p.B
		}
		if p.B > maxB {
			maxB = p.B
		}
	}

	return max(maxR-minR, max(maxG-minG, maxB-minB))
}

func splitBucket(b bucket) (bucket, bucket) {
	if len(b.pixels) == 0 {
		return bucket{}, bucket{}
	}

	rRange := channelRange(b.pixels, func(p RGB) int { return p.R })
	gRange := channelRange(b.pixels, func(p RGB) int { return p.G })
	bRange := channelRange(b.pixels, func(p RGB) int { return p.B })

	switch {
	case rRange >= gRange && rRange >= bRange:
		sort.Slice(b.pixels, func(i, j int) bool {
			return b.pixels[i].R < b.pixels[j].R
		})
	case gRange >= bRange:
		sort.Slice(b.pixels, func(i, j int) bool {
			return b.pixels[i].G < b.pixels[j].G
		})
	default:
		sort.Slice(b.pixels, func(i, j int) bool {
			return b.pixels[i].B < b.pixels[j].B
		})
	}

	mid := len(b.pixels) / 2
	return bucket{pixels: b.pixels[:mid]}, bucket{pixels: b.pixels[mid:]}
}

func channelRange(pixels []RGB, selector func(RGB) int) int {
	min, max := math.MaxInt32, math.MinInt32
	for _, p := range pixels {
		v := selector(p)
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
	}
	return max - min
}

func averageColor(pixels []RGB) RGB {
	if len(pixels) == 0 {
		return RGB{}
	}

	var sumR, sumG, sumB int
	for _, p := range pixels {
		sumR += p.R
		sumG += p.G
		sumB += p.B
	}

	n := len(pixels)
	return RGB{
		R: sumR / n,
		G: sumG / n,
		B: sumB / n,
	}
}

func (c RGB) toHex() string {
	return fmt.Sprintf("#%02x%02x%02x", c.R, c.G, c.B)
}

func (c RGB) toCSSColor() color.RGBA {
	return color.RGBA{R: uint8(c.R), G: uint8(c.G), B: uint8(c.B), A: 255}
}
