package errors

type ErrValidacao string

func (e ErrValidacao) Error() string {
	return "Validacao: " + string(e)
}

type ErrLimite string

func (e ErrLimite) Error() string {
	return "Cliente sem limite: " + string(e)
}

type ErrSql string

func (e ErrSql) Error() string {
	return "Sql: " + string(e)
}
