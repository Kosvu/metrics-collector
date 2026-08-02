package middleware

import (
	"bytes"
	"compress/gzip"
	"encoding/hex"
	"io"
	"metrics/internal/hash"
	"net/http"
	"strings"
	"time"

	"go.uber.org/zap"
)

type signWriter struct {
	http.ResponseWriter
	buf    bytes.Buffer
	status int
}

func (w *signWriter) Write(b []byte) (int, error) {
	return w.buf.Write(b)
}

func (w *signWriter) WriteHeader(status int) {
	w.status = status
}

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w *gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

type responseData struct {
	status int
	size   int
}

type loggingResponseWriter struct {
	http.ResponseWriter
	responseData *responseData
}

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

type Middleware func(http.Handler) http.Handler

func WithLogging(log *zap.SugaredLogger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			responseData := &responseData{
				status: 200,
				size:   0,
			}
			lw := &loggingResponseWriter{
				ResponseWriter: w,
				responseData:   responseData,
			}

			next.ServeHTTP(lw, r)

			duration := time.Since(start)

			log.Infoln(
				"uri", r.RequestURI,
				"method", r.Method,
				"status", responseData.status,
				"duration", duration,
				"size", responseData.size,
			)
		})
	}
}

func SignResponse(key string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			if key == "" {
				next.ServeHTTP(w, r)
				return
			}

			wr := &signWriter{ResponseWriter: w, status: 200}
			next.ServeHTTP(wr, r)
			data := wr.buf.Bytes()
			signData := hash.Sign(data, []byte(key))
			w.Header().Set("HashSHA256", hex.EncodeToString(signData))

			w.WriteHeader(wr.status)
			w.Write(data)
		})
	}
}

func SignCheck(key string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if key == "" {
				next.ServeHTTP(w, r)
				return
			}

			dataSign := r.Header.Get("HashSHA256")

			if dataSign == "" {
				next.ServeHTTP(w, r)
				return
			}

			dataDecSign, err := hex.DecodeString(dataSign)

			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			body, err := io.ReadAll(r.Body)

			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			bodyReader := bytes.NewReader(body)
			doneReader := io.NopCloser(bodyReader)

			r.Body = doneReader

			f := hash.Verify(body, dataDecSign, []byte(key))

			if !f {
				w.WriteHeader(http.StatusBadRequest)
				return
			} else {
				next.ServeHTTP(w, r)
			}

		})
	}
}

func GZipHandle() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
				zr, err := gzip.NewReader(r.Body)
				if err != nil {
					io.WriteString(w, err.Error())
					return

				}
				r.Body = zr
				defer zr.Close()
			}

			if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
				zw, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
				if err != nil {
					io.WriteString(w, err.Error())
					return
				}
				defer zw.Close()

				w.Header().Set("Content-Encoding", "gzip")
				next.ServeHTTP(&gzipWriter{ResponseWriter: w, Writer: zw}, r)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
