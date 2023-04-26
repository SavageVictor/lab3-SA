package test

import (
	"testing"

	"github.com/SavageVictor/lab3-SA/painter"
	"github.com/stretchr/testify/assert"
)

func TestWhiteFill(t *testing.T) {

	fl := &painter.FigureList{}

	op := painter.WhiteFill(fl)

	assert.NotNil(t, op)
}

func TestGreenFill(t *testing.T) {

	fl := &painter.FigureList{}

	op := painter.GreenFill(fl)

	assert.NotNil(t, op)
}

func TestDrawFigure(t *testing.T) {
	fl := &painter.FigureList{}

	op := painter.DrawFigure(fl, 0.5, 0.5)

	assert.NotNil(t, op, fl)
}
func TestBlackSquare(t *testing.T) {
	fl := &painter.FigureList{}

	op := painter.DrawBlackSquare(0.5, 0.5, 0.7, 0.7, fl)

	assert.NotNil(t, op, fl)
}

func TestDrawFigureWithoutAdding(t *testing.T) {

	op := painter.DrawFigureWithoutListing(0.5, 0.5)

	assert.NotNil(t, op)
}

func TestBlackSquareWithoutAdding(t *testing.T) {

	op := painter.DrawBlackSquareWithoutListing(0.5, 0.5, 0.7, 0.7)

	assert.NotNil(t, op)
}
