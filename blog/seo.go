package blog

import (
	"fmt"
	"strings"
)

type SEO struct {
	Title       string
	Description string
	URL         string
	Image       string
	Type        string
	Author      string
	Published   string
	TwitterCard string
	TwitterSite string
}

func (s SEO) OpenGraphTags() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf(`<meta property="og:title" content="%s">`, escapeAttr(s.Title)))
	sb.WriteString("\n")

	if s.Description != "" {
		sb.WriteString(fmt.Sprintf(`<meta property="og:description" content="%s">`, escapeAttr(s.Description)))
		sb.WriteString("\n")
	}

	if s.URL != "" {
		sb.WriteString(fmt.Sprintf(`<meta property="og:url" content="%s">`, escapeAttr(s.URL)))
		sb.WriteString("\n")
	}

	if s.Image != "" {
		sb.WriteString(fmt.Sprintf(`<meta property="og:image" content="%s">`, escapeAttr(s.Image)))
		sb.WriteString("\n")
	}

	ogType := s.Type
	if ogType == "" {
		ogType = "website"
	}
	sb.WriteString(fmt.Sprintf(`<meta property="og:type" content="%s">`, ogType))
	sb.WriteString("\n")

	return sb.String()
}

func (s SEO) TwitterCardTags() string {
	var sb strings.Builder

	card := s.TwitterCard
	if card == "" {
		card = "summary_large_image"
	}
	sb.WriteString(fmt.Sprintf(`<meta name="twitter:card" content="%s">`, card))
	sb.WriteString("\n")

	sb.WriteString(fmt.Sprintf(`<meta name="twitter:title" content="%s">`, escapeAttr(s.Title)))
	sb.WriteString("\n")

	if s.Description != "" {
		sb.WriteString(fmt.Sprintf(`<meta name="twitter:description" content="%s">`, escapeAttr(s.Description)))
		sb.WriteString("\n")
	}

	if s.Image != "" {
		sb.WriteString(fmt.Sprintf(`<meta name="twitter:image" content="%s">`, escapeAttr(s.Image)))
		sb.WriteString("\n")
	}

	if s.TwitterSite != "" {
		sb.WriteString(fmt.Sprintf(`<meta name="twitter:site" content="%s">`, escapeAttr(s.TwitterSite)))
		sb.WriteString("\n")
	}

	return sb.String()
}

func (s SEO) JSONLD() string {
	if s.Type != "article" {
		return ""
	}

	return fmt.Sprintf(`<script type="application/ld+json">
{
    "@context": "https://schema.org",
    "@type": "Article",
    "headline": %q,
    "author": {
        "@type": "Person",
        "name": %q
    },
    "datePublished": %q
}
</script>`, s.Title, s.Author, s.Published)
}

func escapeAttr(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	return s
}
