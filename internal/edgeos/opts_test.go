package edgeos

import (
	"runtime"
	"testing"
	"time"

	. "github.com/britannic/testutils"
)

func TestOption(t *testing.T) {
	vanilla := Parms{
		API:       "",
		Arch:      "",
		Cores:     0,
		Debug:     false,
		Dex:       list{entry: entry(nil)},
		Dir:       "",
		DNSsvc:    "",
		Exc:       list{entry: entry(nil)},
		Ext:       "",
		File:      "",
		FnFmt:     "",
		InCLI:     "",
		Level:     "",
		Method:    "",
		Nodes:     []string(nil),
		Pfx:       "",
		Poll:      0,
		Ltypes:    []string(nil),
		Test:      false,
		Timeout:   0,
		Verbosity: 0,
		Wildcard:  Wildcard{},
	}

	want := "edgeos.Parms{\nWildcard:  \"{*s *}\"\nAPI:       \"/bin/cli-shell-api\"\nArch:      \"amd64\"\nBash:      \"/bin/bash\"\nCores:     \"2\"\nDebug:     \"true\"\nDex:       \"**not initialized**\"\nDir:       \"/tmp\"\nDNSsvc:    \"service dnsmasq restart\"\nExc:       \"**not initialized**\"\nExt:       \"blacklist.conf\"\nFile:      \"/config/config.boot\"\nFnFmt:     \"%v/%v.%v.%v\"\nInCLI:     \"inSession\"\nLevel:     \"service dns forwarding\"\nMethod:    \"GET\"\nNodes:     \"[domains hosts]\"\nPfx:       \"address=\"\nPoll:      \"10\"\nLtypes:    \"[file pre-configured-domain pre-configured-host url]\"\nTest:      \"true\"\nTimeout:   \"30s\"\nVerbosity: \"2\"\n}\n"

	wantRaw := Parms{
		API:       "/bin/cli-shell-api",
		Arch:      "amd64",
		Bash:      "/bin/bash",
		Cores:     2,
		Debug:     true,
		Dex:       list{entry: entry{}},
		Dir:       "/tmp",
		DNSsvc:    "service dnsmasq restart",
		Exc:       list{entry: entry{}},
		Ext:       "blacklist.conf",
		File:      "/config/config.boot",
		FnFmt:     "%v/%v.%v.%v",
		InCLI:     "inSession",
		Level:     "service dns forwarding",
		Method:    "GET",
		Nodes:     []string{domains, hosts},
		Pfx:       "address=",
		Poll:      10,
		Ltypes:    []string{files, PreDomns, PreHosts, urls},
		Test:      true,
		Timeout:   30000000000,
		Verbosity: 2,
		Wildcard:  Wildcard{Node: "*s", Name: "*"},
	}

	Equals(t, vanilla, Parms{})

	c := NewConfig(
		Arch(runtime.GOARCH),
		API("/bin/cli-shell-api"),
		Bash("/bin/bash"),
		Cores(2),
		Debug(true),
		Dir("/tmp"),
		DNSsvc("service dnsmasq restart"),
		Ext("blacklist.conf"),
		File("/config/config.boot"),
		FileNameFmt("%v/%v.%v.%v"),
		InCLI("inSession"),
		Method("GET"),
		Nodes([]string{domains, hosts}),
		Poll(10),
		Prefix("address="),
		Level("service dns forwarding"),
		LTypes([]string{"file", PreDomns, PreHosts, urls}),
		Test(true),
		Timeout(30*time.Second),
		Verbosity(2),
		WCard(Wildcard{Node: "*s", Name: "*"}),
	)

	wantRaw.Dex.RWMutex = c.Dex.RWMutex
	wantRaw.Exc.RWMutex = c.Exc.RWMutex
	Equals(t, wantRaw, *c.Parms)

	Equals(t, want, c.Parms.String())
}
