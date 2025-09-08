package portal

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
)

// loadKeyPairLoose reads PEM files for cert + key and builds a tls.Certificate
// without strict AIA/chain resolution. This is intentionally permissive to
// interop with site-provided client certs that can trip strict loaders.
func loadKeyPairLoose(certPath, keyPath string) (tls.Certificate, error) {
	certPEM, err := os.ReadFile(certPath)
	if err != nil {
		return tls.Certificate{}, fmt.Errorf("read cert: %w", err)
	}
	keyPEM, err := os.ReadFile(keyPath)
	if err != nil {
		return tls.Certificate{}, fmt.Errorf("read key: %w", err)
	}

	var chain [][]byte
	for {
		block, rest := pem.Decode(certPEM)
		if block == nil {
			break
		}
		if block.Type == "CERTIFICATE" && len(block.Bytes) > 0 {
			chain = append(chain, block.Bytes)
		}
		certPEM = rest
	}
	if len(chain) == 0 {
		return tls.Certificate{}, errors.New("no CERTIFICATE blocks found")
	}

	kb, _ := pem.Decode(keyPEM)
	if kb == nil {
		return tls.Certificate{}, errors.New("no PRIVATE KEY block found")
	}

	var parsed any
	switch kb.Type {
	case "RSA PRIVATE KEY":
		parsed, err = x509.ParsePKCS1PrivateKey(kb.Bytes)
	case "EC PRIVATE KEY":
		parsed, err = x509.ParseECPrivateKey(kb.Bytes)
	case "PRIVATE KEY":
		parsed, err = x509.ParsePKCS8PrivateKey(kb.Bytes)
	default:
		return tls.Certificate{}, fmt.Errorf("unsupported key type %q", kb.Type)
	}
	if err != nil {
		return tls.Certificate{}, fmt.Errorf("parse private key: %w", err)
	}
	switch parsed.(type) {
	case *rsa.PrivateKey, *ecdsa.PrivateKey, ed25519.PrivateKey:
	default:
		return tls.Certificate{}, fmt.Errorf("unsupported private key kind %T", parsed)
	}

	return tls.Certificate{Certificate: chain, PrivateKey: parsed}, nil
}
