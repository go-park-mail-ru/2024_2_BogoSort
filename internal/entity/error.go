package entity

import "errors"

func RedisWrap(errs ...error) error {
	return errors.Join(ErrRedis, errors.Join(errs...))
}

func UsecaseWrap(errs ...error) error {
	return errors.Join(ErrInternal, errors.Join(errs...))
}

func PSQLWrap(errs ...error) error {
	return errors.Join(ErrPSQL, errors.Join(errs...))
}

var (
	ErrInternal = errors.New("internal error")
	ErrRedis    = errors.New("redis error")
	ErrPSQL     = errors.New("psql error")
)
