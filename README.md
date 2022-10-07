## Signature
Golang package for generating/adding the required authentication information to a MYRA API call.

[![go report card](https://goreportcard.com/badge/github.com/Myra-Security-GmbH/signature "go report card")](https://goreportcard.com/report/github.com/Myra-Security-GmbH/signature)
[![MIT license](http://img.shields.io/badge/license-MIT-brightgreen.svg)](http://opensource.org/licenses/MIT)
[![Go Reference](https://pkg.go.dev/badge/github.com/Myra-Security-GmbH/signature.svg)](https://pkg.go.dev/github.com/Myra-Security-GmbH/signature)
[![tests](https://github.com/Myra-Security-GmbH/signature/actions/workflows/test.yml/badge.svg)](https://github.com/Myra-Security-GmbH/signature/actions/workflows/test.yml)

### Usage
```
// Create a new Signature instance using your API credentials provided by Myra Security GmbH and the prepared request
s := signature.New(secret, apiKey, request)

// Append signature to the request and return prepared request
req, err := s.Append()
...`

// OR - return the signature string for own interaction:
t := time.Now().Format(time.RFC3339)
signature, err := s.Signature(t)
...

// OR - just generate and return the SigningString (required for the signature):
signingString := signature.SigningString("content data", "GET", "/en/rapi/...", t)
...
```
