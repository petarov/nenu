package engine

import (
	"bufio"
	"errors"
	"fmt"
	"html/template"
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

type postPageData struct {
	*config.YML
	Post post
}

// PostMeta post meta data used for indexing
type PostMeta struct {
	Title     string
	Date      string
	Permalink string
}

type post struct {
	filename     string
	filenameHTML string
	filepath     string
	date         time.Time
	Meta         *PostMeta
	// Title    string
	// Date     string
	Subtitle string
	Content  template.HTML
	// Permalink string
	Summary  string
	ImageURL string
	Prev     *post
}

func writePost(dest *os.File, post *post, data []byte, templates *Templates) (err error) {
	var (
		lines   = make([]string, 0)
		parsing = false
		publish = true
	)

	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	for scanner.Scan() {
		line := scanner.Text()
		switch {
		case strings.HasPrefix(line, "---"):
			parsing = !parsing
		case strings.HasPrefix(line, "title:"):
			post.Meta.Title = strings.TrimSpace(line[7:])
			if len(post.Meta.Title) == 0 {
				return errors.New("post title cannot be empty")
			}
		case strings.HasPrefix(line, "subtitle:"):
			post.Subtitle = strings.TrimSpace(line[10:])
		case strings.HasPrefix(line, "date:"):
			post.date, err = time.Parse(time.RFC3339, line[6:])
			if err != nil {
				return
			}
		case strings.HasPrefix(line, "publish:"):
			publish, err = strconv.ParseBool(line[9:])
			if err != nil {
				return
			}
		default:
			if !parsing {
				lines = append(lines, scanner.Text())
			}
		}
	}

	if publish {
		parser := parser.NewWithExtensions(extensions)
		post.Content = template.HTML(markdown.ToHTML([]byte(strings.Join(lines, "")), parser, nil))
		post.Meta.Date = post.date.In(config.TimeZoneLocation).Local().Format("Mon, 02 Jan 2006")

		dirpath := filepath.Join(config.TempPath, post.filepath)
		if err = os.MkdirAll(dirpath, 0755); err != nil {
			return
		}

		var f *os.File
		f, err = os.Create(filepath.Join(dirpath, post.filenameHTML))
		if err != nil {
			return
		}
		defer f.Close()

		pd := &postPageData{config.YMLConfig, *post}

		// header
		if err = templates.Header.Execute(f, pd); err != nil {
			return
		}
		// body
		if err = templates.Post.Execute(f, pd); err != nil {
			return
		}
		// footer
		if err = templates.Footer.Execute(f, pd); err != nil {
			return
		}
	} else {
		fmt.Printf("| Skipped %s (publish = false)\n", post.filename)
	}

	return
}

// SpewPosts generate all website
func SpewPosts(dest *os.File, templates *Templates) ([]*PostMeta, error) {
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

			fileNoExt := strings.ReplaceAll(file.Name(), filepath.Ext(file.Name()), "")

			post := new(post)
			post.filename = file.Name()
			post.filenameHTML = fmt.Sprintf("%s.html", fileNoExt)
			dirs := strings.Split(post.filename[:10], "-")
			if len(dirs) != 3 {
				return nil, fmt.Errorf("%s unexpected date prefix. YYYY-mm-dd required", post.filename)
			}
			post.filepath = filepath.Join(dirs[0], dirs[1], dirs[2])
			post.date, err = time.Parse(layout, post.filename[:10]+" 00:00:00")
			if err != nil {
				return nil, err
			}

			post.Meta = new(PostMeta)
			post.Meta.Title = strings.ReplaceAll(fileNoExt, "-", " ")

			posts = append(posts, post)
		}
	}

	fmt.Println("| Generating posts...")

	index := make([]*PostMeta, 0, len(posts))

	for _, post := range posts {
		fmt.Println("|--> ", post.filename)

		index = append(index, post.Meta)

		data, err := ioutil.ReadFile(filepath.Join(path, post.filename))
		if err != nil {
			return nil, err
		}

		if err = writePost(dest, post, data, templates); err != nil {
			return nil, err
		}
	}

	sort.Slice(posts, func(a, b int) bool {
		return posts[a].date.After(posts[b].date)
	})

	return index, nil
}
