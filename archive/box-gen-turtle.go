// package main

// import (
// 	"fmt"
// 	"os"

// 	svg "github.com/ajstarks/svgo"
// )

// func main() {
// 	// Box dimensions in mm
// 	Length := 200
// 	Width := 150
// 	Height := 100
// 	Gap := 5
// 	CornerRadius := 10

// 	// Create filename based on dimensions
// 	filename := fmt.Sprintf("box_%d_%d_%d.svg", Length, Width, Height)
// 	f, err := os.Create(filename)
// 	if err != nil {
// 		fmt.Println("Error creating file:", err)
// 		return
// 	}
// 	defer f.Close()

// 	// SVG canvas dimensions (add some margin)
// 	canvasWidth := 3*Width + 2*Length + 4*Gap + 100
// 	canvasHeight := 2*Height + Width + 4*Gap + 100

// 	// Calculate panel positions
// 	leftPanelX := Gap
// 	bottomPanelX := leftPanelX + Width + Gap
// 	rightPanelX := bottomPanelX + Length + Gap
// 	topPanelX := leftPanelX + Width + Gap
// 	lidPanelX := rightPanelX + Width + Gap

// 	topY := Gap
// 	middleY := topY + Height + Gap
// 	bottomY := middleY + Length + Gap

// 	// Start SVG
// 	canvas := svg.New(f)
// 	canvas.Start(canvasWidth, canvasHeight)

// 	// Style definitions
// 	blueSolid := "stroke:blue;stroke-width:2;fill:none"
// 	redDotted := "stroke:red;stroke-width:2;stroke-dasharray:5,5;fill:none"

// 	// Red lines (fold lines) - Starting from bottom left
// 	// Step 1: Bottom horizontal fold lines
// 	canvas.Line(leftPanelX, bottomY, leftPanelX+Width+Length+Width, bottomY, redDotted)

// 	// Step 2: Vertical fold lines between panels
// 	canvas.Line(leftPanelX+Width, middleY, leftPanelX+Width, bottomY, redDotted)
// 	canvas.Line(bottomPanelX+Length, middleY, bottomPanelX+Length, bottomY, redDotted)
// 	canvas.Line(rightPanelX+Width, middleY, rightPanelX+Width, bottomY, redDotted)

// 	// Step 3: Middle horizontal fold line
// 	canvas.Line(leftPanelX, middleY, lidPanelX+Width, middleY, redDotted)

// 	// Step 4: Top horizontal fold line for the lid
// 	canvas.Line(topPanelX, topY+Height, topPanelX+Length, topY+Height, redDotted)

// 	// Blue lines (cut lines) - Starting from top left, moving clockwise
// 	// Step 1: Left side of the template
// 	// Top left panel top edge
// 	canvas.Line(leftPanelX, middleY, leftPanelX, middleY-CornerRadius, blueSolid)
// 	canvas.Qbez(leftPanelX, middleY-CornerRadius, leftPanelX+CornerRadius, middleY-CornerRadius, leftPanelX+CornerRadius, middleY-2*CornerRadius, blueSolid)
// 	canvas.Line(leftPanelX+CornerRadius, middleY-2*CornerRadius, leftPanelX+Width-CornerRadius, middleY-2*CornerRadius, blueSolid)
// 	canvas.Qbez(leftPanelX+Width-CornerRadius, middleY-2*CornerRadius, leftPanelX+Width, middleY-2*CornerRadius, leftPanelX+Width, middleY-CornerRadius, blueSolid)

// 	// Step 2: Middle top edge
// 	canvas.Line(leftPanelX+Width, middleY-CornerRadius, bottomPanelX+Length, middleY-CornerRadius, blueSolid)

// 	// Step 3: Right top panel top edge
// 	canvas.Qbez(bottomPanelX+Length, middleY-CornerRadius, bottomPanelX+Length+CornerRadius, middleY-CornerRadius, bottomPanelX+Length+CornerRadius, middleY-2*CornerRadius, blueSolid)
// 	canvas.Line(bottomPanelX+Length+CornerRadius, middleY-2*CornerRadius, rightPanelX+Width-CornerRadius, middleY-2*CornerRadius, blueSolid)
// 	canvas.Qbez(rightPanelX+Width-CornerRadius, middleY-2*CornerRadius, rightPanelX+Width, middleY-2*CornerRadius, rightPanelX+Width, middleY-CornerRadius, blueSolid)

// 	// Step 4: Top panel (lid)
// 	canvas.Line(rightPanelX+Width, middleY-CornerRadius, rightPanelX+Width, middleY, blueSolid)

// 	// Step 5: Top of lid panel
// 	canvas.Line(topPanelX, topY, topPanelX+Length, topY, blueSolid)

// 	// Step 6: Right edge of the lid panel
// 	canvas.Line(topPanelX+Length, topY, topPanelX+Length, topY+Height, blueSolid)

// 	// Step 7: Left edge of the lid panel
// 	canvas.Line(topPanelX, topY, topPanelX, topY+Height, blueSolid)

// 	// Step 8: Rest of right edge
// 	canvas.Line(rightPanelX+Width, middleY, rightPanelX+Width, bottomY, blueSolid)

// 	// Step 9: Bottom edge of right panel
// 	canvas.Line(rightPanelX+Width, bottomY, rightPanelX+Width+CornerRadius, bottomY+CornerRadius, blueSolid)
// 	canvas.Line(rightPanelX+Width+CornerRadius, bottomY+CornerRadius, rightPanelX+Width+CornerRadius, bottomY+Height-CornerRadius, blueSolid)
// 	canvas.Qbez(rightPanelX+Width+CornerRadius, bottomY+Height-CornerRadius, rightPanelX+Width+CornerRadius, bottomY+Height, rightPanelX+Width, bottomY+Height, blueSolid)

// 	// Step 10: Bottom tabs
// 	canvas.Line(rightPanelX+Width, bottomY+Height, rightPanelX, bottomY+Height, blueSolid)
// 	canvas.Qbez(rightPanelX, bottomY+Height, rightPanelX-CornerRadius, bottomY+Height, rightPanelX-CornerRadius, bottomY+Height-CornerRadius, blueSolid)
// 	canvas.Line(rightPanelX-CornerRadius, bottomY+Height-CornerRadius, rightPanelX-CornerRadius, bottomY+CornerRadius, blueSolid)
// 	canvas.Qbez(rightPanelX-CornerRadius, bottomY+CornerRadius, rightPanelX-CornerRadius, bottomY, rightPanelX, bottomY, blueSolid)

// 	// Middle bottom tabs
// 	canvas.Line(bottomPanelX+Length, bottomY, bottomPanelX+Length+CornerRadius, bottomY+CornerRadius, blueSolid)
// 	canvas.Line(bottomPanelX+Length+CornerRadius, bottomY+CornerRadius, bottomPanelX+Length+CornerRadius, bottomY+Height-CornerRadius, blueSolid)
// 	canvas.Qbez(bottomPanelX+Length+CornerRadius, bottomY+Height-CornerRadius, bottomPanelX+Length+CornerRadius, bottomY+Height, bottomPanelX+Length, bottomY+Height, blueSolid)
// 	canvas.Line(bottomPanelX+Length, bottomY+Height, bottomPanelX, bottomY+Height, blueSolid)
// 	canvas.Qbez(bottomPanelX, bottomY+Height, bottomPanelX-CornerRadius, bottomY+Height, bottomPanelX-CornerRadius, bottomY+Height-CornerRadius, blueSolid)
// 	canvas.Line(bottomPanelX-CornerRadius, bottomY+Height-CornerRadius, bottomPanelX-CornerRadius, bottomY+CornerRadius, blueSolid)
// 	canvas.Qbez(bottomPanelX-CornerRadius, bottomY+CornerRadius, bottomPanelX-CornerRadius, bottomY, bottomPanelX, bottomY, blueSolid)

// 	// Left bottom tabs
// 	canvas.Line(leftPanelX+Width, bottomY, leftPanelX+Width+CornerRadius, bottomY+CornerRadius, blueSolid)
// 	canvas.Line(leftPanelX+Width+CornerRadius, bottomY+CornerRadius, leftPanelX+Width+CornerRadius, bottomY+Height-CornerRadius, blueSolid)
// 	canvas.Qbez(leftPanelX+Width+CornerRadius, bottomY+Height-CornerRadius, leftPanelX+Width+CornerRadius, bottomY+Height, leftPanelX+Width, bottomY+Height, blueSolid)
// 	canvas.Line(leftPanelX+Width, bottomY+Height, leftPanelX, bottomY+Height, blueSolid)
// 	canvas.Qbez(leftPanelX, bottomY+Height, leftPanelX-CornerRadius, bottomY+Height, leftPanelX-CornerRadius, bottomY+Height-CornerRadius, blueSolid)
// 	canvas.Line(leftPanelX-CornerRadius, bottomY+Height-CornerRadius, leftPanelX-CornerRadius, bottomY+CornerRadius, blueSolid)
// 	canvas.Qbez(leftPanelX-CornerRadius, bottomY+CornerRadius, leftPanelX-CornerRadius, bottomY, leftPanelX, bottomY, blueSolid)

// 	// Step 11: Complete left edge to connect back to starting point
// 	canvas.Line(leftPanelX, bottomY, leftPanelX, middleY, blueSolid)

// 	// Add dimension labels
// 	labelStyle := "text-anchor:middle;font-size:14px;font-family:Arial"
// 	canvas.Text(bottomPanelX+Length/2, bottomY+Height+30, fmt.Sprintf("Length: %dmm", Length), labelStyle)
// 	canvas.Text(leftPanelX+Width/2, bottomY+Height+30, fmt.Sprintf("Width: %dmm", Width), labelStyle)
// 	canvas.Text(leftPanelX-20, middleY+Height/2, fmt.Sprintf("Height: %dmm", Height), labelStyle)

// 	canvas.End()

// 	fmt.Printf("SVG box template created: %s\n", filename)
// }
