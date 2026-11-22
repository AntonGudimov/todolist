package application

import "context"

type Component interface {
	Run(ctx context.Context) error
	Shutdown(ctx context.Context) error
}
