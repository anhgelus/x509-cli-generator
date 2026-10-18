# X.509 CLI generator

Simple CLI application to generate X.509 certificates.

Usage:
```
x509 <common_name> [options]

Usage of options:
  -can-sign
        if the certificate can sign
  -ext int
        x509 extension key usage (Go int) (default -1)
  -p224
        use a P-226 curve
  -p256
        use a P-256 curve (default)
  -p384
        use a P-384 curve
  -p521
        use a P-521 curve
  -parent string
        path to the PEM encoded parent
```
