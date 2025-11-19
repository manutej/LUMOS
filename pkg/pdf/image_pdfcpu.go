// +build pdfcpu_available

package pdf

import (
	"bytes"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"os"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

// extractImagesWithPdfcpu uses pdfcpu to extract images from a PDF page
func (d *Document) extractImagesWithPdfcpu(pageNum int, opts ImageExtractionOptions) ([]PageImage, error) {
	file, err := os.Open(d.filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to open PDF: %w", err)
	}
	defer file.Close()

	var images []PageImage
	extractedCount := 0

	// Callback function for each image found
	digestFunc := func(img model.Image, single bool, maxPageDigits int) error {
		// Skip if we've reached the limit
		if opts.MaxImagesPerPage > 0 && extractedCount >= opts.MaxImagesPerPage {
			return nil
		}

		// Filter by size if options specify
		if img.Width < int(opts.MinWidth) || img.Height < int(opts.MinHeight) {
			return nil
		}

		// Skip image masks and thumbnails
		if img.IsImgMask || img.Thumb {
			return nil
		}

		// Read image data into memory
		imageData, err := io.ReadAll(img.Reader)
		if err != nil {
			return fmt.Errorf("failed to read image data: %w", err)
		}

		// Decode image to get actual image.Image object
		imgReader := bytes.NewReader(imageData)
		decodedImg, _, err := image.Decode(imgReader)
		if err != nil {
			// Log but continue - some image formats might not decode properly
			fmt.Printf("Warning: Failed to decode image %s: %v\n", img.Name, err)
			return nil
		}

		// Determine format from filter
		format := "PNG"
		if img.Filter == "DCTDecode" {
			format = "JPEG"
		} else if img.Filter == "JPXDecode" {
			format = "JP2"
		} else if img.Filter == "CCITTFaxDecode" {
			format = "TIFF"
		}

		// Create PageImage struct
		pageImg := PageImage{
			Data:   decodedImg,
			Index:  extractedCount,
			Format: format,
			Title:  img.Name,
			Width:  float64(img.Width),
			Height: float64(img.Height),
		}

		images = append(images, pageImg)
		extractedCount++

		return nil
	}

	// Extract images from the specific page using pdfcpu
	// Page numbers are 1-indexed
	pageStr := fmt.Sprintf("%d", pageNum)
	err = api.ExtractImages(file, []string{pageStr}, digestFunc, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to extract images from page %d: %w", pageNum, err)
	}

	return images, nil
}
