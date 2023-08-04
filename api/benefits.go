package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type getBenefitsRequest struct {
	CPF string `json:"cpf" binding:"required"`
}

func (s *server) getBenefits(ctx *gin.Context) {
	fmt.Println("teste")
}
