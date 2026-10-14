# charm-color

Utilitários de manipulação de cores para UIs de terminal -- escurecer, clarear, misturar, contraste e mais.

> **[Read in English](README.md)**

## Instalação

```bash
go get github.com/junhinhow/charm-color@latest
```

## Uso

O nome do pacote é `color`, o mesmo de `image/color`. Ao importar os dois, dê
um alias para este pacote:

```go
package main

import (
    "fmt"
    "image/color"

    charmcolor "github.com/junhinhow/charm-color"
)

func main() {
    vermelho := color.NRGBA{R: 255, G: 0, B: 0, A: 255}

    escuro := charmcolor.Darken(vermelho, 0.3)
    claro := charmcolor.Lighten(vermelho, 0.3)
    comp := charmcolor.Complementary(vermelho)

    fmt.Println("Escuro:", charmcolor.ToHex(escuro))
    fmt.Println("Claro:", charmcolor.ToHex(claro))
    fmt.Println("Complementar:", charmcolor.ToHex(comp))

    paleta := charmcolor.Palette(vermelho, 5)
    for i, c := range paleta {
        fmt.Printf("  Paleta[%d]: %s\n", i, charmcolor.ToHex(c))
    }
}
```

## Tratamento do alfa

As cores são manipuladas no modelo de **alfa reto** (`color.NRGBA`): o canal
alfa nunca é multiplicado nos canais RGB.

```go
// Os canais RGB são os retos e o alfa é preservado.
charmcolor.ToHex(color.NRGBA{R: 255, G: 0, B: 0, A: 128})            // "#ff0000"
charmcolor.ToHexA(color.NRGBA{R: 255, G: 0, B: 0, A: 128})           // "#ff000080"
charmcolor.ToHex(charmcolor.Darken(color.NRGBA{255, 0, 0, 128}, 0))  // "#ff0000" (identidade)
```

`ToHex` ignora o canal alfa (`#RRGGBB`); use `ToHexA` (`#RRGGBBAA`) para
mantê-lo.

## Leitura de valores hexadecimais

`FromHex` mantém o comportamento original e devolve preto para uma entrada
inválida. `FromHexE` valida o formato e devolve um erro:

```go
c, err := charmcolor.FromHexE("#ff000080") // #RRGGBBAA é aceito
if errors.Is(err, charmcolor.ErrHexInvalido) {
    // trata a entrada inválida
}
```

Formatos aceitos: `#RRGGBB`, `RRGGBB`, `#RRGGBBAA` e `RRGGBBAA`.

## API

| Função | Descrição |
|--------|-----------|
| `Darken(c, amount)` | Escurece uma cor pelo fator (0.0-1.0), preservando o alfa |
| `Lighten(c, amount)` | Clareia uma cor pelo fator (0.0-1.0), preservando o alfa |
| `Blend(c1, c2, t)` | Interpola entre duas cores |
| `Contrast(c)` | Retorna preto ou branco para legibilidade |
| `Complementary(c)` | Retorna a cor complementar (matiz rotacionado em 180°) |
| `Inverse(c)` | Inverte os canais RGB (`255 - canal`) |
| `Palette(base, n)` | Gera n cores harmônicas por rotação de hue |
| `ToHex(c)` | Converte a cor para a string `#RRGGBB` |
| `ToHexA(c)` | Converte a cor para a string `#RRGGBBAA` (com alfa) |
| `FromHex(hex)` | Converte `#RRGGBB` ou `RRGGBB` para cor (preto se inválido) |
| `FromHexE(hex)` | O mesmo, mas devolve `ErrHexInvalido` para formato inválido |
| `FromLipgloss(s)` | Converte uma string hex/ANSI usando `lipgloss.Color` |

## Licença

[MIT](LICENSE)
