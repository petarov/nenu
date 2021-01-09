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
	"github.com/goodsign/monday"
	"github.com/petarov/nenu/config"
)

type postPageData struct {
	*config.YML
	Post Post
}

// Post post meta data
type Post struct {
	filename     string
	filenameHTML string
	filepath     string
	date         time.Time
	Title        string
	Subtitle     string
	Summary      string
	Date         string
	ShortDate    string
	Content      template.HTML
	Permalink    template.URL
	PermalinkURL string
	ImageURL     string
	IsPublish    bool
	IsIndex      bool
	Prev         *Post
}

var (
	markdownExt    = []string{".md", ".markdown", ".mkd", ".mdown", ".mdwn"}
	htmlDateLayout = "Monday, 02 January 2006"
)

func isExtOk(needle string) bool {
	for _, ext := range markdownExt {
		if ext == needle {
			return true
		}
	}
	return false
}

func writePost(post *Post, data []byte, templates *Templates) (err error) {
	var (
		lines   = make([]string, 0)
		parsing = false
	)

	// published unless specified otherwise
	post.IsPublish = true

	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	for scanner.Scan() {
		line := scanner.Text()
		switch {
		case strings.HasPrefix(line, "---"):
			parsing = !parsing
		case strings.HasPrefix(line, "title:"):
			post.Title = strings.TrimSpace(line[7:])
			if len(post.Title) == 0 {
				return errors.New("post title cannot be empty")
			}
		case strings.HasPrefix(line, "subtitle:"):
			post.Subtitle = strings.TrimSpace(line[10:])
		case strings.HasPrefix(line, "summary:"):
			post.Summary = strings.TrimSpace(line[8:])
			// case strings.HasPrefix(line, "date:"):
		// 	post.date, err = time.Parse(time.RFC3339, line[6:])
		// 	if err != nil {
		// 		return
		// 	}
		case strings.HasPrefix(line, "publish:"):
			post.IsPublish, err = strconv.ParseBool(line[9:])
			if err != nil {
				return
			}
		default:
			if !parsing {
				if len(line) > 0 {
					lines = append(lines, line)
					lines = append(lines, "\n")
				} else {
					lines = append(lines, "\n\n")
				}
			}
		}
	}

	if post.IsPublish {
		parser := parser.NewWithExtensions(extensions)
		post.Content = template.HTML(markdown.ToHTML([]byte(strings.Join(lines, "")), parser, nil))

		locale := monday.Locale(config.YMLConfig.Content.Locale)
		post.Date = monday.Format(post.date.In(config.TimeZoneLocation).Local(), htmlDateLayout, locale)
		post.ShortDate = monday.Format(post.date.In(config.TimeZoneLocation).Local(),
			monday.MediumFormatsByLocale[locale], locale)

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

// SpewPosts generate site posts
func SpewPosts(templates *Templates) ([]*Post, error) {
	path := config.PostsPath
	fmt.Printf("| Indexing posts from %s...\n", path)

	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	posts := make([]*Post, 0, len(files))

	for _, file := range files {
		ext := filepath.Ext(file.Name())

		if !file.IsDir() && isExtOk(ext) {
			fmt.Println("|--> ", file.Name())
			filenameClean := strings.ReplaceAll(file.Name(), ext, "")[11:]

			post := new(Post)
			post.filename = file.Name()
			post.filenameHTML = fmt.Sprintf("%s.html", filenameClean)
			dirs := strings.Split(post.filename[:10], "-")
			if len(dirs) != 3 {
				return nil, fmt.Errorf("%s unexpected date prefix. YYYY-mm-dd required", post.filename)
			}
			post.filepath = filepath.Join(dirs[0], dirs[1], dirs[2])
			post.date, err = time.Parse(layout, post.filename[:10]+" 00:00:00")
			if err != nil {
				return nil, err
			}

			post.IsIndex = false
			post.PermalinkURL = fmt.Sprintf("%s/%s/%s", config.YMLConfig.Site.URL, post.filepath,
				post.filenameHTML)
			post.Permalink = template.URL(post.Permalink)
			post.Title = strings.ReplaceAll(filenameClean, "-", " ")

			posts = append(posts, post)
		}
	}

	fmt.Println("| Generating posts...")

	//  sort by ascending date
	sort.Slice(posts, func(a, b int) bool {
		return posts[a].date.Before(posts[b].date)
	})

	var prev *Post

	for _, post := range posts {
		fmt.Println("|--> ", post.filename)

		post.Prev = prev

		data, err := ioutil.ReadFile(filepath.Join(path, post.filename))
		if err != nil {
			return nil, err
		}

		if err = writePost(post, data, templates); err != nil {
			return nil, err
		}

		if post.IsPublish {
			prev = post
		}
	}

	fmt.Println("| Generating index...")

	index := posts[len(posts)-1]
	index.filenameHTML = "index.html"
	index.filepath = ""
	index.IsIndex = true

	data, err := ioutil.ReadFile(filepath.Join(path, index.filename))
	if err != nil {
		return nil, err
	}

	if err = writePost(index, data, templates); err != nil {
		return nil, err
	}

	// sort by descending date
	for i := len(posts)/2 - 1; i >= 0; i-- {
		tmp := len(posts) - 1 - i
		posts[i], posts[tmp] = posts[tmp], posts[i]
	}

	return posts, nil
}
