package main

import (
	"bytes"
	"embed"
	"flag"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"slices"
	"time"
)

var (
	port = flag.Int("port", 8000, "specify port number")
	data = flag.String("data", "data", "data's root path")
)

//go:embed templates
var templateFS embed.FS

type Page struct {
	Path      string
	Title     string
	Handler   func(r *http.Request) map[string]any
	Templates []string
}

type TemplateData struct {
	Active string
	Title  string
	Data   map[string]any
}

type Post struct {
	Title     string
	Slug      string
	CreatedAt time.Time
}

var pages = map[string]Page{
	"index":       {"/", "Главная", handleIndex, []string{"index.html", "ascii.html"}},
	"education":   {"/education", "Образование ", nil, []string{"education.html"}},
	"thingsilike": {"/thingsilike", "Любимые вкусы", nil, []string{"thingsilike.html"}},
	"skills":      {"/skills", "Умения", nil, []string{"skills.html"}},
	"blog":        {"/blog", "Дневник", handleBlog, []string{"blog.html"}},
	"blog.detail": {"/blog/{slug}", "Дневник", handleBlogDetail, []string{"blog.detail.html"}},
}

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

var posts = []Post{
	{Title: "Работаю разнорабочим на линии упаковки", Slug: "work-as-packager", CreatedAt: time.Now()},
	{Title: "Работаю грузчиком в мебельном центре", Slug: "work-as-freight", CreatedAt: time.Now()},
	{Title: "Работаю слесарем-сантехником", Slug: "work-as-plumber", CreatedAt: time.Now().Add(-12 * time.Hour * 24)},
	{Title: "Работаю электриком в «Жилкомсервисе»", Slug: "work-as-electrician", CreatedAt: time.Now()},
}

var secret string

func main() {
	// Пытаемся дешифровать файл с мастер-паролем.
	var err error
	secret, err = decodeSecretPassword(".key.gpg")
	if err != nil {
		log.Fatalln(err)
	}

	flag.Parse()

	// Инициализация БД.
	if err := initDB(fmt.Sprintf("%s/database.db", *data)); err != nil {
		log.Fatal(err)
	}

	// Регистрация всех страниц.
	for name, page := range pages {
		tt := make([]string, 0, len(page.Templates)+1)
		tt = append(tt, "templates/layout.html")
		for _, t := range page.Templates {
			tt = append(tt, fmt.Sprintf("templates/%s", t))
		}
		tmpl := template.Must(template.ParseFS(templateFS, tt...))

		http.HandleFunc(page.Path, func(name string, page Page) http.HandlerFunc {
			return func(w http.ResponseWriter, r *http.Request) {
				data := TemplateData{name, page.Title, nil}
				if page.Handler != nil {
					data.Data = page.Handler(r)
				}
				render(w, tmpl, &data)
			}
		}(name, page))
	}

	// Регистрация админстраторской части.
	registerAdminRoutes()

	// Статический контент.
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	log.Printf("Server is running and working on http://localhost:%d\n", *port)
	log.Fatalln(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}

func decodeSecretPassword(filepath string) (string, error) {
	cmd := exec.Command("gpg", "--quiet", "--decrypt", filepath)

	var out bytes.Buffer
	cmd.Stdin = os.Stdin
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return "", err
	}

	return string(out.Bytes()), nil
}

func handleIndex(r *http.Request) map[string]any {
	// Calculate how's old author.
	birth := time.Date(1994, 02, 14, 0, 0, 0, 0, time.Local)
	years := int(time.Since(birth).Hours() / 24 / 365)

	return map[string]any{
		"Job":   jobs[rand.Intn(len(jobs))],
		"Years": years,
	}
}

func handleBlog(r *http.Request) map[string]any {
	slices.SortFunc(posts, func(a Post, b Post) int {
		return b.CreatedAt.Compare(a.CreatedAt)
	})
	return map[string]any{
		"Posts": posts,
	}
}

func handleBlogDetail(r *http.Request) map[string]any {
	slug := r.PathValue("slug")
	i := slices.IndexFunc(posts, func(p Post) bool {
		return p.Slug == slug
	})
	if i == -1 {
		return nil
	}

	return map[string]any{
		"Post": posts[i],
	}
}

func render(w http.ResponseWriter, tmpl *template.Template, data *TemplateData) error {
	if err := tmpl.ExecuteTemplate(w, "layout", data); err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return err
	}
	return nil
}
