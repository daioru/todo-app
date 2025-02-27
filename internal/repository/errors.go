package repository

import "errors"

var ErrNoRowsUpdated error = errors.New("no rows affected")
var ErrUniqueUser error = errors.New("username already exists")
var ErrUserNotFound = errors.New("username not found")
