package filter

import (
	"compress/flate"
	"compress/gzip"
	"io"
	"net/http"
	"reflect"
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

var writerType = reflect.TypeOf(&gzipWriter{})

func newGzipWriter(w http.ResponseWriter) *gzipWriter {
	w.Header().Set("Content-Encoding", "gzip")
	return &gzipWriter{ResponseWriter: w, gzWriter: func() *gzip.Writer {
		w, _ := gzip.NewWriterLevel(w, flate.BestCompression)
		return w
	}()}
}

func (gw *gzipWriter) Write(bs []byte) (int, error) {
	gw.written = true
	n, err := gw.gzWriter.Write(bs)
	return n, err
}
func (gw *gzipWriter) Close() error {
	if gw.written {
		return gw.gzWriter.Close()
	}
	return nil
}

func (l *GzipFilter) PreHandle(req *cango.WebRequest) interface{} {
	if strings.Contains(req.Request.Header.Get("Accept-Encoding"), "gzip") {
		return newGzipWriter(req.ResponseWriter)
	}
	return true
}

func (l *GzipFilter) PostHandle(req *cango.WebRequest) interface{} {
	if reflect.TypeOf(req.ResponseWriter) == writerType {
		err := req.ResponseWriter.(io.Closer).Close()
		if err != nil {
			canlog.CanError(req.URL.Path+req.URL.RawQuery, err)
		}
	}
	return true
}
