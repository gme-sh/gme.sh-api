package tpl

import (
	"github.com/gofiber/fiber/v2"
	"log"
	"regexp"
	"strings"
)

type Template struct {
	TemplateURL string `json:"template_url" bson:"template_url"`
	FullURL     string `json:"full_url" bson:"full_url"`
	params      []string
}

const paramPattern = `\:([A-Za-z]+)`

var paramRegex *regexp.Regexp

func init() {
	var err error
	paramRegex, err = regexp.Compile(paramPattern)
	if err != nil {
		log.Fatalln("error compiling param regex:", err)
	}
}

func NewTemplate(templateURL, fullURL string) (t *Template) {
	t = &Template{
		TemplateURL: templateURL,
		FullURL:     fullURL,
	}
	t.Params()
	t.Check()
	return
}

func (t *Template) Params() (res []string) {
	if t.params != nil && len(t.params) > 0 {
		log.Println("from cache:", t.params)
		return t.params
	}
	p := paramRegex.FindAllString(t.TemplateURL, -1)
	for _, s := range p {
		if !strings.HasPrefix(s, ":") {
			log.Println("wth?:", s)
			continue
		}
		res = append(res, s[1:])
	}
	log.Println("not from cache:", res)
	t.params = res
	return
}

func (t *Template) Check() (val bool) {
	val = true
	for _, p := range t.Params() {
		if !strings.Contains(t.FullURL, ":"+p) {
			log.Println("WARN :: Parameter", p, "not in full url", t.FullURL)
			val = false
		}
	}
	return
}

func (t *Template) Register(app *fiber.App) {
	log.Println("TPL :: Registering", t.TemplateURL)
	app.Get(t.TemplateURL, func(ctx *fiber.Ctx) error {
		url := t.FullURL
		for _, p := range t.Params() {
			url = strings.ReplaceAll(url, ":"+p, ctx.Params(p))
		}
		log.Println("Redirecting Template", t.TemplateURL, "to", url)
		return ctx.Redirect(url)
	})
}
