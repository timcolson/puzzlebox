package main

import (
	"fmt"
	"math"
	"os"
	"strings"

	"github.com/ajstarks/svgo"
)

type CornerType int

const (
	SquareCorner CornerType = iota
	RoundedCorner
)

type PathSegment struct {
	EndX       int
	EndY       int
	CornerType CornerType
	Radius     int
}

type AdvancedPathBuilder struct {
	commands []string
	lastX    int
	lastY    int
	segments []PathSegment
}

type CornerBuilder struct {
	pathBuilder *AdvancedPathBuilder
	endX        int
	endY        int
}

func NewAdvancedPathBuilder() *AdvancedPathBuilder {
	return &AdvancedPathBuilder{
		commands: make([]string, 0),
		segments: make([]PathSegment, 0),
	}
}

func (apb *AdvancedPathBuilder) MoveTo(x, y int) *AdvancedPathBuilder {
	apb.commands = append(apb.commands, fmt.Sprintf("M%d,%d", x, y))
	apb.lastX = x
	apb.lastY = y
	return apb
}

func (apb *AdvancedPathBuilder) LineTo(x, y int) *CornerBuilder {
	return &CornerBuilder{
		pathBuilder: apb,
		endX:        x,
		endY:        y,
	}
}

func (cb *CornerBuilder) Square() *AdvancedPathBuilder {
	cb.pathBuilder.segments = append(cb.pathBuilder.segments, PathSegment{
		EndX:       cb.endX,
		EndY:       cb.endY,
		CornerType: SquareCorner,
	})
	return cb.pathBuilder
}

func (cb *CornerBuilder) Rounded(radius int) *AdvancedPathBuilder {
	cb.pathBuilder.segments = append(cb.pathBuilder.segments, PathSegment{
		EndX:       cb.endX,
		EndY:       cb.endY,
		CornerType: RoundedCorner,
		Radius:     radius,
	})
	return cb.pathBuilder
}

// generateRoundedCorner creates a smooth curve at a corner
func (apb *AdvancedPathBuilder) generateRoundedCorner(prevX, prevY, cornerX, cornerY, nextX, nextY, radius int) []string {
	var commands []string

	// Calculate vectors from corner point
	// Vector from corner to previous point (incoming direction)
	incomingX := float64(prevX - cornerX)
	incomingY := float64(prevY - cornerY)

	// Vector from corner to next point (outgoing direction)
	outgoingX := float64(nextX - cornerX)
	outgoingY := float64(nextY - cornerY)

	// Normalize vectors
	incomingLen := math.Sqrt(incomingX*incomingX + incomingY*incomingY)
	outgoingLen := math.Sqrt(outgoingX*outgoingX + outgoingY*outgoingY)

	// Handle edge cases
	if incomingLen == 0 || outgoingLen == 0 {
		commands = append(commands, fmt.Sprintf("L%d,%d", cornerX, cornerY))
		return commands
	}

	incomingUnitX := incomingX / incomingLen
	incomingUnitY := incomingY / incomingLen
	outgoingUnitX := outgoingX / outgoingLen
	outgoingUnitY := outgoingY / outgoingLen

	radiusF := float64(radius)

	// Ensure radius doesn't exceed half the length of either segment
	maxRadius := math.Min(incomingLen/2, outgoingLen/2)
	if radiusF > maxRadius {
		radiusF = maxRadius
	}

	// Calculate curve start and end points
	// Start point: move back from corner along incoming vector
	curveStartX := cornerX + int(incomingUnitX*radiusF)
	curveStartY := cornerY + int(incomingUnitY*radiusF)

	// End point: move forward from corner along outgoing vector
	curveEndX := cornerX + int(outgoingUnitX*radiusF)
	curveEndY := cornerY + int(outgoingUnitY*radiusF)

	// Line to curve start
	commands = append(commands, fmt.Sprintf("L%d,%d", curveStartX, curveStartY))

	// Quadratic curve through the corner
	commands = append(commands, fmt.Sprintf("Q%d,%d,%d,%d", cornerX, cornerY, curveEndX, curveEndY))

	return commands
}

func (apb *AdvancedPathBuilder) Build() string {
	if len(apb.segments) == 0 {
		return strings.Join(apb.commands, "")
	}

	currentX := apb.lastX
	currentY := apb.lastY

	for i, segment := range apb.segments {
		if i == len(apb.segments)-1 || segment.CornerType == SquareCorner {
			// Last segment or square corner - draw a straight line
			apb.commands = append(apb.commands, fmt.Sprintf("L%d,%d", segment.EndX, segment.EndY))
		} else if segment.CornerType == RoundedCorner && i < len(apb.segments)-1 {
			// Rounded corner - need next segment for curve calculation
			nextSegment := apb.segments[i+1]
			curves := apb.generateRoundedCorner(
				currentX, currentY,
				segment.EndX, segment.EndY,
				nextSegment.EndX, nextSegment.EndY,
				segment.Radius,
			)
			apb.commands = append(apb.commands, curves...)
		}
		currentX = segment.EndX
		currentY = segment.EndY
	}

	return strings.Join(apb.commands, "")
}

// Example path creation functions
func createSelectivelyRoundedRectangle() string {
	return NewAdvancedPathBuilder().
		MoveTo(50, 50).
		LineTo(150, 50).Square().     // Top edge - sharp corner
		LineTo(150, 100).Rounded(10). // Right edge - rounded corner
		LineTo(50, 100).Rounded(10).  // Bottom edge - rounded corner
		LineTo(50, 50).Square().      // Left edge - sharp corner
		Build()
}

func createZigzagWithRoundedPeaks() string {
	return NewAdvancedPathBuilder().
		MoveTo(20, 80).
		LineTo(40, 40).Rounded(8).  // Peak 1 - rounded
		LineTo(60, 80).Square().    // Valley 1 - sharp
		LineTo(80, 40).Rounded(12). // Peak 2 - rounded
		LineTo(100, 80).Square().   // Valley 2 - sharp
		LineTo(120, 40).Rounded(6). // Peak 3 - rounded
		LineTo(140, 80).Square().   // End - sharp
		Build()
}

func createComplexShape() string {
	return NewAdvancedPathBuilder().
		MoveTo(200, 50).
		LineTo(250, 50).Square().     // Horizontal line
		LineTo(275, 75).Rounded(8).   // Diagonal down-right
		LineTo(275, 125).Square().    // Vertical line
		LineTo(250, 150).Rounded(12). // Diagonal down-left
		LineTo(200, 150).Square().    // Horizontal line
		LineTo(175, 125).Rounded(5).  // Diagonal up-left
		LineTo(175, 75).Square().     // Vertical line
		LineTo(200, 50).Rounded(10).  // Diagonal up-right, back to start
		Build()
}

func createStarWithRoundedPoints() string {
	centerX, centerY := 350, 100
	outerRadius := 40
	innerRadius := 20

	builder := NewAdvancedPathBuilder()

	// Calculate star points
	points := make([][2]int, 10) // 5 outer points + 5 inner points
	for i := 0; i < 10; i++ {
		angle := float64(i) * math.Pi / 5.0 // 36 degrees each
		var radius float64
		if i%2 == 0 {
			radius = float64(outerRadius) // Outer points
		} else {
			radius = float64(innerRadius) // Inner points
		}

		x := centerX + int(radius*math.Cos(angle-math.Pi/2))
		y := centerY + int(radius*math.Sin(angle-math.Pi/2))
		points[i] = [2]int{x, y}
	}

	// Build star path with rounded outer points, sharp inner points
	builder.MoveTo(points[0][0], points[0][1])
	for i := 1; i < len(points); i++ {
		if i%2 == 0 {
			// Outer point - rounded
			builder = builder.LineTo(points[i][0], points[i][1]).Rounded(5)
		} else {
			// Inner point - sharp
			builder = builder.LineTo(points[i][0], points[i][1]).Square()
		}
	}
	// Close the star
	builder = builder.LineTo(points[0][0], points[0][1]).Square()

	return builder.Build()
}

func main() {
	canvas := svg.New(os.Stdout)
	canvas.Start(450, 200)
	canvas.Title("Fluent Path Builder with Selective Rounded Corners")

	// Style for all paths
	pathStyle := "stroke:black;fill:none;stroke-width:2"

	// Example 1: Rectangle with selective rounding
	rect := createSelectivelyRoundedRectangle()
	canvas.Path(rect, pathStyle)

	// Example 2: Zigzag with rounded peaks
	zigzag := createZigzagWithRoundedPeaks()
	canvas.Path(zigzag, "stroke:blue;fill:none;stroke-width:2")

	// Example 3: Complex shape
	complex := createComplexShape()
	canvas.Path(complex, "stroke:red;fill:none;stroke-width:2")

	// Example 4: Star with rounded outer points
	star := createStarWithRoundedPoints()
	canvas.Path(star, "stroke:green;fill:none;stroke-width:2")

	// Add labels
	canvas.Text(100, 25, "Rectangle (mixed corners)", "font-family:Arial;font-size:12;fill:black")
	canvas.Text(80, 170, "Zigzag (rounded peaks)", "font-family:Arial;font-size:12;fill:blue")
	canvas.Text(225, 25, "Complex shape", "font-family:Arial;font-size:12;fill:red")
	canvas.Text(310, 170, "Star (rounded points)", "font-family:Arial;font-size:12;fill:green")

	canvas.End()
}

// Additional utility functions for debugging
func debugPath(path string) {
	fmt.Fprintf(os.Stderr, "Generated path: %s\n", path)

	// Count different command types
	commands := []string{"M", "L", "Q", "C"}
	for _, cmd := range commands {
		count := strings.Count(path, cmd)
		if count > 0 {
			fmt.Fprintf(os.Stderr, "Command %s appears %d times\n", cmd, count)
		}
	}
}

// Validation function
func validatePath(path string) bool {
	if !strings.HasPrefix(path, "M") {
		fmt.Fprintf(os.Stderr, "Warning: Path should start with M command\n")
		return false
	}
	return true
}
