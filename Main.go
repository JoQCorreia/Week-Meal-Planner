package main

import (
	"database/sql"
	"fmt"
	//"os"

	_ "image"
	_ "image/color"
	_ "image/png"
	_ "io"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/theme"

	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	_ "modernc.org/sqlite" //use the side effects of the sqlite driver but not the package
)

var db *sql.DB

type Receitas struct {
	ID       int64
	Receita  string
	Tipo     string
	Proteina string
	Domingo  string
}

var a fyne.App
var w fyne.Window

func main() {
	var err error

	//Getting db handle for queries
	db, err = sql.Open("sqlite", "./database/ReceitasFinal.db")
	if err != nil {
		log.Fatal(err)
	}

	//Pinging the database to be sure it's connected
	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}

	fmt.Println("Connection sucessful")

	// Starting and configuring main window

	a = app.New()
	w = a.NewWindow("Refeições da semana")
	w.CenterOnScreen() // Window to the center of screen
	w.Resize(fyne.NewSize(380, 800))
	w.SetPadded(false)
	w.SetFullScreen(false)
	i := theme.GridIcon()
	w.SetIcon(i)

	//setting main window images and layout

	gui := imageOpen() //Slice with canvas.Image entries for layout

	//Background image
	gui[0].Resize(fyne.NewSize(380, 800))
	gui[0].SetMinSize(fyne.NewSize(380, 800))
	gui[0].FillMode = canvas.ImageFillContain
	backgroundLayout := container.NewCenter(gui[0])

	//Top banner

	//Bottom banner

	//Button on initial screen
	button := widget.NewButton("Criar menu", recipeButton)
	button.Resize(fyne.NewSize(100, 100))
	buttonLayout := container.New(layout.NewCenterLayout(), button)
	layout := container.NewStack(backgroundLayout, buttonLayout) //Content to the center of container with layout

	//UI update

	w.SetContent(layout)
	w.ShowAndRun()

}

func imageOpen() []*canvas.Image {
	//Opening and converting image.Image into background images

	files := []string{"D:/Documents/Ementa da semana/GUIf2.png", "D:/Documents/Ementa da semana/GUI2.svg", "D:/Documents/Ementa da semana/GUI3.svg"}
	var gui []*canvas.Image
	for _, f := range files {
		parsed, err := fyne.LoadResourceFromPath(f)
		if err != nil {
			log.Fatal("I got to loading the resources but I failed because:\n", err)
		}
		parsedImage := canvas.NewImageFromResource(parsed)
		gui = append(gui, parsedImage)
	}
	return gui
}

func recipeButton() {
	//database query
	receitasCarne, err := queryReceitas("Carne")
	if err != nil {
		log.Fatal(err)
	}

	receitasPeixe, err := queryReceitas("Peixe")
	if err != nil {
		log.Fatal(err)
	}

	receitasDomingo, err := queryReceitas("Domingo")
	if err != nil {
		log.Fatal(err)
	}

	//meal lists
	textCarne := lista(receitasCarne)
	textPeixe := lista(receitasPeixe)
	textDomingo := lista(receitasDomingo)

	content := container.New(layout.NewGridLayout(3), textCarne, textPeixe, textDomingo)
	w.SetContent(content)

}

func lista(receitas []Receitas) *widget.List {
	//List generation with recipes
	list := widget.NewList(
		func() int {
			return len(receitas)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("Receitas")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(receitas[i].Receita)
		})
	return list
}

func queryReceitas(tipo string) ([]Receitas, error) {
	var receitas []Receitas
	//recipe query for the different types
	switch tipo {
	case "Domingo":
		{
			rows, err := db.Query("SELECT * FROM receitas WHERE domingo = 'true' ORDER BY random() LIMIT 1;")
			if err != nil {
				fmt.Printf("Error here")
				return receitas, err
			}
			for rows.Next() {
				var rec Receitas
				if err := rows.Scan(&rec.ID, &rec.Receita, &rec.Tipo, &rec.Proteina, &rec.Domingo); err != nil {
					return receitas, err
				}
				receitas = append(receitas, rec)
			}
			if err := rows.Err(); err != nil {
				return receitas, err
			}
			return receitas, nil
		}
	default:
		{
			rows, err := db.Query("SELECT * FROM receitas WHERE tipo = '" + tipo + "' AND domingo = 'false' ORDER BY random() LIMIT 8;")
			if err != nil {
				fmt.Printf("Error here")
				return receitas, err
			}
			for rows.Next() {
				var rec Receitas
				if err := rows.Scan(&rec.ID, &rec.Receita, &rec.Tipo, &rec.Proteina, &rec.Domingo); err != nil {
					return receitas, err
				}
				receitas = append(receitas, rec)
			}
			if err := rows.Err(); err != nil {
				return receitas, err
			}
			return receitas, nil
		}
	}
}
