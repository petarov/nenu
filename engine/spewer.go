package engine

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/parser"
	"github.com/petarov/nenu/config"
)

var (
	layout     = "2006-01-02 15:04:05"
	extensions = parser.CommonExtensions | parser.AutoHeadingIDs
)

type post struct {
	filename string
	date     time.Time
	title    string
	subtitle string
	publish  bool
}

func writeArchives(dest *os.File) {
	// TODO
}

func writePost(dest *os.File, post *post, data []byte) (err error) {
	lines := make([]string, 0)
	parsing := false

	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	for scanner.Scan() {
		line := scanner.Text()
		switch {
		case strings.HasPrefix(line, "---"):
			parsing = !parsing
		case strings.HasPrefix(line, "title:"):
			post.title = strings.TrimSpace(line[7:])
			if len(post.title) == 0 {
				return errors.New("post title cannot be empty")
			}
		case strings.HasPrefix(line, "subtitle:"):
			post.subtitle = strings.TrimSpace(line[10:])
		case strings.HasPrefix(line, "date:"):
			post.date, err = time.Parse(time.RFC3339, line[6:])
			if err != nil {
				return
			}
		case strings.HasPrefix(line, "publish:"):
			post.publish, err = strconv.ParseBool(line[9:])
			if err != nil {
				return
			}
		default:
			if !parsing {
				lines = append(lines, scanner.Text())
			}
		}
	}

	if post.publish {
		parser := parser.NewWithExtensions(extensions)
		md := markdown.ToHTML([]byte(strings.Join(lines, "")), parser, nil)

		fmt.Println("MD---", len(md))
		// TODO
		// f, err := os.Create("/tmp/dat2")
		// if err != nil {
		// 	return
		// }

		// defer f.Close()
	} else {
		fmt.Printf("| Skipped %s (publish = false)\n", post.filename)
	}

	return
}

func writePosts(dest *os.File) ([]*post, error) {
	path := config.PostsPath
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

			post := new(post)
			post.filename = file.Name()
			post.title = strings.ReplaceAll(post.filename[11:strings.LastIndex(post.filename, ".")], "-", " ")
			post.date, err = time.Parse(layout, post.filename[:10]+" 00:00:00")
			post.publish = true
			if err != nil {
				return nil, err
			}
			posts = append(posts, post)
		}
	}

	fmt.Println("| Generating posts...")

	for _, post := range posts {
		fmt.Println("|--> ", post.filename)

		data, err := ioutil.ReadFile(filepath.Join(path, post.filename))
		if err != nil {
			return nil, err
		}

		if err = writePost(dest, post, data); err != nil {
			return nil, err
		}
	}

	sort.Slice(posts, func(a, b int) bool {
		return posts[a].date.After(posts[b].date)
	})

	return posts, nil
}

// Spew generates website
func Spew() (err error) {
	tempDir, err := ioutil.TempFile(config.TempPath, "nenu-gen-")
	if err != nil {
		return
	}
	defer os.Remove(tempDir.Name())

	_, err = writePosts(tempDir)
	if err != nil {
		return err
	}

	return nil
}
