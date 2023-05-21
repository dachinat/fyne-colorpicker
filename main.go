package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/atotto/clipboard"
	"github.com/dachinat/colornameconv"
	"github.com/go-vgo/robotgo"
	hook "github.com/robotn/gohook"
	"image/color"
)

func main() {
	a := app.New()
	w := a.NewWindow("Color Picker")

	var data = []string{}

	rect := canvas.NewRectangle(color.White)

	list := widget.NewList(
		func() int {
			return len(data)
		},
		func() fyne.CanvasObject {
			return container.NewBorder(
				nil,
				nil,
				widget.NewLabel("template"),
				widget.NewButtonWithIcon("", theme.ContentCopyIcon(), nil),
			)
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			item.(*fyne.Container).Objects[0].(*widget.Label).SetText("#" + data[id] + " (" + HexToName(data[id]) + ")")

			item.(*fyne.Container).Objects[1].(*widget.Button).OnTapped = func() {
				clipboard.WriteAll("#" + data[id])
			}
		},
	)

	list.OnSelected = func(id widget.ListItemID) {
		c, _ := ParseHexColor("#" + data[id])
		rect.FillColor = c
		rect.Refresh()
	}

	var btn *widget.Button
	var btn2 *widget.Button

	btn = widget.NewButton("Start Picking", func() {
		mouse(rect, &data, list, btn, btn2)
	})

	btn2 = widget.NewButton("Clear list", func() {
		data = data[:0]
		list.Refresh()
	})

	list.Resize(fyne.NewSize(300, 300))

	content := container.NewGridWithColumns(2, rect, btn)
	w.SetContent(
		container.NewBorder(content, btn2, nil, nil, list),
	)

	w.Resize(fyne.NewSize(580, 380))
	w.ShowAndRun()
}

func mouse(rect *canvas.Rectangle, data *[]string, list *widget.List, btn *widget.Button, btn2 *widget.Button) {

	hook.Register(hook.MouseMove, []string{}, func(e hook.Event) {
		c := robotgo.GetPixelColor(int(e.X), int(e.Y))

		btn.Text = "..."
		btn.Refresh()

		newColor, _ := ParseHexColor("#" + c)
		rect.FillColor = newColor
		rect.Refresh()
	})

	hook.Register(hook.MouseDown, []string{}, func(e hook.Event) {
		c := robotgo.GetPixelColor(int(e.X), int(e.Y))

		*data = append(*data, c)

		btn.Text = "Pick another color"
		btn.Refresh()

		newColor, _ := ParseHexColor("#" + c)
		rect.FillColor = newColor
		rect.Refresh()

		list.Refresh()

		hook.End()
	})

	s := hook.Start()
	<-hook.Process(s)
}

func HexToName(hex string) string {
	name, _ := colornameconv.New(hex)
	return name
}

func ParseHexColor(s string) (c color.RGBA, err error) {
	c.A = 0xff
	switch len(s) {
	case 7:
		_, err = fmt.Sscanf(s, "#%02x%02x%02x", &c.R, &c.G, &c.B)
	case 4:
		_, err = fmt.Sscanf(s, "#%1x%1x%1x", &c.R, &c.G, &c.B)
		// Double the hex digits:
		c.R *= 17
		c.G *= 17
		c.B *= 17
	default:
		err = fmt.Errorf("invalid length, must be 7 or 4")
	}
	return
}
