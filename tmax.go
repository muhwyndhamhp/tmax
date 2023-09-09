package tmax

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type Renderer struct {
	Templates map[string]*template.Template
	RootName  string
	ViewName  string
}

var rg = regexp.MustCompile(`\w+\.html$`)

const (
	separator = "#"
	fileExt   = ".html"
)

func (t *Renderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	str := strings.Split(name, separator)
	if str[0] == t.RootName {
		return t.Templates[str[1]].ExecuteTemplate(w, t.RootName, data)
	} else {
		return t.Templates[name].ExecuteTemplate(w, t.ViewName, data)
	}
}

func NewEchoTemplateRenderer(e *echo.Echo, rootName, viewName, componentPath string, viewPaths ...string) {
	tmpL := map[string]*template.Template{}
	comPaths := extractComponents(componentPath)
	fmt.Println(comPaths)
	for i := range viewPaths {
		c, err := os.ReadDir(viewPaths[i])
		if err != nil {
			log.Error(err)
		}

		for _, entry := range c {
			if !rg.MatchString(entry.Name()) {
				continue
			}
			tmp := &template.Template{}
			comPaths = append(comPaths, filepath.Join(viewPaths[i], entry.Name()))

			template.Must(tmp.ParseFiles(comPaths...))

			name := strings.ReplaceAll(entry.Name(), fileExt, "")
			tmpL[name] = tmp
		}
	}

	t := newTemplate(tmpL, rootName, viewName)
	e.Renderer = t
}

func extractComponents(componentPath string) []string {
	comPaths := []string{}
	var appendComps func(path string)
	appendComps = func(path string) {
		comps, err := os.ReadDir(path)
		if err != nil {
			log.Error(err)
		}
		for i := range comps {
			if comps[i].IsDir() {
				appendComps(filepath.Join(path, comps[i].Name()))
			}
			if rg.MatchString(comps[i].Name()) {
				comPaths = append(comPaths, filepath.Join(path, comps[i].Name()))
			}
		}
	}
	appendComps(componentPath)
	return comPaths
}

func newTemplate(templates map[string]*template.Template, rootName, viewName string) echo.Renderer {
	return &Renderer{
		Templates: templates,
		RootName:  rootName,
		ViewName:  viewName,
	}
}
