package main

import (
	"database/sql"
	_ "embed"
	"fmt"
	"image/color"
	_ "image/png"
	_ "io"
	"log"
	_ "os"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	_ "fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	_ "modernc.org/sqlite" //use the side effects of the sqlite driver but not the package
)

// Database handle
var db *sql.DB

type Receitas struct {
	ID       int64
	Receita  string
	Tipo     string
	Proteina string
	Domingo  string
}

// Creating and implementing the custom theme for the background color
type myTheme struct{}

var _ fyne.Theme = (*myTheme)(nil) //type assertion

func (m myTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {

	switch name {
	case theme.ColorNameBackground:
		return color.RGBA{206, 199, 177, 0} //bege

	case theme.ColorNameButton:
		return color.RGBA{159, 185, 74, 0} //light green

	case theme.ColorNameForeground:
		return color.White

	default:
		return theme.DefaultTheme().Color(name, variant)
	}
}

func (m myTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (m myTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (m myTheme) Size(name fyne.ThemeSizeName) float32 {
	if name == theme.SizeNameText {
		return 20
	}
	return theme.DefaultTheme().Size(name)
}

var a fyne.App
var w fyne.Window

//go:embed data/ReceitasFinal.db
var database []byte

//TODO embed database as []byte

func main() {

	var err error
	a = app.NewWithID("Menu Da Semana")
	// //Getting database handle for queries

	// //creating temp file to store db file
	// tempFile, err := os.CreateTemp(wd.RootURI().Path(), "database")
	// if err != nil {
	// 	fmt.Print(err, tempFile)
	// }
	// tempFile.Write(database)

	// defer os.Remove(tempFile.Name())

	db, err := sql.Open("sqlite", "")
	if err != nil {
		log.Fatal("I did not open the database handle")
	}

	//Pinging the database to be sure it's connected
	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}

	//fmt.Println("Connection sucessful")

	// Starting and configuring main window

	a.Settings().SetTheme(&myTheme{}) // calling the custom theme
	w = a.NewWindow("Refeições da semana")
	w.CenterOnScreen() // Window to the center of screen
	w.Resize(fyne.NewSize(380, 800))
	w.SetPadded(false)
	w.SetFullScreen(false)
	i := theme.GridIcon()
	w.SetIcon(i)

	//setting main window images and layout
	// foo := a.Storage().RootURI()
	// fmt.Print(foo)

	gui := imageOpen() //Slice with canvas.Image entries for layout

	//Background image
	gui[0].Resize(fyne.NewSize(380, 800))
	gui[0].SetMinSize(fyne.NewSize(380, 800))
	gui[0].FillMode = canvas.ImageFillContain
	backgroundLayout := container.NewCenter(gui[0])

	//Top and bottom banner and container
	gui[1].FillMode = canvas.ImageFillContain
	gui[1].Resize(fyne.NewSize(280, 250))
	gui[1].SetMinSize(fyne.NewSize(280, 250))
	bannerLayout := container.New(layout.NewHBoxLayout(), layout.NewSpacer(), gui[1])

	gui[2].FillMode = canvas.ImageFillStretch
	gui[2].SetMinSize(fyne.NewSize(380, 100))
	footerLayout := container.NewVBox(layout.NewSpacer(), gui[2])
	topDownLayout := container.New(layout.NewVBoxLayout(), bannerLayout, layout.NewSpacer(), layout.NewSpacer())

	UIBox := container.NewStack(backgroundLayout, topDownLayout)

	//Button on initial screen
	button := widget.NewButton("Criar menu", recipeButton)
	button.Resize(fyne.NewSize(100, 100))
	buttonLayout := container.New(layout.NewCustomPaddedLayout(0, 55, 10, 0), container.NewVBox(layout.NewSpacer(), container.New(layout.NewHBoxLayout(), layout.NewSpacer(), button, layout.NewSpacer())))
	footerButtonLayout := container.NewStack(footerLayout, buttonLayout)
	graphicLayout := container.NewStack(topDownLayout, footerButtonLayout)

	layout := container.NewCenter(container.NewStack(UIBox, graphicLayout)) //Content to the center of container with layout

	//UI update
	w.SetContent(layout)
	w.ShowAndRun()
}

func calendar() string {
	//Calculating the date of the Monday where the menu has to start
	dateNow := time.Now()
	date := 0
	weekdayFunc := dateNow.AddDate(0, 0, date)

	//Calculating how many days until the next monday (no more infinite loops (҂◡_◡) ᕤ)
	for weekdayFunc.Weekday().String() != "Monday" {
		date += 1
		weekdayFunc = dateNow.AddDate(0, 0, date)
	}
	startDay := dateNow.Day() + date
	monthNow := dateNow.Month().String()

	//Determining if the menu date is in a new month and returning the final
	switch monthNow {
	case "January", "March", "May", "July", "August", "October", "December":
		if startDay <= 31 {
			dateString := strconv.Itoa(startDay) + " " + datePT(monthNow) + " " + strconv.Itoa(dateNow.Year())
			return dateString
		}
		monthInt := dateNow.Month() + 1
		dateNewMonth := strconv.Itoa(31 - startDay)
		dateString := dateNewMonth + " " + datePT(monthInt.String()) + " " + strconv.Itoa(dateNow.Year())

		return dateString

	case "February":

		if dateNow.Year()%4 == 0 && dateNow.Year()%100 != 0 {
			if startDay <= 29 {
				dateString := strconv.Itoa(startDay) + " " + datePT(monthNow) + " " + strconv.Itoa(dateNow.Year())
				return dateString
			}
			monthInt := dateNow.Month() + 1
			dateNewMonth := strconv.Itoa(29 - startDay)
			dateString := dateNewMonth + " " + datePT(monthInt.String()) + " " + strconv.Itoa(dateNow.Year())

			return dateString
		} else {
			if startDay <= 28 {
				dateString := strconv.Itoa(startDay) + " " + datePT(monthNow) + " " + strconv.Itoa(dateNow.Year())
				return dateString
			} else {
				monthInt := dateNow.Month() + 1
				dateNewMonth := strconv.Itoa(28 - startDay)
				dateString := dateNewMonth + " " + datePT(monthInt.String()) + " " + strconv.Itoa(dateNow.Year())

				return dateString
			}
		}

	default:
		if startDay <= 30 {
			dateString := strconv.Itoa(startDay) + " " + datePT(monthNow) + " " + strconv.Itoa(dateNow.Year())
			return dateString
		}
		monthInt := dateNow.Month() + 1
		dateNewMonth := strconv.Itoa(30 - startDay)
		dateString := dateNewMonth + " " + datePT(monthInt.String()) + " " + strconv.Itoa(dateNow.Year())

		return dateString

	}
}

func datePT(m string) string {
	//Translating month names returned by time package
	switch m {
	case "January":
		return "Janeiro"
	case "February":
		return "Fevereiro"
	case "March":
		return "Março"
	case "April":
		return "Abril"
	case "June":
		return "Junho"
	case "July":
		return "Julho"
	case "August":
		return "Agosto"
	case "September":
		return "Setembro"
	case "October":
		return "Outubro"
	case "November":
		return "Novembro"
	case "December":
		return "Dezembro"
	}
	return ""
}

func imageOpen() []*canvas.Image {
	//Criating and converting fyne.Resource into background images

	files := []*fyne.StaticResource{resourceDataGUIf2Png, resourceDataGUIa2CPng, resourceDataGUIa3Png}

	var gui []*canvas.Image

	for _, f := range files {
		parsedImage := canvas.NewImageFromResource(f)
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

	var dias []string = []string{"Segunda", "Terça", "Quarta", "Quinta", "Sexta", "Sábado", "Domingo"}

	// Recipe page layout

	rectHeader := canvas.NewRectangle(color.NRGBA{R: 159, G: 185, B: 74, A: 180})
	rectHeader.Resize(fyne.NewSize(380, 100))
	rectHeader.SetMinSize(fyne.NewSize(380, 100))

	dateHeader := canvas.NewText(calendar(), color.White)

	stringHeader := canvas.NewText("SEMANA", color.White)
	stringHeader.TextStyle = fyne.TextStyle{Bold: true}

	textHeader := container.New(layout.NewHBoxLayout(), layout.NewSpacer(), stringHeader, dateHeader, layout.NewSpacer())

	headerSemana := container.NewStack(rectHeader, textHeader)

	//Monday container
	mondayName := canvas.NewText(dias[0], theme.Color(theme.ColorNameForeground))
	mondayName.TextStyle.Bold = true
	mondayName.TextSize = 18
	mondayPadding := container.NewCenter(mondayName)

	mondayRecipes := container.New(layout.NewCustomPaddedLayout(5, 0, 10, 0), container.New(layout.NewGridWrapLayout(fyne.NewSize(380, 75)), lista([]Receitas{receitasCarne[0], receitasPeixe[0]})))

	mondayCont := container.NewCenter(container.NewVBox(mondayPadding, mondayRecipes))

	//Tuesday container
	tuesdayName := canvas.NewText(dias[1], theme.Color(theme.ColorNameForeground))
	tuesdayName.TextStyle.Bold = true
	tuesdayName.TextSize = 18
	tuesdayPadding := container.NewCenter(tuesdayName)

	tuesdayRecipes := container.New(layout.NewCustomPaddedLayout(5, 0, 10, 0), container.New(layout.NewGridWrapLayout(fyne.NewSize(380, 75)), lista([]Receitas{receitasCarne[1], receitasPeixe[1]})))

	tuesdayCont := container.NewCenter(container.NewVBox(tuesdayPadding, tuesdayRecipes))

	//Wed container
	wedName := canvas.NewText(dias[2], theme.Color(theme.ColorNameForeground))
	wedName.TextStyle.Bold = true
	wedName.TextSize = 18
	wedPadding := container.NewCenter(wedName)

	wedRecipes := container.New(layout.NewCustomPaddedLayout(5, 0, 10, 0), container.New(layout.NewGridWrapLayout(fyne.NewSize(380, 75)), lista([]Receitas{receitasCarne[2], receitasPeixe[2]})))

	wedCont := container.NewCenter(container.NewVBox(wedPadding, wedRecipes))

	//Thu container
	thuName := canvas.NewText(dias[3], theme.Color(theme.ColorNameForeground))
	thuName.TextStyle.Bold = true
	thuName.TextSize = 18
	thuPadding := container.NewCenter(thuName)

	thuRecipes := container.New(layout.NewCustomPaddedLayout(5, 0, 10, 0), container.New(layout.NewGridWrapLayout(fyne.NewSize(380, 75)), lista([]Receitas{receitasCarne[3], receitasPeixe[3]})))

	thuCont := container.NewCenter(container.NewVBox(thuPadding, thuRecipes))

	//Fri container
	friName := canvas.NewText(dias[4], theme.Color(theme.ColorNameForeground))
	friName.TextStyle.Bold = true
	friName.TextSize = 18
	friPadding := container.NewCenter(friName)
	friRecipes := container.New(layout.NewCustomPaddedLayout(5, 0, 10, 0), container.New(layout.NewGridWrapLayout(fyne.NewSize(380, 75)), lista([]Receitas{receitasCarne[4], receitasPeixe[4]})))

	friCont := container.NewCenter(container.NewVBox(friPadding, friRecipes))

	//Sat container
	satName := canvas.NewText(dias[5], theme.Color(theme.ColorNameForeground))
	satName.TextStyle.Bold = true
	satName.TextSize = 18
	satPadding := container.NewCenter(satName)

	satRecipes := container.New(layout.NewCustomPaddedLayout(5, 0, 10, 0), container.New(layout.NewGridWrapLayout(fyne.NewSize(380, 75)), lista([]Receitas{receitasCarne[5], receitasPeixe[5]})))

	satCont := container.NewCenter(container.NewVBox(satPadding, satRecipes))

	//Sunday container
	sundayName := canvas.NewText(dias[6], theme.Color(theme.ColorNameForeground))
	sundayName.TextStyle.Bold = true
	sundayName.TextSize = 18
	sundayPadding := container.NewCenter(sundayName)

	sundayRecipes := container.New(layout.NewCustomPaddedLayout(5, 0, 10, 0), container.New(layout.NewGridWrapLayout(fyne.NewSize(380, 75)), lista([]Receitas{receitasDomingo[0]})))

	sundayCont := container.NewCenter(container.NewVBox(sundayPadding, sundayRecipes))

	//final UI

	UIReceitas := container.NewCenter(container.New(layout.NewCustomPaddedVBoxLayout(4), mondayCont, tuesdayCont, wedCont, thuCont, friCont, satCont, sundayCont))
	content := container.NewCenter(container.New(layout.NewVBoxLayout(), headerSemana, UIReceitas))

	themeOverrideContainer := container.NewThemeOverride(content, theme.DefaultTheme())

	w.SetContent(themeOverrideContainer)

}

func lista(receitas []Receitas) *widget.List {
	//List generation with recipes
	list := widget.NewList(
		func() int {
			return len(receitas)
		},
		func() fyne.CanvasObject {
			return widget.NewLabelWithStyle("Receitas", fyne.TextAlignCenter, fyne.TextStyle{})
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
