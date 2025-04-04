package main

import (
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Dependency struct
type Dependency struct {
	ID        uint      `gorm:"primaryKey"`
	Service1  string    `json:"service_1" gorm:"column:service_1"`
	Service2  string    `json:"service_2" gorm:"column:service_2"`
	Method    string    `json:"method"`
	CreatedAt time.Time `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
}

// Database instance
var db *gorm.DB

func initDB() {
	// Only load .env in local environment
	if _, exists := os.LookupEnv("RENDER"); !exists {
		err := godotenv.Load()
		if err != nil {
			log.Println("⚠️ Warning: No .env file found, using system environment variables.")
		}
	}

	// Database connection string
	dsn := "host=" + os.Getenv("DB_HOST") +
		" user=" + os.Getenv("DB_USER") +
		" password=" + os.Getenv("DB_PASSWORD") +
		" dbname=" + os.Getenv("DB_NAME") +
		" port=" + os.Getenv("DB_PORT") +
		" sslmode=require target_session_attrs=read-write"

	// Connect to database
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database")
	}

	// Auto Migrate
	database.AutoMigrate(&Dependency{})
	db = database
}

func main() {
	// Initialize DB
	initDB()

	// Start Gin server
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://microviz.netlify.app/"}, 
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Default route ("/") for the backend
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to the MicroViz's Backend!",
			"status":  "running",
		})
	})

	// Routes
	r.POST("/api/track", trackDependency)
	r.GET("/api/dependencies", getDependencies)

	// Get port from environment variables (for Render)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default for local development
	}

	// Start server
	r.Run(":" + port)

}

func trackDependency(c *gin.Context) {
	var dependency Dependency
	if err := c.ShouldBindJSON(&dependency); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Save to DB
	db.Create(&dependency)

	c.JSON(200, gin.H{"message": "Dependency tracked successfully"})
}

func getDependencies(c *gin.Context) {
	var dependencies []Dependency
	db.Find(&dependencies)
	c.JSON(200, dependencies)
}
