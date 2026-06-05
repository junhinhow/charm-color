# charm-color

Color manipulation utilities for terminal UIs -- darken, lighten, blend, contrast, and more.

> **[Leia em Portugues](README.pt-br.md)**

## Install

```bash
go get github.com/junhinhow/charm-color@latest
```

## Usage

```go
package main

import (
    "fmt"
    "image/color"

    cc "github.com/junhinhow/charm-color"
)

func main() {
    red := color.NRGBA{R: 255, G: 0, B: 0, A: 255}

    dark := cc.Darken(red, 0.3)
    light := cc.Lighten(red, 0.3)
    comp := cc.Complementary(red)
    text := cc.Contrast(red)
    hex := cc.ToHex(dark)

    fmt.Println("Dark red:", hex)
    fmt.Println("Light red:", cc.ToHex(light))
    fmt.Println("Complement:", cc.ToHex(comp))
    fmt.Println("Text color:", cc.ToHex(text))

    palette := cc.Palette(red, 5)
    for i, c := range palette {
        fmt.Printf("  Palette[%d]: %s\n", i, cc.ToHex(c))
    }
}
```

## API

| Function | Description |
|----------|-------------|
| `Darken(c, amount)` | Darken a color by a factor (0.0-1.0) |
| `Lighten(c, amount)` | Lighten a color by a factor (0.0-1.0) |
| `Blend(c1, c2, t)` | Interpolate between two colors |
| `Contrast(c)` | Returns black or white for text readability |
| `Complementary(c)` | Returns the complementary (opposite) color |
| `Palette(base, n)` | Generates n harmonious colors by hue rotation |
| `ToHex(c)` | Convert color to `#RRGGBB` string |
| `FromHex(hex)` | Parse `#RRGGBB` or `RRGGBB` to color |

## License

[MIT](LICENSE)
