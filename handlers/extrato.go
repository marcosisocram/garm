package handlers

import (
	"fmt"
	"net/http"
	"rinha-de-bk-go/db"
	"strconv"

	"github.com/gin-gonic/gin"
)

type saldo struct {
	Total       int    `json:"total"`
	Limite      int    `json:"limite"`
	DataExtrato string `json:"data_extrato"`
}

type ultimaTransacao struct {
	Valor       int    `json:"valor"`
	Tipo        string `json:"tipo"`
	Descricao   string `json:"descricao"`
	RealizadoEm string `json:"realizado_em"`
}

type extrato struct {
	Saldo             saldo             `json:"saldo"`
	UltimasTransacoes []ultimaTransacao `json:"ultimas_transacoes"`
}

var ultimasTransacoes = []ultimaTransacao{
	{Valor: 0, Tipo: "c", Descricao: "descri", RealizadoEm: "29392893929"},
	{Valor: 0, Tipo: "d", Descricao: "cricao", RealizadoEm: "83823832328"},
}

var extratos = extrato{
	Saldo:             saldo{Total: 0, Limite: 10000, DataExtrato: "3939393"},
	UltimasTransacoes: ultimasTransacoes,
}

func GetExtratos(ctx *gin.Context) {
	idStr := ctx.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{})
		return
	}

	if id > 5 {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{})
		return
	}

	tt, err := db.RecuperaExtratos(id)
	if err != nil {
		// FIX: Ferrou aqui
		fmt.Println(err)
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{})
		return
	}

	ctx.JSON(http.StatusOK, tt) // IndentedJSON
}
