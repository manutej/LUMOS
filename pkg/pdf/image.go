package pdf

import (
	"fmt"
	"image"
)

// PageImage represents an image extracted from a PDF page
type PageImage struct {
	// Image data
	Data image.Image

	// Position on page (in points)
	X      float64
	Y      float64
	Width  float64
	Height float64

	// Index in page's image list
	Index int

	// Metadata
	Format string // "JPEG", "PNG", "TIFF", etc.
	Title  string // Optional title/alt text
}

// RenderSize specifies how to render an image in the terminal
type RenderSize struct {
	Width  int // Characters (not pixels)
	Height int // Lines (not pixels)
	// Aspect ratio preserved
}

// TerminalImageFormat describes how to render image for specific terminal
type TerminalImageFormat string

const (
	// Kitty graphics protocol - modern, best performance
	KittyFormat TerminalImageFormat = "kitty"

	// iTerm2 inline images - macOS, widely supported
	ITermFormat TerminalImageFormat = "iterm2"

	// SIXEL - older terminals, but still supported
	SixelFormat TerminalImageFormat = "sixel"

	// Unicode halfblocks - fallback for all terminals
	HalfblockFormat TerminalImageFormat = "halfblock"

	// Text placeholder - fallback when no graphics available
	TextFormat TerminalImageFormat = "text"
)

// ImageExtractionOptions controls how images are extracted
type ImageExtractionOptions struct {
	// MaxWidth filters images smaller than this (points)
	MinWidth float64

	// MaxHeight filters images smaller than this (points)
	MinHeight float64

	// OnlyInlineImages skips background images
	OnlyInlineImages bool

	// PreserveDPI if true, respects DPI information in PDF
	PreserveDPI bool

	// MaxImagesPerPage limits number of images extracted
	MaxImagesPerPage int
}

// DefaultImageExtractionOptions returns sensible defaults
func DefaultImageExtractionOptions() ImageExtractionOptions {
	return ImageExtractionOptions{
		MinWidth:         10,  // At least 10 points wide (for meaningful content)
		MinHeight:        10,  // At least 10 points tall
		OnlyInlineImages: true,
		PreserveDPI:      true,
		MaxImagesPerPage: 100, // Reasonable limit per page
	}
}

// GetPageImages extracts images from a specific page
// Requires pdfcpu library for extraction
// Returns empty slice if pdfcpu not available or page has no images
func (d *Document) GetPageImages(pageNum int, opts ImageExtractionOptions) ([]PageImage, error) {
	if pageNum < 1 || pageNum > d.pages {
		return nil, fmt.Errorf("page number out of range: %d", pageNum)
	}

	// Check cache first
	if cached, exists := d.imageCache.Get(pageNum); exists {
		return cached, nil
	}

	// TODO: Implement with pdfcpu when library is available
	// For now, return empty slice (images not extracted yet)
	// Code structure is ready for integration

	/*
	// Once pdfcpu is added, this is the implementation pattern:

	f, err := os.Open(d.filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to open PDF: %w", err)
	}
	defer f.Close()

	pdfContext, err := pdfcpu.Read(f, pdfcpu.NewDefaultConfiguration())
	if err != nil {
		return nil, fmt.Errorf("failed to read PDF with pdfcpu: %w", err)
	}

	pageDict := pdfContext.PageDict(pageNum)
	if pageDict == nil {
		return nil, fmt.Errorf("page %d not found in PDF", pageNum)
	}

	// Extract images from page using pdfcpu API
	var images []PageImage
	// ... extraction logic here ...

	// Cache and return
	d.imageCache.Put(pageNum, images)
	return images, nil
	*/

	images := []PageImage{}
	d.imageCache.Put(pageNum, images)
	return images, nil
}

// ImageInfo contains metadata about images on a page
type ImageInfo struct {
	PageNum      int
	ImageCount   int
	TotalSize    int64 // Total bytes of all images
	ImageDims    []image.Rectangle
	HasLargeImg  bool // Any image > 500x500
	HasSmallImg  bool // Any image < 100x100
}

// EstimatePageImageCount provides a rough estimate of images on a page
// without extracting (useful for UI display)
func (d *Document) EstimatePageImageCount(pageNum int) (int, error) {
	if pageNum < 1 || pageNum > d.pages {
		return 0, fmt.Errorf("page number out of range: %d", pageNum)
	}

	// TODO: Implement basic estimation with pdfcpu
	// For now, return 0
	return 0, nil
}

// HasImageCache returns true if image cache is available
func (d *Document) HasImageCache() bool {
	return d.imageCache != nil
}

// ClearImageCache clears the image cache
func (d *Document) ClearImageCache() {
	if d.imageCache != nil {
		d.imageCache.Clear()
	}
}

// ImageCacheStats returns statistics about image caching
func (d *Document) ImageCacheStats() ImageCacheStats {
	if d.imageCache == nil {
		return ImageCacheStats{}
	}
	return d.imageCache.Stats()
}
