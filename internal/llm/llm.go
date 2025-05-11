package llm

import (
	"regexp"
	"strings"
)

// PrepareForLLM nettoie le HTML pour n’en garder que le texte pertinent.
func PrepareForLLM(html string) string {
	// Supprimer les balises HTML
	tagRe := regexp.MustCompile(`<[^>]+>`)
	text := tagRe.ReplaceAllString(html, "")

	// Nettoyage supplémentaire
	text = strings.TrimSpace(text)
	text = strings.ReplaceAll(text, "\n", " ")
	text = strings.Join(strings.Fields(text), " ") // remove excess whitespace

	return text
}
