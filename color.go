// Package color fornece utilitarios de manipulacao de cores para UIs de terminal.
//
// As cores são manipuladas no modelo de alfa reto (straight alpha, NRGBA): o
// canal alfa nunca é aplicado aos canais RGB. Isso evita que operações como
// escurecer, clarear e misturar escureçam silenciosamente as cores com
// transparência.
package color

import (
	"errors"
	"fmt"
	"image/color"
	"math"
	"strconv"
	"strings"

	"charm.land/lipgloss/v2"
)

// ErrHexInvalido indica que uma string hexadecimal não pôde ser convertida em
// cor. Use errors.Is para detectá-lo.
var ErrHexInvalido = errors.New("color: valor hexadecimal inválido")

// Darken escurece a cor pelo fator amount (0.0 a 1.0).
//
// A cor é convertida para alfa reto antes da operação, portanto o canal alfa
// é preservado e Darken(c, 0) devolve a própria cor.
func Darken(c color.Color, amount float64) color.Color {
	n := paraNRGBA(c)
	fator := clamp(1.0 - amount)
	return color.NRGBA{
		R: paraUint8(float64(n.R) * fator),
		G: paraUint8(float64(n.G) * fator),
		B: paraUint8(float64(n.B) * fator),
		A: n.A,
	}
}

// Lighten clareia a cor pelo fator amount (0.0 a 1.0).
//
// A cor é convertida para alfa reto antes da operação, portanto o canal alfa
// é preservado e Lighten(c, 0) devolve a própria cor.
func Lighten(c color.Color, amount float64) color.Color {
	n := paraNRGBA(c)
	amt := clamp(amount)
	return color.NRGBA{
		R: paraUint8(float64(n.R) + (255-float64(n.R))*amt),
		G: paraUint8(float64(n.G) + (255-float64(n.G))*amt),
		B: paraUint8(float64(n.B) + (255-float64(n.B))*amt),
		A: n.A,
	}
}

// Blend interpola entre c1 e c2 pelo fator t (0.0 = c1, 1.0 = c2).
func Blend(c1, c2 color.Color, t float64) color.Color {
	t = clamp(t)
	a := paraNRGBA(c1)
	b := paraNRGBA(c2)
	return color.NRGBA{
		R: paraUint8(lerp(float64(a.R), float64(b.R), t)),
		G: paraUint8(lerp(float64(a.G), float64(b.G), t)),
		B: paraUint8(lerp(float64(a.B), float64(b.B), t)),
		A: paraUint8(lerp(float64(a.A), float64(b.A), t)),
	}
}

// Contrast retorna preto ou branco dependendo da luminancia da cor,
// garantindo legibilidade do texto.
func Contrast(c color.Color) color.Color {
	n := paraNRGBA(c)
	// Luminancia relativa (formula WCAG)
	lum := 0.2126*srgbToLinear(float64(n.R)/255) +
		0.7152*srgbToLinear(float64(n.G)/255) +
		0.0722*srgbToLinear(float64(n.B)/255)
	if lum > 0.179 {
		return color.NRGBA{R: 0, G: 0, B: 0, A: 255}
	}
	return color.NRGBA{R: 255, G: 255, B: 255, A: 255}
}

// Complementary retorna a cor complementar, ou seja, a oposta na roda de
// cores (matiz deslocado em 180 graus). O alfa é preservado.
func Complementary(c color.Color) color.Color {
	n := paraNRGBA(c)
	h, s, l := rgbToHSL(n)
	return hslToNRGBA(math.Mod(h+180, 360), s, l, n.A)
}

// Inverse inverte os canais RGB da cor (255 - canal), preservando o alfa.
//
// Esta é a inversão de componentes, não o complementar: use [Complementary]
// para a cor oposta na roda de cores.
func Inverse(c color.Color) color.Color {
	n := paraNRGBA(c)
	return color.NRGBA{
		R: 255 - n.R,
		G: 255 - n.G,
		B: 255 - n.B,
		A: n.A,
	}
}

// Palette gera uma paleta harmoniosa de n cores a partir de uma cor base,
// rotacionando o hue uniformemente. O alfa da cor base é preservado.
func Palette(base color.Color, n int) []color.Color {
	if n <= 0 {
		return nil
	}
	nb := paraNRGBA(base)
	h, s, l := rgbToHSL(nb)
	palette := make([]color.Color, n)
	for i := 0; i < n; i++ {
		hue := math.Mod(h+float64(i)*(360.0/float64(n)), 360)
		palette[i] = hslToNRGBA(hue, s, l, nb.A)
	}
	return palette
}

// ToHex converte uma cor para representacao hexadecimal (#RRGGBB).
//
// O canal alfa é ignorado: os canais RGB são os valores retos (sem
// multiplicação pelo alfa). Use [ToHexA] para incluir o alfa.
func ToHex(c color.Color) string {
	n := paraNRGBA(c)
	return fmt.Sprintf("#%02x%02x%02x", n.R, n.G, n.B)
}

// ToHexA converte uma cor para representacao hexadecimal (#RRGGBBAA),
// incluindo o canal alfa.
func ToHexA(c color.Color) string {
	n := paraNRGBA(c)
	return fmt.Sprintf("#%02x%02x%02x%02x", n.R, n.G, n.B, n.A)
}

// FromHex converte uma string hexadecimal (#RRGGBB ou RRGGBB) para color.Color.
// Retorna preto se o formato for invalido.
//
// Use [FromHexE] para receber o erro do formato inválido.
func FromHex(hex string) color.Color {
	c, err := FromHexE(hex)
	if err != nil {
		return color.NRGBA{A: 255}
	}
	return c
}

// FromHexE converte uma string hexadecimal para color.Color validando o
// formato. São aceitos #RRGGBB, RRGGBB, #RRGGBBAA e RRGGBBAA.
//
// Devolve [ErrHexInvalido] quando a string não representa um hexadecimal
// válido.
func FromHexE(hex string) (color.Color, error) {
	s := strings.TrimPrefix(strings.TrimSpace(hex), "#")
	if len(s) != 6 && len(s) != 8 {
		return nil, fmt.Errorf("%w: %q não tem 6 nem 8 dígitos hexadecimais", ErrHexInvalido, hex)
	}

	valor, err := strconv.ParseUint(s, 16, 64)
	if err != nil {
		return nil, fmt.Errorf("%w: %q não é um número hexadecimal válido", ErrHexInvalido, hex)
	}

	if len(s) == 6 {
		return color.NRGBA{
			R: uint8(valor >> 16),
			G: uint8(valor >> 8),
			B: uint8(valor),
			A: 255,
		}, nil
	}

	return color.NRGBA{
		R: uint8(valor >> 24),
		G: uint8(valor >> 16),
		B: uint8(valor >> 8),
		A: uint8(valor),
	}, nil
}

// FromLipgloss converte uma string hex/ANSI usando lipgloss.Color para color.Color.
// Util para integrar com estilos lipgloss.
func FromLipgloss(s string) color.Color {
	return lipgloss.Color(s)
}

// --- funcoes auxiliares ---

// paraNRGBA converte qualquer color.Color para color.NRGBA preservando o alfa.
//
// O método RGBA() devolve canais pré-multiplicados pelo alfa; o modelo NRGBA
// devolve os canais retos, que são os que devem ser usados nas operações de
// cor.
func paraNRGBA(c color.Color) color.NRGBA {
	switch v := c.(type) {
	case nil:
		return color.NRGBA{A: 255}
	case color.NRGBA:
		return v
	case color.NRGBA64:
		return color.NRGBA{
			R: uint8(v.R >> 8),
			G: uint8(v.G >> 8),
			B: uint8(v.B >> 8),
			A: uint8(v.A >> 8),
		}
	}
	return color.NRGBAModel.Convert(c).(color.NRGBA)
}

// paraUint8 converte um valor de canal para uint8 arredondando e limitando ao
// intervalo válido.
func paraUint8(v float64) uint8 {
	r := math.Round(v)
	if r <= 0 {
		return 0
	}
	if r >= 255 {
		return 255
	}
	return uint8(r)
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
	n := paraNRGBA(c)
	rf := float64(n.R) / 255
	gf := float64(n.G) / 255
	bf := float64(n.B) / 255

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

// hslToNRGBA converte HSL para uma cor de alfa reto com o alfa informado.
func hslToNRGBA(h, s, l float64, a uint8) color.NRGBA {
	if s == 0 {
		v := paraUint8(l * 255)
		return color.NRGBA{R: v, G: v, B: v, A: a}
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
		R: paraUint8(r * 255),
		G: paraUint8(g * 255),
		B: paraUint8(b * 255),
		A: a,
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
