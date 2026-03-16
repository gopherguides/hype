package hype

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Validate_ValidDocument(t *testing.T) {
	r := require.New(t)

	cab := os.DirFS("testdata/validate/valid")
	p := NewParser(cab)

	doc, err := p.ParseFile("module.md")
	r.NoError(err)

	result := Validate(context.Background(), doc, ValidateOptions{})
	r.Empty(result.Issues)
	r.False(result.HasErrors())
	r.Equal("no issues found", result.Summary())
}

func Test_Validate_MissingImage(t *testing.T) {
	r := require.New(t)

	cab := os.DirFS("testdata/validate/missing-image")
	p := NewParser(cab)

	doc, err := p.ParseFile("module.md")
	r.NoError(err)

	result := &ValidationResult{}
	validateAssets(doc, result)

	r.Len(result.Issues, 1)
	r.Equal(SeverityError, result.Issues[0].Severity)
	r.Equal(CategoryAsset, result.Issues[0].Category)
	r.Contains(result.Issues[0].Message, "image not found")
	r.Contains(result.Issues[0].Message, "images/nonexistent.png")
}

func Test_Validate_MissingSourceCode(t *testing.T) {
	r := require.New(t)

	cab := os.DirFS("testdata/validate/missing-source")
	p := NewParser(cab)

	doc, err := p.ParseFile("module.md")
	r.NoError(err)

	result := &ValidationResult{}
	validateAssets(doc, result)

	r.Len(result.Issues, 1)
	r.Equal(SeverityError, result.Issues[0].Severity)
	r.Equal(CategoryAsset, result.Issues[0].Category)
	r.Contains(result.Issues[0].Message, "source file not found")
	r.Contains(result.Issues[0].Message, "missing.go")
}

func Test_Validate_HeadingSkip(t *testing.T) {
	r := require.New(t)

	cab := os.DirFS("testdata/validate/heading-skip")
	p := NewParser(cab)

	doc, err := p.ParseFile("module.md")
	r.NoError(err)

	result := &ValidationResult{}
	validateHeadingHierarchy(doc, result)

	r.Len(result.Issues, 1)
	r.Equal(SeverityWarning, result.Issues[0].Severity)
	r.Equal(CategoryHeading, result.Issues[0].Category)
	r.Contains(result.Issues[0].Message, "heading skip")
	r.Contains(result.Issues[0].Message, "h1 -> h3")
}

func Test_Validate_BrokenAnchor(t *testing.T) {
	r := require.New(t)

	cab := os.DirFS("testdata/validate/broken-anchor")
	p := NewParser(cab)

	doc, err := p.ParseFile("module.md")
	r.NoError(err)

	result := &ValidationResult{}
	validateLocalLinks(doc, result)

	r.Len(result.Issues, 1)
	r.Equal(SeverityError, result.Issues[0].Severity)
	r.Equal(CategoryLink, result.Issues[0].Category)
	r.Contains(result.Issues[0].Message, "anchor target not found")
	r.Contains(result.Issues[0].Message, "#nonexistent")
}

func Test_Validate_BrokenRef(t *testing.T) {
	r := require.New(t)

	cab := os.DirFS("testdata/validate/broken-ref")
	p := NewParser(cab)

	doc, err := p.ParseFile("module.md")
	r.NoError(err)

	result := &ValidationResult{}
	validateLocalLinks(doc, result)

	r.Len(result.Issues, 1)
	r.Equal(SeverityError, result.Issues[0].Severity)
	r.Equal(CategoryLink, result.Issues[0].Category)
	r.Contains(result.Issues[0].Message, "ref target not found")
	r.Contains(result.Issues[0].Message, "missing-fig")
}

func Test_Validate_DuplicateIDs(t *testing.T) {
	r := require.New(t)

	cab := os.DirFS("testdata/validate/duplicate-id")
	p := NewParser(cab)

	doc, err := p.ParseFile("module.md")
	r.NoError(err)

	result := &ValidationResult{}
	validateDuplicateIDs(doc, result)

	r.Len(result.Issues, 1)
	r.Equal(SeverityError, result.Issues[0].Severity)
	r.Equal(CategoryDuplicateID, result.Issues[0].Category)
	r.Contains(result.Issues[0].Message, "duplicate figure id")
	r.Contains(result.Issues[0].Message, "dup")
}

func Test_Validate_Summary(t *testing.T) {
	r := require.New(t)

	result := &ValidationResult{}
	r.Equal("no issues found", result.Summary())

	result.Add(ValidationIssue{Severity: SeverityError, Category: CategoryAsset, Message: "test"})
	result.Add(ValidationIssue{Severity: SeverityWarning, Category: CategoryHeading, Message: "test"})
	r.Equal("2 issues found (1 errors, 1 warnings)", result.Summary())
	r.True(result.HasErrors())
}

func Test_Validate_IssueString(t *testing.T) {
	r := require.New(t)

	issue := ValidationIssue{
		Severity: SeverityError,
		Category: CategoryAsset,
		Filename: "doc.md",
		Message:  "image not found: missing.png",
	}
	r.Equal("doc.md: ERROR asset: image not found: missing.png", issue.String())

	issue.Filename = ""
	r.Equal("ERROR asset: image not found: missing.png", issue.String())
}
