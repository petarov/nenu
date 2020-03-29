package engine

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/parser"

	"github.com/petarov/nenu/config"
)

func ContentsLoader(contentPath string, yml *config.YML) ([]string, error) {
	fmt.Printf("| Loading contents from %s...\n", contentPath)

	extensions := parser.CommonExtensions | parser.AutoHeadingIDs

	files, err := ioutil.ReadDir(contentPath)
	if err != nil {
		return nil, err
	}

	articles := make([]string, len(files))

	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".md" {
			fmt.Println("|--> ", file.Name())
			md, err := ioutil.ReadFile(filepath.Join(contentPath, file.Name()))
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
