package main

import (
	"fmt"
	"os"

	svg "github.com/ajstarks/svgo"
)

func main() {
	// Box dimensions in mm
	length := 200
	width := 150
	height := 100
	gap := 5
	cornerRadius := 10
	foldInset := 1 // 1mm inset for folding

	// Create filename based on dimensions
	filename := fmt.Sprintf("box_%d_%d_%d.svg", length, width, height)
	f, err := os.Create(filename)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer f.Close()

	// SVG canvas dimensions (add some margin)
	canvasWidth := 3*width + 2*length + 4*gap + 100
	canvasHeight := 2*height + width + 4*gap + 100

	// Calculate panel positions
	leftPanelStart := gap
	bottomPanelStart := leftPanelStart + width + gap
	rightPanelStart := bottomPanelStart + length + gap
	topPanelStart := leftPanelStart + width + gap
	lidPanelStart := rightPanelStart + width + gap

	topEdge := gap
	middleEdge := topEdge + height + gap
	bottomEdge := middleEdge + length + gap

	// Start SVG
	canvas := svg.New(f)
	canvas.Start(canvasWidth, canvasHeight)

	// Draw fold lines and cut lines
	drawFoldLines(canvas, leftPanelStart, bottomPanelStart, rightPanelStart, topPanelStart,
		lidPanelStart, topEdge, middleEdge, bottomEdge, width, length, height)

	drawCutLines(canvas, leftPanelStart, bottomPanelStart, rightPanelStart, topPanelStart,
		topEdge, middleEdge, bottomEdge, width, length, height, cornerRadius, foldInset)

	// Add dimension labels
	labelStyle := "text-anchor:middle;font-size:14px;font-family:Arial"
	canvas.Text(bottomPanelStart+length/2, bottomEdge+height+30, fmt.Sprintf("Length: %dmm", length), labelStyle)
	canvas.Text(leftPanelStart+width/2, bottomEdge+height+30, fmt.Sprintf("Width: %dmm", width), labelStyle)
	canvas.Text(leftPanelStart-20, middleEdge+height/2, fmt.Sprintf("Height: %dmm", height), labelStyle)

	canvas.End()

	fmt.Printf("SVG box template created: %s\n", filename)
}

// drawFoldLines draws all the fold lines (dotted red lines)
func drawFoldLines(canvas *svg.SVG, leftPanelStart, bottomPanelStart, rightPanelStart, topPanelStart,
	lidPanelStart, topEdge, middleEdge, bottomEdge, width, length, height int) {

	redDotted := "stroke:red;stroke-width:2;stroke-dasharray:5,5;fill:none"

	// Step 1: Bottom horizontal fold line
	canvas.Line(leftPanelStart, bottomEdge, leftPanelStart+width+length+width, bottomEdge, redDotted)

	// Step 2: Vertical fold lines between panels
	canvas.Line(leftPanelStart+width, middleEdge, leftPanelStart+width, bottomEdge, redDotted)
	canvas.Line(bottomPanelStart+length, middleEdge, bottomPanelStart+length, bottomEdge, redDotted)
	canvas.Line(rightPanelStart+width, middleEdge, rightPanelStart+width, bottomEdge, redDotted)

	// Step 3: Middle horizontal fold line
	canvas.Line(leftPanelStart, middleEdge, lidPanelStart+width, middleEdge, redDotted)

	// Step 4: Top horizontal fold line for the lid
	canvas.Line(topPanelStart, topEdge+height, topPanelStart+length, topEdge+height, redDotted)
}

// drawCutLines draws all the cut lines (solid blue lines)
func drawCutLines(canvas *svg.SVG, leftPanelStart, bottomPanelStart, rightPanelStart, topPanelStart,
	topEdge, middleEdge, bottomEdge, width, length, height, cornerRadius, foldInset int) {

	blueSolid := "stroke:blue;stroke-width:2;fill:none"

	// Apply inset for folding to top and bottom panels
	topInsetStart := topPanelStart + foldInset
	topInsetLength := length - 2*foldInset

	// Step 1: Left side of the template
	// Top left panel top edge
	canvas.Line(leftPanelStart, middleEdge, leftPanelStart, middleEdge-cornerRadius, blueSolid)
	canvas.Qbez(leftPanelStart, middleEdge-cornerRadius, leftPanelStart+cornerRadius, middleEdge-cornerRadius, leftPanelStart+cornerRadius, middleEdge-2*cornerRadius, blueSolid)
	canvas.Line(leftPanelStart+cornerRadius, middleEdge-2*cornerRadius, leftPanelStart+width-cornerRadius, middleEdge-2*cornerRadius, blueSolid)
	canvas.Qbez(leftPanelStart+width-cornerRadius, middleEdge-2*cornerRadius, leftPanelStart+width, middleEdge-2*cornerRadius, leftPanelStart+width, middleEdge-cornerRadius, blueSolid)

	// Step 2: Middle top edge
	canvas.Line(leftPanelStart+width, middleEdge-cornerRadius, bottomPanelStart+length, middleEdge-cornerRadius, blueSolid)

	// Step 3: Right top panel top edge
	canvas.Qbez(bottomPanelStart+length, middleEdge-cornerRadius, bottomPanelStart+length+cornerRadius, middleEdge-cornerRadius, bottomPanelStart+length+cornerRadius, middleEdge-2*cornerRadius, blueSolid)
	canvas.Line(bottomPanelStart+length+cornerRadius, middleEdge-2*cornerRadius, rightPanelStart+width-cornerRadius, middleEdge-2*cornerRadius, blueSolid)
	canvas.Qbez(rightPanelStart+width-cornerRadius, middleEdge-2*cornerRadius, rightPanelStart+width, middleEdge-2*cornerRadius, rightPanelStart+width, middleEdge-cornerRadius, blueSolid)

	// Step 4: Top panel (lid) with inset
	canvas.Line(rightPanelStart+width, middleEdge-cornerRadius, rightPanelStart+width, middleEdge, blueSolid)

	// Step 5: Top of lid panel with inset
	canvas.Line(topInsetStart, topEdge, topInsetStart+topInsetLength, topEdge, blueSolid)

	// Step 6: Right edge of the lid panel with inset
	canvas.Line(topInsetStart+topInsetLength, topEdge, topInsetStart+topInsetLength, topEdge+height, blueSolid)

	// Step 7: Left edge of the lid panel with inset
	canvas.Line(topInsetStart, topEdge, topInsetStart, topEdge+height, blueSolid)

	// Step 8: Rest of right edge
	canvas.Line(rightPanelStart+width, middleEdge, rightPanelStart+width, bottomEdge, blueSolid)

	// Step 9: Bottom edge of right panel
	// Bottom right corner and tab
	canvas.Line(rightPanelStart+width, bottomEdge, rightPanelStart+width+cornerRadius, bottomEdge+cornerRadius, blueSolid)
	canvas.Line(rightPanelStart+width+cornerRadius, bottomEdge+cornerRadius, rightPanelStart+width+cornerRadius, bottomEdge+height-cornerRadius, blueSolid)
	canvas.Qbez(rightPanelStart+width+cornerRadius, bottomEdge+height-cornerRadius, rightPanelStart+width+cornerRadius, bottomEdge+height, rightPanelStart+width, bottomEdge+height, blueSolid)

	// Step 10: Bottom tabs
	// Right bottom tab
	canvas.Line(rightPanelStart+width, bottomEdge+height, rightPanelStart, bottomEdge+height, blueSolid)
	canvas.Qbez(rightPanelStart, bottomEdge+height, rightPanelStart-cornerRadius, bottomEdge+height, rightPanelStart-cornerRadius, bottomEdge+height-cornerRadius, blueSolid)
	canvas.Line(rightPanelStart-cornerRadius, bottomEdge+height-cornerRadius, rightPanelStart-cornerRadius, bottomEdge+cornerRadius, blueSolid)
	canvas.Qbez(rightPanelStart-cornerRadius, bottomEdge+cornerRadius, rightPanelStart-cornerRadius, bottomEdge, rightPanelStart, bottomEdge, blueSolid)

	// Middle bottom tabs with inset
	canvas.Line(bottomPanelStart+length, bottomEdge, bottomPanelStart+length+cornerRadius, bottomEdge+cornerRadius, blueSolid)
	canvas.Line(bottomPanelStart+length+cornerRadius, bottomEdge+cornerRadius, bottomPanelStart+length+cornerRadius, bottomEdge+height-cornerRadius, blueSolid)
	canvas.Qbez(bottomPanelStart+length+cornerRadius, bottomEdge+height-cornerRadius, bottomPanelStart+length+cornerRadius, bottomEdge+height, bottomPanelStart+length, bottomEdge+height, blueSolid)
	canvas.Line(bottomPanelStart+length, bottomEdge+height, bottomPanelStart+foldInset, bottomEdge+height, blueSolid)
	canvas.Qbez(bottomPanelStart+foldInset, bottomEdge+height, bottomPanelStart+foldInset-cornerRadius, bottomEdge+height, bottomPanelStart+foldInset-cornerRadius, bottomEdge+height-cornerRadius, blueSolid)
	canvas.Line(bottomPanelStart+foldInset-cornerRadius, bottomEdge+height-cornerRadius, bottomPanelStart+foldInset-cornerRadius, bottomEdge+cornerRadius, blueSolid)
	canvas.Qbez(bottomPanelStart+foldInset-cornerRadius, bottomEdge+cornerRadius, bottomPanelStart+foldInset-cornerRadius, bottomEdge, bottomPanelStart+foldInset, bottomEdge, blueSolid)

	// Left bottom tabs
	canvas.Line(leftPanelStart+width, bottomEdge, leftPanelStart+width+cornerRadius, bottomEdge+cornerRadius, blueSolid)
	canvas.Line(leftPanelStart+width+cornerRadius, bottomEdge+cornerRadius, leftPanelStart+width+cornerRadius, bottomEdge+height-cornerRadius, blueSolid)
	canvas.Qbez(leftPanelStart+width+cornerRadius, bottomEdge+height-cornerRadius, leftPanelStart+width+cornerRadius, bottomEdge+height, leftPanelStart+width, bottomEdge+height, blueSolid)
	canvas.Line(leftPanelStart+width, bottomEdge+height, leftPanelStart, bottomEdge+height, blueSolid)
	canvas.Qbez(leftPanelStart, bottomEdge+height, leftPanelStart-cornerRadius, bottomEdge+height, leftPanelStart-cornerRadius, bottomEdge+height-cornerRadius, blueSolid)
	canvas.Line(leftPanelStart-cornerRadius, bottomEdge+height-cornerRadius, leftPanelStart-cornerRadius, bottomEdge+cornerRadius, blueSolid)
	canvas.Qbez(leftPanelStart-cornerRadius, bottomEdge+cornerRadius, leftPanelStart-cornerRadius, bottomEdge, leftPanelStart, bottomEdge, blueSolid)

	// Step 11: Complete left edge to connect back to starting point
	canvas.Line(leftPanelStart, bottomEdge, leftPanelStart, middleEdge, blueSolid)
}
