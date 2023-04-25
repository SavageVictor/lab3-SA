package lang

import (
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/SavageVictor/lab3-SA/painter"
)

// HttpHandler конструює обробник HTTP запитів, який дані з запиту віддає у Parser, а потім відправляє отриманий список
// операцій у painter.Loop.
func HttpHandler(loop *painter.Loop, p *Parser, figureList *painter.FigureList) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.ServeFile(rw, r, "index.html")
			return
		}

		var in io.Reader = r.Body
		if r.Method == http.MethodGet {
			in = strings.NewReader(r.URL.Query().Get("cmd"))
		}

		err := p.Parse(in, loop, figureList)
		if err != nil {
			log.Printf("Bad script: %s", err)
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		//loop.Post(painter.OperationList(cmds))
		rw.WriteHeader(http.StatusOK)
	})
}
