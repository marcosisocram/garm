package main

import (
	"log"
	"os"
	"rinha-de-bk-go/db"
	"rinha-de-bk-go/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	route := gin.Default()
	route.GET("/clientes/:id/extrato", handlers.GetExtratos)
	route.POST("/clientes/:id/transacoes", handlers.PostTransacoes)

	port := os.Getenv("PORT")

	err := db.InitPool()
	if err != nil {
		log.Println(err)
	}

	log.Println("Servidor rodando :", port, "...")

	route.Run(":" + port)
}
