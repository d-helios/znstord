package stmf

import (
	"bytes"
	"io"
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
	newS := make([]string, len(s))

	for i := range s {
		newS[i] = strings.Trim(s[i], trim)
	}
	return newS
}

// GetZvol - return zvol on which logical unit is based
func (logicalunit *LogicalUnit) GetZvol() string {
	return logicalunit.DataFile[15:]
}

// wrapper for exec/Command
func (c *command) Run(filter string, arg ...string) ([][]string, error) {

	cmd := exec.Command(c.Command, arg...)

	var stdout, stderr bytes.Buffer

	filterFunction := func(r rune) bool { return true }

	if filter == ":" {
		filterFunction = fieldFunc
	} else {
		filterFunction = unicode.IsSpace
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

	lines := strings.Split(stdout.String(), "\n")

	// last line is always blank
	lines = lines[0 : len(lines)-1]

	output := make([][]string, len(lines))

	for i, j := range lines {
		output[i] = trimList(strings.FieldsFunc(j, filterFunction), " ")

	}

	return output, nil
}
