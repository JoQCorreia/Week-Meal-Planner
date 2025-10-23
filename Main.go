package main

import (
	"database/sql"
	"fmt"
	//"image/color"
	_ "image/color"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	//"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	_ "modernc.org/sqlite" //use the side effects of the sqlite driver but not the package
	//"os"
)

var db *sql.DB

type Receitas struct {
	ID       int64
	Receita  string
	Tipo     string
	Proteina string
	Domingo  string
}

func main() {
	fmt.Println("Seedling statement (more2com)")
	var err error

	//Getting db handle for queries
	db, err = sql.Open("sqlite", "./database/Receitas2.db")
	if err != nil {
		log.Fatal(err)
	}

	//Pinging the database to be sure it's connected
	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}

	fmt.Println("Connection sucessful")

	a := app.New()
	w := a.NewWindow("Refeições da semana")

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

	//UI update
	w.SetContent(content)
	w.ShowAndRun()

}

func lista(receitas []Receitas) *widget.List {
	//List generation with recipes
	list := widget.NewList(
		func() int {
			return len(receitas)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
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
			rows, err := db.Query("SELECT * FROM receitas WHERE tipo = " + "\"" + tipo + "\"" + " AND domingo = 'false' ORDER BY random() LIMIT 8;")
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
