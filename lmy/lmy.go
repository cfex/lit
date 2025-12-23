package lmy

import (
	"database/sql"
	"strings"

	golightning "github.com/tracewayapp/go-lightning"
)

var NamingStrategy = golightning.DefaultDbNamingStrategy{}

type MySqlInsertUpdateQueryGenerator struct{}

func (MySqlInsertUpdateQueryGenerator) GenerateInsertQuery(tableName string, columnKeys []string, hasIntId bool) (string, []string) {

	var insertQuery strings.Builder

	insertQuery.WriteString("INSERT INTO ")
	insertQuery.WriteString(" ")
	insertQuery.WriteString(tableName)
	insertQuery.WriteString(" ")
	insertQuery.WriteString("(")

	totalKeys := len(columnKeys)
	for i, k := range columnKeys {
		insertQuery.WriteString(k)
		if i != totalKeys-1 {
			insertQuery.WriteString(",")
		}
	}

	insertQuery.WriteString(") VALUES (")

	insertColumns := []string{}
	for i, k := range columnKeys {
		if hasIntId && k == "id" {
			insertQuery.WriteString("NULL")
			if i != totalKeys-1 {
				insertQuery.WriteString(",")
			}
		} else {
			insertColumns = append(insertColumns, k)
			insertQuery.WriteString("?")
			if i != totalKeys-1 {
				insertQuery.WriteString(",")
			}
		}
	}
	insertQuery.WriteString(")")

	return insertQuery.String(), insertColumns
}

func (MySqlInsertUpdateQueryGenerator) GenerateUpdateQuery(tableName string, columnKeys []string) string {
	var updateQuery strings.Builder
	updateQuery.WriteString("UPDATE ")
	updateQuery.WriteString(tableName)
	updateQuery.WriteString(" SET ")

	totalKeys := len(columnKeys)
	for i, k := range columnKeys {
		updateQuery.WriteString(k)
		updateQuery.WriteString(" = ?")
		if i != totalKeys-1 {
			updateQuery.WriteString(",")
		}
	}

	updateQuery.WriteString(" WHERE ")

	return updateQuery.String()
}

var mySqlInsertUpdateQueryGenerator = MySqlInsertUpdateQueryGenerator{}

func Register[T any]() {
	golightning.Register[T](NamingStrategy, mySqlInsertUpdateQueryGenerator)
}

func SelectGenericSingle[T any](tx *sql.Tx, query string, args ...any) (*T, error) {
	return golightning.SelectGenericSingle[T](tx, query, args...)
}
func SelectGeneric[T any](tx *sql.Tx, query string, args ...any) ([]*T, error) {
	return golightning.SelectGeneric[T](tx, query, args...)
}
func InsertGeneric[T any](tx *sql.Tx, t *T) (int, error) {
	return golightning.InsertGeneric[T](tx, t)
}
func InsertGenericExistingUuid[T any](tx *sql.Tx, t *T) error {
	return golightning.InsertGenericExistingUuid[T](tx, t)
}

func InsertGenericUuid[T any](tx *sql.Tx, t *T) (string, error) {
	return golightning.InsertGenericUuid[T](tx, t)
}

func UpdateGeneric[T any](tx *sql.Tx, t *T, where string, args ...any) error {
	return golightning.UpdateGeneric[T](tx, t, where, args...)
}

func SelectMultiple[T any](tx *sql.Tx, mapLine func(*sql.Rows, *T) error, query string, args ...any) ([]*T, error) {
	return golightning.SelectMultiple(tx, mapLine, query, args...)
}
func SelectSingle[T any](tx *sql.Tx, mapLine func(*sql.Rows, *T) error, query string, args ...any) (*T, error) {
	return golightning.SelectSingle(tx, mapLine, query, args...)
}
func Insert(tx *sql.Tx, query string, args ...any) (int, error) {
	return golightning.Insert(tx, query, args...)
}
func Update(tx *sql.Tx, query string, args ...any) error {
	return golightning.Update(tx, query, args...)
}
func Delete(tx *sql.Tx, query string, args ...any) error {
	return golightning.Delete(tx, query, args...)
}
func JoinStringForIn(offset int, params []string) string {
	var sb strings.Builder
	for index := range params {
		sb.WriteString("?")
		if index < len(params)-1 {
			sb.WriteString(",")
		}
	}
	return sb.String()
}
func JoinForIn(ids []int) string {
	return golightning.JoinForIn(ids)
}
