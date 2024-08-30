package parser

import (
	"bufio"
	"fmt"
	"io"
	"kinsyn/pkg/commons"
	"strconv"
	"strings"
	"time"
)

/* Kindle Highlight Format
Mistborn: The Final Empire (Brandon Sanderson)
- Your Highlight on Location 1080-1081 | Added on Friday, August 18, 2023 2:26:04 PM

The only reason to be subservient to those with power is so that you can learn to someday take what they have.
==========
Modern Software Engineering: Doing What Works to Build Better Software Faster
- Your Highlight on Location 613-615 | Added on Thursday, December 14, 2023 11:30:08 PM

Engineering is the application of an empirical, scientific approach to finding efficient, economic solutions to practical problems.
==========
*/

func ParseHighlights(reader io.Reader) ([]commons.Highlight, error) {
	scanner := bufio.NewScanner(reader)
	var highlights []commons.Highlight
	var currentHighlight commons.Highlight

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if isTitleAndAuthorLine(line) {
			err := parseTitleAndAuthor(line, &currentHighlight)
			if err != nil {
				return nil, err // Log or handle error as needed
			}
		} else if isLocationLine(line) {
			err := parseLocationAndDate(line, &currentHighlight)
			if err != nil {
				return nil, err // Log or handle error as needed
			}
		} else if isSeparatorLine(line) {
			highlights = append(highlights, currentHighlight)
			currentHighlight = commons.Highlight{} // Reset for the next highlight
		} else {
			appendContent(line, &currentHighlight)
		}
	}

	// Add the last highlight if it exists
	if currentHighlight.BookTitle != "" && currentHighlight.BookAuthor != "" {
		highlights = append(highlights, currentHighlight)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanner error: %v", err)
	}

	return highlights, nil
}

// isTitleAndAuthorLine checks if the line contains the book title and author.
func isTitleAndAuthorLine(line string) bool {
	return strings.Contains(line, "(") && strings.HasSuffix(line, ")")
}

// isLocationLine checks if the line contains location and date information.
func isLocationLine(line string) bool {
	return strings.HasPrefix(line, "- Your Highlight on Location")
}

// isSeparatorLine checks if the line is the separator ("==========").
func isSeparatorLine(line string) bool {
	return line == "=========="
}

// parseTitleAndAuthor extracts the book title and author from the line.
func parseTitleAndAuthor(line string, highlight *commons.Highlight) error {
	parts := strings.Split(line, "(")
	if len(parts) != 2 {
		return fmt.Errorf("invalid title and author line format")
	}
	highlight.BookTitle = strings.TrimSpace(parts[0])
	highlight.BookAuthor = strings.TrimSuffix(strings.TrimSpace(parts[1]), ")")
	return nil
}

// parseLocationAndDate extracts the location and date from the line.
// - Your Highlight on Location 1885-1887 | Added on Sunday, January 7, 2024 11:57:36 PM
func parseLocationAndDate(line string, highlight *commons.Highlight) error {
	parts := strings.Split(line, "| Added on ")
	if len(parts) != 2 {
		return fmt.Errorf("invalid location and date line format")
	}

	// - Your Highlight on Location 1885-1887
	locationParts := strings.Split(parts[0], " ")
	if len(locationParts) < 3 {
		return fmt.Errorf("invalid location format")
	}

	locationRange := locationParts[len(locationParts)-2]
	// 1885-1887
	rangeParts := strings.Split(locationRange, "-")

	if len(rangeParts) != 2 {
		return fmt.Errorf("invalid location range format")
	}

	start, err := strconv.Atoi(rangeParts[0])
	if err != nil {
		return fmt.Errorf("failed to parse location start: %v", err)
	}
	end, err := strconv.Atoi(rangeParts[1])
	if err != nil {
		return fmt.Errorf("failed to parse location end: %v", err)
	}
	highlight.BookLocationStart = start
	highlight.BookLocationEnd = end

	// Sunday, January 7, 2024 11:57:36 PM
	createdAt, err := time.Parse("Monday, January 2, 2006 3:04:05 PM", strings.TrimSpace(parts[1]))
	if err != nil {
		return fmt.Errorf("failed to parse time: %v", err)
	}
	highlight.CreatedAt = createdAt

	return nil
}

// appendContent adds content to the current highlight.
func appendContent(line string, highlight *commons.Highlight) {
	if highlight.Content != "" {
		highlight.Content += "\n"
	}
	highlight.Content += line
}
