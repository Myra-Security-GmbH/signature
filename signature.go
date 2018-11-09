package signature

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const (
	// RequestTypeJSON ...
	RequestTypeJSON = "application/json"
)

//
// Signature holds information about the request, API key, API Secret
// API key and API secret were provided by MYRA Security GmbH
//
type Signature struct {
	secret  string
	apikey  string
	request *http.Request
}

//
// New creates and returns a new Signature instance
//
func New(secret string, apikey string, request *http.Request) *Signature {
	return &Signature{
		secret:  secret,
		apikey:  apikey,
		request: request,
	}
}

//
// Append builds and adds signature/authentication information to the request and returns the request.
//
func (s *Signature) Append() (*http.Request, error) {

	t := time.Now().Format(time.RFC3339)

	signature, err := s.Signature(t)
	if err != nil {
		return s.request, err
	}
	s.request.Header.Add("Authorization", fmt.Sprintf("MYRA %s:%s", s.apikey, signature))
	s.request.Header.Add("Date", t)
	s.request.Header.Add("Content-Type", RequestTypeJSON)

	return s.request, nil
}

//
// SigningString returns a string for the signature, formatted as descibed in the MYRA API documentation
//
func SigningString(body string, method string, path string, date string) string {
	return fmt.Sprintf(
		"%x#%s#%s#%s#%s",
		md5.Sum([]byte(body)),
		strings.ToUpper(method),
		path,
		RequestTypeJSON,
		date,
	)
}

//
// Signature builds and returns the signature.
// date is a time.RFC3339 formatted date string
//
func (s *Signature) Signature(date string) (string, error) {

	var err error
	body := []byte("")
	if s.request.Body != nil {
		body, err = ioutil.ReadAll(s.request.Body)
		if err != nil {
			return "", err
		}
		s.request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	}

	h1 := hmac.New(sha256.New, []byte("MYRA"+s.secret))
	_, err = h1.Write([]byte(date))
	if err != nil {
		return "", err
	}
	dateKey := fmt.Sprintf("%x", h1.Sum(nil))

	h2 := hmac.New(sha256.New, []byte(dateKey))
	_, err = h2.Write([]byte("myra-api-request"))
	if err != nil {
		return "", err
	}
	signingKey := fmt.Sprintf("%x", h2.Sum(nil))

	path := s.request.URL.Path
	if len(s.request.URL.RawQuery) > 0 {
		path = fmt.Sprintf("%s?%s", path, s.request.URL.RawQuery)
	}

	signingString := SigningString(string(body), s.request.Method, path, date)
	h3 := hmac.New(sha512.New, []byte(signingKey))
	_, err = h3.Write([]byte(signingString))
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(h3.Sum(nil)), nil
}
