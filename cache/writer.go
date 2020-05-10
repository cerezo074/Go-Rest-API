package cache

import "net/http"

//Writer is a wrapper for the response writer that caches the responses
type Writer struct {
	writer   http.ResponseWriter
	response response
	resource string
}

var (
	_ http.ResponseWriter = (*Writer)(nil)
)

//NewWritter creates a cache writer
func NewWritter(writer http.ResponseWriter, request *http.Request) *Writer {
	return &Writer{
		writer:   writer,
		resource: MakeResource(request),
		response: response{
			header: http.Header{},
		},
	}
}

//Header gets the headers from a Writer
func (w *Writer) Header() http.Header {
	return w.response.header
}

//WriteHeader writes headers to the response writer
func (w *Writer) WriteHeader(code int) {
	copyHeader(w.response.header, w.writer.Header())
	w.response.code = code
	w.writer.WriteHeader(code)
}

func (w *Writer) Write(b []byte) (int, error) {
	w.response.body = make([]byte, len(b))
	for index, value := range b {
		w.response.body[index] = value
	}

	copyHeader(w.Header(), w.writer.Header())
	set(w.resource, &w.response)
	return w.writer.Write(b)
}
