package flag

import (
	"testing"
)

func Test_toCase(t *testing.T) {
	testCases := []struct {
		Input    string
		Expected string
	}{
		{
			Input:    "alllower",
			Expected: "alllower",
		},
		{
			Input:    "Title",
			Expected: "title",
		},
		{
			Input:    "weirD",
			Expected: "weirD",
		},
		{
			Input:    "ALLUPPER",
			Expected: "allupper",
		},
		{
			Input:    "FIRSTlower",
			Expected: "firstlower",
		},
		{
			Input:    "FIRSTlowerSuffix",
			Expected: "firstlowerSuffix",
		},
		{
			Input:    "FIRSTlowerUPPERSuffix",
			Expected: "firstlowerUPPERSuffix",
		},
	}

	for i, testCase := range testCases {
		expected := toCase(testCase.Input)
		if expected != testCase.Expected {
			t.Fatal("case", i+1, "expected", testCase.Expected, "got", expected)
		}
	}
}

func Test_Init(t *testing.T) {
	e := ""
	f := testFlag{}
	s := ""

	// Make sure the uninitialized flag structure does not have any values set.
	{
		e = ""
		s = f.Config.Dirs
		if s != e {
			t.Fatal("expected", e, "got", s)
		}
		e = ""
		s = f.Config.Files
		if s != e {
			t.Fatal("expected", e, "got", s)
		}
		e = ""
		s = f.Server.Listen.Address
		if s != e {
			t.Fatal("expected", e, "got", s)
		}
		e = ""
		s = f.Server.TLS.CaFile
		if s != e {
			t.Fatal("expected", e, "got", s)
		}
		e = ""
		s = f.Server.TLS.CrtFile
		if s != e {
			t.Fatal("expected", e, "got", s)
		}
		e = ""
		s = f.Server.TLS.KeyFile
		if s != e {
			t.Fatal("expected", e, "got", s)
		}
		e = ""
		s = f.Foo
		if s != e {
			t.Fatal("expected", e, "got", s)
		}
	}

	Init(&f)

	// Make sure the initialized flag structure does have the proper values set.
	{
		e = "config.dirs"
		s = f.Config.Dirs
		if s != e {
			t.Fatal("expected", e, "got", s)
		}
		e = "config.files"
		s = f.Config.Files
		if s != e {
			t.Fatal("expected", e, "got", s)
		}
		e = "server.listen.address"
		s = f.Server.Listen.Address
		if s != e {
			t.Fatal("expected", e, "got", s)
		}
		e = "server.tls.caFile"
		s = f.Server.TLS.CaFile
		if s != e {
			t.Fatal("expected", e, "got", s)
		}
		e = "server.tls.crtFile"
		s = f.Server.TLS.CrtFile
		if s != e {
			t.Fatal("expected", e, "got", s)
		}
		e = "server.tls.keyFile"
		s = f.Server.TLS.KeyFile
		if s != e {
			t.Fatal("expected", e, "got", s)
		}
		e = "foo"
		s = f.Foo
		if s != e {
			t.Fatal("expected", e, "got", s)
		}
	}
}

type testFlag struct {
	Config testConfig
	Server testServer
	Foo    string
}

type testConfig struct {
	Dirs  string
	Files string
}

type testServer struct {
	Listen testListen
	TLS    testTLS
}

type testListen struct {
	Address string
}

type testTLS struct {
	CaFile  string
	CrtFile string
	KeyFile string
}
