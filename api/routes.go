package api

import (
	"github.com/gin-gonic/gin"
)

func (s *server) setupRoutes() {
	router := gin.Default()

	router.POST("/api/findBenefits", s.findBenefits)

	s.router = router
}
