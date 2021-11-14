package filters

import (
	"compress/flate"
	"compress/gzip"
	"net/http"
	"strings"

	"github.com/JessonChan/cango"
	"github.com/JessonChan/canlog"
)

type GzipFilter struct {
	cango.Filter `value:"/*"`
}

type gzipWriter struct {
	http.ResponseWriter
	gzWriter *gzip.Writer
	written  bool
}

func newGzipWriter(w http.ResponseWriter) *gzipWriter {
	return &gzipWriter{ResponseWriter: w, gzWriter: func() *gzip.Writer {
		w, _ := gzip.NewWriterLevel(w, flate.BestCompression)
		return w
	}()}
}

func (gw *gzipWriter) Write(bs []byte) (int, error) {
	if !gw.written {
		gw.Header().Del("Content-Length")
		gw.Header().Set("Content-Encoding", "gzip")
	}
	gw.written = true
	return gw.gzWriter.Write(bs)
}

func (gw *gzipWriter) Close() error {
	if gw.written {
		return gw.gzWriter.Close()
	}
	return nil
}

func (l *GzipFilter) PreHandle(req *cango.WebRequest) interface{} {
	if strings.Contains(req.Request.Header.Get("Accept-Encoding"), "gzip") {
		req.ResponseWriter = newGzipWriter(req.ResponseWriter)
	}
	return true
}

func (l *GzipFilter) PostHandle(req *cango.WebRequest) interface{} {
	if rw, ok := req.ResponseWriter.(*gzipWriter); ok {
		err := rw.Close()
		if err != nil {
			canlog.CanError(req.URL.Path+req.URL.RawQuery, err)
		}
	}
	return true
}
