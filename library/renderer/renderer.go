package renderer

import (
	"errors"
	"html/template"
	"net/url"
	"strings"
	"time"

	"github.com/andycai/werite/core"
	"github.com/andycai/werite/library/authentication"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/gomarkdown/markdown"
	mdhtml "github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

func ViewEngineStart() *html.Engine {

	viewEngine := html.New("./templates", ".html")

	viewEngine.AddFunc("IsAuthenticated", func(c *fiber.Ctx) bool {
		isAuthenticated, _ := authentication.AuthGet(c)
		return isAuthenticated
	})

	viewEngine.AddFunc("Iterate", func(start int, end int) []int {
		n := end - start + 1
		result := make([]int, n)
		for i := 0; i < n; i++ {
			result[i] = start + i
		}
		return result
	})

	viewEngine.AddFunc("Dict", func(values ...interface{}) (map[string]interface{}, error) {
		if len(values)%2 != 0 {
			return nil, errors.New("invalid dict call")
		}
		dict := make(map[string]interface{}, len(values)/2)
		for i := 0; i < len(values); i += 2 {
			key, ok := values[i].(string)
			if !ok {
				return nil, errors.New("dict keys must be strings")
			}
			dict[key] = values[i+1]
		}
		return dict, nil
	})

	viewEngine.AddFunc("Lang", func(a []string, sep string) string {
		return core.Lang()
	})

	viewEngine.AddFunc("Join", func(a []string, sep string) string {
		return strings.Join(a, sep)
	})

	viewEngine.AddFunc("DateFormat", func(t time.Time, layout string) string {
		return core.DateFormat(t, layout)
	})

	viewEngine.AddFunc("GetErrors", func() []string {
		return core.GetErrors()
	})

	viewEngine.AddFunc("GetMessages", func() []string {
		return core.GetMessages()
	})

	viewEngine.AddFunc("QueryUnescape", func(s string) string {
		query, err := url.QueryUnescape(s)
		if err != nil {
			return s
		}
		return query
	})

	viewEngine.AddFunc("IsNotZero", func(t time.Time) bool {
		return !t.IsZero()
	})

	viewEngine.AddFunc("Str2HTML", func(s string) template.HTML {
		return template.HTML(s)
	})

	viewEngine.AddFunc("MD2HTML", func(s string) template.HTML {
		// create markdown parser with extensions
		extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
		p := parser.NewWithExtensions(extensions)
		doc := p.Parse([]byte(s))

		// create HTML renderer with extensions
		htmlFlags := mdhtml.CommonFlags | mdhtml.HrefTargetBlank
		opts := mdhtml.RendererOptions{Flags: htmlFlags}
		renderer := mdhtml.NewRenderer(opts)

		return template.HTML(string(markdown.Render(doc, renderer)))

	})

	return viewEngine
}
