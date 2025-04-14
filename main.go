package main

import (
	"database/sql"
	"embed"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"text/template"
	"time"
)

var (
	port = flag.Int("port", 8000, "specify port number")
	data = flag.String("data", "data", "data's root path")

	db *sql.DB

	//go:embed templates/*.html
	templateFS embed.FS

	pages map[string]*template.Template = make(map[string]*template.Template)
)

var jobs = []string{
	"подрабатываю курьером – развожу еду и продукты всем подряд. Коплю на мопед.",
	"работаю в «Горэлектротрансе» – сооставляю маршруты, сверяюсь с расписанием, слежу чтоб всё шло по рельсам.",
	"работаю звуковым инженером в концертном зале – настраиваю микшер, ловлю настроение музыкантов в реальном времени.",
	"работаю библиотекарем в спальном районе Петербурга – в основном ничего не делаю, залипаю в ютуб шортсы.",
	"работаю техническим переводчиком с немецкого в любимом издательстве – заказчик присылает текст, я перевожу. И так – снова и снова, пока всем всё не понравится",
	"работаю репетитором – готовлю школьников к экзаменам, объясняю как механически и бездумно решать уравнения.",
	"пишу письма за тех, кто мнётся, теряется и не знает как открыть душу близкому. Проститься, признаться, поблагодарить. В общем, пишу за других, красиво и по делу.",
	"проектирую сборку мебели – черчу схемы, думаю, как сделать так, чтобы табуретка собиралась легко и ловко (на деле – копируем стиль икеи).",
	"верстаю меню для ресторанов и кафе – горячее в начале, напитки в самом конце.",
	"продаю выпечку в популярной пекарне – на ногах весь день, за сутки наматываю километров двадцать.",
	"батрачу на автомойке у своего одногруппника – оттираю грязь с капота, копоть с бампера. В общем, чищу тачки.",
	"работаю электриком в «Жилкомсервисе» – хожу по вызовам, возвращаю людям свет.",
	"я разнорабочий на линии упаковки – укладываю сырки и другую молочку в коробки. По пути думаю: кто их будет есть?",
	"работаю грузчиком в мебельном центре – таскаю диваны, иногда собираю на заказ кухню. По вечерам думаю о сыре...",
	"работаю слесарем-сантехником – чиню то, то течёт, капает или не смывает.",
	"стажируюсь техником в ростелекоме – устанавливаю роутеры, прокидываю витую пару, завожу вайфай.",
	"верстаю афиши для цирковых представлений. В основном, просто переиспользую шаблон с разноцветными шариками.",
}

func main() {
	// Parse templates.
	flag.Parse()

	// Init DB.
	var err error
	db, err = initDB(fmt.Sprintf("%s/database.db", *data))
	if err != nil {
		log.Fatal(err)
	}

	// Parse templates for every page.
	pages["index"] = template.Must(template.ParseFS(templateFS,
		"templates/layout.html", "templates/index.html", "templates/ascii.html",
	))
	pages["education"] = template.Must(template.ParseFS(templateFS,
		"templates/layout.html", "templates/education.html",
	))
	pages["thingsilike"] = template.Must(template.ParseFS(templateFS,
		"templates/layout.html", "templates/thingsilike.html",
	))
	pages["skills"] = template.Must(template.ParseFS(templateFS,
		"templates/layout.html", "templates/skills.html",
	))


	// Serve.
	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/education", handleEducation)
	http.HandleFunc("/thingsilike", handleThingsilike)
	http.HandleFunc("/skills", handleSkills)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	log.Printf("Server is running and working on http://localhost:%d\n", *port)
	log.Fatalln(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}

type templateData struct {
	Active string
	Title  string
	Data   map[string]any
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	// Calculate how's old author.
	birth := time.Date(1994, 02, 14, 0, 0, 0, 0, time.Local)
	years := int(time.Since(birth).Hours() / 24 / 365)

	// Get random job.
	data := templateData{
		Active: "index",
		Title:  "Главная",
		Data: map[string]any{
			"Job":   jobs[rand.Intn(len(jobs))],
			"Years": years,
		},
	}

	render(w, "index", &data)
}

func handleEducation(w http.ResponseWriter, r *http.Request) {
	data := templateData{
		Active: "education",
		Title:  "Образование",
	}

	render(w, "education", &data)
}

func handleThingsilike(w http.ResponseWriter, r *http.Request) {
	data := templateData{
		Active: "thingsilike",
		Title:  "Любимые вкусы",
	}

	render(w, "thingsilike", &data)
}

func handleSkills(w http.ResponseWriter, r *http.Request) {
	data := templateData{
		Active: "skills",
		Title:  "Навыки",
	}

	render(w, "skills", &data)
}

func render(w http.ResponseWriter, name string, data *templateData) {
	p, ok := pages[name]
	if !ok {
		http.Error(w, "Page not found", http.StatusNotFound)
		return
	}
	if err := p.ExecuteTemplate(w, "layout", data); err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
	}
}
