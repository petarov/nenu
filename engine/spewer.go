package engine

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/gomarkdown/markdown/parser"
	"github.com/otiai10/copy"
	"github.com/petarov/nenu/config"
)

var (
	layout        = "2006-01-02 15:04:05"
	extensions    = parser.CommonExtensions | parser.AutoHeadingIDs
	templateNames = []string{"header.html", "footer.html", "post.html", "archive.html"}
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

	// generate posts
	posts, err := SpewPosts(templates)
	if err != nil {
		return err
	}

	doneArch := make(chan error)
	doneRss := make(chan error)

	// generate archive
	go func() {
		doneArch <- SpewArchive(posts, templates)
	}()

	// generate RSS feed
	go func() {
		if config.YMLConfig.Content.Rss {
			doneRss <- SpewAtom(posts)
		} else {
			doneRss <- nil
		}
	}()

	errArch := <-doneArch
	if errArch != nil {
		return errArch
	}

	errRss := <-doneRss
	if errRss != nil {
		return errRss
	}

	// copy all generated content to the specified destination
	config.OutputPath, err = filepath.Abs(config.OutputPath)
	if err != nil {
		return err
	}

	err = copy.Copy(config.TempPath, config.OutputPath)
	if err != nil {
		return err
	}

	// copy web template resources
	err = copy.Copy(config.TemplatePath, config.OutputPath, copy.Options{
		Skip: func(src string) (bool, error) {
			for _, v := range templateNames {
				if strings.HasSuffix(src, v) {
					return true, nil
				}
			}
			return false, nil
		},
	})
	if err != nil {
		return err
	}

	return nil
}
