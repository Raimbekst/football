package main

import "carWash/internal/app"

const (
	configPath = "configs"
)

// @title Car Wash Bismillahirrahmanirrahim
// @version 2.0
// @description API Server for CarWash Application

// @host localhost:8080
// @BasePath /api/v1/

// @securityDefinitions.apikey User_Auth
// @in header
// @name Authorization

func main() {
	app.Run(configPath)
}
