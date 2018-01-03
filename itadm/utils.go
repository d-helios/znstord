package itadm

import (
	"bytes"
	"io"
	_ "log"
	"os/exec"
	"strings"
	"unicode"
)

type command struct {
	Command string
	Stdin   io.Reader
	Stdout  io.Writer
}

func fieldFunc(r rune) bool {
	if r == ':' {
		return true
	}
	return false
}

func trimList(s []string, trim string) []string {
	new_s := make([]string, len(s))

	for i := range s {
		new_s[i] = strings.Trim(s[i], trim)
	}
	return new_s
}

// wrapper for exec/Command
func (c *command) Run(filter string, arg ...string) ([][]string, error) {

	cmd := exec.Command(c.Command, arg...)

	var stdout, stderr bytes.Buffer

	func_filter := func(r rune) bool { return true }

	if filter == ":" {
		func_filter = fieldFunc
	} else {
		func_filter = unicode.IsSpace
	}

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	joinedArgs := strings.Join(cmd.Args, " ")

	err := cmd.Run()

	if err != nil {
		return nil, &Error{
			Err:    err,
			Debug:  strings.Join([]string{cmd.Path, c.Command, joinedArgs}, " "),
			Stderr: stderr.String(),
		}
	}

	// replace all output " = " with "=". Needed for itadm list-targets
	lines := strings.Split(strings.Replace(stdout.String(), " = ", "=", -1), "\n")

	// last line is always blank
	lines = lines[0 : len(lines)-1]

	output := make([][]string, len(lines))

	for i, j := range lines {
		output[i] = trimList(strings.FieldsFunc(j, func_filter), " ")

	}

	return output, nil
}
