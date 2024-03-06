package db

import (
	"context"
	"fmt"
	"rinha-de-bk-go/errors"
	"strconv"

	crdbpgx "github.com/cockroachdb/cockroach-go/v2/crdb/crdbpgxv5"
	"github.com/jackc/pgx/v5"
	pkgerrors "github.com/pkg/errors"
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

func InsertTransacao(clienteId int, transacao *Transacao) (Cliente, error) {
	conn, err := GetConnection()
	if err != nil {
		return Cliente{}, err
	}

	defer conn.Release()

	tx, err := conn.Begin(context.Background())

	defer tx.Rollback(context.Background())

	// Connect to redis.
	//	client := redis.NewClient(&redis.Options{
	//		Network: "tcp",
	//		Addr:    "127.0.0.1:6379",
	//	})
	//	defer client.Close()

	// Create a new lock client.
	//	locker := redislock.New(client)

	//	ctx := context.Background()

	// Retry every 100ms, for up-to 3x
	// backoff := redislock.LimitRetry(redislock.LinearBackoff(100*time.Millisecond), 6)

	// Try to obtain lock.
	//lock, err := locker.Obtain(ctx, strconv.Itoa(clienteId), 2000*time.Millisecond, &redislock.Options{
	//	RetryStrategy: backoff,
	//})
	//if err == redislock.ErrNotObtained {
	//	fmt.Println("Could not obtain lock!")
	//} else if err != nil {
	//	log.Fatalln(err)
	//}

	// Don't forget to defer Release.
	// defer lock.Release(ctx)
	// fmt.Println("I have a lock!")

	row := tx.QueryRow(context.Background(), "SELECT saldo, limite FROM clientes WHERE id = $1", clienteId)

	var saldo int32
	var limite int32

	err = row.Scan(&saldo, &limite)
	if err != nil {
		fmt.Println("Erro no select: ")
		fmt.Println(err)
		return Cliente{}, errors.ErrSql("Select")
	}

	novoSaldo := saldo + transacao.Valor

	// Validar tipo e saldo
	if transacao.Tipo == "d" {
		novoSaldo = saldo - transacao.Valor
		if novoSaldo < (limite * -1) {
			return Cliente{}, errors.ErrLimite(strconv.Itoa(clienteId))
		}
	}

	_, err = tx.Exec(context.Background(), "INSERT INTO transacoes (cliente_id, tipo, valor, descricao) values ($1, $2, $3, $4)", clienteId, transacao.Tipo, transacao.Valor, transacao.Descricao)
	if err != nil {
		fmt.Println(err)
		return Cliente{}, errors.ErrSql("Insert")
	}

	tx.Commit(context.Background())

	err = crdbpgx.ExecuteTx(context.Background(), conn, pgx.TxOptions{}, func(txx pgx.Tx) error {
		// TODO: continuar colocando o update aqui

		_, err = txx.Exec(context.Background(), "UPDATE clientes SET saldo = $1 WHERE id = $2", novoSaldo, clienteId)
		if err != nil {
			return pkgerrors.Wrap(err, "Update")
		}

		//_, err = txx.Exec(context.Background(), "DELETE FROM transacoes WHERE id NOT IN ( SELECT id FROM transacoes WHERE cliente_id = $1 ORDER BY realizada_em DESC LIMIT 10 )", clienteId)
		//if err != nil {
		//	return pkgerrors.Wrap(err, "Delete")
		//}
		//_, err = tx.Exec(ctx, "INSERT INTO transacoes (cliente_id, tipo, valor, descricao) values ($1, $2, $3, $4)", clienteId, transacao.Tipo, transacao.Valor, transacao.Descricao)
		//if err != nil {
		//	return pkgerrors.Wrap(err, "Insert")
		//}
		return nil
	})
	if err != nil {
		fmt.Println(err)
		return Cliente{}, errors.ErrSql("Ferrou!")
	}

	//	_, err = tx.Exec(context.Background(), "UPDATE clientes SET saldo = $1 WHERE id = $2", novoSaldo, clienteId)
	//	if err != nil {
	//		fmt.Println(err)
	//		return Cliente{}, errors.ErrSql("Update")
	//	}

	//_, err = tx.Exec(context.Background(), "INSERT INTO transacoes (cliente_id, tipo, valor, descricao) values ($1, $2, $3, $4)", clienteId, transacao.Tipo, transacao.Valor, transacao.Descricao)
	//if err != nil {
	//	fmt.Println(err)
	//	return Cliente{}, errors.ErrSql("Insert")
	//}

	//	tx.Commit(context.Background())

	return Cliente{Saldo: int(novoSaldo), Limite: int(limite)}, nil
}
