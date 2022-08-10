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
	"log"
	"net/http"
	"reflect"
	"testing"
	"time"

	timerate "golang.org/x/time/rate"
)

func TestServer_NewLogger(t *testing.T) {
	tests := []struct {
		name string
		s    Server
		want *log.Logger
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.NewLogger(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Server.NewLogger() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServer_NewRouter(t *testing.T) {
	tests := []struct {
		name string
		s    Server
		want *http.ServeMux
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.NewRouter(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Server.NewRouter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServer_Serve(t *testing.T) {
	tests := []struct {
		name string
		s    Server
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.s.Serve()
		})
	}
}

func Test_index(t *testing.T) {
	type args struct {
		delay      int
		percentage int
	}
	tests := []struct {
		name string
		args args
		want http.Handler
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := index(tt.args.delay, tt.args.percentage); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("index() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_notFound(t *testing.T) {
	type args struct {
		start time.Time
	}
	tests := []struct {
		name string
		args args
		want http.Handler
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := notFound(tt.args.start); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("notFound() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_healthz(t *testing.T) {
	type args struct {
		percentage int
	}
	tests := []struct {
		name string
		args args
		want http.Handler
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := healthz(tt.args.percentage); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("healthz() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_logResp(t *testing.T) {
	type args struct {
		logger *log.Logger
	}
	tests := []struct {
		name string
		args args
		want func(http.Handler) http.Handler
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := logResp(tt.args.logger); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("logResp() = %p, want %p", got, tt.want)
			}
		})
	}
}

func Test_tracing(t *testing.T) {
	type args struct {
		nextRequestID func() string
	}
	tests := []struct {
		name string
		args args
		want func(http.Handler) http.Handler
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tracing(tt.args.nextRequestID); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("tracing() = %p, want %p", got, tt.want)
			}
		})
	}
}

func Test_newClient(t *testing.T) {
	type args struct {
		rateLimit *timerate.Limiter
	}
	tests := []struct {
		name string
		args args
		want *RLHTTPClient
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newClient(tt.args.rateLimit); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newClient() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRLHTTPClient_Do(t *testing.T) {
	type args struct {
		url        string
		percentage int
		logger     *log.Logger
	}
	tests := []struct {
		name string
		c    *RLHTTPClient
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.c.Do(tt.args.url, tt.args.percentage, tt.args.logger)
		})
	}
}

func TestRequest_logReq(t *testing.T) {
	type args struct {
		logger *log.Logger
	}
	tests := []struct {
		name string
		req  Request
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.req.logReq(tt.args.logger)
		})
	}
}

func Test_randString(t *testing.T) {
	type args struct {
		n int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := randString(tt.args.n); got != tt.want {
				t.Errorf("randString() = %v, want %v", got, tt.want)
			}
		})
	}
}
