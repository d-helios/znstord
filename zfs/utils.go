package zfs

import (
	"bytes"
	"io"
	"os/exec"
	"reflect"
	"strconv"
	"strings"
	"unicode"
)

type command struct {
	Command string
	Stdin   io.Reader
	Stdout  io.Writer
}

// wrapper for exec/Command
func (c *command) Run(filter string, arg ...string) ([][]string, error) {

	cmd := exec.Command(c.Command, arg...)

	var stdout, stderr bytes.Buffer

	functionFilter := func(r rune) bool { return true }

	if filter == ":" {
		functionFilter = unicode.IsPunct
	} else {
		functionFilter = unicode.IsSpace
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
		output[i] = strings.FieldsFunc(j, functionFilter)
	}

	return output, nil
}

// wrapper for cmdZfs calls
func cmdZfs(arg ...string) ([][]string, error) {
	c := command{Command: "zfs"}
	return c.Run("", arg...)
}

// wrapper for cmdZpool calls
func cmdZpool(arg ...string) ([][]string, error) {
	c := command{Command: "zpool"}
	return c.Run("", arg...)
}

func cmdTest(arg ...string) error {
	c := command{Command: "test"}
	_, err := c.Run("", arg...)
	return err
}

func setFloat(field *float64, value string) error {
	if strings.Contains(value, "x") {
		v, err := strconv.ParseFloat(
			strings.Replace(value, "x", "", -1), 64)
		if err != nil {
			return err
		}
		*field = v
	}
	return nil
}

// convert cmdZfs size values into uint64
func setUint(field *uint64, value string) error {
	if value != "-" && value != "" {
		v, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return err
		}
		*field = v
	}
	return nil
}

// Convert Numeric Dataset properties to uint64
func IsItNumericProp(prop string, list []string) bool {
	for _, b := range list {
		if b == prop {
			return true
		}
	}
	return false
}

func FillDataset(dataset interface{}, dict map[string]string) error {
	val := reflect.ValueOf(dataset).Elem()
	var numField uint64
	var floatField float64
	for i := 0; i < val.NumField(); i++ {
		field := strings.ToLower(val.Type().Field(i).Name)
		if IsItNumericProp(field, NumericProps) {
			err := setUint(&numField, dict[field])
			if err != nil {
				return err
			}
			val.Field(i).SetUint(numField)
			continue
		}
		if IsItNumericProp(field, FloatProps) {
			err := setFloat(&floatField, dict[field])
			if err != nil {
				return err
			}
			val.Field(i).SetFloat(floatField)
			continue
		}
		val.Field(i).SetString(dict[field])
	}
	return nil
}

func GetOsRelease() (string, error) {
	args := []string{"list", "-H", "-o", "version", RootPool}
	out, err := cmdZpool(args...)
	if err != nil {
		return "", err
	}
	if out[0][0] == "-" {
		return OpenSolaris, nil
	} else {
		return OracleSolaris, nil
	}
}
