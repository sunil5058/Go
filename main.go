package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Book struct {
	ID          int      `json:"id"`
	Title       string   `json:"title"`
	Author      string   `json:"author"`
	Description string   `json:"description"`
	Category    string   `json:"category"`
	ReadTime    string   `json:"readTime"`
	Tags        []string `json:"tags"`
}

var catalog = []Book{
	{ID: 1, Title: "The Go Way", Author: "Maya Singh", Description: "A practical introduction to writing clear, dependable Go programs.", Category: "Technology", ReadTime: "6h 20m", Tags: []string{"Go", "Programming"}},
	{ID: 2, Title: "Quiet Habits", Author: "Alex Rowan", Description: "Small, sustainable routines for a more focused life.", Category: "Self growth", ReadTime: "4h 10m", Tags: []string{"Habits", "Wellness"}},
	{ID: 3, Title: "The Last Observatory", Author: "Iris Bell", Description: "A scientist receives a signal from the edge of the solar system.", Category: "Science fiction", ReadTime: "8h 45m", Tags: []string{"Space", "Adventure"}},
	{ID: 4, Title: "Letters from Kathmandu", Author: "Nima Karki", Description: "A collection of intimate stories set in a changing city.", Category: "Literary fiction", ReadTime: "5h 30m", Tags: []string{"Nepal", "Stories"}},
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/api/books", books)
	mux.HandleFunc("/api/books/", bookByID)
	server := &http.Server{Addr: ":8080", Handler: logging(mux)}
	log.Println("E-book platform running at http://localhost:8080")
	log.Fatal(server.ListenAndServe())
}

func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(`<!doctype html><html><head><meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1"><title>Leafline</title><style>body{max-width:760px;margin:4rem auto;padding:0 1.5rem;font:16px system-ui;color:#17251d;background:#f7faf6}h1{font-size:2.8rem;margin-bottom:.4rem}p{line-height:1.6}.book{background:#fff;border:1px solid #dce8dc;border-radius:12px;padding:1rem 1.25rem;margin:1rem 0}small{color:#55705b}</style></head><body><h1>Leafline</h1><p>Your small, thoughtful e-book shelf.</p><main id="books">Loading books...</main><script>fetch('/api/books').then(r=>r.json()).then(books=>document.querySelector('#books').innerHTML=books.map(b=>'<article class="book"><small>'+b.category+' / '+b.readTime+'</small><h2>'+b.title+'</h2><strong>'+b.author+'</strong><p>'+b.description+'</p></article>').join('')).catch(()=>document.querySelector('#books').textContent='Unable to load the catalog.')</script></body></html>`))
}

func books(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		methodNotAllowed(w)
		return
	}
	query := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("q")))
	if query == "" {
		respondJSON(w, http.StatusOK, catalog)
		return
	}
	matched := make([]Book, 0)
	for _, book := range catalog {
		text := strings.ToLower(book.Title + " " + book.Author + " " + book.Description + " " + book.Category + " " + strings.Join(book.Tags, " "))
		if strings.Contains(text, query) {
			matched = append(matched, book)
		}
	}
	respondJSON(w, http.StatusOK, matched)
}

func bookByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		methodNotAllowed(w)
		return
	}
	id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/api/books/"))
	if err != nil {
		http.NotFound(w, r)
		return
	}
	for _, book := range catalog {
		if book.ID == id {
			respondJSON(w, http.StatusOK, book)
			return
		}
	}
	respondJSON(w, http.StatusNotFound, map[string]string{"error": "book not found"})
}

func respondJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func methodNotAllowed(w http.ResponseWriter) {
	respondJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
}

func logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
