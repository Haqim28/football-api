package domain

import "errors"

// ErrNotFound adalah helper untuk membuat error not found yang bisa dicek dengan errors.As
type ErrNotFound string

func (e ErrNotFound) Error() string {
	return string(e)
}

var ErrRecordNotFound = errors.New("record tidak ditemukan")
