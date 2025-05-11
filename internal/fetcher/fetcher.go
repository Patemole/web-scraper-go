package fetcher

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
	"time"

	"golang.org/x/net/html"
)

// extractText extrait récursivement le texte d'un nœud HTML
func extractText(n *html.Node, text *strings.Builder) {
	if n.Type == html.TextNode {
		text.WriteString(strings.TrimSpace(n.Data))
		text.WriteString(" ")
	}

	// Ignorer les balises non désirées
	if n.Type == html.ElementNode {
		switch n.Data {
		case "script", "style", "nav", "header", "footer", "iframe", "noscript", "meta", "link":
			return
		}
	}

	// Parcourir récursivement les enfants
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		extractText(c, text)
	}
}

// FetchHTML télécharge le HTML avec support gzip
func FetchHTML(url string) (string, error) {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Accept-Encoding", "gzip")
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; WebScraper/1.0)")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var reader io.ReadCloser
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(resp.Body)
		if err != nil {
			return "", err
		}
		defer reader.Close()
	default:
		reader = resp.Body
	}

	body, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// FetchAndCleanHTML télécharge et nettoie le HTML en streaming
func FetchAndCleanHTML(url string) (string, error) {
	// Télécharger le HTML
	htmlContent, err := FetchHTML(url)
	if err != nil {
		return "", err
	}

	// Parser le HTML
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return "", err
	}

	// Extraire le texte
	var text strings.Builder
	extractText(doc, &text)

	// Nettoyer le texte final
	cleanText := strings.TrimSpace(text.String())
	cleanText = strings.Join(strings.Fields(cleanText), " ") // Normaliser les espaces

	return cleanText, nil
}
