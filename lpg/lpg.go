package lpg

import (
	"database/sql"
	"errors"
	"reflect"
	"strconv"
	"strings"

	golightning "github.com/tracewayapp/go-lightning"
)

var NamingStrategy = golightning.DefaultDbNamingStrategy{}

type PgInsertUpdateQueryGenerator struct{}

func (PgInsertUpdateQueryGenerator) GenerateInsertQuery(tableName string, columnKeys []string, hasIntId bool) (string, []string) {

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

	counter := 1
	insertColumns := []string{}
	for i, k := range columnKeys {
		if hasIntId && k == "id" {
			insertQuery.WriteString("DEFAULT")
			if i != totalKeys-1 {
				insertQuery.WriteString(",")
			}
		} else {
			insertColumns = append(insertColumns, k)
			insertQuery.WriteString("$" + strconv.Itoa(counter))
			if i != totalKeys-1 {
				insertQuery.WriteString(",")
			}

			counter++
		}
	}
	insertQuery.WriteString(") RETURNING id")

	return insertQuery.String(), insertColumns
}

func (PgInsertUpdateQueryGenerator) GenerateUpdateQuery(tableName string, columnKeys []string) string {
	var updateQuery strings.Builder
	updateQuery.WriteString("UPDATE ")
	updateQuery.WriteString(tableName)
	updateQuery.WriteString(" SET ")

	totalKeys := len(columnKeys)
	for i, k := range columnKeys {
		updateQuery.WriteString(k)
		updateQuery.WriteString(" = $" + strconv.Itoa(i+1))
		if i != totalKeys-1 {
			updateQuery.WriteString(",")
		}
	}

	updateQuery.WriteString(" WHERE ")

	return updateQuery.String()
}

var pgInsertUpdateQueryGenerator = PgInsertUpdateQueryGenerator{}

func Register[T any]() {
	golightning.Register[T](NamingStrategy, pgInsertUpdateQueryGenerator)
}

func SelectGenericSingle[T any](tx *sql.Tx, query string, args ...any) (*T, error) {
	return golightning.SelectGenericSingle[T](tx, query, args...)
}
func SelectGeneric[T any](tx *sql.Tx, query string, args ...any) ([]*T, error) {
	return golightning.SelectGeneric[T](tx, query, args...)
}

func InsertGeneric[T any](tx *sql.Tx, t *T) (int, error) {
	// due to the way the driver works we need to read the result set of the query row with insert
	// thus we cannot use the generic version of insert which uses Exec and LastInsertId
	tType := reflect.TypeOf(*t)
	fieldMap, err := golightning.GetFieldMap(tType)

	if err != nil {
		return 0, err
	}

	result := tx.QueryRow(
		fieldMap.InsertQuery,
		*golightning.GetPointersForColumns[T](fieldMap.InsertColumns, fieldMap, t)...,
	)

	var id int
	err = result.Scan(&id)

	if err != nil {
		return 0, err
	}

	return int(id), nil
}
func InsertGenericExistingUuid[T any](tx *sql.Tx, t *T) error {
	return golightning.InsertGenericExistingUuid[T](tx, t)
}

func InsertGenericUuid[T any](tx *sql.Tx, t *T) (string, error) {
	return golightning.InsertGenericUuid[T](tx, t)
}

func UpdateGeneric[T any](tx *sql.Tx, t *T, where string, args ...any) error {
	// we cannot use generic version for this one because we need to replace $ identifiers with reindexed values
	// note: this function won't work if they try to write $1 inside of a string - too bad
	if len(where) == 0 {
		return errors.New("parameter 'where' was not present")
	}
	tType := reflect.TypeOf(*t)
	fieldMap, err := golightning.GetFieldMap(tType)
	if err != nil {
		return err
	}

	params := append(*golightning.GetPointersForColumns(fieldMap.ColumnKeys, fieldMap, t), args...)

	if strings.Contains(where, "$") {
		newWhere := strings.Builder{}

		offset := strings.Count(fieldMap.UpdateQuery, "$")
		parsingIdentifier := false
		for _, c := range where {
			if c == '$' {
				parsingIdentifier = true
				newWhere.WriteRune(c)
			} else if parsingIdentifier {
				if c >= '0' && c <= '9' {
					continue
				} else {
					parsingIdentifier = false
					newWhere.WriteString(strconv.Itoa(offset + 1))
					newWhere.WriteRune(c)
					offset += 1
				}
			} else {
				newWhere.WriteRune(c)
			}
		}
		if parsingIdentifier {
			newWhere.WriteString(strconv.Itoa(offset + 1))
		}

		where = newWhere.String()
	}
	_, err = tx.Exec(
		fieldMap.UpdateQuery+where,
		params...,
	)
	if err != nil {
		return err
	}

	return nil
}

func SelectMultiple[T any](tx *sql.Tx, mapLine func(*sql.Rows, *T) error, query string, args ...any) ([]*T, error) {
	return golightning.SelectMultiple(tx, mapLine, query, args...)
}
func SelectSingle[T any](tx *sql.Tx, mapLine func(*sql.Rows, *T) error, query string, args ...any) (*T, error) {
	return golightning.SelectSingle(tx, mapLine, query, args...)
}
func Insert(tx *sql.Tx, query string, args ...any) (int, error) {
	result := tx.QueryRow(
		query,
		args...,
	)

	var id int
	err := result.Scan(&id)

	if err != nil {
		return 0, err
	}

	return int(id), nil
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
		sb.WriteString("$" + strconv.Itoa(index+1+offset))
		if index < len(params)-1 {
			sb.WriteString(",")
		}
	}
	return sb.String()
}
func JoinForIn(ids []int) string {
	return golightning.JoinForIn(ids)
}
