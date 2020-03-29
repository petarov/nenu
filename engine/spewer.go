package engine

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sort"
	"time"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/parser"

	"github.com/petarov/nenu/config"
)

const layout = "2006-01-02 15:04:05"

type article struct {
	filename string
	date     time.Time
	title    string
	subtitle string
}

func writeHeader() {
	// TODO
}

func writeArchives() {
	// TODO
}

func writeArticle(art *article, html []byte) {
	// TODO
}

func writeArticles() ([]*article, error) {
	path := config.ArticlesPath
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs

	fmt.Printf("| Indexing contents from %s...\n", path)

	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	articles := make([]*article, 0, len(files))

	for _, file := range files {
		ext := filepath.Ext(file.Name())
		if !file.IsDir() && (ext == ".md" || ext == ".markdown") {
			fmt.Println("|--> ", file.Name())

			art := new(article)
			art.filename = file.Name()
			art.date, err = time.Parse(layout, art.filename[:10]+" 12:00:00")
			if err != nil {
				return nil, err
			}
			articles = append(articles, art)
		}
	}

	fmt.Println("| Generating contents...")

	sort.Slice(articles, func(a, b int) bool {
		return articles[a].date.After(articles[b].date)
	})

	for _, art := range articles {
		fmt.Println("|--> ", art.filename)

		md, err := ioutil.ReadFile(filepath.Join(path, art.filename))
		if err != nil {
			return nil, err
		}

		parser := parser.NewWithExtensions(extensions)
		data := markdown.ToHTML(md, parser, nil)

		writeArticle(art, data)
	}

	return articles, nil
}

func Spew() (err error) {
	_, err = writeArticles()
	if err != nil {
		return err
	}

	return nil
}
