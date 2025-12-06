package utils

import (
	"fmt"
	"log/slog"
	"reflect"
	"strings"

	"gorm.io/gorm"
)

func IsCleanDB(db *gorm.DB) bool {
	slog.Warn("IsCleanDB is not yet implemented for dialect", "dialect", db.Dialector.Name())
	return false
}

func HasConstraints(db *gorm.DB) bool {
	slog.Warn("HasForeignKeyConstraints is not yet implemented for dialect", "dialect", db.Dialector.Name())
	return false
}

func WhereNullable(query *gorm.DB, col string, val any) *gorm.DB {
	if val == nil || reflect.ValueOf(val).IsNil() {
		return query.Where(fmt.Sprintf("%s is null", col))
	}
	return query.Where(fmt.Sprintf("%s = ?", col), val)
}

func WithPaging(query *gorm.DB, limit, skip int) *gorm.DB {
	if limit >= 0 {
		query = query.Limit(limit)
	}
	if skip >= 0 {
		query = query.Offset(skip)
	}
	return query
}

type stringWriter struct {
	*strings.Builder
}

func (s stringWriter) WriteByte(c byte) error {
	return s.Builder.WriteByte(c)
}

func (s stringWriter) WriteString(str string) (int, error) {
	return s.Builder.WriteString(str)
}

// QuoteDBIdentifier quotes a column name used in a query.
func QuoteDBIdentifier(db *gorm.DB, identifier string) string {
	builder := stringWriter{Builder: &strings.Builder{}}
	db.Dialector.QuoteTo(builder, identifier)
	return builder.Builder.String()
}

// QuoteSql quotes a SQL statement with the given identifiers.
func QuoteSql(db *gorm.DB, queryTemplate string, identifiers ...string) string {
	quotedIdentifiers := make([]interface{}, len(identifiers))
	for i, identifier := range identifiers {
		quotedIdentifiers[i] = QuoteDBIdentifier(db, identifier)
	}
	return fmt.Sprintf(queryTemplate, quotedIdentifiers...)
}
