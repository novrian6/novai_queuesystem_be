package main

import (
	"log"
	"os"
	"queue-system-backend/database"
	"queue-system-backend/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize database
	database.ConnectDB()

	// Ensure database is properly initialized
	if database.DB == nil {
		log.Fatal("❌ Database connection is not initialized.")
	}

	// Initialize controllers
	//statsController := controllers.NewStatisticsController(database.DB)

	// Initialize Gin
	r := gin.Default()
	//r.Use(cors.Default())

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:8081", "http://localhost", "https://reqbin.com/"}, // Allow Vue frontend
		AllowMethods:     []string{"GET", "POST", "OPTIONS", "PUT", "DELETE"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	// Register routes
	routes.AuthRoutes(r)
	routes.UserRoutes(r)

	// Register Routes
	routes.VenueRoutes(r)

	//Register Company Routes
	routes.CompanyRoutes(r)

	//Register Role Routes
	routes.RoleRoutes(r)

	//Register Role Routes
	routes.ServiceRoutes(r)

	// Add queue ticket routes
	routes.QueueTicketRoutes(r)

	//RegisterViewTicketRoutes
	//routes.QueueTicketPublicRoutes(r) // Added the new public routes
	routes.RegisterViewTicketRoutes(r)

	// Add counter   routes
	routes.RegisterCounterRoutes(r)

	// Register Routes
	routes.RegisterQueueDisplayRoutes(r)

	// Register User Counter Map Routes
	routes.RegisterUserCounterMapRoutes(r)

	// Setup statistics routes
	routes.SetupStatisticsRoutes(r)

	// Define port (with fallback to default port 8081)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	// Start server
	log.Printf("🚀 Server is running on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}
