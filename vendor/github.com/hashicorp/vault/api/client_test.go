package api

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"testing"
)

func init() {
	// Ensure our special envvars are not present
	os.Setenv("VAULT_ADDR", "")
	os.Setenv("VAULT_TOKEN", "")
}

func TestDefaultConfig_envvar(t *testing.T) {
	os.Setenv("VAULT_ADDR", "https://vault.mycompany.com")
	defer os.Setenv("VAULT_ADDR", "")

	config := DefaultConfig()
	if config.Address != "https://vault.mycompany.com" {
		t.Fatalf("bad: %s", config.Address)
	}

	os.Setenv("VAULT_TOKEN", "testing")
	defer os.Setenv("VAULT_TOKEN", "")

	client, err := NewClient(config)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if token := client.Token(); token != "testing" {
		t.Fatalf("bad: %s", token)
	}
}

func TestClientToken(t *testing.T) {
	tokenValue := "foo"
	handler := func(w http.ResponseWriter, req *http.Request) {}

	config, ln := testHTTPServer(t, http.HandlerFunc(handler))
	defer ln.Close()

	client, err := NewClient(config)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	client.SetToken(tokenValue)

	// Verify the token is set
	if v := client.Token(); v != tokenValue {
		t.Fatalf("bad: %s", v)
	}

	client.ClearToken()

	if v := client.Token(); v != "" {
		t.Fatalf("bad: %s", v)
	}
}

func TestClientRedirect(t *testing.T) {
	primary := func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("test"))
	}
	config, ln := testHTTPServer(t, http.HandlerFunc(primary))
	defer ln.Close()

	standby := func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Location", config.Address)
		w.WriteHeader(307)
	}
	config2, ln2 := testHTTPServer(t, http.HandlerFunc(standby))
	defer ln2.Close()

	client, err := NewClient(config2)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	// Set the token manually
	client.SetToken("foo")

	// Do a raw "/" request
	resp, err := client.RawRequest(client.NewRequest("PUT", "/"))
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	// Copy the response
	var buf bytes.Buffer
	io.Copy(&buf, resp.Body)

	// Verify we got the response from the primary
	if buf.String() != "test" {
		t.Fatalf("Bad: %s", buf.String())
	}
}

func TestClientEnvSettings(t *testing.T) {
	cwd, _ := os.Getwd()
	oldCACert := os.Getenv(EnvVaultCACert)
	oldCAPath := os.Getenv(EnvVaultCAPath)
	oldClientCert := os.Getenv(EnvVaultClientCert)
	oldClientKey := os.Getenv(EnvVaultClientKey)
	oldSkipVerify := os.Getenv(EnvVaultInsecure)
	os.Setenv("VAULT_CACERT", cwd+"/test-fixtures/keys/cert.pem")
	os.Setenv("VAULT_CAPATH", cwd+"/test-fixtures/keys")
	os.Setenv("VAULT_CLIENT_CERT", cwd+"/test-fixtures/keys/cert.pem")
	os.Setenv("VAULT_CLIENT_KEY", cwd+"/test-fixtures/keys/key.pem")
	os.Setenv("VAULT_SKIP_VERIFY", "true")
	defer os.Setenv("VAULT_CACERT", oldCACert)
	defer os.Setenv("VAULT_CAPATH", oldCAPath)
	defer os.Setenv("VAULT_CLIENT_CERT", oldClientCert)
	defer os.Setenv("VAULT_CLIENT_KEY", oldClientKey)
	defer os.Setenv("VAULT_SKIP_VERIFY", oldSkipVerify)

	config := DefaultConfig()
	if err := config.ReadEnvironment(); err != nil {
		t.Fatalf("error reading environment: %v", err)
	}

	tlsConfig := config.HttpClient.Transport.(*http.Transport).TLSClientConfig
	if len(tlsConfig.RootCAs.Subjects()) == 0 {
		t.Fatalf("bad: expected a cert pool with at least one subject")
	}
	if len(tlsConfig.Certificates) != 1 {
		t.Fatalf("bad: expected client tls config to have a client certificate")
	}
	if tlsConfig.InsecureSkipVerify != true {
		t.Fatalf("bad: %v", tlsConfig.InsecureSkipVerify)
	}
}
