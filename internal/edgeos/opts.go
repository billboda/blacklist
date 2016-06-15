package edgeos

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
)

// Parms is struct of parameters
type Parms struct {
	Arch      string
	Cores     int
	Debug     bool
	Dex       List
	Dir       string
	Exc       List
	Ext       string
	File      string
	FnFmt     string
	Method    string
	Nodes     []string
	Pfx       string
	Poll      int
	Stypes    []string
	Test      bool
	Verbosity int
}

// Option sets is a recursive function
type Option func(p *Parms) Option

// SetOpt sets the specified options passed as Parms and returns an option to restore the last arg's previous value
func (p *Parms) SetOpt(opts ...Option) (previous Option) {
	// apply all the options, and replace each with its inverse
	for i, opt := range opts {
		opts[i] = opt(p)
	}

	// Reverse the list of inverses, since we want them to be applied in reverse order
	for i, j := 0, len(opts)-1; i <= j; i, j = i+1, j-1 {
		opts[i], opts[j] = opts[j], opts[i]
	}

	return func(p *Parms) Option {
		return p.SetOpt(opts...)
	}
}

// Arch sets target CPU architecture
func Arch(arch string) Option {
	return func(p *Parms) Option {
		previous := p.Arch
		p.Arch = arch
		return Arch(previous)
	}
}

// Cores sets max CPU cores
func Cores(i int) Option {
	return func(p *Parms) Option {
		previous := p.Cores
		runtime.GOMAXPROCS(i)
		p.Cores = i
		return Cores(previous)
	}
}

// Debug toggles debug level on or off
func Debug(b bool) Option {
	return func(p *Parms) Option {
		previous := p.Debug
		p.Debug = b
		return Debug(previous)
	}
}

// Dir sets directory location
func Dir(d string) Option {
	return func(p *Parms) Option {
		previous := p.Dir
		p.Dir = d
		return Dir(previous)
	}
}

// Excludes sets nodes exclusions
func Excludes(l List) Option {
	return func(p *Parms) Option {
		previous := p.Exc
		p.Exc = l
		return Excludes(previous)
	}
}

// Ext sets the blacklist file n extension
func Ext(e string) Option {
	return func(p *Parms) Option {
		previous := p.Ext
		p.Ext = e
		return Ext(previous)
	}
}

// File sets the EdgeOS configuration file n
func File(f string) Option {
	return func(p *Parms) Option {
		previous := p.File
		p.File = f
		return File(previous)
	}
}

// FileNameFmt sets the EdgeOS configuration file name format
func FileNameFmt(f string) Option {
	return func(p *Parms) Option {
		previous := p.FnFmt
		p.FnFmt = f
		return FileNameFmt(previous)
	}
}

// Prefix sets the dnsmasq configuration address line prefix
func Prefix(l string) Option {
	return func(p *Parms) Option {
		previous := p.Pfx
		p.Pfx = l
		return Prefix(previous)
	}
}

// Method sets the HTTP method
func Method(method string) Option {
	return func(p *Parms) Option {
		previous := p.Method
		p.Method = method
		return Method(previous)
	}
}

// NewParms sets a new *Parms instance
func NewParms() *Parms {
	return &Parms{
		Dex: make(List),
		Exc: make(List),
	}
}

// Nodes sets the node ns array
func Nodes(nodes []string) Option {
	return func(p *Parms) Option {
		previous := p.Nodes
		p.Nodes = nodes
		return Nodes(previous)
	}
}

// Poll sets the polling interval in seconds
func Poll(t int) Option {
	return func(p *Parms) Option {
		previous := p.Poll
		p.Poll = t
		return Poll(previous)
	}
}

// String method to implement fmt.Print interface
func (p *Parms) String() string {
	type pArray struct {
		n string
		i int
		v string
	}

	var fields []pArray

	maxLen := func(pA []pArray) int {
		smallest := len(pA[0].n)
		biggest := len(pA[0].n)
		for i := range pA {
			if len(pA[i].n) > biggest {
				biggest = len(pA[i].n)
			} else if len(pA[i].n) < smallest {
				smallest = len(pA[i].n)
			}
		}
		return biggest
	}

	v := reflect.ValueOf(p).Elem()
	for i := 0; i < v.NumField(); i++ {
		fields = append(fields, pArray{n: v.Type().Field(i).Name, v: strings.Replace(fmt.Sprint(v.Field(i).Interface()), "\n", "", -1)})
	}

	max := maxLen(fields)

	pad := func(s string) string {
		i := len(s)
		repeat := max - i + 1
		return strings.Repeat(" ", repeat)
	}

	r := fmt.Sprintln("edgeos.Parms{")
	for _, field := range fields {
		r += fmt.Sprintf("%v:%v%v\n", field.n, pad(field.n), field.v)
	}

	r += fmt.Sprintln("}")

	return r
}

// STypes sets an array of legal types used by Source
func STypes(s []string) Option {
	return func(p *Parms) Option {
		previous := p.Stypes
		p.Stypes = s
		return STypes(previous)
	}
}

// Test toggles testing mode on or off
func Test(b bool) Option {
	return func(p *Parms) Option {
		previous := p.Test
		p.Test = b
		return Test(previous)
	}
}

// Verbosity sets the verbosity level to v
func Verbosity(i int) Option {
	return func(p *Parms) Option {
		previous := p.Verbosity
		p.Verbosity = i
		return Verbosity(previous)
	}
}
