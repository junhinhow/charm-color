package color

import (
	"errors"
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

// --- regressões ---

// TestToHexIgnoraPreMultiplicacao garante que os canais RGB não sejam lidos
// do valor pré-multiplicado devolvido por RGBA().
func TestToHexIgnoraPreMultiplicacao(t *testing.T) {
	// color.NRGBA é alfa reto: RGBA() devolve 32896 (isto é, 128 pré-
	// multiplicado), mas os canais retos são 255, 0 e 0.
	if got := ToHex(color.NRGBA{R: 255, G: 0, B: 0, A: 128}); got != "#ff0000" {
		t.Errorf("ToHex de NRGBA{255,0,0,128} = %s, esperado #ff0000", got)
	}
}

// TestToHexNaoUsaValorPreMultiplicado cobre color.RGBA, que é documentado
// como alfa pré-multiplicado. A tolerância existe porque a precisão da
// conversão depende da implementação de image/color da versão do Go.
func TestToHexNaoUsaValorPreMultiplicado(t *testing.T) {
	got := paraNRGBA(color.RGBA{R: 255, G: 0, B: 0, A: 128})

	if got.G != 0 || got.B != 0 {
		t.Errorf("ToHex de RGBA{255,0,0,128} = %v, esperado canais G e B zerados", got)
	}
	if got.R < 253 || got.R > 255 {
		t.Errorf("ToHex de RGBA{255,0,0,128} = %v, esperado R próximo de 255", got)
	}
	if got.A != 128 {
		t.Errorf("alfa = %d, esperado 128", got.A)
	}
}

// TestDarkenIdentidade verifica que Darken(c, 0) devolve a própria cor.
func TestDarkenIdentidade(t *testing.T) {
	c := color.NRGBA{R: 255, G: 0, B: 0, A: 128}

	escuro := paraNRGBA(Darken(c, 0))
	if escuro != c {
		t.Errorf("Darken(c, 0) = %v, esperado %v", escuro, c)
	}
}

// TestLightenIdentidade verifica que Lighten(c, 0) devolve a própria cor.
func TestLightenIdentidade(t *testing.T) {
	c := color.NRGBA{R: 255, G: 0, B: 0, A: 128}

	claro := paraNRGBA(Lighten(c, 0))
	if claro != c {
		t.Errorf("Lighten(c, 0) = %v, esperado %v", claro, c)
	}
}

func TestDarkenPreservaAlfa(t *testing.T) {
	escuro := paraNRGBA(Darken(color.NRGBA{R: 255, G: 0, B: 0, A: 128}, 0.5))
	if escuro.A != 128 {
		t.Errorf("alfa = %d, esperado 128", escuro.A)
	}
	if escuro.R != 128 {
		t.Errorf("canal R = %d, esperado 128", escuro.R)
	}
}

func TestLightenPreservaAlfa(t *testing.T) {
	claro := paraNRGBA(Lighten(color.NRGBA{R: 255, G: 0, B: 0, A: 128}, 0.5))
	if claro.A != 128 {
		t.Errorf("alfa = %d, esperado 128", claro.A)
	}
}

func TestBlendPreservaAlfa(t *testing.T) {
	a := color.NRGBA{R: 255, G: 0, B: 0, A: 128}
	b := color.NRGBA{R: 0, G: 0, B: 255, A: 0}

	mistura := paraNRGBA(Blend(a, b, 0.5))
	if mistura.A != 64 {
		t.Errorf("alfa = %d, esperado 64", mistura.A)
	}
	if mistura.R != 128 {
		t.Errorf("canal R = %d, esperado 128", mistura.R)
	}
}

// TestArredondamento verifica que as operações arredondam em vez de truncar.
func TestArredondamento(t *testing.T) {
	branco := color.NRGBA{R: 255, G: 255, B: 255, A: 255}

	escuro := paraNRGBA(Darken(branco, 0.5))
	if escuro.R != 128 {
		t.Errorf("Darken(branco, 0.5) = %d, esperado 128 (arredondado)", escuro.R)
	}

	mistura := paraNRGBA(Blend(color.NRGBA{A: 255}, branco, 0.5))
	if mistura.R != 128 {
		t.Errorf("Blend(preto, branco, 0.5) = %d, esperado 128 (arredondado)", mistura.R)
	}
}

// TestContrastIgnoraPreMultiplicacao garante que o contraste use os canais
// retos: um vermelho com transparência não pode virar marrom escuro.
func TestContrastIgnoraPreMultiplicacao(t *testing.T) {
	contraste := paraNRGBA(Contrast(color.NRGBA{R: 255, G: 0, B: 0, A: 128}))
	if contraste.R != 0 || contraste.G != 0 || contraste.B != 0 {
		t.Errorf("Contrast de vermelho = %v, esperado preto", contraste)
	}
}

func TestPalettePreservaAlfa(t *testing.T) {
	paleta := Palette(color.NRGBA{R: 255, G: 0, B: 0, A: 128}, 3)
	for i, c := range paleta {
		if n := paraNRGBA(c); n.A != 128 {
			t.Errorf("alfa da cor %d da paleta = %d, esperado 128", i, n.A)
		}
	}
}

func TestToHexA(t *testing.T) {
	if got := ToHexA(color.NRGBA{R: 255, G: 0, B: 0, A: 128}); got != "#ff000080" {
		t.Errorf("ToHexA = %s, esperado #ff000080", got)
	}
}

func TestFromHexEInvalido(t *testing.T) {
	invalidos := []string{"", "#", "invalid", "#fff", "#12345", "1234567", "#gggggg", "#ff00zz00"}

	for _, entrada := range invalidos {
		if _, err := FromHexE(entrada); !errors.Is(err, ErrInvalidHex) {
			t.Errorf("FromHexE(%q) devolveu %v, esperado ErrInvalidHex", entrada, err)
		}
	}
}

func TestFromHexEComAlfa(t *testing.T) {
	c, err := FromHexE("#ff000080")
	if err != nil {
		t.Fatalf("FromHexE devolveu %v", err)
	}
	n := paraNRGBA(c)
	if n.R != 255 || n.G != 0 || n.B != 0 || n.A != 128 {
		t.Errorf("FromHexE(#ff000080) = %v, esperado (255,0,0,128)", n)
	}
}

func TestFromHexEValido(t *testing.T) {
	c, err := FromHexE("#ff8000")
	if err != nil {
		t.Fatalf("FromHexE devolveu %v", err)
	}
	n := paraNRGBA(c)
	if n.R != 255 || n.G != 128 || n.B != 0 || n.A != 255 {
		t.Errorf("FromHexE(#ff8000) = %v, esperado (255,128,0,255)", n)
	}
}

// TestFromHexInvalidoMantemPreto garante a compatibilidade de FromHex.
func TestFromHexInvalidoMantemPreto(t *testing.T) {
	n := paraNRGBA(FromHex("xyz"))
	if n.R != 0 || n.G != 0 || n.B != 0 {
		t.Errorf("FromHex invalido = %v, esperado preto", n)
	}
}

// TestComplementarReal verifica o complementar no espaço HSL (matiz + 180º).
func TestComplementarReal(t *testing.T) {
	casos := []struct {
		entrada  string
		esperado string
	}{
		{"#402080", "#608020"},
		{"#00a0b0", "#b01000"},
		{"#f0e68c", "#8c96f0"},
	}

	for _, caso := range casos {
		base, err := FromHexE(caso.entrada)
		if err != nil {
			t.Fatalf("FromHexE(%q) devolveu %v", caso.entrada, err)
		}
		if got := ToHex(Complementary(base)); got != caso.esperado {
			t.Errorf("Complementary(%s) = %s, esperado %s", caso.entrada, got, caso.esperado)
		}
	}
}

func TestComplementarPreservaAlfa(t *testing.T) {
	comp := paraNRGBA(Complementary(color.NRGBA{R: 255, G: 0, B: 0, A: 128}))
	if comp.A != 128 {
		t.Errorf("alfa = %d, esperado 128", comp.A)
	}
}

// TestInverse verifica a inversão dos canais RGB.
func TestInverse(t *testing.T) {
	invertida := paraNRGBA(Inverse(color.NRGBA{R: 0x40, G: 0x20, B: 0x80, A: 128}))
	if invertida.R != 0xbf || invertida.G != 0xdf || invertida.B != 0x7f {
		t.Errorf("Inverse = %v, esperado (191,223,127)", invertida)
	}
	if invertida.A != 128 {
		t.Errorf("alfa = %d, esperado 128", invertida.A)
	}
}

func TestOperacoesComCorNaoNula(t *testing.T) {
	// Uma cor nil não deve entrar em panic.
	if got := ToHex(nil); got != "#000000" {
		t.Errorf("ToHex(nil) = %s, esperado #000000", got)
	}
	if got := paraNRGBA(Darken(nil, 0.5)); got.A != 255 {
		t.Errorf("Darken(nil, 0.5) alfa = %d, esperado 255", got.A)
	}
}
