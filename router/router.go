package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/kritpi/agnos-swe-assignment/docs" // generated swagger docs
	"github.com/kritpi/agnos-swe-assignment/handler"
	"github.com/kritpi/agnos-swe-assignment/middleware"
	"github.com/kritpi/agnos-swe-assignment/property"
)

// SetUpRouter registers all routes and the Swagger UI on the Gin engine.
func SetUpRouter(r *gin.Engine, h handler.Handler, staff middleware.StaffVerifier, cfg *property.Config) {
	// Swagger UI at /docs/api
	r.GET("/docs/api", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/docs/api/index.html")
	})
	r.GET("/docs/api/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Keep in sync with @BasePath in cmd/server/main.go.
	api := r.Group("/api")
	v1 := api.Group("/v1")

	// Staff
	staffGroup := v1.Group("/staff")
	{
		staffGroup.POST("/create", h.StaffCreate)
		staffGroup.POST("/login", h.StaffLogin)
	}

	// Patient
	patient := v1.Group("/patient")
	patient.Use(middleware.RequireStaffLogin(cfg.JWT.Secret, staff))
	{
		patient.POST("/search", h.PatientSearch)
	}
}
