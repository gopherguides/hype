package hype

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/net/html"
)

func Test_Element_StartTag(t *testing.T) {
	t.Parallel()

	r := require.New(t)

	hn := &html.Node{
		Data: "div",
	}

	attrs := &Attributes{}
	r.NoError(attrs.Set("class", "foo"))
	r.NoError(attrs.Set("id", "bar"))

	table := []struct {
		name string
		e    *Element
		exp  string
	}{
		{name: "empty", e: &Element{}, exp: ""},
		{name: "with atom", e: &Element{
			HTMLNode: hn,
		}, exp: "<div>"},
		{name: "with attrs", e: &Element{
			HTMLNode:   hn,
			Attributes: attrs,
		}, exp: `<div class="foo" id="bar">`},
	}

	for _, tc := range table {
		t.Run(tc.name, func(t *testing.T) {
			r := require.New(t)

			r.Equal(tc.exp, tc.e.StartTag())
		})
	}

}

func Test_Element_EndTag(t *testing.T) {
	t.Parallel()

	table := []struct {
		name string
		e    *Element
		exp  string
	}{
		{name: "empty", e: &Element{}, exp: ""},
		{name: "with atom", e: &Element{
			HTMLNode: &html.Node{Data: "div"},
		}, exp: "</div>"},
	}

	for _, tc := range table {
		t.Run(tc.name, func(t *testing.T) {
			r := require.New(t)

			r.Equal(tc.exp, tc.e.EndTag())
		})
	}

}

func Test_Element_String(t *testing.T) {
	t.Parallel()

	r := require.New(t)

	hn := &html.Node{
		Data: "div",
	}

	attrs := &Attributes{}
	r.NoError(attrs.Set("class", "foo"))
	r.NoError(attrs.Set("id", "bar"))

	table := []struct {
		name string
		e    *Element
		exp  string
	}{
		{name: "empty", e: &Element{}, exp: ""},
		{name: "with atom", e: &Element{
			HTMLNode: hn,
		}, exp: "<div></div>"},
		{name: "with attrs", e: &Element{
			HTMLNode:   hn,
			Attributes: attrs,
		}, exp: `<div class="foo" id="bar"></div>`},
		{name: "with kids", e: &Element{
			HTMLNode: hn,
			Nodes:    Nodes{TextNode("hello")},
		}, exp: "<div>hello</div>"},
	}

	for _, tc := range table {
		t.Run(tc.name, func(t *testing.T) {
			r := require.New(t)

			r.Equal(tc.exp, tc.e.String())
		})
	}

}
