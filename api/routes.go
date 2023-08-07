package api

import (
	"github.com/gin-gonic/gin"
)

func (s *server) setupRoutes(rootDir string) {
	router := gin.Default()

	router.Static("/static", rootDir+"/client/build")
	router.POST("/api/crawlerProcesses", s.createCrawler)
	router.GET("/api/crawlerProcesses/:id", s.getCrawlerProcess)
	router.GET("/api/benefits/:cpf", s.getBenefits)

	s.router = router
}
