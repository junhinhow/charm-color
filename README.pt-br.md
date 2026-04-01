# charm-color

Utilitarios de manipulacao de cores para UIs de terminal -- escurecer, clarear, misturar, contraste e mais.

> **[Read in English](README.md)**

## Instalacao

```bash
go get github.com/junhinhow/charm-color@latest
```

## Uso

```go
package main

import (
    "fmt"
    "image/color"

    cc "github.com/junhinhow/charm-color"
)

func main() {
    vermelho := color.NRGBA{R: 255, G: 0, B: 0, A: 255}

    escuro := cc.Darken(vermelho, 0.3)
    claro := cc.Lighten(vermelho, 0.3)
    comp := cc.Complementary(vermelho)

    fmt.Println("Escuro:", cc.ToHex(escuro))
    fmt.Println("Claro:", cc.ToHex(claro))
    fmt.Println("Complementar:", cc.ToHex(comp))

    paleta := cc.Palette(vermelho, 5)
    for i, c := range paleta {
        fmt.Printf("  Paleta[%d]: %s\n", i, cc.ToHex(c))
    }
}
```

## API

| Funcao | Descricao |
|--------|-----------|
| `Darken(c, amount)` | Escurece uma cor pelo fator (0.0-1.0) |
| `Lighten(c, amount)` | Clareia uma cor pelo fator (0.0-1.0) |
| `Blend(c1, c2, t)` | Interpola entre duas cores |
| `Contrast(c)` | Retorna preto ou branco para legibilidade |
| `Complementary(c)` | Retorna a cor complementar (oposta) |
| `Palette(base, n)` | Gera n cores harmoniosas por rotacao de hue |
| `ToHex(c)` | Converte cor para string `#RRGGBB` |
| `FromHex(hex)` | Converte `#RRGGBB` ou `RRGGBB` para cor |

## Licenca

[MIT](LICENSE)
