package handlers

import (
	"log"
	"net/http"
	"rinha-de-bk-go/db"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetExtratos(ctx *gin.Context) {
	idStr := ctx.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{})
		return
	}

	if id > 5 {
		ctx.JSON(http.StatusNotFound, gin.H{})
		return
	}

	extrato, err := db.RecuperarExtratos(id)
	if err != nil {
		log.Println("Erro ao recuperar extratos", err)
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{})
		return
	}

	ctx.JSON(http.StatusOK, extrato)
}
