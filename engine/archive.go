package engine

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/petarov/nenu/config"
)

type archivePageData struct {
	*config.YML
	Posts []Post
	// Post an always null hack that helps with the header template if condition
	Post *Post
}

// SpewArchive generate site posts archive
func SpewArchive(posts []*Post, templates *Templates) error {
	fmt.Println("| Generating archive ...")

	f, err := os.Create(filepath.Join(config.TempPath, "archive.html"))
	if err != nil {
		return err
	}
	defer f.Close()

	published := make([]Post, 0)
	for _, p := range posts {
		if p.IsPublish {
			published = append(published, *p)
		}
	}

	apd := &archivePageData{config.YMLConfig, published, nil}

	// header
	if err = templates.Header.Execute(f, apd); err != nil {
		return err
	}
	// body
	if err = templates.Archive.Execute(f, apd); err != nil {
		return err
	}
	// footer
	if err = templates.Footer.Execute(f, apd); err != nil {
		return err
	}

	return nil
}
