# CSS Design System

This project uses a centralized design system with CSS custom properties (variables) defined in `src/index.css`.

## ⚠️ IMPORTANT: Always Use CSS Variables

**DO NOT use hardcoded values** for spacing, font sizes, colors, etc. **ALWAYS use the CSS variables** defined below to ensure consistency across all components.

## CSS Variables Reference

### Spacing Scale (Ultra-Compact Design)
```css
var(--space-xs)   /* 0.25rem = 4px    - minimal gaps between elements */
var(--space-sm)   /* 0.4rem = 6.4px   - small spacing (cards, buttons) */
var(--space-md)   /* 0.6rem = 9.6px   - medium spacing (containers) */
var(--space-lg)   /* 0.8rem = 12.8px  - large spacing (sections) */
```

### Font Sizes (Ultra-Compact Design)
```css
var(--font-xs)    /* 0.7rem = 11.2px  - hints, footnotes, timestamps */
var(--font-sm)    /* 0.8rem = 12.8px  - labels, table headers */
var(--font-md)    /* 0.85rem = 13.6px - body text, inputs */
var(--font-lg)    /* 0.9rem = 14.4px  - emphasized text */
var(--font-xl)    /* 1rem = 16px      - section headers (h3) */
var(--font-2xl)   /* 1.15rem = 18.4px - page titles (h2) */
var(--font-3xl)   /* 1.25rem = 20px   - main headers (h1) */
```

### Border Radius
```css
var(--radius-sm)  /* 2px - buttons, small elements */
var(--radius-md)  /* 3px - cards, containers */
```

### Component Sizes
```css
var(--btn-padding-sm)  /* 0.25rem 0.5rem   - small icon buttons */
var(--btn-padding-md)  /* 0.3rem 0.6rem    - medium buttons */
var(--btn-padding-lg)  /* 0.3rem 0.85rem   - large text buttons */
var(--input-padding)   /* 0.3rem 0.45rem   - input fields */
var(--icon-size-sm)    /* 24px - small icon buttons */
var(--icon-size-md)    /* 28px - medium icon buttons */
```

### Line Heights
```css
var(--line-height-tight)   /* 1.3 - compact lists */
var(--line-height-normal)  /* 1.35 - body text */
var(--line-height-relaxed) /* 1.4 - comfortable reading */
```

### Shadows
```css
var(--shadow-sm)  /* 0 1px 2px rgba(0, 0, 0, 0.08) - subtle depth */
var(--shadow-md)  /* 0 1px 3px rgba(0, 0, 0, 0.1)  - card elevation */
```

### Colors
```css
var(--color-primary)
var(--color-primary-dark)
var(--color-bg)
var(--color-bg-secondary)
var(--color-text)
var(--color-text-secondary)
var(--color-border)
var(--color-card-bg)
var(--color-error)
var(--color-error-dark)
var(--color-error-bg)
```

## Usage Examples

### ✅ CORRECT - Using CSS Variables
```css
.my-component {
  padding: var(--space-md);
  font-size: var(--font-md);
  border-radius: var(--radius-md);
  box-shadow: var(--shadow-sm);
}

.my-button {
  padding: var(--btn-padding-md);
  font-size: var(--font-md);
  border-radius: var(--radius-sm);
}
```

### ❌ WRONG - Hardcoded Values
```css
.my-component {
  padding: 12px;           /* BAD! Use var(--space-md) */
  font-size: 14px;         /* BAD! Use var(--font-md) */
  border-radius: 4px;      /* BAD! Use var(--radius-md) */
  box-shadow: 0 1px 2px;   /* BAD! Use var(--shadow-sm) */
}
```

## Why Use CSS Variables?

1. **Consistency**: All components use the same spacing and sizing
2. **Maintainability**: Change once, update everywhere
3. **Prevents Drift**: No accidental inconsistencies
4. **Easy Theming**: Can adjust entire UI by changing variables
5. **Responsive**: Can override variables at different breakpoints

## When Adding New Components

1. **Check existing variables first** - 99% of cases are covered
2. **Use semantic variable names** - `var(--space-md)` not `12px`
3. **If you need a new size**, discuss adding it to `index.css`
4. **Test at different zoom levels** to ensure rem units scale properly

## Modifying the Design System

To change spacing/sizing globally:

1. Edit variables in `src/index.css`
2. All components update automatically
3. No need to touch individual CSS files

**Example**: To make entire UI more compact, change:
```css
--space-md: 0.75rem;  /* change to 0.6rem */
```

This will update ALL components using `var(--space-md)`.

## Component-Specific Guidelines

### Buttons
- Small icon buttons: `var(--btn-padding-sm)` + `var(--icon-size-sm)`
- Medium buttons: `var(--btn-padding-md)`
- Large buttons: `var(--btn-padding-lg)`

### Cards/Containers
- Padding: `var(--space-md)` or `var(--space-lg)`
- Border-radius: `var(--radius-md)`
- Shadow: `var(--shadow-sm)`

### Forms
- Input padding: `var(--input-padding)`
- Label margin-bottom: `var(--space-xs)` or `var(--space-sm)`
- Form-group margin-bottom: `var(--space-sm)`

### Lists
- Item gap: `var(--space-sm)`
- Item padding: `var(--space-sm)` or `var(--space-md)`

## Design Philosophy

This project follows an **ultra-compact, information-dense design**:
- Minimize whitespace while maintaining readability
- Consistent spacing creates visual rhythm
- Very small font sizes (11-20px) for maximum screen real estate
- Very tight line heights (1.3-1.4) for dense information
- Minimal border radii (2-3px) for modern, clean look
- Every pixel counts - compact design maximizes visible content

---

**Remember**: When in doubt, use a CSS variable. Never hardcode spacing or sizing values!
