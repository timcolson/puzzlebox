package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"42clients.com/puzzlebox/pkg/box"

	"github.com/ajstarks/svgo"
	"github.com/charmbracelet/log"
)

var redDotted = "stroke:red;stroke-width:1;stroke-dasharray:5,1;fill:none"
var blueSolid = "stroke:blue;stroke-width:1;fill:none"
var greenSolid = "stroke:green;stroke-width:1;fill:none"

func main() {
	logger := log.NewWithOptions(os.Stderr, log.Options{
		ReportCaller:    true,
		ReportTimestamp: true,
		TimeFormat:      time.Kitchen,
	})

	// Define command line flags
	outputFile := flag.String("o", "", "Output file name (default: box_<width>_<depth>_<height>.svg)")
	outputDir := flag.String("d", "out", "Output directory (default: out)")

	// Box dimensions in mm
	width := flag.Float64("width", 276, "Box width (front & back with top flap) in mm")
	depth := flag.Float64("depth", 206, "Box depth (sides) in mm")
	height := flag.Float64("height", 196, "Box height in mm")
	foldGap := flag.Float64("gap", 2, "Gap for folds in mm")

	flag.Parse()

	// Create box
	myBox := box.NewBox(*width, *depth, *height, *foldGap)

	if !myBox.IsValid() {
		logger.Error("Invalid box dimensions", "width", *width, "depth", *depth, "height", *height)
		os.Exit(1)
	}

	// Determine filename
	var filename string
	if *outputFile != "" {
		filename = *outputFile
		if filepath.Ext(filename) != ".svg" {
			filename += ".svg"
		}
	} else {
		filename = fmt.Sprintf("box_%.0f_%.0f_%.0f.svg", *width, *depth, *height)
	}

	// Create the output directory
	if *outputDir != "" {
		if err := os.MkdirAll(*outputDir, 0755); err != nil {
			logger.Error("Error creating directory", "dir", *outputDir, "error", err)
			return
		}
		filename = filepath.Join(*outputDir, filename)
	}

	logger.Info("Creating box template",
		"filename", filename,
		"dimensions", fmt.Sprintf("%.0fx%.0fx%.0f mm", *width, *depth, *height))

	// Create the output file
	f, err := os.Create(filename)
	if err != nil {
		logger.Error("Error creating file", "filename", filename, "error", err)
		return
	}
	defer f.Close()

	// Generate all paths
	paths := myBox.GenerateCompleteBox()

	// SVG canvas dimensions (add padding)
	padding := 20
	canvasWidth := myBox.TotalWidth() + padding*2
	canvasHeight := myBox.TotalHeight() + padding*2

	// Start SVG
	canvas := svg.New(f)
	canvas.Start(canvasWidth, canvasHeight)
	canvas.Rect(0, 0, canvasWidth, canvasHeight, "fill:white")

	// Add metadata
	canvas.Title(fmt.Sprintf("Box %.0fx%.0fx%.0f mm", *width, *depth, *height))
	canvas.Desc(fmt.Sprintf("Generated on %s - Box dimensions: %.0fx%.0fx%.0f mm",
		time.Now().Format("2006-01-02 15:04:05"), *width, *depth, *height))

	// Transform the group for proper positioning
	canvas.Group(fmt.Sprintf("transform=\"translate(%d, %d)\"", padding, padding))

	logger.Info("Added fold lines")
	canvas.Path(paths["fold_lines"], redDotted)

	logger.Info("Added cut lines")
	canvas.Path(paths["cut_lines"], blueSolid)

	logger.Info("Add a rectangle")
	canvas.Path(paths["rect"], greenSolid)

	canvas.Gend() // End transform group
	canvas.End()

	logger.Info("Box generated:",
		"file", filename,
		"size", fmt.Sprintf("%dx%d px", canvasWidth, canvasHeight))

	// Print summary
	fmt.Printf("Box generated:\n")
	fmt.Printf("  File: %s\n", filename)
	fmt.Printf("  Dimensions: %.0f×%.0f×%.0f mm\n", *width, *depth, *height)

}
