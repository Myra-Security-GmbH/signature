## Signature
Golang package for generating/adding the required authentication information to a MYRA API call.

### Usage
```
// Create a new Signature instance using your API credentials provided by MYRA Security GmbH and the prepared request
s := signature.New(secret, apiKey, request)

// Append signature to the request and return prepared request
req, err := s.Append()
...`

// OR - return the signature string for own interaction:
t := time.Now().Format(time.RFC3339)
signature, err := s.Signature(t)
...

// OR - just generate and return the SigningString (required for the signature):
signingString := s.SigningString("content data", "GET", "/en/rapi/...", t)
...
```
