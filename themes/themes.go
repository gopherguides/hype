package themes

import (
	"embed"
	"fmt"
	"html/template"
	"io"
	"os"
	"slices"
)

//go:embed css/*.css
var cssFS embed.FS

//go:embed templates/document.html
var documentTemplate string

const DefaultTheme = "github"

var builtinThemes = []string{
	"github",
	"github-dark",
	"solarized-light",
	"solarized-dark",
	"swiss",
	"air",
	"retro",
}

func ListThemes() []string {
	themes := make([]string, len(builtinThemes))
	copy(themes, builtinThemes)
	return themes
}

func IsBuiltinTheme(name string) bool {
	return slices.Contains(builtinThemes, name)
}

func GetCSS(themeName string) (string, error) {
	if !IsBuiltinTheme(themeName) {
		return "", fmt.Errorf("unknown theme: %s", themeName)
	}

	data, err := cssFS.ReadFile("css/" + themeName + ".css")
	if err != nil {
		return "", fmt.Errorf("failed to read theme %s: %w", themeName, err)
	}

	return string(data), nil
}

func LoadCustomCSS(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to read custom CSS file %s: %w", path, err)
	}

	return string(data), nil
}

type RenderData struct {
	Title string
	CSS   template.CSS
	Body  template.HTML
}

func Render(w io.Writer, data RenderData) error {
	tmpl, err := template.New("document").Parse(documentTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse document template: %w", err)
	}

	if err := tmpl.Execute(w, data); err != nil {
		return fmt.Errorf("failed to execute document template: %w", err)
	}

	return nil
}
