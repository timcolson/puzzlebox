# Parametric Rectangle Visualizer - Complete Context

## Overview
A Go application using Fyne GUI toolkit that creates a real-time parametric rectangle drawing tool with live SVG generation and export capabilities. The application provides immediate visual feedback as parameters change and generates clean, standards-compliant SVG output using transform-based positioning.

## Core Architecture

### Key Components
1. **Real-time GUI controls** - Synchronized sliders and text inputs for parameter adjustment
2. **Live visual preview** - Fyne canvas objects that update instantly
3. **Dynamic SVG generation** - Clean SVG output with transform-based positioning
4. **Export functionality** - Save to file and copy to clipboard capabilities

### Dependencies
```bash
go mod init parametric-rectangle
go get fyne.io/fyne/v2@latest
```

Required imports:
```go
import (
    "fmt"
    "strconv"
    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/app"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/widget"
    "fyne.io/fyne/v2/dialog"
    "fyne.io/fyne/v2/canvas"
    "image/color"
)
```

## Data Structure

### SVGRectangle Type
```go
type SVGRectangle struct {
    Width, Height float64
    StrokeWidth   float64
    StrokeColor   string
    FillColor     string
}

func (r SVGRectangle) ToSVG() string {
    return fmt.Sprintf(`<rect x="0" y="0" width="%.1f" height="%.1f" 
        stroke="%s" stroke-width="%.1f" fill="%s" />`,
        r.Width, r.Height, r.StrokeColor, r.StrokeWidth, r.FillColor)
}
```

### Key Design Principles
- **Transform-based positioning**: Rectangle drawn at origin (0,0), positioned via SVG transform
- **Parametric properties**: Only dimensions are variable, position is fixed via transform
- **Clean separation**: Data structure, visual representation, and SVG generation are decoupled

## SVG Generation Pattern

### Transform-Based Positioning
- Rectangle always drawn at origin `(0,0)` in local coordinates
- Uses `<g transform="translate(20, 20)">` wrapper for positioning
- All child elements (labels, dimension lines) contained within transform group
- Easy to modify positioning by changing single transform value

### Complete SVG Structure
```xml
<?xml version="1.0" encoding="UTF-8"?>
<svg width="600" height="400" xmlns="http://www.w3.org/2000/svg">
    <!-- Grid background for visual reference -->
    <defs>
        <pattern id="grid" width="20" height="20" patternUnits="userSpaceOnUse">
            <path d="M 20 0 L 0 0 0 20" fill="none" stroke="#e0e0e0" stroke-width="0.5"/>
        </pattern>
    </defs>
    <rect width="100%" height="100%" fill="url(#grid)" />
    
    <!-- Transform group with origin at (20,20) -->
    <g transform="translate(20, 20)">
        <!-- Parametric rectangle at transformed origin -->
        <rect x="0" y="0" width="200" height="100" ... />
        
        <!-- Dimension labels and annotation lines -->
        <text x="..." y="..." ...>L = 200</text>
        <line x1="..." y1="..." x2="..." y2="..." .../>
    </g>
</svg>
```

### SVG Features
- **Grid background**: 20px grid pattern for spatial reference
- **Dimension annotations**: Dynamic labels showing current parameter values
- **Dimension lines**: Red dashed lines indicating measurements
- **Clean structure**: Proper XML declaration, namespaces, and organization

## Real-time Update System

### Update Function Pattern
```go
updateVisualization := func() {
    // 1. Update data structure
    svgRect.Width = length
    svgRect.Height = width
    
    // 2. Update visual elements
    visualRect.Resize(fyne.NewSize(float32(length), float32(width)))
    
    // 3. Update dimension annotations
    lengthLabel.Text = fmt.Sprintf("L = %.0f", length)
    lengthLabel.Move(fyne.NewPos(float32(20+length/2-15), float32(0)))
    lengthLabel.Refresh()
    
    // 4. Regenerate SVG code
    newSVG := generateSVG(svgRect, 600, 400)
    svgDisplay.ParseMarkdown("```xml\n" + newSVG + "\n```")
    
    // 5. Update info display
    infoLabel.SetText(fmt.Sprintf("Rectangle: %.0f × %.0f\nArea: %.0f", 
        length, width, length*width))
}
```

### Control Binding Pattern
```go
// Bidirectional binding between sliders and text entries
lengthSlider.OnChanged = func(value float64) {
    length = value
    lengthEntry.SetText(fmt.Sprintf("%.0f", length))
    updateVisualization()
}

lengthEntry.OnChanged = func(text string) {
    if val, err := strconv.ParseFloat(text, 64); err == nil && val > 0 {
        length = val
        lengthSlider.SetValue(length)
        updateVisualization()
    }
}
```

## Visual Components

### Canvas Objects
- **Main rectangle**: `canvas.NewRectangle()` positioned at fixed (20,20) to match SVG transform
- **Dimension labels**: `canvas.NewText()` with red color, dynamically positioned
- **Dimension lines**: `canvas.NewLine()` with dashed stroke style
- **Grid background**: Programmatically generated lines with 20px spacing

### Visual Properties
- **Rectangle fill**: Semi-transparent blue (RGBA{70, 130, 180, 80})
- **Rectangle stroke**: Dark gray, 2px width
- **Dimension annotations**: Red color for visibility
- **Grid**: Light gray (#e0e0e0) for subtle reference

### Layout Structure
```go
// Three-panel horizontal split layout
mainContent := container.NewHSplit(
    controlPanel,                    // 25% - Parameter controls
    container.NewHSplit(
        visualCanvas,                // 37.5% - Real-time preview
        svgCodeDisplay,              // 37.5% - Generated SVG code
    ),
)
```

## Current Features

### Parameter Controls
- **Length**: Range 10-400 pixels with slider and text entry
- **Width**: Range 10-300 pixels with slider and text entry
- **Bidirectional binding**: Changes in slider update text field and vice versa
- **Input validation**: Only accepts positive numeric values
- **Real-time updates**: Immediate visual and SVG code changes

### Export Capabilities
- **File Export**: Save SVG with native file dialog
- **Clipboard Copy**: Copy SVG code to system clipboard
- **Live Preview**: Generated SVG code displayed with syntax highlighting
- **Standards Compliant**: Clean, well-formatted SVG output

### Visual Feedback Systems
- **Grid background**: 20px grid for spatial reference
- **Dimension labels**: Show current L and W values
- **Dimension lines**: Red dashed lines indicating measurements
- **Info panel**: Displays dimensions and calculated area
- **Reset functionality**: Return to default values instantly

## Installation & Usage

### Setup Commands
```bash
mkdir parametric-rectangle
cd parametric-rectangle
go mod init parametric-rectangle
go get fyne.io/fyne/v2@latest
```

### Running the Application
```bash
go run main.go
```

### Potential macOS Warning
```bash
# If you see "ld: warning: ignoring duplicate libraries: '-lobjc'"
# This is harmless - the app will work perfectly
# To suppress: CGO_LDFLAGS="-Wl,-no_warn_duplicate_libraries" go run main.go
```

## Extensibility Points

### Easy Additions
1. **Additional shapes**: Circle, polygon, path, line elements
2. **More parameters**: Rotation angle, corner radius, stroke styles
3. **Color controls**: Fill color, stroke color with color pickers
4. **Multiple objects**: Array of shapes with individual parameter sets
5. **Animation**: Time-based parameter changes for motion
6. **Import functionality**: Load existing SVG files for editing
7. **Export formats**: PNG, PDF, or other vector formats

### Pattern for New Shapes
```go
type SVGCircle struct {
    Radius      float64
    StrokeWidth float64
    StrokeColor string
    FillColor   string
}

func (c SVGCircle) ToSVG() string {
    return fmt.Sprintf(`<circle cx="0" cy="0" r="%.1f" 
        stroke="%s" stroke-width="%.1f" fill="%s" />`, 
        c.Radius, c.StrokeColor, c.StrokeWidth, c.FillColor)
}
```

### Architecture for Complex Shapes
```go
type ParametricShape interface {
    ToSVG() string
    GetBounds() (width, height float64)
    UpdateFromParams(params map[string]float64)
}
```

## Architecture Benefits

### Design Advantages
1. **Clean separation of concerns**: Data ↔ Visual ↔ SVG generation are independent
2. **Transform-based positioning**: Easy coordinate system management and modification
3. **Real-time feedback**: Immediate visual and code updates enhance user experience
4. **Extensible patterns**: Simple, consistent approach for adding new shapes and parameters
5. **Export ready**: Generates clean, standards-compliant SVG suitable for any use
6. **Type safety**: Go's type system prevents many common parameter errors

### Code Organization
- **Modular structure**: Each component has clear responsibilities
- **Consistent patterns**: Same approach for all parameter bindings and updates
- **Error handling**: Input validation and graceful error recovery
- **Performance**: Efficient updates only refresh changed elements

## Technical Considerations

### Performance Optimization
- **Selective updates**: Only modified visual elements are refreshed
- **Efficient SVG generation**: Template-based string formatting
- **Minimal allocations**: Reuse of objects where possible

### Cross-platform Compatibility
- **Fyne framework**: Native look and feel on macOS, Windows, Linux
- **Standard SVG**: Output compatible with all modern browsers and vector graphics tools
- **Go ecosystem**: Easy deployment and distribution

The clean architecture makes it straightforward to combine with other functionality while maintaining the real-time visualization and export capabilities.