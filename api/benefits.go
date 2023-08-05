package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/leoomi/benefits-crawler/infra"
	"github.com/leoomi/benefits-crawler/models"
)

type getBenefitsRequest struct {
	CPF string `json:"cpf" binding:"required"`
}

func (s *server) getBenefits(ctx *gin.Context) {
	var req getBenefitsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var benefits models.Benefits
	err := s.elsearch.SearchSingleDocument(infra.BenefitsIndex, "cpf", req.CPF, &benefits)

	if err != nil {
		if err == infra.ErrESNotFound {
			ctx.JSON(http.StatusNotFound, nil)
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, benefits)
}
