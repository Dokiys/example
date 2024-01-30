package middleware

import (
	"fmt"
	"net/http"
	"runtime"
	"testing"
	"time"

	"golang.org/x/time/rate"
)

type ServerMiddleware func(http.Handler) http.Handler

type Router struct {
	middlewareChain []ServerMiddleware
	mux             map[string]http.Handler
}

func NewRouter() *Router {
	return &Router{
		mux: make(map[string]http.Handler),
	}
}
func (r *Router) Use(m ServerMiddleware) {
	r.middlewareChain = append(r.middlewareChain, m)
}
func (r *Router) Add(route string, h http.Handler, middlewares ...ServerMiddleware) {
	var mergedHandler = h

	for i := len(r.middlewareChain) - 1; i >= 0; i-- {
		mergedHandler = r.middlewareChain[i](mergedHandler)
	}
	for i := len(middlewares) - 1; i >= 0; i-- {
		mergedHandler = middlewares[i](mergedHandler)
	}

	r.mux[route] = mergedHandler
}

func NewServerLogMiddleware() ServerMiddleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Println(time.Now().UnixMilli(), ": server start....2")
			time.Sleep(1 * time.Millisecond)
			next.ServeHTTP(w, r)
			time.Sleep(1 * time.Millisecond)
			fmt.Println(time.Now().UnixMilli(), ": server end....3")
		})
	}
}
func NewServerRecoverMiddleware() ServerMiddleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					fmt.Println(rec)
					var stackInfo = make([]byte, 1<<16) // 64k
					idx := runtime.Stack(stackInfo, false)
					fmt.Println(string(stackInfo[0:idx]))
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
func NewServerCorsMiddleware() ServerMiddleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Set CORS headers for all responses
			w.Header().Set("Access-Control-Allow-Origin", "*")

			// Handle pre-flight OPTIONS request
			if r.Method == http.MethodOptions {
				// Include other necessary headers for pre-flight response
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, HEAD, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Length, Content-Type")
				w.Header().Set("Access-Control-Max-Age", "86400") // Cache pre-flight response for 24 hours

				// Respond with HTTP 204 No Content status
				w.WriteHeader(http.StatusNoContent)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
func NewServerLimiterMiddleware() ServerMiddleware {
	var limit = rate.NewLimiter(rate.Every(time.Second), 1)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if limit.Allow() {
				next.ServeHTTP(w, r)
			}
			w.WriteHeader(http.StatusNotAcceptable)
		})
	}
}
func TestHttpServerMiddleware(t *testing.T) {
	r := NewRouter()
	r.Use(NewServerRecoverMiddleware())
	r.Use(NewServerLogMiddleware())
	r.Use(NewServerCorsMiddleware())
	r.Use(NewServerLimiterMiddleware())
	r.Add("/hello", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello Work!"))
		// }), traceServerMiddleWare())
	}))
	r.Add("/panic", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("test recover middleware")
	}))

	s := http.NewServeMux()
	for pattern, handler := range r.mux {
		s.Handle(pattern, handler)
	}
	fmt.Println("Server listening on port 8080")
	err := http.ListenAndServe(":8080", s)
	if err != nil {
		fmt.Println(err)
	}
}
func TestHttpServerMiddleware2(t *testing.T) {
	s := http.NewServeMux()
	s.Handle("/hello", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers for all responses
		w.Header().Set("Access-Control-Allow-Origin", "*")

		// Handle pre-flight OPTIONS request
		if r.Method == http.MethodOptions {
			// Include other necessary headers for pre-flight response
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, HEAD, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Length, Content-Type")
			w.Header().Set("Access-Control-Max-Age", "86400") // Cache pre-flight response for 24 hours

			// Respond with HTTP 204 No Content status
			w.WriteHeader(http.StatusNoContent)
			return
		}
		w.Write([]byte("Hello Work!"))
	}))

	fmt.Println("Server listening on port 8080")
	err := http.ListenAndServe(":8080", s)
	if err != nil {
		fmt.Println(err)
	}
}

// ClientMiddleware is HTTP Client transport ServerMiddleware.
type ClientMiddleware func(http.RoundTripper) http.RoundTripper

// Chain returns a ClientMiddleware that specifies the chained handler for endpoint.
func Chain(rt http.RoundTripper, middlewares ...ClientMiddleware) http.RoundTripper {
	if rt == nil {
		rt = http.DefaultTransport
	}

	for i := len(middlewares) - 1; i >= 0; i-- {
		rt = middlewares[i](rt)
	}

	return rt
}

// RoundTripFunc 类似于http.HandlerFunc。由于 http.RoundTripper 是一个interface，因此需要一个struct
// 来实现RoundTrip()方法。RoundTripFunc类型实现了该方法，以便我们将一个匿名方法
// 转换成 ClientMiddleware。
type RoundTripFunc func(*http.Request) (*http.Response, error)

func (rt RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return rt(req)
}

func NewClientLogMiddleware() ClientMiddleware {
	return func(tripper http.RoundTripper) http.RoundTripper {
		return RoundTripFunc(func(req *http.Request) (*http.Response, error) {
			fmt.Println(time.Now().UnixMilli(), ": client start....1")
			time.Sleep(1 * time.Millisecond)
			resp, err := tripper.RoundTrip(req)
			time.Sleep(1 * time.Millisecond)
			fmt.Println(time.Now().UnixMilli(), ": client end....4")

			return resp, err
		})
	}
}
func TestHttpClientMiddleware(t *testing.T) {
	defaultTransport := http.DefaultTransport
	middlewares := []ClientMiddleware{
		NewClientLogMiddleware(),
	}

	client := &http.Client{
		Transport: Chain(defaultTransport, middlewares...),
	}

	_, err := client.Get("http://localhost:8080/hello")
	if err != nil {
		fmt.Println(err)
	}
}
