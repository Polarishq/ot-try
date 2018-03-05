package handlers

import (
	"net/http"
	"strings"
	"time"

	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http/httputil"
	"net/url"

	"github.com/Polarishq/middleware/framework/log"
	"github.com/satori/go.uuid"
)

// LoggingHandler provides a middleware handler which logs all requests and responses
type LoggingHandler struct {
	handler         http.Handler
	logRequestBody  bool
	logResponseBody bool
}

// NewLoggingHandler creates a middleware handler which logs all requests and responses excluding the response body
func NewLoggingHandler(handler http.Handler) *LoggingHandler {
	return &LoggingHandler{handler: handler, logRequestBody: true, logResponseBody: false}
}

// NewLoggingHandlerWithBody creates a middleware handler which logs all requests and responses,
// with an option to enable logging of the request and response body.
func NewLoggingHandlerWithBody(handler http.Handler, captureRequestBody bool, captureResponseBody bool) *LoggingHandler {
	return &LoggingHandler{handler: handler, logRequestBody: captureRequestBody, logResponseBody: captureResponseBody}
}

type loggingResponseWriter struct {
	headers     http.Header
	w           http.ResponseWriter
	data        []byte
	code        int
	captureBody bool
}

func (lw *loggingResponseWriter) Write(b []byte) (int, error) {
	if lw.captureBody {
		lw.data = append(lw.data, b...)
	}
	return lw.w.Write(b)
}

func (lw *loggingResponseWriter) WriteHeader(code int) {
	lw.headers = lw.Header()
	lw.code = code
	lw.w.WriteHeader(code)
}

func (lw *loggingResponseWriter) Header() http.Header {
	return lw.w.Header()
}

// Return value if nonempty, def otherwise.
func valueOrDefault(value, def string) string {
	if value != "" {
		return value
	}
	return def
}

// drainBody reads all of b to memory and then returns two equivalent
// ReadClosers yielding the same bytes.
//
// It returns an error if the initial slurp of all bytes fails. It does not attempt
// to make the returned ReadClosers have identical error-matching behavior.
func drainBody(b io.ReadCloser) (r1, r2 io.ReadCloser, err error) {
	if b == http.NoBody {
		// No copying needed. Preserve the magic sentinel meaning of NoBody.
		return http.NoBody, http.NoBody, nil
	}
	var buf bytes.Buffer
	if _, err = buf.ReadFrom(b); err != nil {
		return nil, b, err
	}
	if err = b.Close(); err != nil {
		return nil, b, err
	}
	return ioutil.NopCloser(&buf), ioutil.NopCloser(bytes.NewReader(buf.Bytes())), nil
}

//HeaderIs determines if a header has the the given value
func HeaderIs(header *http.Header, key string, value string) bool {
	if (*header)[key] != nil && (*header)[key][0] == value {
		return true
	}
	return false
}

// DumpRequest returns the given request in its HTTP/1.x wire
// representation. It should only be used by servers to debug Client
// requests. The returned representation is an approximation only;
// some details of the initial request are lost while parsing it into
// an http.Request. In particular, the order and case of header field
// names are lost. The order of values in multi-valued headers is kept
// intact. HTTP/2 requests are dumped in HTTP/1.x form, not in their
// original binary representations.
//
// If body is true, DumpRequest also returns the body. To do so, it
// consumes req.Body and then replaces it with a new io.ReadCloser
// that yields the same bytes. If DumpRequest returns an error,
// the state of req is undefined.
//
// The documentation for http.Request.Write details which fields
// of req are included in the dump.
func DumpRequest(req *http.Request, logRequestBody bool) (*log.Fields, error) {
	f := log.Fields{
		"method": valueOrDefault(req.Method, "GET"),
		"proto":  fmt.Sprintf("HTTP/%d.%d", req.ProtoMajor, req.ProtoMinor),
	}

	// By default, print out the unmodified req.RequestURI, which
	// is always set for incoming server requests. But because we
	// previously used req.URL.RequestURI and the docs weren't
	// always so clear about when to use DumpRequest vs
	// DumpRequestOut, fall back to the old way if the caller
	// provides a non-server Request.
	reqURI := req.RequestURI
	if reqURI == "" {
		reqURI = req.URL.RequestURI()
	}

	uri, err := url.ParseRequestURI(reqURI)
	if err != nil {
		return &f, err
	}

	if uri.RawQuery != "" {
		for k, v := range uri.Query() {
			f[k] = strings.Join(v[:], ",")
		}
	}
	f["path"] = uri.Path

	host := req.Host
	absRequestURI := strings.HasPrefix(req.RequestURI, "http://") ||
		strings.HasPrefix(req.RequestURI, "https://")
	if !absRequestURI {
		if host == "" && req.URL != nil {
			host = req.URL.Host
		}
	}
	f["host"] = host

	if logRequestBody == true {
		body, err := extractBody(req)
		if err != nil {
			return &f, err
		}
		for k, v := range *body {
			s, _ := v.(string)
			f[k] = s
		}
	}

	// Extract headers
	for k, v := range req.Header {
		f[k] = strings.Join(v[:], ",")
	}
	if err != nil {
		return &f, err
	}

	return &f, nil
}

func extractBody(req *http.Request) (*log.Fields, error) {
	var err error
	f := log.Fields{}

	save := req.Body
	if req.Body == nil {
		req.Body = nil
	} else {
		save, req.Body, err = drainBody(req.Body)
		if err != nil {
			return &f, err
		}
	}

	chunked := len(req.TransferEncoding) > 0 && req.TransferEncoding[0] == "chunked"
	if req.Body != nil {
		b := bytes.Buffer{}
		var dest io.Writer = &b
		if chunked {
			dest = httputil.NewChunkedWriter(dest)
		}
		_, err = io.Copy(dest, req.Body)
		if err != nil {
			return &f, err
		}
		if chunked {
			dest.(io.Closer).Close()
		}
		sbody := b.String()
		if sbody != "" {
			if HeaderIs(&req.Header, "Content-Type", "application/json") {
				jbody, err := JSONtoFields(&sbody)
				if err != nil {
					return &f, err
				}
				for k, v := range *jbody {
					f[k] = v
				}
			} else {
				f["body"] = sbody
			}
		}
	}

	req.Body = save
	return &f, nil
}

//JSONtoFields turns a JSON object into a set of log Fields
func JSONtoFields(raw *string) (*log.Fields, error) {
	var data map[string]interface{}
	err := json.Unmarshal([]byte(*raw), &data)
	if err != nil {
		return nil, err
	}
	f := log.Fields{}
	for k, v := range data {
		f[k] = v
	}
	return &f, nil
}

func (l *LoggingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	// Setup request ID tracking if it's present
	reqID := r.Header.Get("X-POLARIS-REQ-ID")
	if reqID == "" {
		reqID = uuid.NewV4().String()
		log.WithField("reqID", reqID).Info("Assigning new request id")
	}
	log.SetReqID(reqID)

	freq, err := DumpRequest(r, l.logRequestBody)
	if err != nil {
		log.Errorf("Failed to dump request: %+v", err)
	}
	log.WithFields(*freq).Info("start request")

	// Process the request
	lwr := loggingResponseWriter{w: w, captureBody: l.logResponseBody}
	l.handler.ServeHTTP(&lwr, r)
	endTime := time.Now()

	fresp := log.Fields{
		// if the code is 0, it means that an outer handler will write the code
		"code":                  fmt.Sprintf("%d", lwr.code),
		"duration_microseconds": fmt.Sprintf("%d", endTime.Sub(startTime).Nanoseconds()/1000),
	}
	log.WithFields(fresp).Info("end request")
}
