package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

var db *pgxpool.Pool

func InitPool() error {
	config, err := pgxpool.ParseConfig("postgresql://root@127.0.0.1:26257/defaultdb?sslmode=disable")
	if err != nil {
		return err
	}

	config.MaxConns = 16
	config.MinConns = 13

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return err
	}

	db = pool
	return nil
}

func GetConnection() (*pgxpool.Conn, error) {
	if db == nil {
		InitPool()
	}

	conn, err := db.Acquire(context.Background())
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func InitTables(ctx context.Context) error {
	conn, err := GetConnection()

	defer conn.Release()

	if err != nil {
		return err
	}

	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}

	defer tx.Commit(ctx)

	//	if _, err := tx.Exec(ctx, "SET CLUSTER SETTING sql.txn.read_committed_isolation.enabled = 'true'"); err != nil {
	//		return err
	//	}

	if _, err := tx.Exec(ctx, "DROP TABLE IF EXISTS transacoes"); err != nil {
		return err
	}

	if _, err := tx.Exec(ctx, "DROP TABLE IF EXISTS clientes"); err != nil {
		return err
	}

	if _, err := tx.Exec(ctx, `CREATE TABLE clientes (
		  id INT2 PRIMARY KEY,
		  limite INT NOT NULL,
		  saldo INT NOT NULL DEFAULT 0
		) `); err != nil {
		return err
	}

	if _, err := tx.Exec(ctx, `
		CREATE TABLE transacoes (
			id SERIAL PRIMARY KEY,
			cliente_id INT2 NOT NULL REFERENCES clientes(id),
			valor INT4 NOT NULL,
			tipo STRING(1) NOT NULL,
			descricao STRING(10) NOT NULL,
			realizada_em TIMESTAMP NOT NULL DEFAULT NOW(),
			INDEX(cliente_id),
			INDEX(realizada_em desc)
		)
		`); err != nil {
		return err
	}

	if _, err := tx.Exec(ctx, "DELETE FROM transacoes WHERE 1 = 1"); err != nil {
		return err
	}

	if _, err := tx.Exec(ctx, "DELETE FROM clientes WHERE 1 = 1"); err != nil {
		return err
	}

	if _, err := tx.Exec(ctx, `
		insert into clientes (id, limite)
			values (1, 100000),
      	 (2, 80000),
      	(3, 1000000),
      	(4, 10000000),
      	(5, 500000)
		`); err != nil {
		return err
	}

	return nil
}
