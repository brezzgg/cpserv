//go:build linux

package clipboard

import "errors"

func GetClipboard() (Clipboard, error) {
	return nil, errors.New("cpserv server is not working on linux")
}
