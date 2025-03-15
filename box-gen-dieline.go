package main

import (
	"fmt"
	"os"

	svg "github.com/ajstarks/svgo"
)

// BoxDimensions stores the primary measurements of the box
type BoxDimensions struct {
	Width  float64 // Width of the box
	Height float64 // Height of the box
	Depth  float64 // Depth of the box
}

// Additional measurements and constants
const (
	marginMM    = 10.0  // Margin around the dieline
	glueFlapMM  = 15.0  // Width of glue flap
	tuckFlapMM  = 20.0  // Height of tuck flap
	cornerGapMM = 2.0   // Small gap at corners for folding
	strokeWidth = 1.0   // Line thickness
	dashPattern = "5,5" // Dash pattern for fold lines
)

func mainDieline() {
	// Example dimensions in millimeters
	box := BoxDimensions{
		Width:  100.0,
		Height: 150.0,
		Depth:  50.0,
	}

	// Create SVG file
	file, err := os.Create("box_dieline.svg")
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		return
	}
	defer file.Close()

	// Calculate total dimensions including margins
	totalWidth := 2*marginMM + 2*box.Depth + box.Width + glueFlapMM
	totalHeight := 2*marginMM + 2*box.Height + 2*box.Depth + 2*tuckFlapMM

	// Initialize SVG canvas
	canvas := svg.New(file)
	canvas.Start(int(totalWidth), int(totalHeight))

	// Set initial position for drawing
	startX := marginMM
	startY := marginMM

	// Draw main body (central rectangle)
	drawMainBody(canvas, startX, startY, box)

	// Draw side panels
	drawSidePanels(canvas, startX, startY, box)

	// Draw top and bottom panels
	drawTopBottomPanels(canvas, startX, startY, box)

	// Draw tuck flaps
	drawTuckFlaps(canvas, startX, startY, box)

	// Draw glue flap
	drawGlueFlap(canvas, startX, startY, box)

	canvas.End()
}

func drawMainBody(canvas *svg.SVG, x, y float64, box BoxDimensions) {
	// Main rectangle - cut lines
	canvas.Rect(
		int(x+box.Depth),
		int(y+box.Depth+tuckFlapMM),
		int(box.Width),
		int(box.Height),
		"fill:none;stroke:black;stroke-width:1")
}

func drawSidePanels(canvas *svg.SVG, x, y float64, box BoxDimensions) {
	// Left panel
	canvas.Rect(
		int(x),
		int(y+box.Depth+tuckFlapMM),
		int(box.Depth),
		int(box.Height),
		"fill:none;stroke:black;stroke-width:1")

	// Right panel
	canvas.Rect(
		int(x+box.Depth+box.Width),
		int(y+box.Depth+tuckFlapMM),
		int(box.Depth),
		int(box.Height),
		"fill:none;stroke:black;stroke-width:1")

	// Add fold lines
	canvas.Line(
		int(x+box.Depth),
		int(y+box.Depth+tuckFlapMM),
		int(x+box.Depth),
		int(y+box.Depth+tuckFlapMM+box.Height),
		"stroke:black;stroke-width:1;stroke-dasharray:5,5")

	canvas.Line(
		int(x+box.Depth+box.Width),
		int(y+box.Depth+tuckFlapMM),
		int(x+box.Depth+box.Width),
		int(y+box.Depth+tuckFlapMM+box.Height),
		"stroke:black;stroke-width:1;stroke-dasharray:5,5")
}

func drawTopBottomPanels(canvas *svg.SVG, x, y float64, box BoxDimensions) {
	// Top panel
	canvas.Rect(
		int(x+box.Depth),
		int(y+tuckFlapMM),
		int(box.Width),
		int(box.Depth),
		"fill:none;stroke:black;stroke-width:1")

	// Bottom panel
	canvas.Rect(
		int(x+box.Depth),
		int(y+box.Depth+tuckFlapMM+box.Height),
		int(box.Width),
		int(box.Depth),
		"fill:none;stroke:black;stroke-width:1")

	// Add fold lines
	canvas.Line(
		int(x+box.Depth),
		int(y+box.Depth+tuckFlapMM),
		int(x+box.Depth+box.Width),
		int(y+box.Depth+tuckFlapMM),
		"stroke:black;stroke-width:1;stroke-dasharray:5,5")

	canvas.Line(
		int(x+box.Depth),
		int(y+box.Depth+tuckFlapMM+box.Height),
		int(x+box.Depth+box.Width),
		int(y+box.Depth+tuckFlapMM+box.Height),
		"stroke:black;stroke-width:1;stroke-dasharray:5,5")
}

func drawTuckFlaps(canvas *svg.SVG, x, y float64, box BoxDimensions) {
	// Top tuck flap (trapezoid shape)
	points := fmt.Sprintf("%d,%d %d,%d %d,%d %d,%d",
		int(x+box.Depth+cornerGapMM), int(y),
		int(x+box.Depth+box.Width-cornerGapMM), int(y),
		int(x+box.Depth+box.Width-cornerGapMM-5), int(y+tuckFlapMM),
		int(x+box.Depth+cornerGapMM+5), int(y+tuckFlapMM))
	canvas.Polygon(points, "fill:none;stroke:black;stroke-width:1")

	// Bottom tuck flap
	bottomY := y + box.Depth + tuckFlapMM + box.Height + box.Depth
	points = fmt.Sprintf("%d,%d %d,%d %d,%d %d,%d",
		int(x+box.Depth+cornerGapMM), int(bottomY+tuckFlapMM),
		int(x+box.Depth+box.Width-cornerGapMM), int(bottomY+tuckFlapMM),
		int(x+box.Depth+box.Width-cornerGapMM-5), int(bottomY),
		int(x+box.Depth+cornerGapMM+5), int(bottomY))
	canvas.Polygon(points, "fill:none;stroke:black;stroke-width:1")
}

func drawGlueFlap(canvas *svg.SVG, x, y float64, box BoxDimensions) {
	// Glue flap (slightly curved for better adhesion)
	startX := x + box.Depth + box.Width + box.Depth
	startY := y + box.Depth + tuckFlapMM

	path := fmt.Sprintf("M %d,%d C %d,%d %d,%d %d,%d L %d,%d C %d,%d %d,%d %d,%d Z",
		int(startX), int(startY),
		int(startX+5), int(startY+10),
		int(startX+glueFlapMM), int(startY+20),
		int(startX+glueFlapMM), int(startY+box.Height-20),
		int(startX), int(startY+box.Height),
		int(startX+5), int(startY+box.Height-10),
		int(startX+glueFlapMM), int(startY+box.Height-20),
		int(startX), int(startY))

	canvas.Path(path, "fill:none;stroke:black;stroke-width:1;stroke-dasharray:5,5")
}
