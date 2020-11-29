package engine

import (
	"os"
	"path/filepath"
	"time"

	"github.com/gorilla/feeds"
	"github.com/petarov/nenu/config"
)

// SpewAtom generate Atom feed from posts
func SpewAtom(posts []*Post) error {
	now := time.Now()
	feed := &feeds.Feed{
		Title:       config.YMLConfig.Site.Title,
		Link:        &feeds.Link{Href: config.YMLConfig.Site.URL},
		Description: config.YMLConfig.Site.Description,
		Author:      &feeds.Author{Name: config.YMLConfig.Site.Author, Email: ""},
		Copyright:   config.YMLConfig.Footer.Copyright,
		Created:     now,
	}

	for _, p := range posts {
		feed.Items = append(feed.Items, &feeds.Item{
			Title: p.Title,
			//Image:       &feeds.Image{Url: p.ImageURL},
			Link:        &feeds.Link{Href: p.PermalinkURL},
			Description: p.Summary,
			Author:      &feeds.Author{Name: config.YMLConfig.Site.Author, Email: ""},
			Created:     p.date.UTC(),
		})
	}

	atom, err := feed.ToAtom()
	if err != nil {
		return err
	}

	f, err := os.Create(filepath.Join(config.TempPath, "feed-atom.xml"))
	if err != nil {
		return err
	}
	defer f.Close()

	f.WriteString(atom)

	return nil
}
