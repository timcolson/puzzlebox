## Puzzle Box Functional Spec

I need a tool to help me build parameterized paths for a box with cut lines and fold lines. 

I need the tool to provide rounded corners between SOME but not all line segments. I also need to create more than one path.

I have found that some lines are easier to describe with absolute coordinates, relative to a fixed 0,0 origin. For some sections, such as the bottom flaps, I found that relative coordinates are easier to specify than absolute. For example, move to the side panel bottom left, then draw a diagonal line to 1/2 the height of the side panel flap to start the structure. 

### Resources
* box-gen.go - maps input parameters to a Box struct, initiates logging, and begins to create the box with fold lines, then cut lines for one "side flap" of the box.
* gui_context.md -  a graphical tool built to create and visualize an SVG in real time as variables are changing.
* svg_path_builder_context.md - a pathbuilder pattern code which provides a way to create paths from MoveTo and LineTo functions with a fluent approach. 

### Tasks
[ ] Incorporate the Box struct into the GUI tool. Note how internal calculations need to be in float64, but external calls will almost always expect an integer.

[ ] GUI - Set up inputs for Width, Depth, Height, and FoldGap

[ ] Design a simple way to edit the path descriptions and see them in real-time without needing to re-run the program. If impossible or overly complex, I can instead run "air" and make the path changes manually in the code.

[ ] Add an editable text field to the GUI where I can specify the paths using the PathBuilder syntax.

[ ] LineTo currently takes absolute coordinates. Create RelativeLine(offsetX,offsetY). Also create convencience functions: HorizontalLine(offset) and VerticalLine(offset) which simply call RelativeLine with 0 for Y and X offset respectively. 


### Technology
* Use golang
* For UI, use Fyne library
* Use SVGO library: github.com/ajstarks/svgo 
