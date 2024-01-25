package server

// md module
//
// Copyright (c) 2023 - Valentin Kuznetsov <vkuznet@gmail.com>
//
import (
	"io"
	"io/fs"
	"log"

	"github.com/gomarkdown/markdown"
	mhtml "github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

// MDToHTML function to generate HTML from given markdown file
func MDToHTML(fsys fs.FS, fname string) (string, error) {
	/*
		filesFS, err := fs.Sub(StaticFs, "static/markdown")
		if err != nil {
			log.Println("ERROR: unable to open static/markdown", err)
			return "", err
		}
		log.Printf("### fileFS %+v", filesFS)
		file, err := filesFS.Open(fname)
	*/
	file, err := fsys.Open(fname)
	if err != nil {
		log.Println("ERROR: unable to open", fname, err)
		return "", err
	}
	/*
	   file, err := os.Open(fname)
	   if err != nil {
	       log.Println("ERROR: unable to open", fname, err)
	       return "", err
	   }
	*/
	defer file.Close()
	var md []byte
	md, err = io.ReadAll(file)
	if err != nil {
		return "", err
	}

	// create markdown parser with extensions
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(md)

	// create HTML renderer with extensions
	//     htmlFlags := mhtml.CommonFlags | mhtml.HrefTargetBlank
	htmlFlags := mhtml.CommonFlags
	opts := mhtml.RendererOptions{Flags: htmlFlags}
	renderer := mhtml.NewRenderer(opts)
	content := markdown.Render(doc, renderer)
	//     return html.EscapeString(string(content)), nil
	return string(content), nil
}
