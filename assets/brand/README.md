# Brand Assets

## Twitter/X Profile Images

### Avatar (`twitter-avatar.png` — 400x400px)

**Gemini prompt:**
> Generate a 400x400 pixel Twitter/X profile avatar image. The design should be: a bold monospace text ">" followed by a space and "Hype" centered on a solid dark background. The ">" character must be in bright emerald green color (#10b981). The word "Hype" must be in off-white (#f8f8f2). The background must be solid very dark navy-black (#0d1117). Use a bold monospace font like Menlo or Consolas. The text should be sized to fill about 60-70% of the image width, perfectly centered both horizontally and vertically. No other elements, no gradients, no shadows, no borders, no decorations. Pure flat minimal design. This will be displayed as a circular crop on Twitter so keep all text well within the center.

### Header (`twitter-header.png` — 1500x500px)

**Gemini prompt:**
> Generate a 1500x500 pixel Twitter/X header banner image. Clean, centered design on a solid dark background (#0d1117). In the center, large bold monospace text "> Hype" where ">" is emerald green (#10b981) and "Hype" is off-white (#f8f8f2). Below the logo, smaller monospace text "Your docs are never out of date." in muted gray (#9ca3af). Both elements centered horizontally, positioned slightly above vertical center. No gradients, no shadows, no borders. Pure flat minimal dark mode design.

### Gemini settings

- **Model:** `gemini-3.1-flash-image-preview` (recommended) or `gemini-3-pro-image-preview` (highest quality)
- **Aspect ratio:** `1:1` for avatar, `16:9` closest for header (or use SVG approach below)
- **Resolution:** `1K` or higher

### SVG source files

The PNGs were generated from the SVG source files in this directory using `rsvg-convert`:

```bash
rsvg-convert -w 400 -h 400 twitter-avatar.svg -o twitter-avatar.png
rsvg-convert -w 1500 -h 500 twitter-header.svg -o twitter-header.png
```

### Brand colors

| Color | Hex | Usage |
|-------|-----|-------|
| Emerald green | `#10b981` | Primary accent, `>` symbol |
| Off-white | `#f8f8f2` | Logo text |
| Dark background | `#0d1117` | Background |
| Muted gray | `#9ca3af` | Tagline text |

### Twitter/X upload specs

| Asset | Dimensions | Max size | Format | Notes |
|-------|-----------|----------|--------|-------|
| Avatar | 400x400px | 2MB | PNG | Displayed as circle; keep content centered |
| Header | 1500x500px | 5MB | PNG | Bottom-left 400x200px overlapped by profile pic |
