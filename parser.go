package hype

import (
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"path/filepath"
	"sync"

	"github.com/gopherguides/hype/atomx"
	"github.com/jmoiron/sqlx"
	"golang.org/x/net/html"
)

// Parser will convert HTML documents into, easy to use, nice types.
type Parser struct {
	fs.FS                                               // the filesystem to use
	Root           string                               `json:"root,omitempty"`      // the root directory of the parser
	FileName       string                               `json:"file_name,omitempty"` // the file being parsed (might be empty)
	DB             *sqlx.DB                             `json:"-"`                   // the database to use (optional)
	Client         *http.Client                         `json:"-"`                   // the http client to use (optional)
	PreProcessor   func(r io.Reader) (io.Reader, error) `json:"-"`                   // the preprocessor to use (optional)
	PreParseHooks  PreParseHooks                        `json:"-"`                   // the pre-parse hooks to use (optional)
	PostParseHooks PostParseHooks                       `json:"-"`                   // the post-parse hooks to use (optional)

	customTags   TagMap
	snippetRules map[string]string
	once         sync.Once
	sync.RWMutex
}

func (p *Parser) init() {
	p.once.Do(func() {

		if p.snippetRules == nil {
			p.snippetRules = map[string]string{}
		}

		if p.customTags == nil {
			p.customTags = TagMap{}
		}

		p.customTags[atomx.Meta] = func(node *Node) (Tag, error) {
			return NewMeta(node)
		}

		img := func(node *Node) (Tag, error) {
			return NewImage(p.FS, node)
		}

		p.customTags[atomx.Img] = img
		p.customTags[atomx.Image] = img

		p.customTags[atomx.Code] = func(node *Node) (Tag, error) {
			return NewCode(p, node)
		}

		p.customTags[atomx.Body] = func(node *Node) (Tag, error) {
			return NewBody(node)
		}

		p.customTags[atomx.File] = func(node *Node) (Tag, error) {
			return NewFile(p.FS, node)
		}

		p.customTags[atomx.Filegroup] = func(node *Node) (Tag, error) {
			return NewFileGroup(node)
		}

		p.customTags[atomx.Include] = func(node *Node) (Tag, error) {
			return NewInclude(node, p)
		}

		p.customTags[atomx.Page] = func(node *Node) (Tag, error) {
			return NewPage(node)
		}

		p.customTags[atomx.Head] = func(node *Node) (Tag, error) {
			return NewHead(node)
		}

		for _, h := range atomx.Headings() {
			p.customTags[h] = func(node *Node) (Tag, error) {
				return NewHeading(node)
			}
		}
	})
}

// SetCustomTag will set a custom tag for the parser.
func (p *Parser) SetCustomTag(atom atomx.Atom, fn CustomTagFn) {
	p.init()

	p.Lock()
	p.customTags[atom] = fn
	p.Unlock()
}

// SubParser will create a new Parser that will use the same FS as the parent.
func (p *Parser) SubParser(path string) (*Parser, error) {
	p.Lock()
	defer p.Unlock()

	cab, err := fs.Sub(p.FS, path)
	if err != nil {
		return nil, err
	}

	p2, err := NewParser(cab)
	if err != nil {
		return nil, err
	}
	p2.PreProcessor = p.PreProcessor
	p2.DB = p.DB

	p2.Root = filepath.Join(p.Root, path)

	for k, v := range p.snippetRules {
		p2.snippetRules[k] = v
	}

	for k, v := range p.customTags {
		p2.customTags[k] = v
	}

	return p2, nil
}

// NewParser will create a new Parser.
func NewParser(cab fs.FS) (*Parser, error) {
	if cab == nil {
		return nil, fmt.Errorf("cab can not be nil")
	}

	p := &Parser{
		FS:         cab,
		customTags: TagMap{},
		snippetRules: map[string]string{
			".go":   "// %s",
			".html": "<!-- %s -->",
			".md":   "// %s",
			".rb":   "# %s",
		},
	}

	return p, nil
}

// ParseFile will parse the requested file and return a Document.
// The file MUST be present in the parser's FS.
func (p *Parser) ParseFile(name string) (*Document, error) {
	if p == nil {
		return nil, fmt.Errorf("parser can not be nil")
	}
	p.FileName = name

	var r io.ReadCloser
	r, err := p.Open(name)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	if filepath.Ext(name) == ".md" {
		src, err := io.ReadAll(r)
		if err != nil {
			return nil, err
		}

		return p.ParseMD(src)
	}

	doc, err := p.ParseReader(r)
	if err != nil {
		return nil, err
	}

	return doc, nil
}

// ParseReader will parse the given reader and return a Document.
func (p *Parser) ParseReader(r io.Reader) (*Document, error) {
	p.init()

	var err error
	if p.PreProcessor != nil {
		r, err = p.PreProcessor(r)
		if err != nil {
			return nil, err
		}
	}

	if err := p.PreParseHooks.Run(p); err != nil {
		return nil, err
	}

	node, err := html.Parse(r)
	if err != nil {
		return nil, err
	}

	return p.NewDocument(node)
}

// CustomTag will return the custom tag for the given atom,
// or nil if there is no custom tag.
func (p *Parser) CustomTag(atom atomx.Atom) (CustomTagFn, bool) {
	p.init()

	p.Lock()
	defer p.Unlock()

	fn, ok := p.customTags[atom]
	return fn, ok
}
