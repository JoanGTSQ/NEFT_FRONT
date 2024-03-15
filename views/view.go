package views

import (
	"bytes"
	"html/template"
	"io"
	"net/http"
	"path/filepath"

	"github.com/gorilla/csrf"
	"jgt.solutions/context"
	"jgt.solutions/logController"
)

var (
	LayoutDir   string = "views/layouts/"
	TemplateDir string = "views/"
	TemplateExt string = ".gohtml"
)

func NewView(layout string, files ...string) *View {

	addTemplatePath(files)
	addTemplateExt(files)
	files = append(files, layoutFiles()...)
	t, err := template.New("").Funcs(funcMap).ParseFiles(files...)
	if err != nil {
		logController.ErrorLogger.Println(err)
		return nil
	}

	return &View{
		Template: t,
		Layout:   layout,
	}
}

type View struct {
	Template *template.Template
	Layout   string
}

func (v *View) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	v.Render(w, r, nil)
}

func (v *View) Render(w http.ResponseWriter, r *http.Request, data interface{}) {
	w.Header().Set("Content-type", "text/html")
	var vd Data

	switch d := data.(type) {

	case *Data:
		vd = *d
	default:
		vd.Yield = d
	}
	vd.Active = r.URL.Path
	vd.User = context.User(r.Context())
	var buf bytes.Buffer
	csrfField := csrf.TemplateField(r)
	tpl := v.Template.Funcs(template.FuncMap{
		"csrfField": func() template.HTML {
			return csrfField
		},
	})

	if err := tpl.ExecuteTemplate(&buf, v.Layout, vd); err != nil {
		logController.ErrorLogger.Println(err)
		http.Redirect(w, r, "/505", http.StatusFound)
		return
	}
	io.Copy(w, &buf)
}

func (v *View) Flush(w http.ResponseWriter, r *http.Request, data interface{}) {
	if r.RequestURI == "/" {
		r.RequestURI = "/#contact"
	}
	w.Header().Set("Content-type", "text/html")
	var vd Data

	switch d := data.(type) {

	case *Data:
		vd = *d
	default:
		vd.Yield = d
	}
	vd.Active = r.URL.Path
	vd.User = context.User(r.Context())
	var buf bytes.Buffer
	csrfField := csrf.TemplateField(r)
	tpl := v.Template.Funcs(template.FuncMap{
		"csrfField": func() template.HTML {
			return csrfField
		},
	})

	if err := tpl.ExecuteTemplate(&buf, v.Layout, vd); err != nil {
		logController.ErrorLogger.Println(err)
		http.Redirect(w, r, "/505", http.StatusFound)
		return
	}
	io.Copy(w, &buf)
	w.(http.Flusher).Flush()
}
func layoutFiles() []string {
	files, err := filepath.Glob(LayoutDir + "*" + TemplateExt)
	if err != nil {
		logController.ErrorLogger.Println(err)
		return nil
	}
	return files
}

func addTemplatePath(files []string) {
	for i, f := range files {
		files[i] = TemplateDir + f
	}
}
func addTemplateExt(files []string) {
	for i, f := range files {
		files[i] = f + TemplateExt
	}
}
