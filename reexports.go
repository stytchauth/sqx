package sqx

import sq "github.com/stytchauth/squirrel"

// And  represents a SQL AND expression. It is a re-export of squirrel.And.
type And = sq.And

// Eq represents a SQL equality = expression. It is a re-export of squirrel.Eq.
type Eq = sq.Eq

// NotEq represents a SQL inequality <> expression. It is a re-export of squirrel.NotEq.
type NotEq = sq.NotEq

// Or represents a SQL OR expression. It is a re-export of squirrel.Or.
type Or = sq.Or

// Sqlizer is an interface containing the ToSql method. It is a re-export of squirrel.Sqlizer.
type Sqlizer = sq.Sqlizer

// Gt represents a SQL > expression. It is a re-export of squirrel.Gt.
type Gt = sq.Gt

// GtOrEq represents a SQL >= expression. It is a re-export of squirrel.GtOrEq.
type GtOrEq = sq.GtOrEq

// Lt represents a SQL < expression. It is a re-export of squirrel.Lt.
type Lt = sq.Lt

// LtOrEq represents a SQL <= expression. It is a re-export of squirrel.LtOrEq.
type LtOrEq = sq.LtOrEq
