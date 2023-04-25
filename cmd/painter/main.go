package main

import (
	"net/http"

	"github.com/SavageVictor/lab3-SA/painter"
	"github.com/SavageVictor/lab3-SA/painter/lang"
	"github.com/SavageVictor/lab3-SA/ui"
)

func main() {
	var (
		pv ui.Visualizer // Візуалізатор створює вікно та малює у ньому.

		// Потрібні для частини 2.
		opLoop painter.Loop // Цикл обробки команд.
		parser lang.Parser  // Парсер команд.
	)

	//pv.Debug = true
	pv.Title = "Simple painter"

	pv.OnScreenReady = opLoop.Start
	opLoop.Receiver = &pv

	figureList := painter.FigureList{}

	go func() {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "cmd/painter/index.html")
		})
		http.Handle("/parse", lang.HttpHandler(&opLoop, &parser, &figureList))
		_ = http.ListenAndServe("localhost:17000", nil)
	}()

	pv.Main()
	opLoop.StopAndWait()
}
