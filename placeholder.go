package sqx

import sq "github.com/stytchauth/squirrel"

// PlaceholderFormat is the interface that wraps the ReplacePlaceholders method.
//
// ReplacePlaceholders takes a SQL statement and replaces each question mark
// placeholder with a (possibly different) SQL placeholder.
type PlaceholderFormat sq.PlaceholderFormat

var (
	// Question is a PlaceholderFormat instance that leaves placeholders as
	// question marks.
	Question = sq.Question

	// Dollar is a PlaceholderFormat instance that replaces placeholders with
	// dollar-prefixed positional placeholders (e.g. $1, $2, $3).
	Dollar = sq.Dollar

	// Colon is a PlaceholderFormat instance that replaces placeholders with
	// colon-prefixed positional placeholders (e.g. :1, :2, :3).
	Colon = sq.Colon

	// AtP is a PlaceholderFormat instance that replaces placeholders with
	// "@p"-prefixed positional placeholders (e.g. @p1, @p2, @p3).
	AtP = sq.AtP

	defaultPlaceholderFormat sq.PlaceholderFormat = Question
)

func SetPlaceholder(placeholderFormat PlaceholderFormat) {
	defaultPlaceholderFormat = placeholderFormat
}

func SetPostgres() {
	SetPlaceholder(Dollar)
}
