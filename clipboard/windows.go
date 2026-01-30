//go:build windows

package clipboard

import (
	"context"

	C "golang.design/x/clipboard"
)

func GetClipboard() (Clipboard, error) {
	err := C.Init()
	if err != nil {
		return nil, err
	}
	return &clip{}, nil
}

type clip struct{}

func (c *clip) Read() (string, error) {
	buf := C.Read(C.FmtText)
	return string(buf), nil
}

func (c *clip) Write(text string) error {
	C.Write(C.FmtText, []byte(text))
	return nil
}

func (c *clip) Watch(ctx context.Context, f func(text string)) {
	ch := C.Watch(ctx, C.FmtText)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case data := <-ch:
				f(string(data))
			}
		}
	}()
}
