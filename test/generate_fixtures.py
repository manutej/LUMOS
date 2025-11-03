#!/usr/bin/env python3
"""
Generate test PDF fixtures for LUMOS testing.
Requires: pip install reportlab
"""

from reportlab.lib.pagesizes import letter
from reportlab.lib.styles import getSampleStyleSheet
from reportlab.platypus import SimpleDocTemplate, Paragraph, Spacer, PageBreak
from reportlab.pdfgen import canvas
import os


def generate_simple_pdf(filename):
    """Generate a simple single-page PDF for basic testing."""
    c = canvas.Canvas(filename, pagesize=letter)
    width, height = letter

    # Set metadata
    c.setTitle("LUMOS Test Document")
    c.setAuthor("LUMOS Test Suite")
    c.setSubject("Testing PDF Reading Functionality")
    c.setCreator("LUMOS Fixture Generator")

    # Add title
    c.setFont("Helvetica-Bold", 16)
    c.drawString(50, height - 50, "LUMOS Test PDF")

    # Add body text
    c.setFont("Helvetica", 12)
    y = height - 100

    lines = [
        "This is a simple test PDF for the LUMOS PDF reader.",
        "",
        "Title: LUMOS Test Document",
        "Author: LUMOS Test Suite",
        "Subject: Testing PDF Reading Functionality",
        "Creator: LUMOS Fixture Generator",
        "",
        "This document contains basic text for testing:",
        "- Document loading",
        "- Page counting",
        "- Text extraction",
        "- Metadata reading",
        "",
        "Page Count: 1",
        "Word Count: Approximately 50 words",
        "Line Count: Multiple lines for testing",
    ]

    for line in lines:
        c.drawString(50, y, line)
        y -= 20

    c.save()
    print(f"✅ Generated {filename}")


def generate_multipage_pdf(filename):
    """Generate a multi-page PDF for testing pagination."""
    c = canvas.Canvas(filename, pagesize=letter)
    width, height = letter

    # Set metadata
    c.setTitle("LUMOS Multi-Page Test")
    c.setAuthor("LUMOS Test Suite")
    c.setSubject("Multi-page Testing")
    c.setCreator("LUMOS Fixture Generator")

    pages = [
        {
            "title": "LUMOS Multi-Page Test PDF - Page 1",
            "content": [
                "This is the first page of a multi-page test document.",
                "",
                "Content for testing:",
                "- Page navigation",
                "- Multi-page caching",
                "- Page range operations",
                "",
                "Total Pages: 5"
            ]
        },
        {
            "title": "LUMOS Multi-Page Test PDF - Page 2",
            "content": [
                "This is the second page.",
                "",
                "Testing features:",
                "- Sequential page access",
                "- Cache behavior with multiple pages",
                "- Page-by-page text extraction",
                "",
                "Current Page: 2 of 5"
            ]
        },
        {
            "title": "LUMOS Multi-Page Test PDF - Page 3",
            "content": [
                "This is the third page, the middle of the document.",
                "",
                "Key testing scenarios:",
                "- Middle page access",
                "- Cache eviction with limited cache size",
                "- Random page access patterns",
                "",
                "Current Page: 3 of 5"
            ]
        },
        {
            "title": "LUMOS Multi-Page Test PDF - Page 4",
            "content": [
                "This is the fourth page.",
                "",
                "Additional test cases:",
                "- Near-end page access",
                "- Cache statistics validation",
                "- Page range boundary testing",
                "",
                "Current Page: 4 of 5"
            ]
        },
        {
            "title": "LUMOS Multi-Page Test PDF - Page 5",
            "content": [
                "This is the final page of the test document.",
                "",
                "Final test validations:",
                "- Last page access",
                "- Complete document traversal",
                "- End-to-end workflow testing",
                "",
                "Current Page: 5 of 5 (END)"
            ]
        }
    ]

    for page_data in pages:
        # Title
        c.setFont("Helvetica-Bold", 16)
        c.drawString(50, height - 50, page_data["title"])

        # Content
        c.setFont("Helvetica", 12)
        y = height - 100

        for line in page_data["content"]:
            c.drawString(50, y, line)
            y -= 20

        c.showPage()

    c.save()
    print(f"✅ Generated {filename}")


def generate_search_test_pdf(filename):
    """Generate a PDF with specific text patterns for search testing."""
    c = canvas.Canvas(filename, pagesize=letter)
    width, height = letter

    # Set metadata
    c.setTitle("LUMOS Search Test")
    c.setAuthor("LUMOS Test Suite")
    c.setSubject("Search Testing")
    c.setCreator("LUMOS Fixture Generator")

    # Add title
    c.setFont("Helvetica-Bold", 16)
    c.drawString(50, height - 50, "LUMOS Search Test PDF")

    # Add test content
    c.setFont("Helvetica", 11)
    y = height - 100

    lines = [
        "This document contains specific text patterns for search testing.",
        "",
        "Case Sensitivity Tests:",
        "- The word 'test' appears in lowercase",
        "- The word 'Test' appears with capital T",
        "- The word 'TEST' appears in all caps",
        "",
        "Word Boundary Tests:",
        "- testing (contains 'test' but is a different word)",
        "- test (exact word match)",
        "- retest (contains 'test' as suffix)",
        "",
        "Multiple Occurrences:",
        "The word search appears here.",
        "Another search instance appears here.",
        "And a third search occurrence here.",
        "",
        "Special Characters:",
        "test@example.com",
        "hello-world",
        "under_score",
        "",
        "Long text for context extraction testing:",
        "Lorem ipsum dolor sit amet, consectetur adipiscing elit.",
        "This sentence contains the word MATCH right here in the middle.",
        "Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.",
        "",
        "Line boundary testing:",
        "First line with keyword",
        "Second line with keyword",
        "Third line with keyword",
    ]

    for line in lines:
        if y < 50:  # Start new page if needed
            c.showPage()
            c.setFont("Helvetica", 11)
            y = height - 50
        c.drawString(50, y, line)
        y -= 18

    c.save()
    print(f"✅ Generated {filename}")


def main():
    fixtures_dir = "fixtures"

    # Ensure fixtures directory exists
    os.makedirs(fixtures_dir, exist_ok=True)

    # Generate all test fixtures
    generate_simple_pdf(os.path.join(fixtures_dir, "simple.pdf"))
    generate_multipage_pdf(os.path.join(fixtures_dir, "multipage.pdf"))
    generate_search_test_pdf(os.path.join(fixtures_dir, "search_test.pdf"))

    print("\n✅ All test fixtures generated successfully!")
    print(f"\nGenerated files in {fixtures_dir}/:")
    for f in os.listdir(fixtures_dir):
        if f.endswith('.pdf'):
            size = os.path.getsize(os.path.join(fixtures_dir, f))
            print(f"  - {f} ({size:,} bytes)")


if __name__ == "__main__":
    main()
