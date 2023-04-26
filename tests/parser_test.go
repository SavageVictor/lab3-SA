package test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/SavageVictor/lab3-SA/painter"
	"github.com/SavageVictor/lab3-SA/painter/lang"
	"github.com/stretchr/testify/assert"
)

func TestParseCommands(t *testing.T) {
	testCases := []struct {
		name       string
		input      string
		expectedOp string
	}{
		{
			name:       "White Fill",
			input:      "cmd=white",
			expectedOp: "OperationFunc",
		},
		{
			name:       "Green Fill",
			input:      "cmd=green",
			expectedOp: "OperationFunc",
		},
		{
			name:       "Update",
			input:      "cmd=update",
			expectedOp: "updateOp",
		},
		{
			name:       "Draw Black Square",
			input:      "cmd=bgrect 0.2 0.2 0.6 0.6",
			expectedOp: "OperationFunc",
		},
		{
			name:       "Draw Figure",
			input:      "cmd=figure 0.3 0.7",
			expectedOp: "OperationFunc",
		},
		{
			name:       "Move Figures",
			input:      "cmd=move 0.5 0.6",
			expectedOp: "OperationFunc",
		},
		{
			name:       "Reset Screen",
			input:      "cmd=reset",
			expectedOp: "OperationFunc",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			parser := &lang.Parser{}
			figureList := painter.FigureList{}
			ops, err := parser.Parse(strings.NewReader(tc.input), &figureList)
			assert.NoError(t, err, "Unexpected error")

			assert.Len(t, ops, 1, "Expected 1 operation")

			opName := strings.Split(fmt.Sprintf("%T", ops[0]), ".")[1]

			assert.Equal(t, tc.expectedOp, opName, "Expected operation")
		})
	}
}
