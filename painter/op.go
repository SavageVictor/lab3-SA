package painter

import (
	"image"
	"image/color"

	"golang.org/x/exp/shiny/screen"
)

// Operation modifies the input texture.
type Operation interface {
	Do(t screen.Texture) bool
}

// OperationList groups a list of operations into one.
type OperationList []Operation

func (ol OperationList) Do(t screen.Texture) (ready bool) {
	for _, o := range ol {
		ready = o.Do(t) || ready
	}
	return
}

// UpdateOp is an operation that does not modify the texture, but signals that the texture should be considered ready.
var UpdateOp = updateOp{}

type updateOp struct{}

func (op updateOp) Do(t screen.Texture) bool { return true }

// OperationFunc is used to convert a texture update function into an Operation.
type OperationFunc func(t screen.Texture)

func (f OperationFunc) Do(t screen.Texture) bool {
	f(t)
	return false
}

// WhiteFill fills the texture with white color. Can be used as an Operation through OperationFunc(WhiteFill).
func WhiteFill(fl *FigureList) OperationFunc {
	return func(t screen.Texture) {
		t.Fill(t.Bounds(), color.White, screen.Src)
		for _, figure := range fl.FigureSq {
			OperationFunc.Do(DrawBlackSquareWithoutListing(figure.X1, figure.Y1, figure.X2, figure.Y2), t)
		}
		for _, figure := range fl.Figures {
			OperationFunc.Do(DrawFigureWithoutListing(figure.X, figure.Y), t)
		}
	}
}

// GreenFill fills the texture with green color. Can be used as an Operation through OperationFunc(GreenFill).
func GreenFill(fl *FigureList) OperationFunc {
	return func(t screen.Texture) {
		t.Fill(t.Bounds(), color.RGBA{G: 0xff, A: 0xff}, screen.Src)
		for _, figure := range fl.FigureSq {
			OperationFunc.Do(DrawBlackSquareWithoutListing(figure.X1, figure.Y1, figure.X2, figure.Y2), t)
		}
		for _, figure := range fl.Figures {
			OperationFunc.Do(DrawFigureWithoutListing(figure.X, figure.Y), t)
		}
	}
}

func ResetScreen(t screen.Texture) {
	t.Fill(t.Bounds(), color.RGBA{A: 0xff}, screen.Src)
}

// DrawBlackSquare draws a black square on the texture using the specified floating-point coordinates.
func DrawBlackSquare(x1, y1, x2, y2 float64, fl *FigureList) OperationFunc {
	return func(t screen.Texture) {
		b := t.Bounds()
		rect := image.Rect(int(float64(b.Dx())*x1), int(float64(b.Dy())*y1), int(float64(b.Dx())*x2), int(float64(b.Dy())*y2))
		black := color.RGBA{A: 0xff}
		t.Fill(rect, black, screen.Src)

		figure := FigureSq{
			X1: x1,
			Y1: y1,
			X2: x2,
			Y2: y2,
		}
		fl.FigureSq = append(fl.FigureSq, figure)
	}
}

type FigureSq struct {
	X1, Y1, X2, Y2 float64
}

type Figure struct {
	X, Y          float64
	Width, Height int
}

type FigureList struct {
	Figures  []Figure
	FigureSq []FigureSq
}

// DrawFigure creates a yellow cross and draws it on the texture using the specified floating-point coordinates.
func DrawFigure(fl *FigureList, x, y float64) OperationFunc {
	return func(t screen.Texture) {
		b := t.Bounds()

		// Create a yellow cross using two rectangles
		rectWidth, rectHeight := 30, 10
		yellow := color.RGBA{R: 255, G: 255, B: 0, A: 0}

		// Draw horizontal rectangle
		horizRect := image.Rect(int(float64(b.Dx())*x)-rectWidth/2, int(float64(b.Dy())*y)-rectHeight/2, int(float64(b.Dx())*x)+rectWidth/2, int(float64(b.Dy())*y)+rectHeight/2)
		t.Fill(horizRect, yellow, screen.Src)

		// Draw vertical rectangle
		vertRect := image.Rect(int(float64(b.Dx())*x)-rectHeight/2, int(float64(b.Dy())*y)-rectWidth/2, int(float64(b.Dx())*x)+rectHeight/2, int(float64(b.Dy())*y)+rectWidth/2)
		t.Fill(vertRect, yellow, screen.Src)

		figure := Figure{
			X:      x,
			Y:      y,
			Width:  rectWidth,
			Height: rectHeight,
		}
		fl.Figures = append(fl.Figures, figure)
	}
}

// MoveFigures moves all the figures in the FigureList to the specified floating-point coordinates.
// З цим можна ще багато чого було б зробитиб але воно вимогам лабораторної підходить тому залишу так
func MoveFigures(fl *FigureList, x, y float64) OperationFunc {
	return func(t screen.Texture) {
		b := t.Bounds()
		for i := range fl.Figures {
			figure := fl.Figures[i]

			// Clear previous figure position
			clearRect1 := image.Rect(int(float64(b.Dx())*figure.X)-figure.Height/2, int(float64(b.Dx())*figure.Y)-figure.Width/2, int(float64(b.Dx())*figure.X)+figure.Height/2, int(float64(b.Dx())*figure.X)+figure.Width/2)
			t.Fill(clearRect1, color.Transparent, screen.Src)
			clearRect2 := image.Rect(int(float64(b.Dx())*figure.X)-figure.Width/2, int(float64(b.Dx())*figure.Y)-figure.Height/2, int(float64(b.Dx())*figure.X)+figure.Width/2, int(float64(b.Dx())*figure.X)+figure.Height/2)
			t.Fill(clearRect2, color.Transparent, screen.Src)

			// Update figure position
			fl.Figures[i].X = x
			fl.Figures[i].Y = y
		}

		// Draw figure at new position
		drawFigureOp := DrawFigureWithoutListing(x, y)
		drawFigureOp(t)
	}
}

func DrawFigureWithoutListing(x, y float64) OperationFunc {
	return func(t screen.Texture) {
		b := t.Bounds()

		// Create a yellow cross using two rectangles
		rectWidth, rectHeight := 30, 10
		yellow := color.RGBA{R: 255, G: 255, B: 0, A: 0}

		// Draw horizontal rectangle
		horizRect := image.Rect(int(float64(b.Dx())*x)-rectWidth/2, int(float64(b.Dy())*y)-rectHeight/2, int(float64(b.Dx())*x)+rectWidth/2, int(float64(b.Dy())*y)+rectHeight/2)
		t.Fill(horizRect, yellow, screen.Src)

		// Draw vertical rectangle
		vertRect := image.Rect(int(float64(b.Dx())*x)-rectHeight/2, int(float64(b.Dy())*y)-rectWidth/2, int(float64(b.Dx())*x)+rectHeight/2, int(float64(b.Dy())*y)+rectWidth/2)
		t.Fill(vertRect, yellow, screen.Src)
	}
}

func DrawBlackSquareWithoutListing(x1, y1, x2, y2 float64) OperationFunc {
	return func(t screen.Texture) {
		b := t.Bounds()
		rect := image.Rect(int(float64(b.Dx())*x1), int(float64(b.Dy())*y1), int(float64(b.Dx())*x2), int(float64(b.Dy())*y2))
		black := color.RGBA{A: 0xff}
		t.Fill(rect, black, screen.Src)
	}
}
