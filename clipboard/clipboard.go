package clipboard

import (
	"context"
)

type Clipboard interface {
	Write(string) error
	Read() (string, error)
	Watch(ctx context.Context, f func(string))
}
