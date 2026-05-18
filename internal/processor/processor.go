package processor

import (
	"bytes"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
)

type Result struct {
	Palette    []string `json:"palette"`
	Dominant   string   `json:"dominant"`
	MoodX      float64  `json:"mood_x"`
	MoodY      float64  `json:"mood_y"`
	Brightness float64  `json:"brightness"`
	Saturation float64  `json:"saturation"`
}

type Processor struct {
	maxClusters int
}

func New(maxClusters int) *Processor {
	return &Processor{maxClusters: maxClusters}
}

func (p *Processor) Analyze(data []byte) (*Result, error) {
	img, err := decodeImage(data)
	if err != nil {
		return nil, fmt.Errorf("не удалось декодировать изображение: %w", err)
	}

	small := resize(img)
	palette := extractPalette(small, p.maxClusters)
	mood := calculateMood(palette)

	return buildResult(palette, mood), nil
}

func decodeImage(data []byte) (image.Image, error) {
	reader := bytes.NewReader(data)
	img, _, err := image.Decode(reader)

	if err != nil {
		return nil, err
	}

	return img, nil
}

func buildResult(palette []RGB, mood MoodCoordinates) *Result {
	hexPalette := make([]string, len(palette))
	for i, c := range palette {
		hexPalette[i] = c.toHex()
	}

	dominant := ""
	if len(hexPalette) > 0 {
		dominant = hexPalette[0]
	}

	return &Result{
		Palette:    hexPalette,
		Dominant:   dominant,
		MoodX:      mood.MoodX,
		MoodY:      mood.MoodY,
		Brightness: mood.Brightness,
		Saturation: mood.Saturation,
	}
}
