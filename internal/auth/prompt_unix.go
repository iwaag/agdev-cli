//go:build linux || darwin

package auth

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"golang.org/x/sys/unix"
)

func PromptLine(prompt string) (string, error) {
	fmt.Fprint(os.Stderr, prompt)

	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(line), nil
}

func PromptPassword(prompt string) (string, error) {
	fd := int(os.Stdin.Fd())
	info, err := os.Stdin.Stat()
	if err != nil {
		return "", err
	}
	if (info.Mode() & os.ModeCharDevice) == 0 {
		return "", fmt.Errorf("password prompt requires a terminal")
	}

	fmt.Fprint(os.Stderr, prompt)

	state, err := unix.IoctlGetTermios(fd, unix.TCGETS)
	if err != nil {
		return "", err
	}

	updated := *state
	updated.Lflag &^= unix.ECHO
	if err := unix.IoctlSetTermios(fd, unix.TCSETS, &updated); err != nil {
		return "", err
	}
	defer func() {
		_ = unix.IoctlSetTermios(fd, unix.TCSETS, state)
		fmt.Fprintln(os.Stderr)
	}()

	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(line), nil
}
