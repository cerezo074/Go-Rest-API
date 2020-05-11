package cache

import (
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

type mockWriter response

func newMockWriter() *mockWriter {
	return &mockWriter{
		body:   []byte{},
		header: http.Header{},
	}
}

func (mw *mockWriter) Write(b []byte) (int, error) {
	mw.body = make([]byte, len(b))

	for index, value := range b {
		mw.body[index] = value
	}

	return len(b), nil
}

func (mw *mockWriter) Header() http.Header {
	return mw.header
}

func (mw *mockWriter) WriteHeader(code int) {
	mw.code = code
}

func TestWriter(testSuite *testing.T) {
	mockWriter := newMockWriter()
	resource := "/test/url?with=params"

	url, error := url.Parse(resource)

	if error != nil {
		testSuite.Fatal("Invalid url")
	}

	request := &http.Request{URL: url}

	///////////////////////////////////
	testSuite.Log("Test NewWritter")

	customWriter := NewWritter(mockWriter, request)

	if customWriter.resource != resource {
		testSuite.Errorf("Resources are different. Expected: %s / Actual: %s", resource, customWriter.resource)
	}

	if customWriter.writer != mockWriter {
		testSuite.Error("Writer not assigned")
	}

	///////////////////////////////////
	testSuite.Log("Test Header")

	header := customWriter.Header()
	header.Add("test", "value")
	header2 := customWriter.response.header

	if header2.Get("test") != "value" {
		testSuite.Error("Value not stored inside the header")
	}

	///////////////////////////////////
	testSuite.Log("Test WriteHeader")

	statusCode := 201
	customWriter.WriteHeader(statusCode)

	if customWriter.response.code != statusCode {
		testSuite.Error("Status code not stored")
	}

	if mockWriter.code != statusCode {
		testSuite.Error("Status code not stored")
	}

	header2 = mockWriter.header

	if header2.Get("test") != "value" {
		testSuite.Error("Header not written")
	}

	///////////////////////////////////
	testSuite.Log("Test Write")
	body := []byte{1, 2, 3, 4, 5}
	lenght, error := customWriter.Write(body)

	if error != nil {
		testSuite.Error("Body not written")
	}

	if lenght != len(body) {
		testSuite.Error("Body lenght is invalid")
	}

	if &customWriter.response.body == &body {
		testSuite.Error("Body assigned, not copied")
	}

	if !reflect.DeepEqual(customWriter.response.body, body) {
		testSuite.Error("Body not copied")
	}

	if &mockWriter.body == &body {
		testSuite.Error("Body assigned, not copied")
	}

	if !reflect.DeepEqual(mockWriter.body, body) {
		testSuite.Error("Body not copied")
	}
}
