package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"flag"
	"os"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		println(
			"Invalid usage.\n" +
				"Usage: x509 <common_name> [options]\n\n" +
				"To see options: x509 test -h\n",
		)
		os.Exit(1)
	}
	cn := os.Args[1]
	flags := flag.NewFlagSet("options", flag.ExitOnError)
	pathParent := flags.String("parent", "", "path to the PEM encoded parent")
	ca := flags.Bool("can-sign", false, "if the certificate can sign")
	var ext Slice[int]
	flags.Var(&ext, "ext", "x509 extension key usage (Go int)")
	is224 := flags.Bool("p224", false, "use a P-226 curve")
	_ = flags.Bool("p256", false, "use a P-256 curve (default)")
	is384 := flags.Bool("p384", false, "use a P-384 curve")
	is521 := flags.Bool("p521", false, "use a P-521 curve")
	var notAfter Duration
	flags.Var(&notAfter, "not-after", "certificate cannot be used after this relative date (default: 30d)")
	var notBefore Duration
	flags.Var(&notBefore, "not-before", "certificate cannot be used before this relative date")
	var dns Slice[string]
	flags.Var(&dns, "dns", "DNS Names authorized for this certificate")
	if len(os.Args) != 2 {
		_ = flags.Parse(os.Args[2:])
	}
	curve := elliptic.P256()
	switch {
	case *is521:
		curve = elliptic.P521()
	case *is384:
		curve = elliptic.P384()
	case *is224:
		curve = elliptic.P224()
	}
	var parent *x509.Certificate
	var parentPriv *ecdsa.PrivateKey
	if len(*pathParent) != 0 {
		b, err := os.ReadFile(*pathParent)
		if err != nil {
			panic(err)
		}
		for len(b) != 0 {
			var bl *pem.Block
			bl, b = pem.Decode(b)
			switch bl.Type {
			case "CERTIFICATE":
				parent, err = x509.ParseCertificate(bl.Bytes)
				if err != nil {
					panic(err)
				}
			case "PRIVATE KEY":
				k, err := x509.ParsePKCS8PrivateKey(bl.Bytes)
				if err != nil {
					panic(err)
				}
				parentPriv = k.(*ecdsa.PrivateKey)
			default:
				panic("unknown PEM type " + bl.Type)
			}
		}
	}
	if !notAfter.has {
		notAfter.dur = 30 * 24 * time.Hour
	}
	now := time.Now()
	templ := x509.Certificate{
		Subject: pkix.Name{
			CommonName: cn,
		},
		IsCA:                  *ca,
		BasicConstraintsValid: *ca,
		KeyUsage:              x509.KeyUsageDigitalSignature,
		NotBefore:             now.Add(notBefore.dur),
		NotAfter:              now.Add(notAfter.dur),
		DNSNames:              dns.Content,
	}
	if templ.IsCA {
		templ.KeyUsage |= x509.KeyUsageCertSign
	}
	if ext.Content != nil {
		templ.ExtKeyUsage = make([]x509.ExtKeyUsage, 0, len(ext.Content))
		for _, v := range ext.Content {
			templ.ExtKeyUsage = append(templ.ExtKeyUsage, x509.ExtKeyUsage(v))
		}
	}
	priv, err := ecdsa.GenerateKey(curve, nil)
	if err != nil {
		panic(err)
	}
	if parent == nil {
		parent = &templ
		parentPriv = priv
	}
	b, err := x509.CreateCertificate(rand.Reader, &templ, parent, &priv.PublicKey, parentPriv)
	if err != nil {
		panic(err)
	}
	var buf bytes.Buffer
	err = pem.Encode(&buf, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: b,
	})
	if err != nil {
		panic(err)
	}
	b, err = x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		panic(err)
	}
	err = pem.Encode(&buf, &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: b,
	})
	if err != nil {
		panic(err)
	}
	os.Stdout.Write(buf.Bytes())
}
