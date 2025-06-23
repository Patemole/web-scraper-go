# web-scraper-go
A fast Go-based tool to query SERP API, scrape HTML pages, and prepare them for LLM input.
The full pipeline (SERP + parallel fetch + stream-clean) runs in ~1.5 seconds, approximately 4× faster compared to typical Python scraping stacks (~6 seconds).
Reducing processing time and ensuring high-quality data retrieval are critical to providing fast, accurate responses from LLM-based tools.


# web-scraper-go

**web-scraper-go** est un outil ultra-rapide écrit en Go pour interroger Google via SERP API, télécharger en parallèle les pages HTML ciblées (ex. : `upenn.edu`), les nettoyer efficacement, et les sauvegarder localement. Le pipeline est optimisé pour la performance réseau, la compression (gzip), et la scalabilité.

---

## 🚀 Fonctionnalités

- 🔍 Interroge **SERP API** avec une requête ciblée (`site:upenn.edu keyword`)
- ⚡ Télécharge en parallèle les pages HTML des résultats
- 📦 Utilise la compression GZIP pour optimiser la bande passante
- 🧹 Nettoie le HTML pour n'en garder que le texte utile (prêt pour LLM)
- 💾 Sauvegarde le HTML brut et le texte nettoyé dans `/data`
- 📊 Affiche des statistiques d'exécution précises pour chaque étape

---

## 📦 Structure du projet

```
web-scraper-go/
├── cmd/                    # Main CLI
│   └── main.go
├── internal/
│   ├── serpapi/            # Requêtes SERP API
│   ├── fetcher/            # Téléchargement & sauvegarde HTML
│   ├── freshness/          # (à venir) estimation de fraîcheur
│   └── llm/                # Nettoyage HTML → texte brut
├── data/                   # Fichiers HTML et/ou texte nettoyé
├── .env.example            # Modèle de config avec SERP_API_KEY
├── go.mod
```

---

## 🛠️ Configuration

1. Clone le repo :
```bash
git clone https://github.com/<your-username>/web-scraper-go.git
cd web-scraper-go
```

2. Installe les dépendances :
```bash
go mod tidy
```

3. Crée un fichier `.env` à la racine avec ta clé SERP API :
```
SERP_API_KEY=your_key_here
```

---

## ▶️ Exécution

### Construction des requêtes

Le flag `--query` permet de spécifier à la fois le domaine et les mots-clés de recherche. La syntaxe est :
```
site:domaine.com mot-clé1 mot-clé2 ...
```

Exemples de requêtes :
```bash
# Recherche sur upenn.edu
./web-scraper-go --query="site:upenn.edu academic calendar"

# Recherche sur harvard.edu
./web-scraper-go --query="site:harvard.edu admission requirements"

# Recherche sur mit.edu avec plusieurs mots-clés
./web-scraper-go --query="site:mit.edu computer science graduate program"
```

### Utilisation en ligne de commande

Le programme accepte plusieurs flags pour personnaliser son comportement :

```bash
# Utilisation avec les valeurs par défaut
./web-scraper-go

# Personnalisation de la requête
./web-scraper-go --query="site:upenn.edu academic calendar"

# Personnalisation complète
./web-scraper-go --query="site:upenn.edu academic calendar" --results=5 --output="results"
```

Flags disponibles :
- `--query` : Requête de recherche (défaut: "site:upenn.edu tuition international students")
  - Format : `site:domaine.com mot-clé1 mot-clé2 ...`
  - Exemple : `site:mit.edu artificial intelligence research`
- `--results` : Nombre de résultats à récupérer (défaut: 3)
- `--output` : Dossier de sortie pour les fichiers (défaut: "data")

### Intégration avec Python

Vous pouvez facilement intégrer web-scraper-go dans vos scripts Python :

```python
import subprocess
import json

def scrape_web_pages(domain, keywords, num_results=3, output_dir="data"):
    """
    Exécute web-scraper-go avec les paramètres spécifiés.
    
    Args:
        domain (str): Domaine à rechercher (ex: "upenn.edu")
        keywords (str): Mots-clés de recherche
        num_results (int): Nombre de résultats à récupérer
        output_dir (str): Dossier de sortie
    
    Returns:
        dict: Résultats de l'exécution
    """
    # Construction de la requête
    query = f"site:{domain} {keywords}"
    
    # Construction de la commande
    cmd = [
        "./web-scraper-go",
        f"--query={query}",
        f"--results={num_results}",
        f"--output={output_dir}"
    ]
    
    # Exécution de la commande
    result = subprocess.run(
        cmd,
        capture_output=True,
        text=True
    )
    
    # Vérification des erreurs
    if result.returncode != 0:
        raise Exception(f"Erreur lors de l'exécution: {result.stderr}")
    
    return {
        "stdout": result.stdout,
        "stderr": result.stderr,
        "returncode": result.returncode
    }

# Exemple d'utilisation
try:
    results = scrape_web_pages(
        domain="upenn.edu",
        keywords="academic calendar",
        num_results=5,
        output_dir="results"
    )
    print("Résultats:", results["stdout"])
except Exception as e:
    print(f"Erreur: {e}")
```

Exemple de sortie :
```
URL 1: https://srfs.upenn.edu/...
Temps de fetch: 58 ms
HTML length: 212519 bytes
Texte nettoyé length: 8053 bytes
HTML saved successfully
...
Temps total d'exécution: 1.8 s
```

---

## 📈 Performances

- Temps total pour 3 pages : **~1.8 secondes**
- Temps moyen de téléchargement HTML : **~620 ms/page**
- Traitement full pipeline (fetch + nettoyage + sauvegarde) optimisé à chaque étape

---

---

## 📌 Optimisations avancées (présentes ou à implémenter)

Voici les optimisations mises en place ou planifiées pour atteindre des performances maximales sur le scraping universitaire :

### 🔍 SERP API
- `engine=google_light` → ne récupère que les résultats organiques (plus rapide, plus léger)
- `ludicrous=true` → exécution prioritaire sur infrastructure accélérée SerpApi
- ⏱️ Objectif : **temps de réponse < 1 seconde**

### 📥 HTML Fetching
- ✅ **Client HTTP réutilisable (Keep-Alive)** → évite le coût du handshake TLS à chaque requête
- ✅ **Compression GZIP** activée → réduit le poids des pages HTML
- 🧠 **DNS caching local** recommandé → limite les latences de résolution
- 🌍 **Exécution depuis un VPS proche géographiquement (ex: AWS US)** → réduit la latence réseau
- 🚦 **Limitation volontaire des requêtes en parallèle (3–5)** → évite la congestion réseau

### 🧹 Nettoyage HTML
- 🚫 Évite `io.ReadAll` quand possible
- ✅ Lecture en **streaming**, ligne par ligne → réduit la mémoire et accélère le traitement
- 🧠 Nettoyage ciblé : uniquement les balises utiles (`<p>`, `<h1>`, etc.)

### 💾 Sauvegarde
- ✅ Écriture directe sur disque (`WriteFile`) sans surcharge
- 💡 Possibilité d’ajouter :
  - compression `.txt.gz` des textes nettoyés
  - cache intelligent pour éviter de re-fetch une page déjà sauvegardée

---

### 🧠 Idée future : scorings & résumés intelligents
- Intégrer un LLM pour :
  - résumer chaque page
  - classer la pertinence
  - enrichir les métadonnées de chaque fichier sauvegardé

---

## 📄 Licence

MIT
