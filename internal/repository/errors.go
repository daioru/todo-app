package repository

import "errors"

var ErrNoRowsUpdated error = errors.New("no rows affected")
