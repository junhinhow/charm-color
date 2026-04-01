// Package color fornece utilitarios de manipulacao de cores para UIs de terminal.
package color

import (
	"fmt"
	"image/color"
	"math"
	"strings"

	"charm.land/lipgloss/v2"
)

// Darken escurece a cor pelo fator amount (0.0 a 1.0).
func Darken(c color.Color, amount float64) color.Color {
	r, g, b, a := c.RGBA()
	factor := clamp(1.0 - amount)
	return color.NRGBA{
		R: uint8(float64(r>>8) * factor),
		G: uint8(float64(g>>8) * factor),
		B: uint8(float64(b>>8) * factor),
		A: uint8(a >> 8),
	}
}

// Lighten clareia a cor pelo fator amount (0.0 a 1.0).
func Lighten(c color.Color, amount float64) color.Color {
	r, g, b, a := c.RGBA()
	rf := float64(r >> 8)
	gf := float64(g >> 8)
	bf := float64(b >> 8)
	amt := clamp(amount)
	return color.NRGBA{
		R: uint8(rf + (255-rf)*amt),
		G: uint8(gf + (255-gf)*amt),
		B: uint8(bf + (255-bf)*amt),
		A: uint8(a >> 8),
	}
}

// Blend interpola entre c1 e c2 pelo fator t (0.0 = c1, 1.0 = c2).
func Blend(c1, c2 color.Color, t float64) color.Color {
	t = clamp(t)
	r1, g1, b1, a1 := c1.RGBA()
	r2, g2, b2, a2 := c2.RGBA()
	return color.NRGBA{
		R: uint8(lerp(float64(r1>>8), float64(r2>>8), t)),
		G: uint8(lerp(float64(g1>>8), float64(g2>>8), t)),
		B: uint8(lerp(float64(b1>>8), float64(b2>>8), t)),
		A: uint8(lerp(float64(a1>>8), float64(a2>>8), t)),
	}
}

// Contrast retorna preto ou branco dependendo da luminancia da cor,
// garantindo legibilidade do texto.
func Contrast(c color.Color) color.Color {
	r, g, b, _ := c.RGBA()
	// Luminancia relativa (formula WCAG)
	lum := 0.2126*srgbToLinear(float64(r>>8)/255) +
		0.7152*srgbToLinear(float64(g>>8)/255) +
		0.0722*srgbToLinear(float64(b>>8)/255)
	if lum > 0.179 {
		return color.NRGBA{R: 0, G: 0, B: 0, A: 255}
	}
	return color.NRGBA{R: 255, G: 255, B: 255, A: 255}
}

// Complementary retorna a cor complementar (oposta na roda de cores).
func Complementary(c color.Color) color.Color {
	r, g, b, a := c.RGBA()
	return color.NRGBA{
		R: uint8(255 - r>>8),
		G: uint8(255 - g>>8),
		B: uint8(255 - b>>8),
		A: uint8(a >> 8),
	}
}

// Palette gera uma paleta harmoniosa de n cores a partir de uma cor base,
// rotacionando o hue uniformemente.
func Palette(base color.Color, n int) []color.Color {
	if n <= 0 {
		return nil
	}
	h, s, l := rgbToHSL(base)
	palette := make([]color.Color, n)
	for i := 0; i < n; i++ {
		hue := math.Mod(h+float64(i)*(360.0/float64(n)), 360)
		palette[i] = hslToRGB(hue, s, l)
	}
	return palette
}

// ToHex converte uma cor para representacao hexadecimal (#RRGGBB).
func ToHex(c color.Color) string {
	r, g, b, _ := c.RGBA()
	return fmt.Sprintf("#%02x%02x%02x", r>>8, g>>8, b>>8)
}

// FromHex converte uma string hexadecimal (#RRGGBB ou RRGGBB) para color.Color.
// Retorna preto se o formato for invalido.
func FromHex(hex string) color.Color {
	hex = strings.TrimPrefix(hex, "#")
	if len(hex) != 6 {
		return color.NRGBA{A: 255}
	}
	var r, g, b uint8
	_, err := fmt.Sscanf(hex, "%02x%02x%02x", &r, &g, &b)
	if err != nil {
		return color.NRGBA{A: 255}
	}
	return color.NRGBA{R: r, G: g, B: b, A: 255}
}

// FromLipgloss converte uma string hex/ANSI usando lipgloss.Color para color.Color.
// Util para integrar com estilos lipgloss.
func FromLipgloss(s string) color.Color {
	return lipgloss.Color(s)
}

// --- funcoes auxiliares ---

func clamp(v float64) float64 {
	if v < 0 {
		return 0
	}
	if v > 1 {
		return 1
	}
	return v
}

func lerp(a, b, t float64) float64 {
	return a + (b-a)*t
}

func srgbToLinear(v float64) float64 {
	if v <= 0.04045 {
		return v / 12.92
	}
	return math.Pow((v+0.055)/1.055, 2.4)
}

func rgbToHSL(c color.Color) (h, s, l float64) {
	r, g, b, _ := c.RGBA()
	rf := float64(r>>8) / 255
	gf := float64(g>>8) / 255
	bf := float64(b>>8) / 255

	max := math.Max(rf, math.Max(gf, bf))
	min := math.Min(rf, math.Min(gf, bf))
	l = (max + min) / 2

	if max == min {
		return 0, 0, l
	}

	d := max - min
	if l > 0.5 {
		s = d / (2 - max - min)
	} else {
		s = d / (max + min)
	}

	switch max {
	case rf:
		h = (gf - bf) / d
		if gf < bf {
			h += 6
		}
	case gf:
		h = (bf-rf)/d + 2
	case bf:
		h = (rf-gf)/d + 4
	}
	h *= 60
	return h, s, l
}

func hslToRGB(h, s, l float64) color.Color {
	if s == 0 {
		v := uint8(l * 255)
		return color.NRGBA{R: v, G: v, B: v, A: 255}
	}

	var q float64
	if l < 0.5 {
		q = l * (1 + s)
	} else {
		q = l + s - l*s
	}
	p := 2*l - q

	hNorm := h / 360

	r := hueToRGB(p, q, hNorm+1.0/3.0)
	g := hueToRGB(p, q, hNorm)
	b := hueToRGB(p, q, hNorm-1.0/3.0)

	return color.NRGBA{
		R: uint8(math.Round(r * 255)),
		G: uint8(math.Round(g * 255)),
		B: uint8(math.Round(b * 255)),
		A: 255,
	}
}

func hueToRGB(p, q, t float64) float64 {
	if t < 0 {
		t++
	}
	if t > 1 {
		t--
	}
	if t < 1.0/6.0 {
		return p + (q-p)*6*t
	}
	if t < 1.0/2.0 {
		return q
	}
	if t < 2.0/3.0 {
		return p + (q-p)*(2.0/3.0-t)*6
	}
	return p
}
