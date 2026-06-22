package app

import (
	"github.com/Prateet-Github/streamit-api/internal/routes"
	"github.com/gin-gonic/gin"
)

func New() *gin.Engine {
	router := gin.Default()

	routes.RegisterHealthRoutes(router)

	return router
}
