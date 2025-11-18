package ui

import "github.com/luxor/lumos/pkg/pdf"

// TOCLoadedMsg indicates table of contents has been loaded
type TOCLoadedMsg struct {
	TOC *pdf.TableOfContents
}

// TOCSelectedMsg indicates a TOC entry was selected
type TOCSelectedMsg struct {
	Page int
}
