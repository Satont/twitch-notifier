package repository

import (
	"errors"

	"github.com/Masterminds/squirrel"
)

var Sq = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

var ErrBadQuery = errors.New("bad query")
