package main

import (
	"flag"
	"fmt"
	"log"
	"runtime"
	"syscall"
	"time"

	e "github.com/britannic/blacklist/internal/edgeos"
)

const (
	all   = "all"
	files = "file"
	pre   = "pre-configured"
	urls  = "url"
)

var (
	// Version vars updated by go build -ldflags
	build   = "UNKNOWN"
	githash = "UNKNOWN"
	version = "UNKNOWN"
)

func main() {

	o := getOpts()
	o.Init("blacklist", flag.ExitOnError)
	o.setArgs(func(code int) {
		syscall.Exit(code)
	})

	c := o.initEdgeOS()
	c.ReadCfg(o.getCFG(c))
	fmt.Println(c.String())

	c.SetOpt(
		e.Dexcludes(c.Excludes("blacklist", "domains")),
		e.Excludes(c.Excludes("hosts")),
	)

	if err := c.GetAll().Files().Remove(); err != nil {
		log.Printf("c.GetAll().Files().Remove() error: %v\n", err)
	}

	f, err := c.CreateObject(e.FileObj)
	if err != nil {
		log.Fatal(err)
	}

	p, err := c.CreateObject(e.PreObj)
	if err != nil {
		log.Fatal(err)
	}

	u, err := c.CreateObject(e.URLObj)
	if err != nil {
		log.Fatal(err)
	}

	c.ProcessContent(p)
	c.ProcessContent(f)
	c.ProcessContent(u)

	b, err := c.ReloadDNS()
	if err != nil {
		log.Printf("ReloadDNS(): %v\n error: %v\n", string(b), err)
	}
	log.Printf("ReloadDNS(): %v\n", string(b))
}

// basename removes directory components and file extensions.
func basename(s string) string {
	// Discard last '/' and everything before.
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == '/' {
			s = s[i+1:]
			break
		}
	}

	// Preserve everything before last '.'.
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == '.' {
			s = s[:i]
			break
		}
	}
	return s
}

func (o *Opts) initEdgeOS() *e.Config {
	return e.NewConfig(
		e.API("/bin/cli-shell-api"),
		e.Arch(runtime.GOARCH),
		e.Bash("/bin/bash"),
		e.Cores(runtime.NumCPU()),
		e.Debug(*o.Debug),
		e.Dir(o.SetDir(*o.ARCH)),
		e.DNSsvc("service dnsmasq restart"),
		e.Ext("blacklist.conf"),
		e.File(*o.File),
		e.FileNameFmt("%v/%v.%v.%v"),
		e.InCLI("inSession"),
		e.Level("service dns forwarding"),
		e.Method("GET"),
		e.Nodes([]string{"domains", "hosts"}),
		e.Poll(*o.Poll),
		e.Prefix("address="),
		e.STypes([]string{"file", pre, "url"}),
		e.Timeout(30*time.Second),
		e.WCard(e.Wildcard{Node: "*s", Name: "*"}),
	)

}
