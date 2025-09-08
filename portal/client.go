package portal

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/kolo/xmlrpc"
)

// Options configure the TLS and transport behavior.
type Options struct {
	Server    string
	Port      int
	Path      string
	CertPEM   string
	KeyPEM    string
	CACertPEM string
	Verify    bool
	Timeout   time.Duration

	DisableHTTP2    bool
	InsecureCiphers bool
}

// Client is a thin XML-RPC client wrapper for the portal.
type Client struct {
	raw  *xmlrpc.Client
	opts Options
}

// New builds a Client with sane TLS defaults and timeouts.
func New(opts Options) (*Client, error) {
	if opts.Server == "" {
		opts.Server = "boss.emulab.net"
	}
	if opts.Port == 0 {
		opts.Port = 3069
	}
	if opts.Path == "" {
		opts.Path = "/usr/testbed"
	}
	if opts.CertPEM == "" {
		home, _ := os.UserHomeDir()
		opts.CertPEM = filepath.Join(home, ".ssl", "emulab.pem")
	}
	if opts.KeyPEM == "" {
		opts.KeyPEM = opts.CertPEM
	}

	cert, err := loadKeyPairLoose(opts.CertPEM, opts.KeyPEM)
	if err != nil {
		return nil, fmt.Errorf("load client cert/key: %w", err)
	}

	tcfg := &tls.Config{Certificates: []tls.Certificate{cert}}
	if opts.Verify {
		if opts.CACertPEM == "" {
			return nil, fmt.Errorf("verify=true requires CACertPEM")
		}
		ca, err := os.ReadFile(opts.CACertPEM)
		if err != nil {
			return nil, fmt.Errorf("read CA: %w", err)
		}
		pool := x509.NewCertPool()
		if !pool.AppendCertsFromPEM(ca) {
			return nil, fmt.Errorf("append CA failed")
		}
		tcfg.RootCAs = pool
	} else {
		tcfg.InsecureSkipVerify = true //nolint:gosec
	}
	if opts.InsecureCiphers {
		tcfg.CipherSuites = []uint16{
			tls.TLS_RSA_WITH_AES_128_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
		}
	}

	tr := &http.Transport{TLSClientConfig: tcfg}
	if opts.Timeout > 0 {
		tr.TLSHandshakeTimeout = opts.Timeout
		tr.ResponseHeaderTimeout = opts.Timeout
		tr.IdleConnTimeout = opts.Timeout
		tr.DialContext = (&net.Dialer{Timeout: opts.Timeout}).DialContext
	}
	if opts.DisableHTTP2 {
		tr.ForceAttemptHTTP2 = false
	}

	url := fmt.Sprintf("https://%s:%d%s", opts.Server, opts.Port, opts.Path)
	xcli, err := xmlrpc.NewClient(url, tr)
	if err != nil {
		return nil, fmt.Errorf("xmlrpc client: %w", err)
	}
	return &Client{raw: xcli, opts: opts}, nil
}

// doMethod invokes module.method with [version, params] args.
func (c *Client) doMethod(module, method string, params map[string]any) (*EmulabResponse, error) {
	full := module + "." + method
	var resp EmulabResponse
	args := []any{packageVersion, params}
	if err := c.raw.Call(full, args, &resp); err != nil {
		return nil, err
	}
	if resp.Code != ResponseSuccess {
		if resp.Value != nil {
			return &resp, fmt.Errorf("emulab error (%d): %v - %s", resp.Code, resp.Value, resp.Output)
		}
		return &resp, fmt.Errorf("emulab error (%d): %s", resp.Code, resp.Output)
	}
	return &resp, nil
}

/* ----- Public high-level helpers (pass-through) ----- */

func (c *Client) StartExperiment(params map[string]any) (*EmulabResponse, error) {
	return c.doMethod("portal", "startExperiment", params)
}

func (c *Client) ModifyExperiment(params map[string]any) (*EmulabResponse, error) {
	return c.doMethod("portal", "modifyExperiment", params)
}

func (c *Client) TerminateExperiment(experiment string) (*EmulabResponse, error) {
	return c.doMethod("portal", "terminateExperiment", map[string]any{"experiment": experiment})
}

func (c *Client) ExtendExperiment(experiment, hours, reason string) (*EmulabResponse, error) {
	return c.doMethod("portal", "extendExperiment", map[string]any{
		"experiment": experiment, "wanted": hours, "reason": reason,
	})
}

func (c *Client) ExperimentStatus(experiment string, asJSON, withCert, refresh bool) (*EmulabResponse, error) {
	p := map[string]any{"experiment": experiment}
	if asJSON {
		p["asjson"] = 1
	}
	if withCert {
		p["withcert"] = 1
	}
	if refresh {
		p["refresh"] = 1
	}
	return c.doMethod("portal", "experimentStatus", p)
}

func (c *Client) ExperimentManifests(experiment string) (*EmulabResponse, error) {
	return c.doMethod("portal", "experimentManifests", map[string]any{"experiment": experiment})
}

func (c *Client) RebootNodes(experiment string, nodes []string, power bool) (*EmulabResponse, error) {
	p := map[string]any{"experiment": experiment, "nodes": strings.Join(nodes, ",")}
	if power {
		p["power"] = 1
	}
	return c.doMethod("portal", "reboot", p)
}

func (c *Client) ConnectSharedLan(srcExp, srcLan, dstExp, dstLan string) (*EmulabResponse, error) {
	return c.doMethod("portal", "connectSharedLan", map[string]any{
		"experiment": srcExp, "sourcelan": srcLan, "targetexp": dstExp, "targetlan": dstLan,
	})
}

func (c *Client) DisconnectSharedLan(exp, lan string) (*EmulabResponse, error) {
	return c.doMethod("portal", "disconnectSharedLan", map[string]any{"experiment": exp, "sourcelan": lan})
}
