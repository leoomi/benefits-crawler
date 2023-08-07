package api

import (
	"os"

	"github.com/gin-gonic/gin"
)

func (s *server) setupRoutes(rootDir string) {
	router := gin.Default()

	router.StaticFile("/", rootDir+"/client/build/index.html")
	files, _ := os.ReadDir(rootDir + "/client/build")
	for _, file := range files {
		if file.Type().IsDir() || file.Name() == "index.html" {
			continue
		}

		router.StaticFile("/"+file.Name(), rootDir+"/client/build/"+file.Name())
	}
	router.Static("/static", rootDir+"/client/build/static")
	router.POST("/api/crawlerProcesses", s.createCrawler)
	router.GET("/api/crawlerProcesses/:id", s.getCrawlerProcess)
	router.GET("/api/benefits/:cpf", s.getBenefits)

	s.router = router
}
