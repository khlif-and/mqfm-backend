package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math"
	"strings"
)

func GetInitials(name string) string {
	if name == "" {
		return ""
	}

	parts := strings.Fields(name)
	if len(parts) == 0 {
		return ""
	}

	initials := string([]rune(parts[0])[0])
	if len(parts) > 1 {
		initials += string([]rune(parts[1])[0])
	}
	
	return strings.ToUpper(initials)
}

func GenerateAmbientColor(input string) string {
	hash := md5.Sum([]byte(input))
	hashStr := hex.EncodeToString(hash[:])

	hInt := int(hashStr[0]) + int(hashStr[1])*256
	hue := float64(hInt % 360)

	sInt := int(hashStr[2])
	saturation := 40.0 + (float64(sInt)/255.0)*30.0

	lInt := int(hashStr[3])
	lightness := 10.0 + (float64(lInt)/255.0)*15.0

	r, g, b := hslToRGB(hue, saturation/100.0, lightness/100.0)
	return fmt.Sprintf("#%02x%02x%02x", r, g, b)
}

func hslToRGB(h, s, l float64) (uint8, uint8, uint8) {
	c := (1 - math.Abs(2*l-1)) * s
	x := c * (1 - math.Abs(math.Mod(h/60.0, 2)-1))
	m := l - c/2

	var r, g, b float64

	switch {
	case 0 <= h && h < 60:
		r, g, b = c, x, 0
	case 60 <= h && h < 120:
		r, g, b = x, c, 0
	case 120 <= h && h < 180:
		r, g, b = 0, c, x
	case 180 <= h && h < 240:
		r, g, b = 0, x, c
	case 240 <= h && h < 300:
		r, g, b = x, 0, c
	case 300 <= h && h < 360:
		r, g, b = c, 0, x
	}

	return uint8((r + m) * 255), uint8((g + m) * 255), uint8((b + m) * 255)
}
