//go:build linux

package ui

import (
	"fmt"
	"os"
	"strings"
	"syscall"
	"unsafe"
)

func SecretInput(label string) string {
	fmt.Print(Cyan + label + ": " + Reset)

	fd := int(os.Stdin.Fd())
	state, err := readTermios(fd)
	if err != nil {
		text, _ := reader.ReadString('\n')
		return strings.TrimSpace(text)
	}

	next := *state
	next.Lflag &^= syscall.ECHO
	if err := writeTermios(fd, &next); err != nil {
		text, _ := reader.ReadString('\n')
		return strings.TrimSpace(text)
	}

	defer func() {
		_ = writeTermios(fd, state)
		fmt.Println()
	}()

	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(text)
}

func readTermios(fd int) (*syscall.Termios, error) {
	state := &syscall.Termios{}
	_, _, errno := syscall.Syscall6(
		syscall.SYS_IOCTL,
		uintptr(fd),
		uintptr(syscall.TCGETS),
		uintptr(unsafe.Pointer(state)),
		0,
		0,
		0,
	)
	if errno != 0 {
		return nil, errno
	}
	return state, nil
}

func writeTermios(fd int, state *syscall.Termios) error {
	_, _, errno := syscall.Syscall6(
		syscall.SYS_IOCTL,
		uintptr(fd),
		uintptr(syscall.TCSETS),
		uintptr(unsafe.Pointer(state)),
		0,
		0,
		0,
	)
	if errno != 0 {
		return errno
	}
	return nil
}
