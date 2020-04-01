package engine

import (
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/gomarkdown/markdown/parser"
	"github.com/petarov/nenu/config"
)

var (
	layout     = "2006-01-02 15:04:05"
	extensions = parser.CommonExtensions | parser.AutoHeadingIDs
)

// Templates website html templates
type Templates struct {
	Header *template.Template
	Post   *template.Template
	Footer *template.Template
}

func loadTemplates() *Templates {
	templates := new(Templates)
	templates.Header = template.Must(template.ParseFiles(filepath.Join(config.TemplatePath, "header.html")))
	templates.Post = template.Must(template.ParseFiles(filepath.Join(config.TemplatePath, "post.html")))
	templates.Footer = template.Must(template.ParseFiles(filepath.Join(config.TemplatePath, "footer.html")))
	return templates
}

// Spew generates website
func Spew() (err error) {
	tempDir, err := ioutil.TempFile(config.TempPath, "nenu-gen-")
	if err != nil {
		return
	}
	defer os.Remove(tempDir.Name())

	templates := loadTemplates()

	_, err = SpewPosts(tempDir, templates)
	if err != nil {
		return err
	}

	return nil
}
