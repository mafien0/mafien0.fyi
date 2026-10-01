package main

import (
	"embed"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strconv"
)

//go:embed public/* public/badges/* templates/*
var staticFiles embed.FS

var tmpl = template.Must(template.ParseFS(staticFiles, "templates/*.html"))

// -- Links --

type Link struct {
	Name string
	Link string
}

var links = []Link{
	{"Page source", "https://github.com/mafien0/mafien0.fyi"},
	{"GitHub", "https://github.com/mafien0"},
	{"Lichess", "https://lichess.org/@/mafien0"},
	{"Steam", "https://steamcommunity.com/id/mafien0"},
	{"Twitter", "https://x.com/mafien0"},
}

// -- Badges --

type Badge struct {
	Name string
	Link string
	Path string
}

var badges = []Badge{
	{"mafien0", "https://mafien0.fyi", "/static/badges/mafien0.png"},
	{"nixos", "https://nixos.org", "/static/badges/nixos.png"},
	{"nowindows", "https://nixos.org", "/static/badges/nowindows.gif"},
	{"linuxnow", "https://distrowatch.com", "/static/badges/linux.gif"},
	{"eeto", "https://eightyeightthirty.one", "/static/badges/eeto.png"},
	{"ublock", "https://github.com/gorhill/uBlock", "/static/badges/ublock.png"},
	{"helium", "https://helium.computer", "/static/badges/helium.png"},
	{"ltt", "https://www.youtube.com/LinusTechTips", "/static/badges/ltt.png"},
	{"steam", "https://store.steampowered.com", "/static/badges/steam.gif"},
	{"oneshot", "https://store.steampowered.com/app/420530", "/static/badges/oneshot.png"},
	{"jsab", "https://store.steampowered.com/app/531510", "/static/badges/jsab.gif"},
	{"celeste", "https://www.celestegame.com", "/static/badges/celeste.gif"},
	{"ltg", "https://lowtierfailure.com", "/static/badges/ltg.png"},
	{"qbit", "https://www.qbittorrent.org", "/static/badges/qbit.png"},
	{"minecraft", "https://prismlauncher.org", "/static/badges/minecraft.gif"},
	{"rust", "https://rust-lang.org", "/static/badges/rust.gif"},
	{"bnix", "https://nixos.org", "/static/badges/bnix.gif"},
	{"neovim", "https://neovim.io", "/static/badges/neovim.png"},
	{"blender", "https://www.blender.org", "/static/badges/blender.gif"},
}

// -- rest --
// actually this is funny because rest is like the rest of code,
// but at the same time, its REST api, you get it?

func main() {
	mux := http.NewServeMux()

	// Serve embedded static files
	publicFS, _ := fs.Sub(staticFiles, "public")
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(publicFS))))

	// Serve index.html at root
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		data, err := staticFiles.ReadFile("public/index.html")
		if err != nil {
			http.Error(w, "Not Found", http.StatusNotFound)
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(data)
	})

	// Fake favicon.ico
	mux.HandleFunc("GET /favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		data, err := staticFiles.ReadFile("public/favicon.png")
		if err != nil {
			http.Error(w, "Not Found", http.StatusNotFound)
		}
		w.Header().Set("Content-Type", "image/png")
		w.Write(data)
	})

	// htmx
	mux.HandleFunc("GET /badges", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := tmpl.ExecuteTemplate(w, "badges", badges); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	})

	// htmx
	mux.HandleFunc("GET /links", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		width, _ := strconv.Atoi(r.Header.Get("X-Screen-Width"))
		name := "links"
		if width > 0 && width < 700 {
			name = "links-mobile"
		}
		if err := tmpl.ExecuteTemplate(w, name, links); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	})

	// Get port from environment variable (Vercel sets this)
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	// Serve
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
