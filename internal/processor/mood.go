package processor

import "math"

type MoodCoordinates struct {
	MoodX      float64
	MoodY      float64
	Brightness float64
	Saturation float64
}

func calculateMood(palette []RGB) MoodCoordinates {
	if len(palette) == 0 {
		return MoodCoordinates{MoodX: 0.5, MoodY: 0.5}
	}

	var totalBrightness, totalSaturation float64

	for _, color := range palette {
		b, s := brightnessAndSaturation(color)
		totalBrightness += b
		totalSaturation += s
	}

	avgBrightness := totalBrightness / float64(len(palette))
	avgSaturation := totalSaturation / float64(len(palette))

	moodY := avgBrightness
	moodX := calculateMoodX(palette, avgSaturation)

	return MoodCoordinates{
		MoodX:      clamp(moodX),
		MoodY:      clamp(moodY),
		Brightness: avgBrightness,
		Saturation: avgSaturation,
	}
}

func brightnessAndSaturation(c RGB) (brightness, saturation float64) {
	r := float64(c.R) / 255.0
	g := float64(c.G) / 255.0
	b := float64(c.B) / 255.0

	maxC := math.Max(r, math.Max(g, b))
	minC := math.Min(r, math.Min(g, b))

	brightness = (maxC + minC) / 2.0

	if maxC == minC {
		saturation = 0
		return
	}

	delta := maxC - minC
	if brightness > 0.5 {
		saturation = delta / (2.0 - maxC - minC)
	} else {
		saturation = delta / (maxC + minC)
	}

	return
}

func calculateMoodX(palette []RGB, avgSaturation float64) float64 {
	significantColors := countSignificantColors(palette)
	diversityScore := float64(significantColors-1) / float64(len(palette)-1)

	saturationWeight := 0.2 + avgSaturation*0.8

	contrastScore := paletteContrast(palette)

	return (diversityScore*saturationWeight)*0.40 +
		contrastScore*0.35 +
		avgSaturation*0.25
}

func paletteContrast(palette []RGB) float64 {
	if len(palette) == 0 {
		return 0
	}

	minBrightness := 1.0
	maxBrightness := 0.0

	for _, c := range palette {
		b, _ := brightnessAndSaturation(c)
		if b < minBrightness {
			minBrightness = b
		}
		if b > maxBrightness {
			maxBrightness = b
		}
	}

	return maxBrightness - minBrightness
}

func countSignificantColors(palette []RGB) int {
	if len(palette) == 0 {
		return 0
	}

	significant := []RGB{palette[0]}
	for _, candidate := range palette[1:] {
		isSignificant := true
		for _, existing := range significant {
			if colorDistance(candidate, existing) < 50 {
				isSignificant = false
				break
			}
		}

		if isSignificant {
			significant = append(significant, candidate)
		}
	}

	return len(significant)
}

func colorDistance(a, b RGB) float64 {
	dr := float64(a.R - b.R)
	dg := float64(a.G - b.G)
	db := float64(a.B - b.B)

	return math.Sqrt(dr*dr + dg*dg + db*db)
}

func clamp(v float64) float64 {
	if v < 0 {
		return 0
	}

	if v > 1 {
		return 1
	}

	return v
}
