package main

import (
	"flag"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	air_router_db "air_router/db"
	air_router_handlers "air_router/handlers"

	"github.com/gin-gonic/gin"
)

var (
	// Version information (set by build flags)
	Version   = "dev"
	BuildTime = "unknown"
	GitCommit = "unknown"
)

func main() {
	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	// Set gin to release mode
	gin.SetMode(gin.ReleaseMode)

	// Parse command-line flags
	configDir := flag.String("config", "./data", "Directory for data files")
	frontendPath := flag.String("path", "../frontend", "Path to frontend directory")
	port := flag.String("port", "8080", "Proxy API port")
	webPort := flag.String("web-port", "9000", "Web interface port")
	flag.Parse()

	// Initialize configuration directory
	if err := initConfigDir(*configDir); err != nil {
		log.Fatal("Error creating config directory: ", err)
	}

	// Setup database path
	dbPath := filepath.Join(*configDir, "accounts.db")

	// Convert frontend path to absolute path
	absFrontendPath, err := filepath.Abs(*frontendPath)
	if err != nil {
		log.Fatal("Error resolving frontend path: ", err)
	}

	// Initialize database
	dbConn, err := air_router_db.InitDB(dbPath)
	if err != nil {
		log.Fatal("Error initializing database: ", err)
	}
	defer dbConn.Close()

	// Initialize account database handler
	accountDB := &air_router_db.AccountDB{DB: dbConn}

	// Initialize model database handler
	modelDB := &air_router_db.ModelDB{DB: dbConn}

	// Initialize handlers
	handlers := air_router_handlers.NewHandlers(absFrontendPath, accountDB, modelDB)

	// Setup routers
	webRouter := air_router_handlers.SetupWebRouter(handlers.IndexHandler, handlers.AccountHandler, handlers.ModelHandler, handlers.ProxyHandler, absFrontendPath)
	proxyRouter := air_router_handlers.SetupProxyRouter(handlers.ProxyHandler)

	// Start web server
	webAddr := ":" + *webPort
	go func() {
		log.Printf("Web server starting on http://localhost%s", webAddr)
		if err := webRouter.Run(webAddr); err != nil {
			log.Fatal("Web server error: ", err)
		}
	}()

	// Start proxy server
	proxyAddr := ":" + *port
	printStartupInfo(webAddr, proxyAddr, absFrontendPath, dbPath)
	if err := proxyRouter.Run(proxyAddr); err != nil {
		log.Fatal("Proxy server error: ", err)
	}
}

// initConfigDir creates the configuration directory if it doesn't exist
func initConfigDir(dir string) error {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	return nil
}

// printStartupInfo prints server startup information
func printStartupInfo(webAddr, proxyAddr, frontendPath, dbPath string) {
	log.Println("================================")
	log.Println("  AI Router Server")
	log.Println("================================")
	log.Printf("  Version:    %s", Version)
	log.Printf("  Build:      %s", BuildTime)
	log.Printf("  Git Commit: %s", GitCommit)
	log.Println("--------------------------------")
	log.Printf("  Web Interface:    http://127.0.0.1%s", webAddr)
	log.Println("  Web endpoints:     /, /debug, /api/*")
	log.Printf("  Frontend:   %s", frontendPath)
	log.Printf("  Database:   %s", dbPath)
	log.Println("--------------------------------")
	log.Printf("  Proxy API:        http://127.0.0.1%s", proxyAddr)
	log.Println("  Proxy endpoints:   /v1/*")
	log.Println("================================")
}
