package services

import "tbot/internal/model"

type Service interface {
	GetTag() string
	Query(query string) (model.Resp, error)
	IsReady() bool
}
