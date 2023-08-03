package api

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/leoomi/benefits-crawler/infra"
	"github.com/leoomi/benefits-crawler/models"
)

type findBenefitsRequest struct {
	CPF      string `json:"cpf" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (s *server) findBenefits(ctx *gin.Context) {
	var req findBenefitsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	process := models.CrawlerProcess{
		CPF:      req.CPF,
		Username: req.Username,
		Password: req.Password,
		State:    models.Created,
	}

	data, _ := json.Marshal(process)
	elsearchRes, err := s.elsearch.CreateIndex(infra.CrawlerProcessIndex, data)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := models.CrawlerProcessWithId{
		ID:             elsearchRes.Id,
		CrawlerProcess: process,
	}

	message, err := json.Marshal(res)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = s.rabbitmq.PublishMessage(message)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, res)
}
