package api

import (
	"github.com/gin-gonic/gin"
)

func (s *server) setupRoutes() {
	router := gin.Default()

	router.POST("/api/crawlerProcesses", s.createCrawler)
	router.GET("/api/crawlerProcesses/:id", s.getCrawlerProcess)

	s.router = router
}
