package edgeos

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

// CFGcli loads configurations using the EdgeOS CFGcli
type CFGcli struct {
	*Config
	Cfg string
}

// CFGstatic loads static configurations for testing
type CFGstatic struct {
	*Config
	Cfg string
}

func active(a string, inCLI bool) string {
	switch inCLI {
	case true:
		switch a {
		case "exists":
			a = "existsActive"
		case "listNodes":
			a = "listActiveNodes"
		case "returnValue":
			a = "returnActiveValue"
		case "returnValues":
			a = "returnActiveValues"
		case "showCfg":
			a = "showConfig"
		}

	default:
		switch a {
		case "existsActive":
			a = "exists"
		case "listActiveNodes":
			a = "listNodes"
		case "returnActiveValue":
			a = "returnValue"
		case "returnActiveValues":
			a = "returnValues"
		case "showConfig":
			a = "showCfg"
		}
	}
	return a
}

// apiCMD returns a map of CLI commands
func apiCMD(a string, inCLI bool) string {
	var apiCMDs = map[string]string{
		"cfExists":           "cfExists",
		"cfReturnValue":      "cfReturnValue",
		"cfReturnValues":     "cfReturnValues",
		"echo":               "true",
		"exists":             "exists",
		"existsActive":       "existsActive",
		"getNodeType":        "getNodeType",
		"inSession":          "inSession",
		"isLeaf":             "isLeaf",
		"isMulti":            "isMulti",
		"isTag":              "isTag",
		"listActiveNodes":    "listActiveNodes",
		"listNodes":          "listNodes",
		"returnActiveValue":  "returnActiveValue",
		"returnActiveValues": "returnActiveValues",
		"returnValue":        "returnValue",
		"returnValues":       "returnValues",
		"showCfg":            "showCfg",
		"showConfig":         "showConfig",
	}
	return apiCMDs[active(a, inCLI)]
}

// deleteFile removes a file if it exists
func deleteFile(f string) bool {
	if err := os.Remove(f); err != nil {
		return false
	}
	return true
}

// GetFile reads a file and returns an io.Reader
func GetFile(f string) (io.Reader, error) {
	return os.Open(f)
}

// mode returns a contextual VYOS API argument
// func mode(insession bool) string {
// if insession {
// 	return "--show-working-only"
// }
// 	return "--show-active-only"
// }

// purgeFiles removes any orphaned blacklist files that don't have sources
func purgeFiles(files []string) error {
	var errs []string
	for _, f := range files {
		if _, err := os.Stat(f); !os.IsNotExist(err) {
			if !deleteFile(f) {
				errs = append(errs, fmt.Sprintf("could not remove %q", f))
			}
		}
	}

	if errs != nil {
		return errors.New(strings.Join(errs, "\n"))
	}
	return nil
}

// read returns an EdgeOS API configuration io.Reader
func (c *CFGcli) read() io.Reader {
	b, err := c.load("showCfg", c.Level)
	if err != nil {
		log.Print(err.Error())
	}
	return bytes.NewReader(b)
}

// read returns an EdgeOS config file io.Reader
func (c *CFGstatic) read() io.Reader {
	return strings.NewReader(c.Cfg)
}

// writeFile saves hosts/domains data to disk
func (b *bList) writeFile() error {
	if b.size == 0 {
		return nil
	}

	w, err := os.Create(b.file)
	if err != nil {
		return err
	}

	// defer w.Close()
	_, err = io.Copy(w, b.r)
	w.Close()
	return err
}
