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
	"strconv"
	"time"
	"os"

	timerate "golang.org/x/time/rate"

	"github.com/spf13/cobra"
	"gopkg.in/DataDog/dd-trace-go.v1/profiler"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

// workerCmd represents the worker command
var workerCmd = &cobra.Command{
	Use:   "worker",
	Short: "Starts a worker instance",
	Run: func(cmd *cobra.Command, args []string) {
		localPort, _ := cmd.Flags().GetInt("health-port")
		url, _ := cmd.Flags().GetString("url")
		port, _ := cmd.Flags().GetInt("port")
		rate, _ := cmd.Flags().GetInt("rate")
		fail, _ := cmd.Flags().GetInt("fail")
		failHealth, _ := cmd.Flags().GetInt("health-fail")
		urlString := url + ":" + strconv.Itoa(port) + "/"
		datadog, _ := cmd.Flags().GetBool("datadog")

		server := &Server{
			port: localPort,
			name: "Worker",
			datadog: datadog,
		}

		// instantiate server
		server.logger = server.NewLogger()
		server.router = server.NewRouter()
		// Initialize DataDog tracing
		if datadog {
			tracer.Start(
				tracer.WithLogStartup(true),
				tracer.WithService(os.Getenv("DD_SERVICE")),
				tracer.WithUniversalVersion(os.Getenv("DD_VERSION")),
				tracer.WithEnv(os.Getenv("DD_ENV")),
			)
			defer tracer.Stop()
			err := profiler.Start(
				profiler.WithLogStartup(true),
				profiler.WithService(os.Getenv("DD_SERVICE")),
				profiler.WithVersion(os.Getenv("DD_VERSION")),
				profiler.WithEnv(os.Getenv("DD_ENV")),
			)
			if err != nil {
				server.logger.Fatal(err)
			}
			defer profiler.Stop()
		}
		server.logger.Printf("Starting %v on port :%v", server.name, server.port)
		if datadog {
			server.router.Handle("/", datadogTraceMiddleware(server.router, notFound(time.Now()), os.Getenv("DD_SERVICE")))
			server.router.Handle("/healthz", datadogTraceMiddleware(server.router, healthz(failHealth), os.Getenv("DD_SERVICE")))
		} else {
			server.router.Handle("/", notFound(time.Now()))
			server.router.Handle("/healthz", healthz(failHealth))
		}


		// allow rate of `rate` requests per second and disallow initial burst
		rateLimit := timerate.NewLimiter(timerate.Limit(rate), 1)

		// instantiate client
		client := newClient(rateLimit)
		ctx := context.Background()

		go func() {
			for {
				err := client.Ratelimiter.Wait(ctx)
				if err == nil { // This is a blocking call. Honors the rate limit
					client.Do(urlString, fail, server.logger)
				}
			}
		}()

		server.Serve()
	},
}

func init() {
	rootCmd.AddCommand(workerCmd)

	// Define flags
	workerCmd.Flags().StringP("url", "u", "http://localhost", "target URL")
	workerCmd.Flags().IntP("health-port", "P", 8081, "worker healthcheck Port")
	workerCmd.Flags().IntP("rate", "r", 1, "rate of requests per second")
	workerCmd.Flags().IntP("port", "p", 8080, "target port")
	workerCmd.Flags().IntP("fail", "f", 0, "% of requests to fail, ex 10 = 10%")
	workerCmd.Flags().IntP("health-fail", "F", 0, "% of requests to /healthz to fail, ex 10 = 10%")
}
