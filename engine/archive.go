package engine

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/petarov/nenu/config"
)

type archivePageData struct {
	*config.YML
	Posts []PostMeta
	// Post an always null hack that helps with the header template if condition
	Post *PostMeta
}

// SpewArchive generate site posts archive
func SpewArchive(meta []*PostMeta, dest *os.File, templates *Templates) error {
	path := config.PostsPath
	fmt.Printf("| Generating archive %s...\n", path)

	var f *os.File
	f, err := os.Create(filepath.Join(config.TempPath, "archive.html"))
	if err != nil {
		return err
	}
	defer f.Close()

	postMetas := make([]PostMeta, 0)
	for _, v := range meta {
		postMetas = append(postMetas, *v)
	}

	pd := &archivePageData{config.YMLConfig, postMetas, nil}

	// header
	if err = templates.Header.Execute(f, pd); err != nil {
		return err
	}
	// body
	if err = templates.Archive.Execute(f, pd); err != nil {
		return err
	}
	// footer
	if err = templates.Footer.Execute(f, pd); err != nil {
		return err
	}

	return nil
}
