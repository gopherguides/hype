package hype

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

// YouTube represents an embedded YouTube video.
// Usage: <youtube id="VIDEO_ID"></youtube>
// Or with optional title: <youtube id="VIDEO_ID" title="Video Title"></youtube>
type YouTube struct {
	*Element
}

// youtubeIDPattern validates YouTube video IDs (11 characters, alphanumeric with - and _)
var youtubeIDPattern = regexp.MustCompile(`^[a-zA-Z0-9_-]{11}$`)

func (yt *YouTube) MarshalJSON() ([]byte, error) {
	if yt == nil {
		return nil, ErrIsNil("youtube")
	}

	yt.RLock()
	defer yt.RUnlock()

	m, err := yt.JSONMap()
	if err != nil {
		return nil, err
	}

	m["type"] = toType(yt)

	id, err := yt.VideoID()
	if err != nil {
		return nil, err
	}
	m["video_id"] = id

	if title, ok := yt.Get("title"); ok {
		m["title"] = title
	}

	return json.MarshalIndent(m, "", "  ")
}

func (yt *YouTube) VideoID() (string, error) {
	if yt == nil {
		return "", ErrIsNil("youtube")
	}

	return yt.ValidAttr("id")
}

func (yt *YouTube) Title() string {
	if yt == nil {
		return ""
	}

	title, _ := yt.Get("title")
	return title
}

func (yt *YouTube) MD() string {
	if yt == nil {
		return ""
	}

	return yt.String()
}

func (yt *YouTube) String() string {
	if yt == nil {
		return ""
	}

	id, err := yt.VideoID()
	if err != nil {
		return ""
	}

	title := yt.Title()
	if title == "" {
		title = "YouTube video player"
	}

	bb := &strings.Builder{}
	bb.WriteString(`<div class="youtube-embed">`)
	bb.WriteString("\n  ")
	fmt.Fprintf(bb, `<iframe src="https://www.youtube.com/embed/%s"`, id)
	fmt.Fprintf(bb, ` title=%q`, title)
	bb.WriteString(` frameborder="0"`)
	bb.WriteString(` allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share"`)
	bb.WriteString(` referrerpolicy="strict-origin-when-cross-origin"`)
	bb.WriteString(` allowfullscreen`)
	bb.WriteString(`></iframe>`)
	bb.WriteString("\n</div>")

	return bb.String()
}

func NewYouTube(el *Element) (*YouTube, error) {
	if el == nil {
		return nil, ErrIsNil("element")
	}

	yt := &YouTube{
		Element: el,
	}

	id, err := yt.ValidAttr("id")
	if err != nil {
		return nil, err
	}

	if !youtubeIDPattern.MatchString(id) {
		return nil, yt.WrapErr(fmt.Errorf("invalid YouTube video ID %q: must be 11 alphanumeric characters", id))
	}

	return yt, nil
}

func NewYouTubeNodes(p *Parser, el *Element) (Nodes, error) {
	yt, err := NewYouTube(el)
	if err != nil {
		return nil, err
	}

	return Nodes{yt}, nil
}
