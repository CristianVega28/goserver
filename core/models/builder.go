package models

import (
	"database/sql"
	"fmt"
	"strings"
)

type (
	Builder struct {
		conn    *sql.DB
		table   string
		str     strings.Builder
		columns []string
	}
)

func (builder *Builder) Where(column string, value string) *Builder {

	query := fmt.Sprintf("WHERE %s = %s ", column, value)
	builder.str.WriteString(query)

	return builder
}

func (builder *Builder) Or(column string, operator string, value string) *Builder {

	query := fmt.Sprintf("OR %s %s %s ", column, operator, value)
	builder.str.WriteString(query)

	return builder
}

func (builder *Builder) Columns(columns []string) *Builder {
	// builder.str.
}
