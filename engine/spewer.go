package engine

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/gomarkdown/markdown/parser"
	"github.com/otiai10/copy"
	"github.com/petarov/nenu/config"
)

var (
	layout     = "2006-01-02 15:04:05"
	extensions = parser.CommonExtensions | parser.AutoHeadingIDs
)

// Templates website html templates
type Templates struct {
	Header  *template.Template
	Footer  *template.Template
	Post    *template.Template
	Archive *template.Template
}

func loadTemplates() *Templates {
	templates := new(Templates)
	templates.Header = template.Must(template.ParseFiles(filepath.Join(config.TemplatePath, "header.html")))
	templates.Footer = template.Must(template.ParseFiles(filepath.Join(config.TemplatePath, "footer.html")))
	templates.Post = template.Must(template.ParseFiles(filepath.Join(config.TemplatePath, "post.html")))
	templates.Archive = template.Must(template.ParseFiles(filepath.Join(config.TemplatePath, "archive.html")))
	return templates
}

// Spew generates website
func Spew() (err error) {
	tempDir, err := ioutil.TempDir(config.TempPath, "nenu-gen-")
	if err != nil {
		return
	}
	defer os.RemoveAll(tempDir)

	config.TempPath, err = filepath.Abs(tempDir)
	if err != nil {
		return err
	}
	fmt.Printf("| Using temp dir: %s\n", config.TempPath)

	templates := loadTemplates()

	meta, err := SpewPosts(templates)
	if err != nil {
		return err
	}

	err = SpewArchive(meta, templates)
	if err != nil {
		return err
	}

	config.OutputPath, err = filepath.Abs(config.OutputPath)
	if err != nil {
		return err
	}

	err = copy.Copy(config.TempPath, config.OutputPath)
	if err != nil {
		return err
	}

	return nil
}
