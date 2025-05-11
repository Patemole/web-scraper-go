package fetcher

import (
	"os"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// CleanFilename génère un nom de fichier safe à partir d'une URL
func CleanFilename(url string) string {
	re := regexp.MustCompile(`[^a-zA-Z0-9]+`)
	return re.ReplaceAllString(url, "_")
}

// CleanHTML nettoie le HTML pour ne garder que le contenu principal
func CleanHTML(html string) (string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return "", err
	}

	// Supprimer les éléments non désirés
	doc.Find("script, style, nav, header, footer, iframe, noscript, meta, link").Remove()

	// Supprimer les attributs non essentiels
	doc.Find("*").Each(func(i int, s *goquery.Selection) {
		s.RemoveAttr("class")
		s.RemoveAttr("id")
		s.RemoveAttr("style")
		s.RemoveAttr("onclick")
		s.RemoveAttr("onload")
	})

	// Extraire le contenu principal (généralement dans main ou article)
	content := doc.Find("main, article, .content, #content, .main, #main")
	if content.Length() == 0 {
		// Si aucun conteneur principal n'est trouvé, prendre le body
		content = doc.Find("body")
	}

	// Nettoyer le texte
	text := content.Text()
	text = strings.TrimSpace(text)
	text = regexp.MustCompile(`\s+`).ReplaceAllString(text, " ")

	return text, nil
}

// SaveHTMLToFile sauvegarde le contenu dans /data
func SaveHTMLToFile(url, content string) error {
	filename := "data/" + CleanFilename(url) + ".txt"
	return os.WriteFile(filename, []byte(content), 0644)
}
