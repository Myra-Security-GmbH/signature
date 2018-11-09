// +build !testing

package signature

import (
	"crypto/md5"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

//
// buildRequest returns a new request for the given parameters
//
func buildRequest(url string, method string, payload map[string]string) *http.Request {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		panic(err)
	}
	return req
}

//
// dpTestNew provides data for TestNew
//
func dpTestNew() [][]string {
	return [][]string{
		{"1234", "4321"},
		{"", ""},
		{" ", " "},
		{"test", ""},
		{"", "test"},
		{"test", "test"},
	}
}

func TestNew(t *testing.T) {
	req := buildRequest("https://api.myracloud.com/en/rapi/dnsRecords/", "GET", map[string]string{})
	for ds, data := range dpTestNew() {
		s := New(data[0], data[1], req)

		assert.Equal(t, data[0], s.secret, fmt.Sprintf("Expected secret to be %s, got %s (dataset: %d)", data[0], s.secret, ds))
		assert.Equal(t, data[1], s.apikey, fmt.Sprintf("Expected apikey to be %s, got %s (dataset: %d)", data[1], s.apikey, ds))
	}
}

//
// dpTestAppend provides data for TestAppend
//
func dpTestAppend() []*Signature {
	return []*Signature{
		New("1234", "4321", buildRequest("https://api.myracloud.com/en/rapi/dnsRecords/", "GET", map[string]string{})),
		New("1234", "4321", buildRequest("https://api.myracloud.com/en/rapi/dnsRecords/", "GET", map[string]string{"test": "me", "some": "parameter"})),
		New("1234", "4321", buildRequest("https://api.myracloud.com/en/rapi/dnsRecords/", "POST", map[string]string{})),
		New("1234", "4321", buildRequest("https://api.myracloud.com/en/rapi/dnsRecords/", "POST", map[string]string{"test": "me", "some": "parameter"})),
	}
}

func TestAppend(t *testing.T) {
	for ds, sig := range dpTestAppend() {
		req, err := sig.Append()
		if err != nil {
			t.Errorf("Unexpected error for dataset %d: %s", ds, err.Error())
		}

		authorization := req.Header.Get("Authorization")
		assert.NotEmpty(t, authorization, fmt.Sprintf("Missing 'Authorization' Header for dataset %d", ds))

		contentType := req.Header.Get("Content-Type")
		assert.NotEmpty(t, contentType, fmt.Sprintf("Missing 'Content-Type' Header for dataset %d", ds))

		date := req.Header.Get("Date")
		assert.NotEmpty(t, date, fmt.Sprintf("Missing 'Date' Header for dataset %d", ds))
		_, err = time.Parse(time.RFC3339, date)
		assert.NoError(t, err)
	}
}

//
// dpTestSigningString provides data for TestSigningString
//
func dpTestSigningString() [][]string {
	return [][]string{
		{"this is the body", "GET", "/test/me", time.Now().Format(time.RFC3339)},
		{"this is the body", "get", "/test/me", time.Now().Format(time.RFC3339)},
		{"", "get", "/test/me?test=me", time.Now().Format(time.RFC3339)},
		{"{\"test\":\"data\"}", "PoSt", "/test/me", time.Now().Format(time.RFC3339)},
		{"", "PUT", "/test/me", time.Now().Format(time.RFC3339)},
		{"{}", "PUT", "/test/me", time.Now().Format(time.RFC3339)},
	}
}

func TestSigningString(t *testing.T) {
	for ds, data := range dpTestSigningString() {
		result := SigningString(data[0], data[1], data[2], data[3])
		assert.NotEmpty(t, result, fmt.Sprintf("Missing signing string for dataset %d", ds))

		assert.Contains(t, result, fmt.Sprintf("%x", md5.Sum([]byte(data[0]))), fmt.Sprintf("Missing body information in signing string for dataset %d", ds))
		assert.Contains(t, result, strings.ToUpper(data[1]), fmt.Sprintf("Missing method information in signing string for dataset %d", ds))
		assert.Contains(t, result, data[2], fmt.Sprintf("Missing path information in signing string for dataset %d", ds))
		assert.Contains(t, result, data[3], fmt.Sprintf("Missing date information in signing string for dataset %d", ds))
	}
}

//
// dpTestSignature provides data for TestSignature
//
func dpTestSignature() []*Signature {
	return []*Signature{
		New("1234", "4321", buildRequest("https://api.myracloud.com/en/rapi/dnsRecords/", "GET", map[string]string{})),
		New("1234", "4321", buildRequest("https://api.myracloud.com/en/rapi/dnsRecords/", "GET", map[string]string{"test": "me", "some": "parameter"})),
		New("1234", "4321", buildRequest("https://api.myracloud.com/en/rapi/dnsRecords/", "POST", map[string]string{})),
		New("1234", "4321", buildRequest("https://api.myracloud.com/en/rapi/dnsRecords/", "POST", map[string]string{"test": "me", "some": "parameter"})),
	}
}

func TestSignature(t *testing.T) {
	for ds, sig := range dpTestSignature() {
		date := time.Now().Format(time.RFC3339)
		result, err := sig.Signature(date)

		assert.NoError(t, err)
		assert.NotEmpty(t, result, "Missing signature result for dataset %d", ds)
	}
}
