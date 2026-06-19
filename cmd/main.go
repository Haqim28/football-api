package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/yourname/football-api/internal/config"
	"github.com/yourname/football-api/internal/database"
	"github.com/yourname/football-api/internal/middleware"

	// Import docs agar Swagger spec ter-register via init()
	_ "github.com/yourname/football-api/docs"

	teamModule   "github.com/yourname/football-api/internal/modules/team"
	playerModule "github.com/yourname/football-api/internal/modules/player"
	matchModule  "github.com/yourname/football-api/internal/modules/match"
	goalModule   "github.com/yourname/football-api/internal/modules/goal"
	reportModule "github.com/yourname/football-api/internal/modules/report"
)

// @title                       Football API - AYO Technical Test 2026
// @version                     1.0
// @description                Backend API untuk manajemen tim sepakbola amatir perusahaan XYZ. Menyediakan CRUD untuk Tim, Pemain, Pertandingan, pencatatan Gol, dan Laporan hasil pertandingan.
// @description                Semua endpoint GET bersifat publik. Endpoint POST/PUT/DELETE memerlukan header `X-API-Key` (lihat tombol Authorize di kanan atas).
// @host                        localhost:8080
// @BasePath                    /api/v1
// @schemes                     http https
// @securityDefinitions.apikey  ApiKeyAuth
// @in                          header
// @name                        X-API-Key
func main() {
	// Load .env
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, using environment variables")
	}

	cfg := config.Load()
	db := database.Connect(cfg)

	// ─── Repositories ─────────────────────────────────────────────────────────
	teamRepo   := teamModule.NewRepository(db)
	playerRepo := playerModule.NewRepository(db)
	matchRepo  := matchModule.NewRepository(db)
	goalRepo   := goalModule.NewRepository(db)

	// ─── Services ─────────────────────────────────────────────────────────────
	teamSvc   := teamModule.NewService(teamRepo)
	playerSvc := playerModule.NewService(playerRepo, teamSvc)
	matchSvc  := matchModule.NewService(matchRepo, teamSvc)
	goalSvc   := goalModule.NewService(goalRepo, matchSvc, playerSvc)
	reportSvc := reportModule.NewService(matchRepo)

	// ─── Handlers ─────────────────────────────────────────────────────────────
	teamHandler   := teamModule.NewHandler(teamSvc)
	playerHandler := playerModule.NewHandler(playerSvc)
	matchHandler  := matchModule.NewHandler(matchSvc)
	goalHandler   := goalModule.NewHandler(goalSvc)
	reportHandler := reportModule.NewHandler(reportSvc)

	// ─── Router ───────────────────────────────────────────────────────────────
	r := gin.Default()
	r.Use(gin.Recovery())

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Swagger UI — akses di http://localhost:8080/swagger/index.html
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	authMiddleware := middleware.APIKeyAuth(cfg.APIKey)

	v1 := r.Group("/api/v1")
	{
		teamHandler.RegisterRoutes(v1, authMiddleware)
		playerHandler.RegisterRoutes(v1, authMiddleware)
		matchHandler.RegisterRoutes(v1, authMiddleware)
		goalHandler.RegisterRoutes(v1, authMiddleware)
		reportHandler.RegisterRoutes(v1)
	}

	log.Printf("server running on :%s", cfg.Port)
	log.Printf("swagger UI: http://localhost:%s/swagger/index.html", cfg.Port)

	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
