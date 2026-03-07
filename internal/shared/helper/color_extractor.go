package helper

import (
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"math"
	"os"
	"sort"
)

type colorRGB struct {
	R, G, B float64
}

type colorBucket struct {
	colors []colorRGB
}

func ExtractDominantColor(imagePath string) (string, error) {
	f, err := os.Open(imagePath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return "", err
	}

	pixels := samplePixels(img, 100)
	if len(pixels) == 0 {
		return "#000000", nil
	}

	palette := medianCut(pixels, 5)
	dominant := selectVibrant(palette)

	return rgbToHex(dominant), nil
}

func samplePixels(img image.Image, maxDim int) []colorRGB {
	bounds := img.Bounds()
	w := bounds.Max.X - bounds.Min.X
	h := bounds.Max.Y - bounds.Min.Y

	stepX := 1
	stepY := 1
	if w > maxDim {
		stepX = w / maxDim
	}
	if h > maxDim {
		stepY = h / maxDim
	}

	var pixels []colorRGB
	for y := bounds.Min.Y; y < bounds.Max.Y; y += stepY {
		for x := bounds.Min.X; x < bounds.Max.X; x += stepX {
			r, g, b, a := img.At(x, y).RGBA()
			if a < 0x8000 {
				continue
			}
			pixels = append(pixels, colorRGB{
				R: float64(r >> 8),
				G: float64(g >> 8),
				B: float64(b >> 8),
			})
		}
	}
	return pixels
}

func medianCut(pixels []colorRGB, depth int) []colorRGB {
	if depth == 0 || len(pixels) == 0 {
		return []colorRGB{averageColor(pixels)}
	}

	rRange, gRange, bRange := channelRanges(pixels)

	var sortFn func(i, j int) bool
	maxRange := math.Max(rRange, math.Max(gRange, bRange))

	switch {
	case maxRange == rRange:
		sortFn = func(i, j int) bool { return pixels[i].R < pixels[j].R }
	case maxRange == gRange:
		sortFn = func(i, j int) bool { return pixels[i].G < pixels[j].G }
	default:
		sortFn = func(i, j int) bool { return pixels[i].B < pixels[j].B }
	}

	sort.Slice(pixels, sortFn)

	mid := len(pixels) / 2
	left := medianCut(pixels[:mid], depth-1)
	right := medianCut(pixels[mid:], depth-1)

	return append(left, right...)
}

func channelRanges(pixels []colorRGB) (float64, float64, float64) {
	minR, maxR := 255.0, 0.0
	minG, maxG := 255.0, 0.0
	minB, maxB := 255.0, 0.0

	for _, p := range pixels {
		minR = math.Min(minR, p.R)
		maxR = math.Max(maxR, p.R)
		minG = math.Min(minG, p.G)
		maxG = math.Max(maxG, p.G)
		minB = math.Min(minB, p.B)
		maxB = math.Max(maxB, p.B)
	}

	return maxR - minR, maxG - minG, maxB - minB
}

func averageColor(pixels []colorRGB) colorRGB {
	if len(pixels) == 0 {
		return colorRGB{}
	}

	var sumR, sumG, sumB float64
	for _, p := range pixels {
		sumR += p.R
		sumG += p.G
		sumB += p.B
	}

	n := float64(len(pixels))
	return colorRGB{R: sumR / n, G: sumG / n, B: sumB / n}
}

func selectVibrant(palette []colorRGB) colorRGB {
	if len(palette) == 0 {
		return colorRGB{}
	}

	bestScore := -1.0
	bestColor := palette[0]

	for _, c := range palette {
		h, s, l := rgbToHSL(c)
		_ = h

		satScore := s
		lightPenalty := 1.0 - math.Abs(l-0.45)*2.0
		if lightPenalty < 0 {
			lightPenalty = 0
		}

		score := satScore*0.6 + lightPenalty*0.4

		if score > bestScore {
			bestScore = score
			bestColor = c
		}
	}

	return bestColor
}

func rgbToHSL(c colorRGB) (float64, float64, float64) {
	r := c.R / 255.0
	g := c.G / 255.0
	b := c.B / 255.0

	maxVal := math.Max(r, math.Max(g, b))
	minVal := math.Min(r, math.Min(g, b))
	l := (maxVal + minVal) / 2.0

	if maxVal == minVal {
		return 0, 0, l
	}

	d := maxVal - minVal
	s := d / (1.0 - math.Abs(2.0*l-1.0))

	var h float64
	switch maxVal {
	case r:
		h = math.Mod((g-b)/d, 6.0)
	case g:
		h = (b-r)/d + 2.0
	case b:
		h = (r-g)/d + 4.0
	}
	h *= 60.0
	if h < 0 {
		h += 360.0
	}

	return h, s, l
}

func rgbToHex(c colorRGB) string {
	r := clampByte(c.R)
	g := clampByte(c.G)
	b := clampByte(c.B)

	hex := [7]byte{'#'}
	const hextable = "0123456789abcdef"
	hex[1] = hextable[r>>4]
	hex[2] = hextable[r&0x0f]
	hex[3] = hextable[g>>4]
	hex[4] = hextable[g&0x0f]
	hex[5] = hextable[b>>4]
	hex[6] = hextable[b&0x0f]
	return string(hex[:])
}

func clampByte(v float64) uint8 {
	if v < 0 {
		return 0
	}
	if v > 255 {
		return 255
	}
	return uint8(math.Round(v))
}
