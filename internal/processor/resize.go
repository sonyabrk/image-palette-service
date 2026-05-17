package processor

import (
	"image"

	"github.com/disintegration/imaging"
)

const targetSize = 150

func resize(img image.Image) image.Image {
	return imaging.Fit(img, targetSize, targetSize, imaging.Lanczos)
}
