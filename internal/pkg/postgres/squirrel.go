package postgres

import (
	"fmt"
	"strings"
	"text/template"

	sq "github.com/Masterminds/squirrel"
)

// Squirrel provides wrapper around squirrel package
type Squirrel struct {
	Builder sq.StatementBuilderType
}

func NewSquirrel() *Squirrel {
	return &Squirrel{sq.StatementBuilder.PlaceholderFormat(sq.Dollar)}
}

func (s *Squirrel) Equal(key string, value interface{}) sq.Eq {
	return sq.Eq{key: value}
}

func (s *Squirrel) EqualStr(key string) EqualStr {
	return EqualStr(key)
}

func (s *Squirrel) ILike(key string, value interface{}) sq.ILike {
	return sq.ILike{key: value}
}

func (s *Squirrel) NotEqual(key string, value interface{}) sq.NotEq {
	return sq.NotEq{key: value}
}

func (s *Squirrel) Or(cond ...sq.Sqlizer) sq.Or {
	sl := make([]sq.Sqlizer, 0, len(cond))
	sl = append(sl, cond...)
	return sl
}

func (s *Squirrel) And(cond ...sq.Sqlizer) sq.And {
	sl := make([]sq.Sqlizer, 0, len(cond))
	sl = append(sl, cond...)
	return sl
}

func (s *Squirrel) Alias(expr sq.Sqlizer, alias string) sq.Sqlizer {
	return sq.Alias(expr, alias)
}

func (s *Squirrel) EqualMany(clauses map[string]interface{}) sq.Eq {
	eqMany := make(sq.Eq, len(clauses))
	for key, value := range clauses {
		eqMany[key] = value
	}
	return eqMany
}

func (s *Squirrel) Gt(key string, value interface{}) sq.Gt {
	return sq.Gt{key: value}
}

func (s *Squirrel) Lt(key string, value interface{}) sq.Lt {
	return sq.Lt{key: value}
}

func (s *Squirrel) Expr(sql string, args ...interface{}) sq.Sqlizer {
	return sq.Expr(sql, args)
}

func (s *Squirrel) JSONPathWhere(fieldName, jsonbOp, searchField, value string) (string, error) {
	var b strings.Builder
	value = template.HTMLEscapeString(value)
	_, err := fmt.Fprintf(&b, "%s %s? '$.%s ?? (@ == %q)'", fieldName, jsonbOp, searchField, value)
	if err != nil {
		return "", fmt.Errorf("in JSONPathWhere squirrel method: %w", err)
	}

	return b.String(), nil
}

type EqualStr string

func (e EqualStr) ToSql() (sql string, args []interface{}, err error) {
	sql = string(e)
	return
}
