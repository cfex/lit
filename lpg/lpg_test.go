package lpg

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	golightning "github.com/tracewayapp/go-lightning"
)

// Test structs for generic functions
type TestUser struct {
	Id        int
	FirstName string
	LastName  string
	Email     string
}

func TestRegister(t *testing.T) {
	// Clear any previous registrations
	delete(golightning.StructToFieldMap, reflect.TypeFor[TestUser]())

	Register[TestUser]()

	fieldMap, err := golightning.GetFieldMap(reflect.TypeFor[TestUser]())
	assert.NoError(t, err)
	assert.NotNil(t, fieldMap)

	// Verify basic structure
	assert.True(t, fieldMap.HasIntId)
	assert.Contains(t, fieldMap.ColumnKeys, "id")
	assert.Contains(t, fieldMap.ColumnKeys, "first_name")
	assert.Contains(t, fieldMap.ColumnKeys, "last_name")
	assert.Contains(t, fieldMap.ColumnKeys, "email")
}

func TestJoinForIn(t *testing.T) {
	tests := []struct {
		name     string
		ids      []int
		expected string
	}{
		{"empty", []int{}, ""},
		{"single", []int{1}, "1"},
		{"multiple", []int{1, 2, 3}, "1,2,3"},
		{"negative", []int{-1, 0, 1}, "-1,0,1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := JoinForIn(tt.ids)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestJoinStringForIn(t *testing.T) {
	tests := []struct {
		name     string
		offset   int
		params   []string
		expected string
	}{
		{"empty", 0, []string{}, ""},
		{"no offset", 0, []string{"a", "b"}, "$1,$2"},
		{"with offset", 2, []string{"a", "b"}, "$3,$4"},
		{"large offset", 10, []string{"x"}, "$11"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := JoinStringForIn(tt.offset, tt.params)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestPgInsertUpdateQueryGenerator_GenerateInsertQuery(t *testing.T) {
	gen := PgInsertUpdateQueryGenerator{}

	tests := []struct {
		name             string
		tableName        string
		columnKeys       []string
		hasIntId         bool
		expectedContains []string
		expectedColumns  []string
	}{
		{
			name:             "with int id",
			tableName:        "users",
			columnKeys:       []string{"id", "first_name", "last_name"},
			hasIntId:         true,
			expectedContains: []string{"INSERT INTO", "users", "DEFAULT", "RETURNING id", "$1", "$2"},
			expectedColumns:  []string{"first_name", "last_name"},
		},
		{
			name:             "without int id",
			tableName:        "products",
			columnKeys:       []string{"product_id", "name", "price"},
			hasIntId:         false,
			expectedContains: []string{"INSERT INTO", "products", "$1", "$2", "$3", "RETURNING id"},
			expectedColumns:  []string{"product_id", "name", "price"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query, columns := gen.GenerateInsertQuery(tt.tableName, tt.columnKeys, tt.hasIntId)

			for _, s := range tt.expectedContains {
				assert.Contains(t, query, s)
			}
			assert.Equal(t, tt.expectedColumns, columns)
		})
	}
}

func TestPgInsertUpdateQueryGenerator_GenerateUpdateQuery(t *testing.T) {
	gen := PgInsertUpdateQueryGenerator{}

	columnKeys := []string{"id", "first_name", "last_name"}
	query := gen.GenerateUpdateQuery("users", columnKeys)

	assert.Contains(t, query, "UPDATE users")
	assert.Contains(t, query, "SET")
	assert.Contains(t, query, "id = $1")
	assert.Contains(t, query, "first_name = $2")
	assert.Contains(t, query, "last_name = $3")
	assert.Contains(t, query, "WHERE")
}
