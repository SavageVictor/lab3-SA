package lang

import (
	"bufio"
	"fmt"
	"io"
	"net/url"
	"strconv"
	"strings"

	"github.com/SavageVictor/lab3-SA/painter"
)

type Parser struct{}

func (p *Parser) Parse(in io.Reader, loop *painter.Loop, figureList *painter.FigureList) error {
	scanner := bufio.NewScanner(in)

	for scanner.Scan() {
		encodedCommand := strings.TrimSpace(scanner.Text())
		decodedCommand, err := url.QueryUnescape(encodedCommand)
		if err != nil {
			return fmt.Errorf("Error decoding command: %v", err)
		}
		cmdFields := strings.Fields(decodedCommand)

		switch cmdFields[0] {
		case "cmd=white":
			loop.Post(painter.OperationFunc(painter.WhiteFill(figureList)))
		case "cmd=green":
			loop.Post(painter.OperationFunc(painter.GreenFill(figureList)))
		case "cmd=update":
			loop.Post(painter.UpdateOp)
		case "cmd=bgrect":
			if len(cmdFields) == 5 {
				x1, _ := strconv.ParseFloat(cmdFields[1], 32)
				y1, _ := strconv.ParseFloat(cmdFields[2], 32)
				x2, _ := strconv.ParseFloat(cmdFields[3], 32)
				y2, _ := strconv.ParseFloat(cmdFields[4], 32)
				loop.Post(painter.OperationFunc(painter.DrawBlackSquare(x1, y1, x2, y2, figureList)))
			} else {
				return fmt.Errorf("Incorrect number of arguments for bgrect command")
			}
		case "cmd=figure":
			if len(cmdFields) == 3 {
				x, _ := strconv.ParseFloat(cmdFields[1], 32)
				y, _ := strconv.ParseFloat(cmdFields[2], 32)
				loop.Post(painter.OperationFunc(painter.DrawFigure(figureList, x, y)))
			} else {
				return fmt.Errorf("Incorrect number of arguments for figure command")
			}
		case "cmd=move":
			if len(cmdFields) == 3 {
				x, _ := strconv.ParseFloat(cmdFields[1], 32)
				y, _ := strconv.ParseFloat(cmdFields[2], 32)
				loop.Post(painter.OperationFunc(painter.MoveFigures(figureList, x, y)))
			} else {
				return fmt.Errorf("Incorrect number of arguments for move command")
			}
		case "cmd=reset":
			loop.Post(painter.OperationFunc(painter.ResetScreen))
		}
	}
	return scanner.Err()
}
