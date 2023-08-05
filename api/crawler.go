package api

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/leoomi/benefits-crawler/infra"
	"github.com/leoomi/benefits-crawler/models"
)

type createCrawlerRequest struct {
	CPF      string `json:"cpf" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (s *server) createCrawler(ctx *gin.Context) {
	var req createCrawlerRequest
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

	id, err := s.elsearch.CreateDocument(infra.CrawlerProcessIndex, process)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := models.CrawlerProcessWithId{
		ID:             id,
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

type getCrawlerProcessRequest struct {
	ID string `uri:"id" binding:"required"`
}

func (s *server) getCrawlerProcess(ctx *gin.Context) {
	var req getCrawlerProcessRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var process models.CrawlerProcess
	err := s.elsearch.GetDocument(infra.CrawlerProcessIndex, req.ID, &process)
	if err != nil {
		if err == infra.ErrESNotFound {
			ctx.JSON(http.StatusNotFound, nil)
		}
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, process)
}
