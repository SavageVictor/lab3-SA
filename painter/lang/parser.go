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

func (p *Parser) Parse(in io.Reader, figureList *painter.FigureList) ([]painter.Operation, error) {
	scanner := bufio.NewScanner(in)
	scanner.Split(bufio.ScanLines)

	var res []painter.Operation
	for scanner.Scan() {
		commandLine := scanner.Text()
		op, err := p.parse(commandLine, figureList)
		if err != nil {
			return nil, err
		}
		res = append(res, op)
	}
	return res, scanner.Err()
}

func (p *Parser) parse(commandLine string, figureList *painter.FigureList) (painter.Operation, error) {
	encodedCommand := strings.TrimSpace(commandLine)
	decodedCommand, err := url.QueryUnescape(encodedCommand)
	if err != nil {
		return nil, fmt.Errorf("Error decoding command: %v", err)
	}
	cmdFields := strings.Fields(decodedCommand)

	switch cmdFields[0] {
	case "cmd=white":
		return painter.OperationFunc(painter.WhiteFill(figureList)), nil
	case "cmd=green":
		return painter.OperationFunc(painter.GreenFill(figureList)), nil
	case "cmd=update":
		return painter.UpdateOp, nil
	case "cmd=bgrect":
		if len(cmdFields) == 5 {
			x1, _ := strconv.ParseFloat(cmdFields[1], 32)
			y1, _ := strconv.ParseFloat(cmdFields[2], 32)
			x2, _ := strconv.ParseFloat(cmdFields[3], 32)
			y2, _ := strconv.ParseFloat(cmdFields[4], 32)
			return painter.OperationFunc(painter.DrawBlackSquare(x1, y1, x2, y2, figureList)), nil
		} else {
			return nil, fmt.Errorf("Incorrect number of arguments for bgrect command")
		}
	case "cmd=figure":
		if len(cmdFields) == 3 {
			x, _ := strconv.ParseFloat(cmdFields[1], 32)
			y, _ := strconv.ParseFloat(cmdFields[2], 32)
			return painter.OperationFunc(painter.DrawFigure(figureList, x, y)), nil
		} else {
			return nil, fmt.Errorf("Incorrect number of arguments for figure command")
		}
	case "cmd=move":
		if len(cmdFields) == 3 {
			x, _ := strconv.ParseFloat(cmdFields[1], 32)
			y, _ := strconv.ParseFloat(cmdFields[2], 32)
			return painter.OperationFunc(painter.MoveFigures(figureList, x, y)), nil
		} else {
			return nil, fmt.Errorf("Incorrect number of arguments for move command")
		}
	case "cmd=reset":
		return painter.OperationFunc(painter.ResetScreen), nil
	default:
		return nil, fmt.Errorf("Invalid command: %s", cmdFields[0])
	}
}
