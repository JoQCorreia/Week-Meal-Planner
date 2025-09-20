package main

import (
	"database/sql"
	"fmt"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
	"log"
	"modernc.org/sqlite"
	"os"
)

func main() {
	fmt.Println("Seedling statement (more2com)")

	a := app.New()
	w := a.NewWindow("Hello World")
	w.SetContent(widget.NewLabel("Hello World!"))
	w.ShowAndRun()

}
