package color

import (
	"image/color"
	"testing"
)

func TestDarken(t *testing.T) {
	c := color.NRGBA{R: 200, G: 100, B: 50, A: 255}
	d := Darken(c, 0.5)
	r, g, b, _ := d.RGBA()
	if r>>8 > 200 || g>>8 > 100 || b>>8 > 50 {
		t.Error("Darken deveria reduzir os valores RGB")
	}
}

func TestLighten(t *testing.T) {
	c := color.NRGBA{R: 100, G: 50, B: 25, A: 255}
	l := Lighten(c, 0.5)
	r, g, b, _ := l.RGBA()
	if r>>8 < 100 || g>>8 < 50 || b>>8 < 25 {
		t.Error("Lighten deveria aumentar os valores RGB")
	}
}

func TestBlend(t *testing.T) {
	black := color.NRGBA{R: 0, G: 0, B: 0, A: 255}
	white := color.NRGBA{R: 255, G: 255, B: 255, A: 255}
	mid := Blend(black, white, 0.5)
	r, _, _, _ := mid.RGBA()
	v := r >> 8
	if v < 120 || v > 135 {
		t.Errorf("Blend(black, white, 0.5) esperado ~127, obtido %d", v)
	}
}

func TestContrastDark(t *testing.T) {
	white := color.NRGBA{R: 255, G: 255, B: 255, A: 255}
	c := Contrast(white)
	r, g, b, _ := c.RGBA()
	if r>>8 != 0 || g>>8 != 0 || b>>8 != 0 {
		t.Error("Contrast de branco deveria retornar preto")
	}
}

func TestContrastLight(t *testing.T) {
	black := color.NRGBA{R: 0, G: 0, B: 0, A: 255}
	c := Contrast(black)
	r, g, b, _ := c.RGBA()
	if r>>8 != 255 || g>>8 != 255 || b>>8 != 255 {
		t.Error("Contrast de preto deveria retornar branco")
	}
}

func TestComplementary(t *testing.T) {
	c := color.NRGBA{R: 255, G: 0, B: 0, A: 255}
	comp := Complementary(c)
	r, g, b, _ := comp.RGBA()
	if r>>8 != 0 || g>>8 != 255 || b>>8 != 255 {
		t.Errorf("Complementar de vermelho deveria ser ciano, obtido (%d,%d,%d)", r>>8, g>>8, b>>8)
	}
}

func TestPalette(t *testing.T) {
	base := color.NRGBA{R: 255, G: 0, B: 0, A: 255}
	p := Palette(base, 4)
	if len(p) != 4 {
		t.Errorf("Palette deveria ter 4 cores, obtido %d", len(p))
	}
}

func TestPaletteZero(t *testing.T) {
	base := color.NRGBA{R: 255, G: 0, B: 0, A: 255}
	p := Palette(base, 0)
	if p != nil {
		t.Error("Palette(0) deveria retornar nil")
	}
}

func TestToHex(t *testing.T) {
	c := color.NRGBA{R: 255, G: 128, B: 0, A: 255}
	hex := ToHex(c)
	if hex != "#ff8000" {
		t.Errorf("ToHex esperado #ff8000, obtido %s", hex)
	}
}

func TestFromHex(t *testing.T) {
	c := FromHex("#ff8000")
	r, g, b, _ := c.RGBA()
	if r>>8 != 255 || g>>8 != 128 || b>>8 != 0 {
		t.Errorf("FromHex #ff8000 esperado (255,128,0), obtido (%d,%d,%d)", r>>8, g>>8, b>>8)
	}
}

func TestFromHexWithoutHash(t *testing.T) {
	c := FromHex("00ff00")
	r, g, b, _ := c.RGBA()
	if r>>8 != 0 || g>>8 != 255 || b>>8 != 0 {
		t.Errorf("FromHex 00ff00 esperado (0,255,0), obtido (%d,%d,%d)", r>>8, g>>8, b>>8)
	}
}

func TestFromHexInvalid(t *testing.T) {
	c := FromHex("invalid")
	r, g, b, _ := c.RGBA()
	if r>>8 != 0 || g>>8 != 0 || b>>8 != 0 {
		t.Error("FromHex invalido deveria retornar preto")
	}
}
