package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
	"testing"

	"github.com/britannic/blacklist/internal/edgeos"
	"github.com/britannic/blacklist/internal/tdata"
	. "github.com/britannic/testutils"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGetOpts(t *testing.T) {
	Convey("Testing commandline output", t, func() {
		exp := vanillaArgs

		exitCmd = func(int) { return }

		out := new(bytes.Buffer)
		o := getOpts()
		o.Init("blacklist", flag.ContinueOnError)
		o.SetOutput(out)

		os.Args = append(os.Args, "-h")
		o.setArgs()

		act, err := ioutil.ReadAll(out)

		So(err, ShouldBeNil)

		So(string(act), ShouldEqual, exp)

	})
}

func TestBasename(t *testing.T) {
	tests := []struct {
		s   string
		exp string
	}{
		{s: "e.txt", exp: "e"},
		{s: "/github.com/britannic/blacklist/internal/edgeos", exp: "edgeos"},
	}
	for _, tt := range tests {
		Equals(t, tt.exp, basename(tt.s))
	}
}

func TestBuild(t *testing.T) {
	want := map[string]string{
		"build":   build,
		"githash": githash,
		"version": version,
	}

	for k := range want {
		Equals(t, "UNKNOWN", want[k])
	}
}

func TestCommandLineArgs(t *testing.T) {
	exitCmd = func(int) { return }
	out := new(bytes.Buffer)
	want := vanillaArgs
	fmt.Println(os.Args)
	os.Args = append(os.Args, "-h")
	o := getOpts()
	o.Init("blacklist", flag.ContinueOnError)

	o.SetOutput(out)
	o.Parse(cleanArgs(os.Args[1:]))
	o.setArgs()

	got, err := ioutil.ReadAll(out)

	Convey("Testing command line arguments", t, func() {
		So(err, ShouldBeNil)
		So(string(got), ShouldEqual, want)
	})
}

func TestGetCFG(t *testing.T) {
	exitCmd = func(int) { return }
	exp := mainGetConfig
	o := getOpts()
	c := o.initEdgeOS()

	c.ReadCfg(o.getCFG(c))
	Equals(t, exp, c.String())
	o.Set("mips64", "amd64")
	// *o.MIPS64 = "amd64"

	c = o.initEdgeOS()
	c.ReadCfg(o.getCFG(c))
	Equals(t, "{\n  \"nodes\": [{\n  }]\n}", c.String())
}

// func TestGetOpts(t *testing.T) {
// 	exitCmd = func(int) { return }
// 	o := getOpts()
// 	want := "FlagSet\nARCH:    \"amd64\"\nDEBUG:   \"false\"\nDIR:     \"/etc/dnsmasq.d\"\nF:       \"**not initialized**\"\nH:       \"false\"\nI:       \"5\"\nMIPS64:  \"mips64\"\nOS:      \"" + runtime.GOOS + "\"\nTEST:    \"false\"\nTMP:     \"/tmp\"\nV:       \"false\"\nVERSION: \"false\"\n"
//
// 	Equals(t, want, o.String())
//
// 	tests := []struct {
// 		name string
// 		test interface{}
// 		exp  interface{}
// 	}{
// 		{
// 			name: "o.Debug",
// 			test: o.Debug,
// 			exp:  true,
// 		},
// 		{
// 			name: "o.File",
// 			test: o.File,
// 			exp:  "",
// 		},
// 		{
// 			name: "o.Poll",
// 			test: o.Poll,
// 			exp:  8,
// 		},
// 		{
// 			name: "o.Test",
// 			test: o.Test,
// 			exp:  true,
// 		},
// 		{
// 			name: "o.Verb",
// 			test: o.Verb,
// 			exp:  true,
// 		},
// 		{
// 			name: "o.Version",
// 			test: o.Version,
// 			exp:  true,
// 		},
// 	}
//
// 	for _, run := range tests {
// 		switch run.test.(type) {
// 		case bool:
// 			Equals(t, run.exp.(bool), run.test.(bool))
//
// 		case string:
// 			Equals(t, run.exp.(string), run.test.(string))
//
// 		case int:
// 			Equals(t, run.exp.(int), run.test.(int))
// 		}
// 	}
// }

func TestSetArch(t *testing.T) {
	exitCmd = func(int) { return }
	o := getOpts()

	tests := []struct {
		arch string
		exp  string
	}{
		{arch: "mips64", exp: "/etc/dnsmasq.d"},
		{arch: "linux", exp: "/tmp"},
		{arch: "darwin", exp: "/tmp"},
	}

	for _, test := range tests {
		Equals(t, test.exp, o.setDir(test.arch))
	}
}

type cfgCLI struct {
	edgeos.CFGcli
}

func (c cfgCLI) Load() io.Reader {
	return strings.NewReader(tdata.Cfg)
}

func TestInitEdgeOS(t *testing.T) {
	exitCmd = func(int) { return }
	o := getOpts()
	p := o.initEdgeOS()
	exp := "edgeos.Parms{\nWildcard:  \"{*s *}\"\nAPI:       \"/bin/cli-shell-api\"\nArch:      \"amd64\"\nBash:      \"/bin/bash\"\nCores:     \"2\"\nDebug:     \"false\"\nDex:       \"**not initialized**\"\nDir:       \"/tmp\"\nDNSsvc:    \"service dnsmasq restart\"\nExc:       \"**not initialized**\"\nExt:       \"blacklist.conf\"\nFile:      \"**not initialized**\"\nFnFmt:     \"%v/%v.%v.%v\"\nInCLI:     \"inSession\"\nLevel:     \"service dns forwarding\"\nMethod:    \"GET\"\nNodes:     \"[domains hosts]\"\nPfx:       \"address=\"\nPoll:      \"5\"\nLtypes:    \"[file pre-configured-domain pre-configured-host url]\"\nTest:      \"false\"\nTimeout:   \"30s\"\nVerbosity: \"0\"\n}\n"

	Equals(t, exp, fmt.Sprint(p.Parms))
}

var (
	JSONcfg = "{\n  \"nodes\": [{\n    \"blacklist\": {\n      \"disabled\": \"false\",\n      \"ip\": \"0.0.0.0\",\n      \"excludes\": [\n        \"122.2o7.net\",\n        \"1e100.net\",\n        \"adobedtm.com\",\n        \"akamai.net\",\n        \"amazon.com\",\n        \"amazonaws.com\",\n        \"apple.com\",\n        \"ask.com\",\n        \"avast.com\",\n        \"bitdefender.com\",\n        \"cdn.visiblemeasures.com\",\n        \"cloudfront.net\",\n        \"coremetrics.com\",\n        \"edgesuite.net\",\n        \"freedns.afraid.org\",\n        \"github.com\",\n        \"githubusercontent.com\",\n        \"google.com\",\n        \"googleadservices.com\",\n        \"googleapis.com\",\n        \"googleusercontent.com\",\n        \"gstatic.com\",\n        \"gvt1.com\",\n        \"gvt1.net\",\n        \"hb.disney.go.com\",\n        \"hp.com\",\n        \"hulu.com\",\n        \"images-amazon.com\",\n        \"msdn.com\",\n        \"paypal.com\",\n        \"rackcdn.com\",\n        \"schema.org\",\n        \"skype.com\",\n        \"smacargo.com\",\n        \"sourceforge.net\",\n        \"ssl-on9.com\",\n        \"ssl-on9.net\",\n        \"static.chartbeat.com\",\n        \"storage.googleapis.com\",\n        \"windows.net\",\n        \"yimg.com\",\n        \"ytimg.com\"\n        ]\n    },\n    \"domains\": {\n      \"disabled\": \"false\",\n      \"ip\": \"0.0.0.0\",\n      \"excludes\": [],\n      \"includes\": [\n        \"adsrvr.org\",\n        \"adtechus.net\",\n        \"advertising.com\",\n        \"centade.com\",\n        \"doubleclick.net\",\n        \"free-counter.co.uk\",\n        \"intellitxt.com\",\n        \"kiosked.com\"\n        ],\n      \"sources\": [{\n        \"malc0de\": {\n          \"disabled\": \"false\",\n          \"description\": \"List of zones serving malicious executables observed by malc0de.com/database/\",\n          \"prefix\": \"zone \",\n          \"file\": \"\",\n          \"url\": \"http://malc0de.com/bl/ZONES\"\n        }\n    }]\n    },\n    \"hosts\": {\n      \"disabled\": \"false\",\n      \"ip\": \"192.168.168.1\",\n      \"excludes\": [],\n      \"includes\": [\"beap.gemini.yahoo.com\"],\n      \"sources\": [{\n        \"adaway\": {\n          \"disabled\": \"false\",\n          \"description\": \"Blocking mobile ad providers and some analytics providers\",\n          \"prefix\": \"127.0.0.1 \",\n          \"file\": \"\",\n          \"url\": \"http://adaway.org/hosts.txt\"\n        },\n        \"malwaredomainlist\": {\n          \"disabled\": \"false\",\n          \"description\": \"127.0.0.1 based host and domain list\",\n          \"prefix\": \"127.0.0.1 \",\n          \"file\": \"\",\n          \"url\": \"http://www.malwaredomainlist.com/hostslist/hosts.txt\"\n        },\n        \"openphish\": {\n          \"disabled\": \"false\",\n          \"description\": \"OpenPhish automatic phishing detection\",\n          \"prefix\": \"http\",\n          \"file\": \"\",\n          \"url\": \"https://openphish.com/feed.txt\"\n        },\n        \"someonewhocares\": {\n          \"disabled\": \"false\",\n          \"description\": \"Zero based host and domain list\",\n          \"prefix\": \"0.0.0.0\",\n          \"file\": \"\",\n          \"url\": \"http://someonewhocares.org/hosts/zero/\"\n        },\n        \"tasty\": {\n          \"disabled\": \"false\",\n          \"description\": \"File source\",\n          \"prefix\": \"\",\n          \"file\": \"../testdata/blist.hosts.src\",\n          \"url\": \"\"\n        },\n        \"volkerschatz\": {\n          \"disabled\": \"false\",\n          \"description\": \"Ad server blacklists\",\n          \"prefix\": \"http\",\n          \"file\": \"\",\n          \"url\": \"http://www.volkerschatz.com/net/adpaths\"\n        },\n        \"winhelp2002\": {\n          \"disabled\": \"false\",\n          \"description\": \"Zero based host and domain list\",\n          \"prefix\": \"0.0.0.0 \",\n          \"file\": \"\",\n          \"url\": \"http://winhelp2002.mvps.org/hosts.txt\"\n        },\n        \"yoyo\": {\n          \"disabled\": \"false\",\n          \"description\": \"Fully Qualified Domain Names only - no prefix to strip\",\n          \"prefix\": \"\",\n          \"file\": \"\",\n          \"url\": \"http://pgl.yoyo.org/as/serverlist.php?hostformat=nohtml&showintro=1&mimetype=plaintext\"\n        }\n    }]\n    }\n  }]\n}"

	mainGetConfig = "{\n  \"nodes\": [{\n    \"blacklist\": {\n      \"disabled\": \"false\",\n      \"ip\": \"0.0.0.0\",\n      \"excludes\": [\n        \"1e100.net\",\n        \"2o7.net\",\n        \"adobedtm.com\",\n        \"akamai.net\",\n        \"akamaihd.net\",\n        \"amazon.com\",\n        \"amazonaws.com\",\n        \"apple.com\",\n        \"ask.com\",\n        \"avast.com\",\n        \"bitdefender.com\",\n        \"cdn.visiblemeasures.com\",\n        \"cloudfront.net\",\n        \"coremetrics.com\",\n        \"edgesuite.net\",\n        \"freedns.afraid.org\",\n        \"github.com\",\n        \"githubusercontent.com\",\n        \"google.com\",\n        \"googleadservices.com\",\n        \"googleapis.com\",\n        \"googletagmanager.com\",\n        \"googleusercontent.com\",\n        \"gstatic.com\",\n        \"gvt1.com\",\n        \"gvt1.net\",\n        \"hb.disney.go.com\",\n        \"hp.com\",\n        \"hulu.com\",\n        \"images-amazon.com\",\n        \"live.com\",\n        \"microsoft.com\",\n        \"msdn.com\",\n        \"msecnd.net\",\n        \"paypal.com\",\n        \"rackcdn.com\",\n        \"schema.org\",\n        \"shopify.com\",\n        \"skype.com\",\n        \"smacargo.com\",\n        \"sourceforge.net\",\n        \"ssl-on9.com\",\n        \"ssl-on9.net\",\n        \"sstatic.net\",\n        \"static.chartbeat.com\",\n        \"storage.googleapis.com\",\n        \"windows.net\",\n        \"xboxlive.com\",\n        \"yimg.com\",\n        \"ytimg.com\"\n        ]\n    },\n    \"domains\": {\n      \"disabled\": \"false\",\n      \"ip\": \"192.168.100.1\",\n      \"excludes\": [],\n      \"includes\": [\n        \"adsrvr.org\",\n        \"adtechus.net\",\n        \"advertising.com\",\n        \"centade.com\",\n        \"doubleclick.net\",\n        \"free-counter.co.uk\",\n        \"intellitxt.com\",\n        \"kiosked.com\",\n        \"patoghee.in\"\n        ],\n      \"sources\": [{\n        \"malc0de\": {\n          \"disabled\": \"false\",\n          \"description\": \"List of zones serving malicious executables observed by malc0de.com/database/\",\n          \"ip\": \"192.168.168.1\",\n          \"prefix\": \"zone \",\n          \"url\": \"http://malc0de.com/bl/ZONES\",\n        },\n        \"malwaredomains.com\": {\n          \"disabled\": \"false\",\n          \"description\": \"Just domains\",\n          \"ip\": \"10.0.0.1\",\n          \"url\": \"http://mirror1.malwaredomains.com/files/justdomains\",\n        },\n        \"simple_tracking\": {\n          \"disabled\": \"false\",\n          \"description\": \"Basic tracking list by Disconnect\",\n          \"url\": \"https://s3.amazonaws.com/lists.disconnect.me/simple_tracking.txt\",\n        },\n        \"zeus\": {\n          \"disabled\": \"false\",\n          \"description\": \"abuse.ch ZeuS domain blocklist\",\n          \"url\": \"https://zeustracker.abuse.ch/blocklist.php?download=domainblocklist\",\n        }\n    }]\n    },\n    \"hosts\": {\n      \"disabled\": \"false\",\n      \"excludes\": [],\n      \"includes\": [\"beap.gemini.yahoo.com\"],\n      \"sources\": [{\n        \"openphish\": {\n          \"disabled\": \"false\",\n          \"description\": \"OpenPhish automatic phishing detection\",\n          \"prefix\": \"http\",\n          \"url\": \"https://openphish.com/feed.txt\",\n        },\n        \"raw.github.com\": {\n          \"disabled\": \"false\",\n          \"description\": \"This hosts file is a merged collection of hosts from reputable sources\",\n          \"prefix\": \"0.0.0.0 \",\n          \"url\": \"https://raw.githubusercontent.com/StevenBlack/hosts/master/hosts\",\n        },\n        \"sysctl.org\": {\n          \"disabled\": \"false\",\n          \"description\": \"This hosts file is a merged collection of hosts from cameleon\",\n          \"ip\": \"172.16.16.1\",\n          \"prefix\": \"127.0.0.1\\t \",\n          \"url\": \"http://sysctl.org/cameleon/hosts\",\n        },\n        \"tasty\": {\n          \"disabled\": \"false\",\n          \"description\": \"File source\",\n          \"ip\": \"10.10.10.10\",\n          \"file\": \"../testdata/blist.hosts.src\",\n        },\n        \"volkerschatz\": {\n          \"disabled\": \"false\",\n          \"description\": \"Ad server blacklists\",\n          \"prefix\": \"http\",\n          \"url\": \"http://www.volkerschatz.com/net/adpaths\",\n        },\n        \"yoyo\": {\n          \"disabled\": \"false\",\n          \"description\": \"Fully Qualified Domain Names only - no prefix to strip\",\n          \"url\": \"http://pgl.yoyo.org/as/serverlist.php?hostformat=nohtml&showintro=1&mimetype=plaintext\",\n        }\n    }]\n    }\n  }]\n}"

	vanillaArgs = `  -arch string
    	Set EdgeOS CPU architecture (default "amd64")
  -debug
    	Enable debug mode
  -dir string
    	Override dnsmasq directory (default "/etc/dnsmasq.d")
  -f string
    	<file> # Load a configuration file
  -h	Display help
  -i int
    	Polling interval (default 5)
  -mips64 string
    	Override target EdgeOS CPU architecture (default "mips64")
  -os string
    	Override native EdgeOS OS (default "` + runtime.GOOS + `")
  -test
    	Run config and data validation tests
  -tmp string
    	Override dnsmasq temporary directory (default "/tmp")
  -v	Verbose display
  -version
    	Show version
`
)
