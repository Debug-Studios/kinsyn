package parser

import (
	"strings"
	"testing"
)

func TestParseHighlightsSuccess(t *testing.T) {
	input := `
The Software Engineer's Guidebook: Navigating senior, tech lead, and staff engineer positions at tech companies and startups (Orosz, Gergely)
- Your Highlight on Location 6868-6871 | Added on Sunday, June 16, 2024 1:16:19 AM

Architecture should not exist in a vacuum; in fact, good architecture is always coupled to the business. What are the goals of the business? The underlying architecture should enable these goals to be met, and if the current systems get in the way of these goals, then it's a worthwhile - and necessary! - task to change them.
==========
The Personal MBA (Kaufman, Josh)
- Your Highlight on Location 853-854 | Added on Monday, June 24, 2024 1:19:20 PM

Life's tough. It's tougher if you're stupid. —JOHN WAYNE, WESTERN FILM ICON
==========
The Personal MBA (Kaufman, Josh)
- Your Highlight on Location 1005-1008 | Added on Monday, June 24, 2024 2:26:11 PM

It is important that students bring a certain ragamuffin, barefoot irreverence to their studies; they are not here to worship what is known, but to question it. —JACOB BRONOWSKI, WRITER AND PRESENTER OF THE ASCENT OF MAN
==========
The Personal MBA (Kaufman, Josh)
- Your Highlight on Location 1086-1089 | Added on Monday, June 24, 2024 2:30:55 PM

I've long believed that a certain system—which almost any intelligent person can learn—works way better than the systems most people use [to understand the world]. What you need is a latticework of mental models in your head. And, with that system, things gradually fit together in a way that enhances cognition.
==========
`
	reader := strings.NewReader(input)
	highlights, err := ParseHighlights(reader)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if len(highlights) != 4 {
		t.Fatalf("Expected 4 highlights, got: %d", len(highlights))
	}

}

func TestParseHighlightsEmptyInput(t *testing.T) {
	input := ``
	reader := strings.NewReader(input)
	highlights, err := ParseHighlights(reader)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if len(highlights) != 0 {
		t.Fatalf("Expected 0 highlights, got: %d", len(highlights))
	}
}

func TestParseHighlightsInvalidTitleAndAuthor(t *testing.T) {
	input := `
Invalid Title and Author Line
- Your Highlight on Location 1080-1081 | Added on Friday, August 18, 2023 2:26:04 PM

Some text content.
==========
`
	reader := strings.NewReader(input)
	_, err := ParseHighlights(reader)

	if err == nil {
		t.Fatalf("Expected error, got none")
	}

	expectedErr := "author and title line missing"
	if !strings.Contains(err.Error(), expectedErr) {
		t.Errorf("Expected error to contain '%s', got: %v", expectedErr, err)
	}
}

func TestParseHighlightsInvalidLocation(t *testing.T) {
	input := `
Mistborn: The Final Empire (Brandon Sanderson)
- Your Highlight on Location invalid-range | Added on Friday, August 18, 2023 2:26:04 PM

Some text content.
==========
`
	reader := strings.NewReader(input)
	_, err := ParseHighlights(reader)

	if err == nil {
		t.Fatalf("Expected error, got none")
	}

	expectedErr := "failed to parse location start: strconv.Atoi: parsing \"invalid\": invalid syntax"
	if !strings.Contains(err.Error(), expectedErr) {
		t.Errorf("Expected error to contain '%s', got: %v", expectedErr, err)
	}
}

func TestParseHighlightsMissingDate(t *testing.T) {
	input := `
Mistborn: The Final Empire (Brandon Sanderson)
- Your Highlight on Location 1080-1081 | Added on

Some text content.
==========
`
	reader := strings.NewReader(input)
	_, err := ParseHighlights(reader)

	if err == nil {
		t.Fatalf("Expected error, got none")
	}

	expectedErr := "invalid location and date line format"
	if !strings.Contains(err.Error(), expectedErr) {
		t.Errorf("Expected error to contain '%s', got: %v", expectedErr, err)
	}
}
