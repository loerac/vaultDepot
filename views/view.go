package views

import (
    "bytes"
    "html/template"
    "io"
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
    view.Render(writer, nil)
}

/**
 * @brief:  Render the view with the predined layout
 *
 * @param:  writer - Render the template
 * @param:  data - Info that is being rendered
 **/
func (view *View) Render(writer http.ResponseWriter, data interface{}) {
    writer.Header().Set("Content-Type", "text/html")

    switch data.(type) {
    case Data:
        // Do nothing
    default:
        data = Data {
            Yield: data,
        }
    }

    var buf bytes.Buffer
    err := view.Template.ExecuteTemplate(&buf, view.Layout, data)
    if err != nil {
        http.Error(writer,
          "OOPSIE WOOPSIE! uwu There's been a fucko wucko.... please email support@vaultdepot.com",
          http.StatusInternalServerError)
        return
    }
    io.Copy(writer, &buf)
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
