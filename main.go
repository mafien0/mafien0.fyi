package main

import (
	"embed"
	"html/template"
	"io/fs"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// Embed static files at compile time
//
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

// TODO: drop gin, i need pure go
func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// Serve embedded static files
	publicFS, _ := fs.Sub(staticFiles, "public")
	r.StaticFS("/static", http.FS(publicFS))

	// Serve index.html at root
	r.GET("/", func(c *gin.Context) {
		data, err := staticFiles.ReadFile("public/index.html")
		if err != nil {
			c.String(http.StatusNotFound, "Not found")
			return
		}
		c.Data(http.StatusOK, "text/html; charset=utf-8", data)
	})

	// Serve favicon at root
	r.GET("/favicon.png", func(c *gin.Context) {
		data, err := staticFiles.ReadFile("public/favicon.ico")
		if err != nil {
			c.Status(http.StatusNotFound)
			return
		}
		c.Data(http.StatusOK, "image/x-icon", data)
	})

	// htmx
	r.GET("/badges", func(c *gin.Context) {
		c.Header("Content-Type", "text/html; charset=utf-8")
		if err := tmpl.ExecuteTemplate(c.Writer, "badges", badges); err != nil {
			c.String(http.StatusInternalServerError, err.Error())
		}
	})

	// htmx
	r.GET("/links", func(c *gin.Context) {
		c.Header("Content-Type", "text/html; charset=utf-8")
		if err := tmpl.ExecuteTemplate(c.Writer, "links", links); err != nil {
			c.String(http.StatusInternalServerError, err.Error())
		}
	})

	// Get port from environment variable (Vercel sets this)
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	r.Run(":" + port)
}
