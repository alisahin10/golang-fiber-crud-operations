package main

import (
	"Fiber/database"
	"Fiber/handlers"
	"Fiber/middleware"
	"Fiber/repository"
	"Fiber/routes"
	"Fiber/services"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"os"
)

var logger *zap.Logger

func initLogger() {
	// Configure the JSON encoder
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "time"
	encoderConfig.LevelKey = "level"
	encoderConfig.NameKey = "logger"
	encoderConfig.MessageKey = "message"
	encoderConfig.StacktraceKey = "stacktrace"
	encoderConfig.LineEnding = zapcore.DefaultLineEnding

	// Create a new core
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(getLogFile()), // Use the log file for output
		zapcore.InfoLevel,
	)

	// Create a new logger with the core
	logger = zap.New(core)
}

func getLogFile() *os.File {
	// Create or open the log file
	file, err := os.OpenFile("app_logs.json", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Can't open log file: %v", err)
	}
	return file
}

func main() {
	// Load environment variables from the .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Initialize logger
	initLogger()
	defer logger.Sync()

	// Get the database path from the environment variable
	databasePath := os.Getenv("LOCAL_DATABASE_PATH")
	if databasePath == "" {
		log.Fatalf("LOCAL_DATABASE_PATH is not set in the environment")
	}

	// Initialize the database
	database.InitDB(databasePath)
	defer database.CloseDB()

	// Create the Fiber app
	app := fiber.New()

	// Use custom logging middleware
	app.Use(middleware.ZapLoggerMiddleware(logger))

	// Initialize the repository
	userRepo := repository.NewBuntDBUserRepository(database.DB, logger)

	// Initialize the service with the repository and logger
	userService := services.NewUserService(userRepo, logger)

	// Initialize the handler with the service and logger
	userHandler := handlers.NewUserHandler(userService, logger)

	// Set up routes
	routes.SetupRoutes(app, userHandler)

	// Start the server using the port defined in the environment
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "3000"
	}

	log.Fatal(app.Listen(":" + port))
}
