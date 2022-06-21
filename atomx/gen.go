// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build ignore
// +build ignore

//go:generate go run gen.go
package main

import (
	"bytes"
	"fmt"
	"go/format"
	"log"
	"os"
	"sort"

	"github.com/gobuffalo/flect"
)

var hype = []string{
	"cmd",
	"code",
	"file",
	"filegroup",
	"go",
	"include",
	"metadata",
	"page",
	"ref",
	"unknown",
}

func main() {

	all := append(elements, attributes...)
	all = append(all, headings...)
	all = append(all, inlines...)
	all = append(all, eventHandlers...)
	all = append(all, extra...)
	all = append(all, hype...)

	// sort.Strings(all)

	dedup := map[string]string{}

	for _, k := range all {
		dedup[k] = flect.Pascalize(k)
	}

	uniq := make([]string, 0, len(dedup))
	for k := range dedup {
		uniq = append(uniq, k)
	}

	sort.Strings(uniq)

	bb := &bytes.Buffer{}

	fmt.Fprintln(bb, "// Code generated by gen.go. DO NOT EDIT.")
	fmt.Fprintf(bb, "\npackage atomx\n\n")

	fmt.Fprintf(bb, "// Headings returns all the heading elements.\n")
	fmt.Fprintln(bb, "\nfunc Headings() Atoms {")
	fmt.Fprintln(bb, "\treturn Atoms{")
	for _, k := range headings {
		fmt.Fprintf(bb, "\t\t%s,\n", dedup[k])
	}
	fmt.Fprintln(bb, "\t}")
	fmt.Fprintln(bb, "}")

	fmt.Fprintf(bb, "// Inlines returns all the inline elements.\n")
	fmt.Fprintln(bb, "\nfunc Inlines() Atoms {")
	fmt.Fprintln(bb, "\treturn Atoms{")
	for _, k := range inlines {
		fmt.Fprintf(bb, "\t\t%s,\n", dedup[k])
	}
	fmt.Fprintln(bb, "\t}")
	fmt.Fprintln(bb, "}")

	fmt.Fprintln(bb, "const(")

	for _, k := range uniq {
		fmt.Fprintf(bb, "\t%s Atom = \"%s\"\n", dedup[k], k)
	}

	fmt.Fprintln(bb, ")")

	b, err := format.Source(bb.Bytes())
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Create("atoms.go")
	if err != nil {
		log.Fatal(err)
	}

	f.Write(b)
	f.Close()
}

var headings = []string{
	"h1",
	"h2",
	"h3",
	"h4",
	"h5",
	"h6",
}

var inlines = []string{
	"a",
	"b",
	"br",
	"image",
	"img",
	"link",
	"ref",
}

var elements = []string{
	"a",
	"abbr",
	"address",
	"area",
	"article",
	"aside",
	"audio",
	"b",
	"base",
	"bdi",
	"bdo",
	"blockquote",
	"body",
	"br",
	"button",
	"canvas",
	"caption",
	"cite",
	"code",
	"col",
	"colgroup",
	"command",
	"data",
	"datalist",
	"dd",
	"del",
	"details",
	"dfn",
	"dialog",
	"div",
	"dl",
	"dt",
	"em",
	"embed",
	"fieldset",
	"file",
	"figcaption",
	"figure",
	"footer",
	"form",
	"head",
	"header",
	"hgroup",
	"hr",
	"html",
	"i",
	"iframe",
	"img",
	"input",
	"ins",
	"kbd",
	"keygen",
	"label",
	"legend",
	"li",
	"link",
	"main",
	"map",
	"mark",
	"menu",
	"menuitem",
	"meta",
	"meter",
	"nav",
	"noscript",
	"object",
	"ol",
	"optgroup",
	"option",
	"output",
	"p",
	"param",
	"picture",
	"pre",
	"progress",
	"q",
	"rp",
	"rt",
	"ruby",
	"s",
	"samp",
	"script",
	"section",
	"select",
	"slot",
	"small",
	"source",
	"span",
	"strong",
	"style",
	"sub",
	"summary",
	"sup",
	"table",
	"tbody",
	"td",
	"template",
	"textarea",
	"tfoot",
	"th",
	"thead",
	"time",
	"title",
	"tr",
	"track",
	"u",
	"ul",
	"var",
	"video",
	"wbr",
}

// https://html.spec.whatwg.org/multipage/indices.html#attributes-3
//
// "challenge", "command", "contextmenu", "dropzone", "icon", "keytype", "mediagroup",
// "radiogroup", "spellcheck", "scoped", "seamless", "sortable" and "sorted" have been removed from the spec,
// but are kept here for backwards compatibility.
var attributes = []string{
	"abbr",
	"accept",
	"accept-charset",
	"accesskey",
	"action",
	"allowfullscreen",
	"allowpaymentrequest",
	"allowusermedia",
	"alt",
	"as",
	"async",
	"autocomplete",
	"autofocus",
	"autoplay",
	"challenge",
	"charset",
	"checked",
	"cite",
	"class",
	"color",
	"cols",
	"colspan",
	"command",
	"content",
	"contenteditable",
	"contextmenu",
	"controls",
	"coords",
	"crossorigin",
	"data",
	"datetime",
	"default",
	"defer",
	"dir",
	"dirname",
	"disabled",
	"download",
	"draggable",
	"dropzone",
	"enctype",
	"for",
	"form",
	"formaction",
	"formenctype",
	"formmethod",
	"formnovalidate",
	"formtarget",
	"headers",
	"height",
	"hidden",
	"high",
	"href",
	"hreflang",
	"http-equiv",
	"icon",
	"id",
	"inputmode",
	"integrity",
	"is",
	"ismap",
	"itemid",
	"itemprop",
	"itemref",
	"itemscope",
	"itemtype",
	"keytype",
	"kind",
	"label",
	"lang",
	"list",
	"loop",
	"low",
	"manifest",
	"max",
	"maxlength",
	"media",
	"mediagroup",
	"method",
	"min",
	"minlength",
	"multiple",
	"muted",
	"name",
	"nomodule",
	"nonce",
	"novalidate",
	"open",
	"optimum",
	"pattern",
	"ping",
	"placeholder",
	"playsinline",
	"poster",
	"preload",
	"radiogroup",
	"readonly",
	"referrerpolicy",
	"rel",
	"required",
	"reversed",
	"rows",
	"rowspan",
	"sandbox",
	"spellcheck",
	"scope",
	"scoped",
	"seamless",
	"selected",
	"shape",
	"size",
	"sizes",
	"sortable",
	"sorted",
	"slot",
	"span",
	"spellcheck",
	"src",
	"srcdoc",
	"srclang",
	"srcset",
	"start",
	"step",
	"style",
	"tabindex",
	"target",
	"title",
	"translate",
	"type",
	"typemustmatch",
	"updateviacache",
	"usemap",
	"value",
	"width",
	"workertype",
	"wrap",
}

// "onautocomplete", "onautocompleteerror", "onmousewheel",
// "onshow" and "onsort" have been removed from the spec,
// but are kept here for backwards compatibility.
var eventHandlers = []string{
	"onabort",
	"onautocomplete",
	"onautocompleteerror",
	"onauxclick",
	"onafterprint",
	"onbeforeprint",
	"onbeforeunload",
	"onblur",
	"oncancel",
	"oncanplay",
	"oncanplaythrough",
	"onchange",
	"onclick",
	"onclose",
	"oncontextmenu",
	"oncopy",
	"oncuechange",
	"oncut",
	"ondblclick",
	"ondrag",
	"ondragend",
	"ondragenter",
	"ondragexit",
	"ondragleave",
	"ondragover",
	"ondragstart",
	"ondrop",
	"ondurationchange",
	"onemptied",
	"onended",
	"onerror",
	"onfocus",
	"onhashchange",
	"oninput",
	"oninvalid",
	"onkeydown",
	"onkeypress",
	"onkeyup",
	"onlanguagechange",
	"onload",
	"onloadeddata",
	"onloadedmetadata",
	"onloadend",
	"onloadstart",
	"onmessage",
	"onmessageerror",
	"onmousedown",
	"onmouseenter",
	"onmouseleave",
	"onmousemove",
	"onmouseout",
	"onmouseover",
	"onmouseup",
	"onmousewheel",
	"onwheel",
	"onoffline",
	"ononline",
	"onpagehide",
	"onpageshow",
	"onpaste",
	"onpause",
	"onplay",
	"onplaying",
	"onpopstate",
	"onprogress",
	"onratechange",
	"onreset",
	"onresize",
	"onrejectionhandled",
	"onscroll",
	"onsecuritypolicyviolation",
	"onseeked",
	"onseeking",
	"onselect",
	"onshow",
	"onsort",
	"onstalled",
	"onstorage",
	"onsubmit",
	"onsuspend",
	"ontimeupdate",
	"ontoggle",
	"onunhandledrejection",
	"onunload",
	"onvolumechange",
	"onwaiting",
}

// extra are ad-hoc values not covered by any of the lists above.
var extra = []string{
	"acronym",
	"align",
	"annotation",
	"annotation-xml",
	"applet",
	"basefont",
	"bgsound",
	"big",
	"blink",
	"center",
	"color",
	"desc",
	"face",
	"font",
	"foreignObject", // HTML is case-insensitive, but SVG-embedded-in-HTML is case-sensitive.
	"foreignobject",
	"frame",
	"frameset",
	"image",
	"isindex",
	"listing",
	"malignmark",
	"marquee",
	"math",
	"mglyph",
	"mi",
	"mn",
	"mo",
	"ms",
	"mtext",
	"nobr",
	"noembed",
	"noframes",
	"plaintext",
	"prompt",
	"public",
	"rb",
	"rtc",
	"spacer",
	"strike",
	"svg",
	"system",
	"tt",
	"xmp",
}
