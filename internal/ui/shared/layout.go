package shared

import (
	"fmt"
	"strings"
)

// RenderPage renders a consistent page layout with a title, body, and help section.
func RenderPage(title string, body string, help string) string {
	var out strings.Builder

	out.WriteString(TitleStyle.Render(fmt.Sprintf("GitHub Actions Manager / %s", title)))
	out.WriteString("\n\n")

	out.WriteString(body)

	if help != "" {
		out.WriteString("\n")
		out.WriteString(help)
	}

	return ContainerStyle.Render(out.String())
}
