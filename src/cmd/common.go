/*
Copyright © 2022 Cameron Larsen <cameron.larsen@nielseniq.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"context"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync/atomic"
	"time"
	"unsafe"

	timerate "golang.org/x/time/rate"
)

type (
	key          int
	RLHTTPClient struct {
		client      *http.Client
		Ratelimiter *timerate.Limiter
	}
	Request struct {
		R  *http.Request
		r  *http.Response
		e  error
		id string
	}
	Server struct {
		name   string
		port   int
		logger *log.Logger
		router *http.ServeMux
	}
	App interface {
		Start()
		Serve()
		NewLogger() *log.Logger
		NewRouter() *http.ServeMux
	}
)

const (
	requestIDKey  key = 0
	letterBytes       = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	letterIdxBits     = 6                    // 6 bits to represent a letter index
	letterIdxMask     = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax      = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var (
	healthy int32
	src     = rand.NewSource(time.Now().UnixNano())
)

func (s Server) NewLogger() *log.Logger {
	return log.New(os.Stdout, "["+s.name+"] ", log.LstdFlags)
}

func (s Server) NewRouter() *http.ServeMux {
	return http.NewServeMux()
}

func (s Server) Serve() {
	nextRequestID := func() string {
		return strconv.FormatInt(time.Now().UnixNano(), 10)
	}

	// instantiate server
	listenAddr := ":" + strconv.Itoa(s.port)
	server := &http.Server{
		Addr:         listenAddr,
		Handler:      tracing(nextRequestID)(logResp(s.logger)(s.router)),
		ErrorLog:     s.logger,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	done := make(chan bool)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	go func() {
		<-quit
		s.logger.Println("Server is shutting down...")
		atomic.StoreInt32(&healthy, 0)

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		server.SetKeepAlivesEnabled(false)
		if err := server.Shutdown(ctx); err != nil {
			s.logger.Fatalf("Unable to gracefully shutdown the server: %v\n", err)
		}
		close(done)
	}()

	s.logger.Printf("Server is ready to handle requests at %v", listenAddr)
	atomic.StoreInt32(&healthy, 1)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		s.logger.Fatalf("Unable to start server on %s: %v\n", listenAddr, err)
	}

	<-done
	s.logger.Println("Server stopped")
}

func index(delay int, percentage int) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		if r.URL.Path != "/" {
			w.Header().Set("X-Response-Code", "404")
			w.Header().Set("X-Request-Duration", time.Since(start).String())
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 - Not Found"))
			return
		}
		if rand.Intn(100) < percentage {
			time.Sleep(time.Duration(delay) * time.Millisecond)
			w.Header().Set("X-Response-Code", "500")
			w.Header().Set("X-Request-Duration", time.Since(start).String())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 - Internal Server Error"))
		} else {
			time.Sleep(time.Duration(delay) * time.Millisecond)
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.Header().Set("X-Content-Type-Options", "nosniff")
			w.Header().Set("X-Response-Code", "200")
			w.Header().Set("X-Request-Duration", time.Since(start).String())
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("200 - OK"))
		}
	})
}

func notFound(start time.Time) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Response-Code", "404")
		w.Header().Set("X-Request-Duration", time.Since(start).String())
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 - Not Found"))
	})
}

func healthz(percentage int) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		if rand.Intn(100) > percentage {
			if atomic.LoadInt32(&healthy) == 1 {
				w.Header().Set("Content-Type", "text/plain; charset=utf-8")
				w.Header().Set("X-Content-Type-Options", "nosniff")
				w.Header().Set("X-Response-Code", "204")
				w.Header().Set("X-Request-Duration", time.Since(start).String())
				w.WriteHeader(http.StatusNoContent)
				return
			}
		} else {
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.Header().Set("X-Content-Type-Options", "nosniff")
			w.Header().Set("X-Response-Code", "503")
			w.Header().Set("X-Request-Duration", time.Since(start).String())
			w.WriteHeader(http.StatusServiceUnavailable)
		}
	})
}

func logResp(logger *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				requestID, ok := r.Context().Value(requestIDKey).(string)
				if !ok {
					requestID = "unknown"
				}
				logger.Printf("%v - [%v][%v][%v] -> [%s] %s", requestID, r.RemoteAddr, r.Method, r.URL.Path, w.Header().Get("X-Response-Code"), w.Header().Get("X-Request-Duration"))
			}()
			next.ServeHTTP(w, r)
		})
	}
}

func tracing(nextRequestID func() string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get("X-Request-Id")
			if requestID == "" {
				requestID = nextRequestID()
			}
			ctx := context.WithValue(r.Context(), requestIDKey, requestID)
			w.Header().Set("X-Request-Id", requestID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// creates a single rate-limited client
func newClient(rateLimit *timerate.Limiter) *RLHTTPClient {
	return &RLHTTPClient{
		client:      &http.Client{
			Timeout:     5 * time.Second,
		},
		Ratelimiter: rateLimit,
	}
}

func (c *RLHTTPClient) Do(url string, percentage int, logger *log.Logger) {
	var req = &Request{
		id: randString(6),
	}
	if rand.Intn(100) < percentage {
		url += req.id + "/"
	}
	req.R, _ = http.NewRequest("GET", url, nil)
	req.r, req.e = c.client.Do(req.R)
	req.logReq(logger)
}

func (req Request) logReq(logger *log.Logger) {
	if req.e != nil {
		logger.Println(req.e.Error())
		return
	}
	logger.Printf("[%v][%v] -> [%s] %s", req.r.Request.Method, req.r.Request.URL, strconv.Itoa(req.r.StatusCode), req.r.Header.Get("X-Request-Duration"))
	defer req.r.Body.Close()
}

func randString(n int) string {
	b := make([]byte, n)
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return *(*string)(unsafe.Pointer(&b))
}
