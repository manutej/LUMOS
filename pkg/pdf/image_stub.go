// +build !pdfcpu_available

package pdf

import (
	"fmt"
)

// extractImagesWithPdfcpu is a stub implementation used when pdfcpu is not available
// This keeps the application functional for text-only PDFs
func (d *Document) extractImagesWithPdfcpu(pageNum int, opts ImageExtractionOptions) ([]PageImage, error) {
	// Pdfcpu not available - return empty slice
	// The caller will handle this gracefully by returning empty images
	return []PageImage{}, nil
}

func init() {
	fmt.Println("Image support not available - pdfcpu library not found.")
	fmt.Println("To enable image extraction, ensure pdfcpu is installed:")
	fmt.Println("  go get github.com/pdfcpu/pdfcpu@latest")
}
