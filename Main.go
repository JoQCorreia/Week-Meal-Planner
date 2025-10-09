package main

import (
	"database/sql"
	"fmt"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
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

	//database query
	receitas, err := queryReceitas()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(receitas[0])

	//UI update
	lista(receitas)

}

func lista(receitas []Receitas) {
	//Fyne UI management
	a := app.New()
	w := a.NewWindow("Refeições da semana")

	fmt.Println(receitas[0].Receita)
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

	w.SetContent(list)
	w.ShowAndRun()
}

func queryReceitas() ([]Receitas, error) {
	var receitas []Receitas

	rows, err := db.Query("SELECT * FROM receitas")
	if err != nil {
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
