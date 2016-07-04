package web

import (
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/gzip"
	"github.com/martini-contrib/render"

	"../logging"
	"./routes"
)

var errorChan chan error

func Initialize(lock chan error, logger *logging.Router, port int, directory string) {
	errorChan = lock
	logger.Banner("Starting Server")

	m := martini.Classic()
	m.Map(logger.Web())

	RegisterMiddleware(m, directory)
	RegisterDefaults(m)

	go func() {
		lock <- http.ListenAndServe(":"+strconv.Itoa(port), m)
	}()
}

func RegisterMiddleware(m *martini.ClassicMartini, directory string) {

	m.Use(gzip.All())

	m.Use(render.Renderer(render.Options{
		Directory:       filepath.Join(directory, "views"),
		Extensions:      []string{".tmpl", ".html"},
		IndentJSON:      true,
		HTMLContentType: "text/html",
		Delims:          render.Delims{"(%", "%)"},
	}))

	m.Use(martini.Static(filepath.Join(directory, "css"), martini.StaticOptions{
		Prefix: "/css",
	}))

	m.Use(martini.Static(filepath.Join(directory, "js"), martini.StaticOptions{
		Prefix: "/js",
	}))
}

func RegisterDefaults(m *martini.ClassicMartini) {
	m.Get("/", routes.Home)
}
