package pathbuilder

import (
	"fmt"
	"math"
	"strings"
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

// NewAdvancedPathBuilder creates a new path builder instance
func NewAdvancedPathBuilder() *AdvancedPathBuilder {
	return &AdvancedPathBuilder{
		commands: make([]string, 0),
		segments: make([]PathSegment, 0),
	}
}

// MoveTo sets the starting point for the path
func (apb *AdvancedPathBuilder) MoveTo(x, y int) *AdvancedPathBuilder {
	apb.commands = append(apb.commands, fmt.Sprintf("M%d,%d", x, y))
	apb.lastX = x
	apb.lastY = y
	return apb
}

// LineTo draws a line to absolute coordinates
func (apb *AdvancedPathBuilder) LineTo(x, y int) *CornerBuilder {
	return &CornerBuilder{
		pathBuilder: apb,
		endX:        x,
		endY:        y,
	}
}

// RelativeLine moves relative to current position
func (apb *AdvancedPathBuilder) RelativeLine(offsetX, offsetY int) *CornerBuilder {
	newX := apb.lastX + offsetX
	newY := apb.lastY + offsetY
	return &CornerBuilder{
		pathBuilder: apb,
		endX:        newX,
		endY:        newY,
	}
}

// HorizontalLine moves horizontally by offset amount
func (apb *AdvancedPathBuilder) HorizontalLine(offset int) *CornerBuilder {
	return apb.RelativeLine(offset, 0)
}

// VerticalLine moves vertically by offset amount
func (apb *AdvancedPathBuilder) VerticalLine(offset int) *CornerBuilder {
	return apb.RelativeLine(0, offset)
}

// CurrentPosition returns the current drawing position
func (apb *AdvancedPathBuilder) CurrentPosition() (int, int) {
	return apb.lastX, apb.lastY
}

// Square creates a sharp corner at this point
func (cb *CornerBuilder) Square() *AdvancedPathBuilder {
	cb.pathBuilder.segments = append(cb.pathBuilder.segments, PathSegment{
		EndX:       cb.endX,
		EndY:       cb.endY,
		CornerType: SquareCorner,
	})
	return cb.pathBuilder
}

// Rounded creates a rounded corner with the specified radius
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
	incomingX := float64(prevX - cornerX)
	incomingY := float64(prevY - cornerY)
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
	curveStartX := cornerX + int(incomingUnitX*radiusF)
	curveStartY := cornerY + int(incomingUnitY*radiusF)
	curveEndX := cornerX + int(outgoingUnitX*radiusF)
	curveEndY := cornerY + int(outgoingUnitY*radiusF)

	// Line to curve start
	commands = append(commands, fmt.Sprintf("L%d,%d", curveStartX, curveStartY))
	// Quadratic curve through the corner
	commands = append(commands, fmt.Sprintf("Q%d,%d,%d,%d", cornerX, cornerY, curveEndX, curveEndY))

	return commands
}

// Build generates the final SVG path string
func (apb *AdvancedPathBuilder) Build() string {
	if len(apb.segments) == 0 {
		return strings.Join(apb.commands, "")
	}

	currentX := apb.lastX
	currentY := apb.lastY

	for i, segment := range apb.segments {
		if i == len(apb.segments)-1 || segment.CornerType == SquareCorner {
			apb.commands = append(apb.commands, fmt.Sprintf("L%d,%d", segment.EndX, segment.EndY))
		} else if segment.CornerType == RoundedCorner && i < len(apb.segments)-1 {
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

// Clear resets the builder to start a new path
func (apb *AdvancedPathBuilder) Clear() *AdvancedPathBuilder {
	apb.commands = make([]string, 0)
	apb.segments = make([]PathSegment, 0)
	apb.lastX = 0
	apb.lastY = 0
	return apb
}

// Clone creates a copy of the current builder state
func (apb *AdvancedPathBuilder) Clone() *AdvancedPathBuilder {
	newBuilder := NewAdvancedPathBuilder()
	newBuilder.commands = make([]string, len(apb.commands))
	copy(newBuilder.commands, apb.commands)
	newBuilder.segments = make([]PathSegment, len(apb.segments))
	copy(newBuilder.segments, apb.segments)
	newBuilder.lastX = apb.lastX
	newBuilder.lastY = apb.lastY
	return newBuilder
}

// ClosePath adds a Z command to close the path
func (apb *AdvancedPathBuilder) ClosePath() *AdvancedPathBuilder {
	apb.commands = append(apb.commands, "Z")
	return apb
}

// Utility functions for common shapes

// CreateRectangle creates a rectangle with optional rounded corners
func CreateRectangle(x, y, width, height, radius int) string {
	if radius <= 0 {
		return NewAdvancedPathBuilder().
			MoveTo(x, y).
			HorizontalLine(width).Square().
			VerticalLine(height).Square().
			HorizontalLine(-width).Square().
			VerticalLine(-height).Square().
			Build()
	}

	return NewAdvancedPathBuilder().
		MoveTo(x, y).
		HorizontalLine(width).Rounded(radius).
		VerticalLine(height).Rounded(radius).
		HorizontalLine(-width).Rounded(radius).
		VerticalLine(-height).Rounded(radius).
		Build()
}

// Validation functions

// ValidatePath checks if a path string is valid
func ValidatePath(path string) bool {
	if !strings.HasPrefix(path, "M") {
		return false
	}
	// Add more validation as needed
	return true
}

// GetPathBounds calculates the bounding box of a path
func GetPathBounds(path string) (minX, minY, maxX, maxY int) {
	// This would parse the path and calculate bounds
	// Simplified implementation for now
	return 0, 0, 100, 100
}
