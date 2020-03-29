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

type post struct {
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

func writePost(art *post, html []byte) {
	// TODO
}

func writePosts() ([]*post, error) {
	path := config.PostsPath
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs

	fmt.Printf("| Indexing posts from %s...\n", path)

	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	posts := make([]*post, 0, len(files))

	for _, file := range files {
		ext := filepath.Ext(file.Name())
		if !file.IsDir() && (ext == ".md" || ext == ".markdown") {
			fmt.Println("|--> ", file.Name())

			art := new(post)
			art.filename = file.Name()
			art.date, err = time.Parse(layout, art.filename[:10]+" 12:00:00")
			if err != nil {
				return nil, err
			}
			posts = append(posts, art)
		}
	}

	fmt.Println("| Generating posts...")

	sort.Slice(posts, func(a, b int) bool {
		return posts[a].date.After(posts[b].date)
	})

	for _, art := range posts {
		fmt.Println("|--> ", art.filename)

		md, err := ioutil.ReadFile(filepath.Join(path, art.filename))
		if err != nil {
			return nil, err
		}

		parser := parser.NewWithExtensions(extensions)
		data := markdown.ToHTML(md, parser, nil)

		writePost(art, data)
	}

	return posts, nil
}

func Spew() (err error) {
	_, err = writePosts()
	if err != nil {
		return err
	}

	return nil
}
