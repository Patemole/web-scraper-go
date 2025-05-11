package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Patemole/web-scraper-go/internal/fetcher"
	"github.com/Patemole/web-scraper-go/internal/serpapi"
	"github.com/joho/godotenv"
)

type URLResult struct {
	URL       string
	HTML      string
	CleanText string
	FetchTime time.Duration
	CleanTime time.Duration
	SaveTime  time.Duration
	Error     error
}

func processURL(url string, resultChan chan<- URLResult) {
	result := URLResult{URL: url}

	// Fetch
	fetchStart := time.Now()
	html, err := fetcher.FetchHTML(url)
	if err != nil {
		result.Error = err
		resultChan <- result
		return
	}
	result.HTML = html
	result.FetchTime = time.Since(fetchStart)

	// Clean
	cleanStart := time.Now()
	cleanText, err := fetcher.CleanHTML(html)
	if err != nil {
		result.Error = err
		resultChan <- result
		return
	}
	result.CleanText = cleanText
	result.CleanTime = time.Since(cleanStart)

	// Save
	saveStart := time.Now()
	if err := fetcher.SaveHTMLToFile(url, cleanText); err != nil {
		result.Error = err
		resultChan <- result
		return
	}
	result.SaveTime = time.Since(saveStart)

	resultChan <- result
}

func main() {
	// Définition des flags
	query := flag.String("query", "site:upenn.edu tuition international students", "Search query to use with SERP API")
	numResults := flag.Int("results", 3, "Number of results to fetch from SERP API")
	outputDir := flag.String("output", "data", "Directory to save the results")
	flag.Parse()

	startTime := time.Now()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Créer le dossier de sortie s'il n'existe pas
	if err := os.MkdirAll(*outputDir, 0755); err != nil {
		log.Fatal("Error creating output directory:", err)
	}

	apiKey := os.Getenv("SERP_API_KEY")
	if apiKey == "" {
		log.Fatal("SERP_API_KEY is not set")
	}

	// Mesure du temps pour la recherche SERP
	serpStartTime := time.Now()
	urls, err := serpapi.GetTopUrls(apiKey, *query, *numResults)
	if err != nil {
		log.Fatal(err)
	}
	serpDuration := time.Since(serpStartTime)
	fmt.Printf("\nTemps de recherche SERP: %v\n", serpDuration)
	fmt.Printf("Requête: %s\n", *query)
	fmt.Printf("Nombre de résultats demandés: %d\n", *numResults)

	fmt.Println("\nFetching HTML for top URLs:")

	// Créer un canal pour recevoir les résultats
	resultChan := make(chan URLResult, len(urls))

	// Lancer les goroutines pour chaque URL
	for _, url := range urls {
		go processURL(url, resultChan)
	}

	// Collecter les résultats
	var results []URLResult
	var totalFetchTime, totalCleanTime, totalSaveTime time.Duration

	for i := 0; i < len(urls); i++ {
		result := <-resultChan
		results = append(results, result)

		if result.Error != nil {
			fmt.Printf("\nURL %d: %s\n", i+1, result.URL)
			fmt.Printf("Error: %v\n", result.Error)
			continue
		}

		fmt.Printf("\nURL %d: %s\n", i+1, result.URL)
		fmt.Printf("Temps de fetch: %v\n", result.FetchTime)
		fmt.Printf("Temps de nettoyage: %v\n", result.CleanTime)
		fmt.Printf("Temps de sauvegarde: %v\n", result.SaveTime)
		fmt.Printf("HTML length: %d bytes\n", len(result.HTML))
		fmt.Printf("Texte nettoyé length: %d bytes\n", len(result.CleanText))
		fmt.Printf("HTML saved successfully\n")

		totalFetchTime += result.FetchTime
		totalCleanTime += result.CleanTime
		totalSaveTime += result.SaveTime
	}

	// Affichage des statistiques finales
	totalDuration := time.Since(startTime)
	fmt.Printf("\n=== Statistiques d'exécution ===\n")
	fmt.Printf("Temps total d'exécution: %v\n", totalDuration)
	fmt.Printf("Temps moyen par URL:\n")
	fmt.Printf("  - Fetch: %v\n", totalFetchTime/time.Duration(len(urls)))
	fmt.Printf("  - Nettoyage: %v\n", totalCleanTime/time.Duration(len(urls)))
	fmt.Printf("  - Sauvegarde: %v\n", totalSaveTime/time.Duration(len(urls)))
	fmt.Printf("Temps total par opération:\n")
	fmt.Printf("  - Recherche SERP: %v\n", serpDuration)
	fmt.Printf("  - Fetch total: %v\n", totalFetchTime)
	fmt.Printf("  - Nettoyage total: %v\n", totalCleanTime)
	fmt.Printf("  - Sauvegarde total: %v\n", totalSaveTime)
}
