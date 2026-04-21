//go:build !linux

package ui

func SecretInput(label string) string {
	return Input(label)
}
