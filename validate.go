package hype

import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"strings"
)

type IssueSeverity int

const (
	SeverityError IssueSeverity = iota
	SeverityWarning
)

func (s IssueSeverity) String() string {
	switch s {
	case SeverityError:
		return "ERROR"
	case SeverityWarning:
		return "WARN"
	default:
		return "UNKNOWN"
	}
}

func (s IssueSeverity) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

type IssueCategory string

const (
	CategoryAsset       IssueCategory = "asset"
	CategoryHeading     IssueCategory = "heading"
	CategoryLink        IssueCategory = "link"
	CategoryDuplicateID IssueCategory = "duplicate-id"
	CategoryExecution   IssueCategory = "execution"
)

type ValidationIssue struct {
	Severity IssueSeverity `json:"severity"`
	Category IssueCategory `json:"category"`
	Filename string        `json:"filename"`
	Element  string        `json:"element"`
	Message  string        `json:"message"`
}

func (vi ValidationIssue) String() string {
	if vi.Filename != "" {
		return fmt.Sprintf("%s: %s %s: %s", vi.Filename, vi.Severity, vi.Category, vi.Message)
	}
	return fmt.Sprintf("%s %s: %s", vi.Severity, vi.Category, vi.Message)
}

type ValidationResult struct {
	Issues []ValidationIssue `json:"issues"`
}

func (vr *ValidationResult) Add(issue ValidationIssue) {
	vr.Issues = append(vr.Issues, issue)
}

func (vr *ValidationResult) Errors() []ValidationIssue {
	var errs []ValidationIssue
	for _, i := range vr.Issues {
		if i.Severity == SeverityError {
			errs = append(errs, i)
		}
	}
	return errs
}

func (vr *ValidationResult) Warnings() []ValidationIssue {
	var warns []ValidationIssue
	for _, i := range vr.Issues {
		if i.Severity == SeverityWarning {
			warns = append(warns, i)
		}
	}
	return warns
}

func (vr *ValidationResult) HasErrors() bool {
	for _, i := range vr.Issues {
		if i.Severity == SeverityError {
			return true
		}
	}
	return false
}

func (vr *ValidationResult) Summary() string {
	errs := len(vr.Errors())
	warns := len(vr.Warnings())
	total := errs + warns
	if total == 0 {
		return "no issues found"
	}
	return fmt.Sprintf("%d issues found (%d errors, %d warnings)", total, errs, warns)
}

type ValidateOptions struct {
	Exec bool
}

func Validate(ctx context.Context, doc *Document, opts ValidateOptions) *ValidationResult {
	result := &ValidationResult{}

	validateAssets(doc, result)
	validateHeadingHierarchy(doc, result)
	validateLocalLinks(doc, result)
	validateDuplicateIDs(doc, result)

	if opts.Exec {
		validateExecution(ctx, doc, result)
	}

	return result
}

func validateAssets(doc *Document, result *ValidationResult) {
	images := ByType[*Image](doc.Nodes)
	for _, img := range images {
		src, err := img.ValidAttr("src")
		if err != nil {
			continue
		}
		if strings.HasPrefix(src, "http") {
			continue
		}
		if _, err := fs.Stat(doc.FS, src); err != nil {
			result.Add(ValidationIssue{
				Severity: SeverityError,
				Category: CategoryAsset,
				Filename: img.Filename,
				Element:  img.StartTag(),
				Message:  fmt.Sprintf("image not found: %s", src),
			})
		}
	}

	codes := ByType[*SourceCode](doc.Nodes)
	for _, code := range codes {
		src, ok := code.Get("src")
		if !ok {
			continue
		}
		filePath := strings.SplitN(src, "#", 2)[0]
		if _, err := fs.Stat(doc.FS, filePath); err != nil {
			result.Add(ValidationIssue{
				Severity: SeverityError,
				Category: CategoryAsset,
				Filename: code.Filename,
				Element:  code.StartTag(),
				Message:  fmt.Sprintf("source file not found: %s", filePath),
			})
		}
	}
}

func validateHeadingHierarchy(doc *Document, result *ValidationResult) {
	headings := ByType[*Heading](doc.Nodes)
	if len(headings) == 0 {
		return
	}

	prevLevel := headings[0].Level()
	for _, h := range headings[1:] {
		level := h.Level()
		if level > prevLevel+1 {
			result.Add(ValidationIssue{
				Severity: SeverityWarning,
				Category: CategoryHeading,
				Filename: h.Filename,
				Element:  h.StartTag(),
				Message:  fmt.Sprintf("heading skip: h%d -> h%d (expected h%d)", prevLevel, level, prevLevel+1),
			})
		}
		prevLevel = level
	}
}

func validateLocalLinks(doc *Document, result *ValidationResult) {
	ids := make(map[string]bool)

	figures := ByType[*Figure](doc.Nodes)
	for _, f := range figures {
		if id, ok := f.Get("id"); ok {
			ids[id] = true
		}
	}

	links := ByType[*Link](doc.Nodes)
	for _, l := range links {
		href, err := l.Href()
		if err != nil || !strings.HasPrefix(href, "#") {
			continue
		}
		anchor := strings.TrimPrefix(href, "#")
		if !ids[anchor] {
			result.Add(ValidationIssue{
				Severity: SeverityError,
				Category: CategoryLink,
				Filename: l.Filename,
				Element:  l.StartTag(),
				Message:  fmt.Sprintf("anchor target not found: %s", href),
			})
		}
	}
}

func validateDuplicateIDs(doc *Document, result *ValidationResult) {
	seen := make(map[string]bool)

	figures := ByType[*Figure](doc.Nodes)
	for _, f := range figures {
		id, ok := f.Get("id")
		if !ok {
			continue
		}
		if seen[id] {
			result.Add(ValidationIssue{
				Severity: SeverityWarning,
				Category: CategoryDuplicateID,
				Filename: f.Filename,
				Element:  f.StartTag(),
				Message:  fmt.Sprintf("duplicate figure id: %s", id),
			})
		}
		seen[id] = true
	}
}

func validateExecution(ctx context.Context, doc *Document, result *ValidationResult) {
	if err := doc.Execute(ctx); err != nil {
		result.Add(ValidationIssue{
			Severity: SeverityError,
			Category: CategoryExecution,
			Filename: doc.Filename,
			Message:  err.Error(),
		})
	}
}
