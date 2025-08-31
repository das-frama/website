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
	"time"
)

var (
	port = flag.Int("port", 8000, "Порт для запуска сервера")
	data = flag.String("data", "data", "Каталог с данными")
	rpid = flag.String("rpid", "localhost", "Разделительный символ для портов")
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

var pages = map[string]Page{
	"index":       {"GET /{$}", "Главная", handleIndex, []string{"index.html", "ascii.html"}},
	"education":   {"GET /education", "Образование ", nil, []string{"education.html"}},
	"thingsilike": {"GET /thingsilike", "Любимые вкусы", nil, []string{"thingsilike.html"}},
	"skills":      {"GET /skills", "Умения", nil, []string{"skills.html"}},
	"blog":        {"GET /blog", "Дневник", handleBlog, []string{"blog.html"}},
	"blog.detail": {"GET /blog/{slug}", "Дневник", handleBlogDetail, []string{"blog.detail.html"}},
}

var jobs = []string{
	"подрабатываю курьером – развожу еду и продукты всем подряд. Коплю на мопед.",
	"работаю в «Горэлектротрансе» – сооставляю маршруты, сверяюсь с расписанием, слежу чтоб всё шло по рельсам.",
	"работаю звуковым инженером в ДК – настраиваю микшер, ставлю свет, мечтаю о Галь Гадот.",
	"работаю библиотекарем в спальном районе Петербурга – в основном ничего не делаю, часами залипаю в ютуб шортсы.",
	"работаю техническим переводчиком с немецкого в любимом издательстве – заказчик присылает текст, я перевожу. И так – снова и снова, пока всем всё не понравится.",
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
	"стажируюсь техником в ростелекоме, устанавливаю роутеры, прокидываю витую пару, белооранжевый-оранжевый-белозелёный-синий-белосиний-зелёный-белокоричневый-коричневый.",
	"верстаю афиши для цирковых представлений. В основном, просто переиспользую шаблон с разноцветными шариками.",
	"слежу за цветами в крупном строительном магазине. Пересаживаю растения из одного горшка в другой.",
	"работаю швеёй в ИК-14. Делаем носки, футболки-маечки, хлопчатые трусы и роскошные шарфы.",
}

var secret string

func init() {
	flag.Parse()
	if *port == 0 {
		log.Fatalln("Порт не может быть нулевым")
	}
	if *data == "" {
		log.Fatalln("Каталог с данными не может быть пустым")
	}
	if *rpid == "" {
		log.Fatalln("Разделительный символ для портов не может быть пустым")
	}
}

func main() {
	// Пытаемся дешифровать файл с мастер-паролем.
	var err error
	secret, err = decodeSecretPassword(".key.gpg")
	if err != nil {
		log.Fatalln(err)
	}

	// Покажем какой rpid
	fmt.Println("RPID:", *rpid)

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

	return out.String(), nil
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
	posts, err := listPosts(r.Context())
	if err != nil {
		return nil
	}

	return map[string]any{
		"Posts": posts,
	}
}

func handleBlogDetail(r *http.Request) map[string]any {
	slug := r.PathValue("slug")
	post, err := getPostBySlug(r.Context(), slug)
	if err != nil {
		log.Printf("Error getting post by slug: %v", err)
		return nil
	}
	if !post.Active {
		return nil
	}

	return map[string]any{
		"Post": post,
	}
}

func render(w http.ResponseWriter, tmpl *template.Template, data *TemplateData) error {
	if err := tmpl.ExecuteTemplate(w, "layout", data); err != nil {
		http.Error(w, fmt.Sprintf("Template error: %v", err), http.StatusInternalServerError)
		return err
	}
	return nil
}
