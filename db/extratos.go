package db

import (
	"context"
	"time"
)

type Saldo struct {
	Total       int       `json:"total"`
	Limite      int       `json:"limite"`
	DataExtrato time.Time `json:"data_extrato"`
}

type UltimaTransacao struct {
	Valor       int       `json:"valor"`
	Tipo        string    `json:"tipo"`
	Descricao   string    `json:"descricao"`
	RealizadaEm time.Time `json:"realizada_em"`
}

type Extrato struct {
	Saldo             Saldo             `json:"saldo"`
	UltimasTransacoes []UltimaTransacao `json:"ultimas_transacoes"`
}

func RecuperarExtratos(clientId int) (Extrato, error) {
	conn, err := GetConnection()
	if err != nil {
		return Extrato{}, err
	}

	defer conn.Release()

	tx, err := conn.Begin(context.Background())
	if err != nil {
		return Extrato{}, err
	}

	defer tx.Commit(context.Background())

	row := tx.QueryRow(context.Background(), "SELECT saldo, limite, now() FROM clientes WHERE id = $1", clientId)

	extrato := Extrato{}

	saldo := Saldo{}

	err = row.Scan(&saldo.Total, &saldo.Limite, &saldo.DataExtrato)
	if err != nil {
		return Extrato{}, err
	}

	extrato.Saldo = saldo

	rows, err := tx.Query(context.Background(), "SELECT valor, tipo, descricao, realizada_em FROM transacoes WHERE cliente_id = $1 ORDER BY realizada_em desc limit 10", clientId)
	if err != nil {
		return Extrato{}, err
	}

	extrato.UltimasTransacoes = []UltimaTransacao{}

	for rows.Next() {
		var valor int
		var tipo string
		var descricao string
		var realizadaEm time.Time

		err := rows.Scan(&valor, &tipo, &descricao, &realizadaEm)
		if err != nil {
			return Extrato{}, err
		}

		extrato.UltimasTransacoes = append(extrato.UltimasTransacoes, UltimaTransacao{Valor: valor, Tipo: tipo, Descricao: descricao, RealizadaEm: realizadaEm})

	}

	if err = rows.Err(); err != nil {
		return Extrato{}, err
	}

	return extrato, nil
}
