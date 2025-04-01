package main

import (
	"log"
	"os"
	"time"

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
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Database connection string
	dsn := "host=" + os.Getenv("DB_HOST") +
		" user=" + os.Getenv("DB_USER") +
		" password=" + os.Getenv("DB_PASSWORD") +
		" dbname=" + os.Getenv("DB_NAME") +
		" port=" + os.Getenv("DB_PORT") +
		" sslmode=disable"

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

	// Routes
	r.POST("/api/track", trackDependency)
	r.GET("/api/dependencies", getDependencies)

	// Start server
	r.Run(":8080")
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
