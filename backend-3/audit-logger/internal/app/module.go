package app

import (
	"go.uber.org/fx"

	"github.com/xakepp35/pkg/ufx"
)

func NewApp() *fx.App {
	return fx.New(
		ufx.WithZeroLogger(),
		config.Module,
		http.Module,
		usecase.Module,
		repository.Module,
		pkg.Module,
		worker.Module,
	)
}
