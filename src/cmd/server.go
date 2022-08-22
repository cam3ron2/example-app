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
	"github.com/spf13/cobra"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Starts a server instance",
	Run: func(cmd *cobra.Command, args []string) {
		port, _ := cmd.Flags().GetInt("port")
		delay, _ := cmd.Flags().GetInt("delay")
		fail, _ := cmd.Flags().GetInt("fail")
		failHealth, _ := cmd.Flags().GetInt("health-fail")
		datadog, _ := cmd.Flags().GetBool("datadog")
		// Initialize DataDog tracing
		if datadog {
			tracer.Start()
		}
		server := &Server{
			port: port,
			name: "Server",
			datadog: datadog,
		}

		server.logger = server.NewLogger()
		server.router = server.NewRouter()
		server.logger.Printf("Starting %v on port :%v", server.name, server.port)
		if datadog {
			server.router.Handle("/", datadogTraceMiddleware(server.router, index(delay, fail)))
			server.router.Handle("/healthz", datadogTraceMiddleware(server.router, healthz(failHealth)))
		} else {
			server.router.Handle("/", index(delay, fail))
			server.router.Handle("/healthz", healthz(failHealth))
		}
		
		server.Serve()
		if datadog {
			defer tracer.Stop()
		}
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	// Define flags
	serverCmd.Flags().IntP("delay", "d", 0, "response delay in ms")
	serverCmd.Flags().IntP("port", "p", 8080, "port to listen on")
	serverCmd.Flags().IntP("fail", "f", 0, "% of requests to fail, ex 10 = 10%")
	serverCmd.Flags().IntP("health-fail", "F", 0, "% of requests to /healthz to fail, ex 10 = 10%")
}
