package box

import (
	"42clients.com/puzzlebox/pkg/pathbuilder"
)

// Box represents dimensions as floats to simplify calculations.
// Internal calculations return floats.
// Public calculations output integers for ease of use by external drawing commands which expect ints.
type Box struct {
	Width            float64
	Depth            float64
	Height           float64
	FoldGap          float64
	BottomTabPercent float64
}

// NewBox creates a new Box with default values
func NewBox(width, depth, height, foldGap float64) *Box {
	return &Box{
		Width:            width,
		Depth:            depth,
		Height:           height,
		FoldGap:          foldGap,
		BottomTabPercent: 0.50,
	}
}

// TotalWidth dimension calculations
func (b Box) TotalWidth() int {
	totalWidth := int((2.0 * b.Width) + (2.0 * b.Depth) + b.SideFlapWidth())
	return totalWidth
}
func (b Box) W() int { return int(b.Width) }
func (b Box) D() int { return int(b.Depth) }
func (b Box) H() int { return int(b.Height) }

func (b Box) SideFlapWidth() float64 {
	return .25 * b.Depth
}

func (b Box) TotalHeight() int {
	return int(b.TopFlapHeight() + b.Depth + b.Height + b.BottomFlapMaxHeight())
}

func (b Box) BottomFlapMaxHeight() float64 {
	return b.BottomTabPercent * b.Depth
}

func (b Box) TopFlapHeight() float64 {
	return .2 * b.Depth
}

// Position calculations
func (b Box) BackRight() int {
	return int((2 * b.Width) + (2 * b.Depth))
}

func (b Box) Bottom() int {
	return int(b.TopFlapHeight() + b.Depth + b.Height)
}

func (b Box) Top() int {
	return int(b.TopFlapHeight() + b.Depth)
}

func (b Box) SideALeft() int {
	return 0
}

func (b Box) SideFlapHeight() float64 {
	return b.Depth * 0.66
}

// Panel positions for more complex layouts
func (b Box) FrontPanelLeft() int {
	return int(b.Depth)
}

func (b Box) FrontPanelRight() int {
	return int(b.Depth + b.Width)
}

func (b Box) BackPanelLeft() int {
	return int(b.Depth + b.Width + b.Depth)
}

func (b Box) BackPanelRight() int {
	return int(b.Depth + b.Width + b.Depth + b.Width)
}

// PathBuilder integration for generating various box components

// GenerateFoldLines creates the main fold lines for the box
func (b Box) GenerateFoldLines() string {
	builder := pathbuilder.NewAdvancedPathBuilder()

	originX := b.SideALeft()
	originY := b.Bottom()
	mainBoxWidth := b.TotalWidth() - int(b.SideFlapWidth())

	// Main box outline fold lines
	// Bottom horizontal fold line
	builder = builder.
		MoveTo(originX, originY).
		HorizontalLine(mainBoxWidth).Square().
		VerticalLine(-1 * b.H()).Square(). // Right vertical fold line
		HorizontalLine(-1 * mainBoxWidth).Square()

	return builder.Build()
}

// GenerateCutLines creates cut lines for side flaps using relative coordinates
func (b Box) GenerateCutLines() string {
	builder := pathbuilder.NewAdvancedPathBuilder()

	// Left side flap - simple rectangular flap
	leftFlapPath := builder.
		MoveTo(0, b.Top()).
		VerticalLine(int(b.Height)).Square().
		HorizontalLine(int(b.Depth)).Square().
		VerticalLine(-int(b.Height)).Square().
		Build()

	return leftFlapPath
}

// GenerateBottomFlaps creates the bottom flaps for the box
func (b Box) GenerateBottomFlaps() string {
	builder := pathbuilder.NewAdvancedPathBuilder()

	// Front bottom flap - simple rectangular flap
	frontBottomPath := builder.
		MoveTo(int(b.Depth), b.Bottom()).
		VerticalLine(int(b.BottomFlapMaxHeight())).Square().
		HorizontalLine(int(b.Width)).Square().
		VerticalLine(-int(b.BottomFlapMaxHeight())).Square().
		Build()

	return frontBottomPath
}

// GenerateTopFlaps creates the top flaps for the box
func (b Box) GenerateTopFlaps() string {
	builder := pathbuilder.NewAdvancedPathBuilder()

	// Front top flap - simple rectangular flap
	frontTopPath := builder.
		MoveTo(int(b.Depth), b.Top()).
		VerticalLine(-int(b.TopFlapHeight())).Square().
		HorizontalLine(int(b.Width)).Square().
		VerticalLine(int(b.TopFlapHeight())).Square().
		Build()

	return frontTopPath
}

// GenerateSideTabs creates tabs for gluing sides
func (b Box) GenerateSideTabs() string {
	builder := pathbuilder.NewAdvancedPathBuilder()

	// Right side tab - simple rectangular tab
	rightTabPath := builder.
		MoveTo(int(b.Depth+b.Width+b.Depth), b.Top()).
		HorizontalLine(int(b.SideFlapWidth())).Square().
		VerticalLine(int(b.Height)).Square().
		HorizontalLine(-int(b.SideFlapWidth())).Square().
		Build()

	return rightTabPath
}

// GenerateCompleteBox creates all paths for a complete box template
func (b Box) GenerateCompleteBox() map[string]string {
	return map[string]string{
		"fold_lines": b.GenerateFoldLines(),
		"cut_lines":  b.GenerateCutLines(),
		"rect":       pathbuilder.CreateRectangle(30, 20, 50, 30, 1),
	}
}

// Validation methods
func (b Box) IsValid() bool {
	return b.Width > 0 && b.Depth > 0 && b.Height > 0 && b.FoldGap >= 0
}
