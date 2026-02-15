# Self-Hosted Fonts

This directory contains self-hosted font files for offline support.

## Required Font Files

Download the following woff2 files from Google Fonts and place them in this directory:

### DM Serif Display (Editorial serif for headings)
- `dm-serif-display-v15-latin-regular.woff2`
- Download from: https://fonts.google.com/specimen/DM+Serif+Display
- Weight: 400 (Regular)
- Format: woff2

### Plus Jakarta Sans (Geometric sans for UI)
- `plus-jakarta-sans-v8-latin-regular.woff2` (Weight: 400)
- `plus-jakarta-sans-v8-latin-500.woff2` (Weight: 500)
- `plus-jakarta-sans-v8-latin-600.woff2` (Weight: 600)
- `plus-jakarta-sans-v8-latin-700.woff2` (Weight: 700)
- Download from: https://fonts.google.com/specimen/Plus+Jakarta+Sans
- Format: woff2

## Quick Download Method

You can use [google-webfonts-helper](https://gwfh.mranftl.com/fonts) to easily download woff2 files:

1. Visit https://gwfh.mranftl.com/fonts
2. Search for "DM Serif Display"
3. Select "latin" charset
4. Select "regular (400)" weight
5. Download the woff2 file
6. Repeat for "Plus Jakarta Sans" with weights 400, 500, 600, 700

## Fallback Fonts

If the font files are not available, the application will fall back to system fonts:
- Headings: Georgia, serif
- Body: system-ui, -apple-system, sans-serif
- Monospace: JetBrains Mono, Fira Code, Consolas, monospace
