package engine

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/parser"

	"github.com/petarov/nenu/config"
)

func loadFiles(path string, yml *config.YML) ([]string, error) {
	fmt.Printf("| Loading contents from %s...\n", path)

	extensions := parser.CommonExtensions | parser.AutoHeadingIDs

	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	articles := make([]string, len(files))

	for _, file := range files {
		ext := filepath.Ext(file.Name())
		if !file.IsDir() && (ext == ".md" || ext == ".markdown") {
			fmt.Println("|--> ", file.Name())
			md, err := ioutil.ReadFile(filepath.Join(path, file.Name()))
			if err != nil {
				return nil, err
			}
			parser := parser.NewWithExtensions(extensions)
			parsed := markdown.ToHTML(md, parser, nil)
			articles = append(articles, string(parsed))
		}
	}

	return articles, nil
}

func Spew(yml *config.YML) (err error) {
	articles, err := loadFiles(config.ArticlesPath, yml)
	if err != nil {
		return err
	}

	fmt.Printf("%v", articles)

	return nil
}
