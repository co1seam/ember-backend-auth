package ports

import (
	"context"
)

type Creator interface {
	Create(ctx context.Context, entity ...interface{}) (interface{}, error)
}

type Reader interface {
	Read(ctx context.Context, entity ...interface{}) (interface{}, error)
}

type Updater interface {
	Update(ctx context.Context, entity ...interface{}) (interface{}, error)
}

type Deleter interface {
	Delete(ctx context.Context, entity ...interface{}) (interface{}, error)
}

type Geter interface {
	Get(ctx context.Context, entity ...interface{}) (interface{}, error)
}

type Seter interface {
	Set(ctx context.Context, entity ...interface{}) (interface{}, error)
}

type CRUD interface {
	Creator
	Reader
	Updater
	Deleter
}

type Cache interface {
	Geter
	Seter
}
