package views

import (
    "html/template"
    "net/http"
    "path/filepath"
)

var (
    TemplateDir string = "views/"
    LayoutDir   string = TemplateDir + "layouts/"
    TemplateExt string = ".html"
)

type View struct {
    Template    *template.Template
    Layout      string
}

func layoutFiles() []string {
    files, err := filepath.Glob(LayoutDir + "*" + TemplateExt)
    if err != nil {
        panic(err)
    }

    return files
}

func NewView(layout string, files ...string) *View {
    addTemplatePath(files)
    addTemplateExt(files)
    files = append(files, layoutFiles()...)
    tmpl, err := template.ParseFiles(files...)
    if err != nil {
        panic(err)
    }

    return &View{
        Template: tmpl,
        Layout: layout,
    }
}

func (view *View) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
    if err := view.Render(writer, nil); err != nil {
        panic(err)
    }
}

func (view *View) Render(writer http.ResponseWriter, data interface{}) error {
    writer.Header().Set("Content-Type", "text/html")
    return view.Template.ExecuteTemplate(writer, view.Layout, data)
}

/**
 * @brief:  Prepends the template directory
 *          to each string in the slice
 *
 * @param:  files - File paths for templates
 **/
func addTemplatePath(files []string) {
    for i, file := range files {
        files[i] = TemplateDir + file
    }
}

/**
 * @brief:  Appends the template directory
 *          to each string in the slice
 *
 * @param:  files - File paths for templates
 **/
func addTemplateExt(files []string) {
    for i, file := range files {
        files[i] = file + TemplateExt
    }
}
