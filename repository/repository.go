package repository

type Authorization interface {
}

type Chat interface {
}

type Repository struct {
	Authorization
	Chat
}

func NewRepository()
