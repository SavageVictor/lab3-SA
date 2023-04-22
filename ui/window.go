package ui

import (
	"image"
	"image/color"
	"log"

	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/imageutil"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/image/draw"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/mouse"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
)

type Visualizer struct {
	Title         string
	Debug         bool
	OnScreenReady func(s screen.Screen)

	w    screen.Window
	tx   chan screen.Texture
	done chan struct{}

	sz  size.Event
	pos image.Rectangle
}

func (pw *Visualizer) Main() {
	pw.tx = make(chan screen.Texture)
	pw.done = make(chan struct{})
	pw.pos.Max.X = 200
	pw.pos.Max.Y = 200
	driver.Main(pw.run)
}

func (pw *Visualizer) Update(t screen.Texture) {
	pw.tx <- t
}

func (pw *Visualizer) run(s screen.Screen) {
	if pw.OnScreenReady != nil {
		pw.OnScreenReady(s)
	}

	w, err := s.NewWindow(&screen.NewWindowOptions{
		Title:  pw.Title,
		Width:  800, // Set the initial window width to 800
		Height: 800, // Set the initial window height to 800
	})
	if err != nil {
		log.Fatal("Failed to initialize the app window:", err)
	}
	defer func() {
		w.Release()
		close(pw.done)
	}()

	pw.w = w

	events := make(chan any)
	go func() {
		for {
			e := w.NextEvent()
			if pw.Debug {
				log.Printf("new event: %v", e)
			}
			if detectTerminate(e) {
				close(events)
				break
			}
			events <- e
		}
	}()

	var t screen.Texture

	for {
		select {
		case e, ok := <-events:
			if !ok {
				return
			}
			pw.handleEvent(e, t)

		case t = <-pw.tx:
			w.Send(paint.Event{})
		}
	}
}

func detectTerminate(e any) bool {
	switch e := e.(type) {
	case lifecycle.Event:
		if e.To == lifecycle.StageDead {
			return true // Window destroy initiated.
		}
	case key.Event:
		if e.Code == key.CodeEscape {
			return true // Esc pressed.
		}
	}
	return false
}

func (pw *Visualizer) handleEvent(e any, t screen.Texture) {
	switch e := e.(type) {

	case size.Event: // Update window size data.
		pw.sz = e

	case error:
		log.Printf("ERROR: %s", e)

	case mouse.Event:
		if t == nil && e.Button == mouse.ButtonLeft && e.Direction == mouse.DirPress {
			// Move the figure to the coords of a mouse click
			pw.pos = image.Rect(int(e.X)-100, int(e.Y)-100, int(e.X)+100, int(e.Y)+100)
			pw.w.Send(paint.Event{})
		}

	case paint.Event:
		// Draw window content.
		if t == nil {
			pw.drawDefaultUI()
		} else {
			// Use texture received through the Update call.
			pw.w.Scale(pw.sz.Bounds(), t, t.Bounds(), draw.Src, nil)
		}
		pw.w.Publish()
	}
}

func (pw *Visualizer) drawDefaultUI() {
	pw.w.Fill(pw.sz.Bounds(), color.Black, draw.Src) // Background.

	// Render yellow cross figure.
	pw.drawYellowCross()

	// Draw white border.
	for _, br := range imageutil.Border(pw.sz.Bounds(), 10) {
		pw.w.Fill(br, color.White, draw.Src)
	}
}

func (pw *Visualizer) drawYellowCross() {
	crossColor := color.RGBA{255, 255, 0, 255}
	lineWidth := 10

	// Horizontal line of the cross.
	horizontalLine := image.Rect(pw.pos.Min.X, pw.pos.Min.Y+(pw.pos.Dy()/2)-lineWidth/2, pw.pos.Max.X, pw.pos.Min.Y+(pw.pos.Dy()/2)+lineWidth/2)
	pw.w.Fill(horizontalLine, crossColor, draw.Src)

	// Vertical line of the cross.
	verticalLine := image.Rect(pw.pos.Min.X+(pw.pos.Dx()/2)-lineWidth/2, pw.pos.Min.Y, pw.pos.Min.X+(pw.pos.Dx()/2)+lineWidth/2, pw.pos.Max.Y)
	pw.w.Fill(verticalLine, crossColor, draw.Src)
}
