package command

import (
	"context"
	"tickets/domain"
)

type CreateShowHandler struct {
	repo domain.Repository
}

func NewCreateShowHandler(repo domain.Repository) CreateShowHandler {
	return CreateShowHandler{repo: repo}
}

func (h *CreateShowHandler) Handle(ctx context.Context, show domain.Show) error {
	return h.repo.CreateShow(ctx, show)
}
