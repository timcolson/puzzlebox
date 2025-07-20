package main

import (
	"flag"
	"fmt"
	"github.com/charmbracelet/log"
	"os"
	"path/filepath"
	"time"

	svg "github.com/ajstarks/svgo"
)

var redDotted = "stroke:red;stroke-width:1;stroke-dasharray:5,1;fill:none"
var blueSolid = "stroke:blue;stroke-width:1;fill:none"

// Box represents dimensions as floats to simplify calculations.
// Internal calculations return floats.
// Public calculations output integers for easy of use by external drawing commands which expect ints.
type Box struct {
	Width            float64
	Depth            float64
	Height           float64
	FoldGap          float64
	BottomTabPercent float64
}

func (b Box) TotalWidth() int {
	totalWidth := int((2.0 * b.Width) + (2.0 * b.Depth) + b.SideFlapWidth())
	return totalWidth
}

func (b Box) SideFlapWidth() float64 {
	return .25 * b.Depth
}

func (b Box) TotalHeight() int {
	return int(b.topFlapHeight() + b.Depth + b.Height + b.bottomFlapMaxHeight())
}

func (b Box) bottomFlapMaxHeight() float64 {
	return 0.75 * b.Depth
}

func (b Box) topFlapHeight() float64 {
	return .2 * b.Depth
}

func (b Box) BackRight() int {
	return int((2 * b.Width) + (2 * b.Depth))
}

func (b Box) Bottom() int {
	return int(b.topFlapHeight() + b.Depth + b.Height)
}
func (b Box) Top() int {
	return int(b.topFlapHeight() + b.Depth)
}

func (b Box) SideALeft() int {
	return 0
}

func (b Box) SideFlapHeight() float64 {
	return b.Depth * 0.66
}

func main() {

	// Create a log file
	// Commented out because file is unused at this time
	//logFile, err := os.OpenFile("box-gen.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	//if err != nil {
	//	log.Fatalf("Failed to open log file: %v", err)
	//}
	//defer logFile.Close()

	// Writes to both console and file
	// Commented out because using this multiwriter seems to only output plain text to STDERR, so no format/color!
	//multiWriter := io.MultiWriter(os.Stderr, logFile)

	logger := log.NewWithOptions(os.Stderr, log.Options{
		ReportCaller:    true,
		ReportTimestamp: true,
		TimeFormat:      time.Kitchen,
	})

	// Define command line flags
	outputFile := flag.String("o", "box.svg", "Output file name (default: box_<width>_<depth>_<height>.svg)")
	outputDir := flag.String("d", "out", "Output directory (default: current directory)")

	// Box dimensions in mm
	width := flag.Int("width", 36, "Box width (front & back with top flap) in mm")
	depth := flag.Int("depth", 29, "Box depth (sides) in mm")
	height := flag.Int("height", 48, "Box height in mm")
	foldGap := flag.Int("gap", 5, "Gap for folds in mm")

	flag.Parse()

	// Dereference for convenience
	W := float64(*width)
	D := float64(*depth)
	H := float64(*height)
	G := float64(*foldGap)

	myBox := Box{
		Width:   W,
		Depth:   D,
		Height:  H,
		FoldGap: G,
	}

	// Determine filename
	var filename string
	if *outputFile != "" {
		filename = *outputFile
		// Add .svg extension if not present
		if filepath.Ext(filename) != ".svg" {
			filename += ".svg"
		}
	} else {
		// Create filename based on dimensions
		filename = fmt.Sprintf("box_%d_%d_%d.svg", *width, *depth, *height)
	}

	// Prepend output directory if specified
	if *outputDir != "" {
		// Create a directory if it doesn't exist

		if err := os.MkdirAll(*outputDir, 0755); err != nil {
			fmt.Printf("Error creating directory %s: %v\n", *outputDir, err)
			return
		}
		filename = filepath.Join(*outputDir, filename)
	}

	// Create an output file
	logger.Info("Created output filename: ", "filename", filename)
	f, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Error creating file %s: %v\n", filename, err)
		return
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
		}
	}(f)

	// SVG canvas dimensions (add some margin)
	padding := 10
	canvasWidth := myBox.TotalWidth() + padding
	canvasHeight := myBox.TotalHeight() + padding

	// Start SVG
	canvas := svg.New(f)
	canvas.Start(canvasWidth, canvasHeight)
	canvas.Rect(0, 0, canvasWidth, canvasHeight, "fill:white")

	// Add metadata
	canvas.Title(fmt.Sprintf("Box Template %dx%dx%d mm", *width, *depth, *height))
	canvas.Desc(fmt.Sprintf("Generated on %s - Box dimensions: %dx%dx%d mm",
		time.Now().Format("2006-01-02 15:04:05"), *width, *depth, *height))

	// Draw fold lines first
	drawFoldLines(canvas, myBox)
	// Draw cut lines second
	drawCutLines(canvas, myBox)
	canvas.End()

	fmt.Printf("SVG box created: %s\n", filename)
}

// drawFoldLines draws all the fold lines (dotted red lines)
func drawFoldLines(canvas *svg.SVG, b Box) {

	// Step 1: Bottom horizontal fold line
	originX := b.SideALeft()
	originY := b.Bottom()
	canvas.Line(originX, originY, b.BackRight(), b.Bottom(), redDotted)

	// Right side
	canvas.Line(b.BackRight(), b.Bottom(), b.BackRight(), b.Top(), redDotted)
	// Box Top
	canvas.Line(b.BackRight(), b.Top(), originX, b.Top(), redDotted)
}

// drawCutLines draws all the cut lines (solid blue lines)
func drawCutLines(canvas *svg.SVG, b Box) {

	// Step 1: Panel1 Tab
	originX := b.SideALeft()
	originY := b.Bottom()
	// Draw angle line to 50% of the side distance and half of the side flap height
	canvas.Line(originX, originY, originX+int(b.Depth*.5), originY+int(b.SideFlapHeight()*.5), blueSolid)
	// Draw straight vertical line down to SideFlapHeight
	canvas.Line(originX+int(b.Depth*.5), originY+int(b.SideFlapHeight()*.5), originX+int(b.Depth*.5), originY+int(b.SideFlapHeight()), blueSolid)
	// Draw straight horizontal line to right side of the panel
	canvas.Line(originX+int(b.Depth*.5), originY+int(b.SideFlapHeight()), originX+int(b.Depth), originY+int(b.SideFlapHeight()), blueSolid)
	// Draw straight vertical line up to the bottom of the box.
	canvas.Line(originX+int(b.Depth), originY+int(b.SideFlapHeight()), originX+int(b.Depth), b.Bottom(), blueSolid)

}
