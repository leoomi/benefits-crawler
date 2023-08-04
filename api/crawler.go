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

type getCrawlerProcessRequest struct {
	ID string `uri:"id" binding:"required"`
}

func (s *server) getCrawlerProcess(ctx *gin.Context) {
	var req getCrawlerProcessRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	data, err := s.elsearch.GetDocument(infra.CrawlerProcessIndex, req.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	docRes := infra.DocResponse[models.CrawlerProcess]{}
	json.Unmarshal(data, &docRes)

	if !docRes.Found {
		ctx.JSON(http.StatusNotFound, nil)
		return
	}

	ctx.JSON(http.StatusOK, docRes.Source)
}
