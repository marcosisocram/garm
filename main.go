package main

import (
	"context"
	"fmt"
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
		fmt.Println(err)
	}

	iniciarBanco := os.Getenv("INIT_TABLES")

	if iniciarBanco == "true" {
		fmt.Println("Criando tabelas")
		err := db.InitTables(context.Background())
		if err != nil {
			fmt.Println("Erro ao criar as tabelas")
			fmt.Println(err)
		} else {
			fmt.Println("Tabelas criadas")
		}
	}

	fmt.Println("Servidor rodando :" + port + "...")

	route.Run(":" + port)
}
