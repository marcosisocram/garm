package db

import (
	"context"
	"rinha-de-bk-go/errors"
	"strconv"
)

type Transacao struct {
	Valor     int32  `json:"valor"`
	Tipo      string `json:"tipo"`
	Descricao string `json:"descricao"`
}

type Cliente struct {
	Saldo  int `json:"saldo"`
	Limite int `json:"limite"`
}

func GravarTransacao(clienteId int, transacao *Transacao) (Cliente, error) {
	conn, err := GetConnection()
	if err != nil {
		return Cliente{}, err
	}

	defer conn.Release()

	tx, err := conn.Begin(context.Background())

	defer tx.Rollback(context.Background())

	if transacao.Tipo == "d" {

		row := tx.QueryRow(context.Background(), "SELECT saldo_r, limite_r, status_r FROM debitar($1, $2, $3)", clienteId, transacao.Valor, transacao.Descricao)

		var status int
		var saldo int
		var limite int

		err = row.Scan(&saldo, &limite, &status)
		if err != nil {
			return Cliente{}, err
		}

		if status == 1 {
			return Cliente{}, errors.ErrSql("Cliente " + strconv.Itoa(clienteId) + " sem limite")
		}

		tx.Commit(context.Background())

		return Cliente{Saldo: saldo, Limite: limite}, nil

	}

	// Não sendo D então
	row := tx.QueryRow(context.Background(), "SELECT saldo_r, limite_r FROM creditar($1, $2, $3)", clienteId, transacao.Valor, transacao.Descricao)

	var saldo int
	var limite int

	err = row.Scan(&saldo, &limite)
	if err != nil {
		return Cliente{}, err
	}

	tx.Commit(context.Background())

	return Cliente{Saldo: saldo, Limite: limite}, nil
}
