package freshness

import (
	"regexp"
	"strings"
)

// EstimateFreshness tente d'estimer si une page est "récente".
// Elle retourne true si elle semble être à jour (date récente dans l'URL ou dans le HTML).
func EstimateFreshness(url string, html string) bool {
	// Heuristique 1 : présence d'une année récente dans l'URL
	if strings.Contains(url, "2024") || strings.Contains(url, "2025") {
		return true
	}

	// Heuristique 2 : chercher une date dans le contenu HTML (format AAAA-MM ou AAAA)
	yearPattern := regexp.MustCompile(`20\d{2}`)
	years := yearPattern.FindAllString(html, -1)

	for _, y := range years {
		if y == "2024" || y == "2025" {
			return true
		}
	}

	// Sinon on considère la page comme "non fraîche"
	return false
}
