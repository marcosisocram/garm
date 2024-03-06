package handlers

import (
	"fmt"
	"net/http"
	"regexp"
	"rinha-de-bk-go/db"
	"rinha-de-bk-go/errors"
	"strconv"

	"github.com/gin-gonic/gin"
)

func validar(transacao *db.Transacao) error {
	if transacao.Descricao == "" || len(transacao.Descricao) > 10 {
		return errors.ErrValidacao("Descricao")
	}

	pattern, err := regexp.Compile(`^[cd]$`)
	if err != nil {
		return err
	}

	if !pattern.MatchString(transacao.Tipo) {
		return errors.ErrValidacao("Tipo")
	}

	return nil
}

func PostTransacoes(ctx *gin.Context) {
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

	var transacao db.Transacao

	if err = ctx.ShouldBindJSON(&transacao); err != nil {
		fmt.Println("Não foi possivel receber body")
		fmt.Println(err)
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{})
		return
	}

	if err = validar(&transacao); err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{})
		return
	}

	cliente, err := db.InsertTransacao(id, &transacao)
	if err != nil {
		fmt.Println("Erro ao inserir...")
		fmt.Println(err)

		ctx.JSON(http.StatusUnprocessableEntity, gin.H{})
		return
	}

	ctx.JSON(http.StatusOK, cliente)
}
