/*!
 * go-rs/cors
 * Copyright(c) 2019 Roshan Gade
 * MIT Licensed
 */
package cors

// Reference to: https://en.m.wikipedia.org/wiki/Cross-origin_resource_sharing
// Reference to: https://fetch.spec.whatwg.org/#http-cors-protocol
import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/go-rs/rest-api-framework"
)

var (
	OriginNotAllowed  = errors.New("ORIGIN_NOT_ALLOWED")
	HeadersNotAllowed = errors.New("HEADERS_NOT_ALLOWED")
	MethodNotAllowed  = errors.New("METHOD_NOT_ALLOWED")
)

/**
 * An HTTP response to a CORS request can include the following headers:
 *
 * `Access-Control-Allow-Origin`
 * Indicates whether the response can be shared, via returning the literal value of the `Origin` request header (which can be `null`) or `*` in a response.
 *
 * `Access-Control-Allow-Credentials`
 * Indicates whether the response can be shared when request’s credentials mode is "include".
 *
 * For a CORS-preflight request, request’s credentials mode is always "omit", but for any subsequent CORS requests it might not be. Support therefore needs to be indicated as part of the HTTP response to the CORS-preflight request as well.
 *
 * An HTTP response to a CORS-preflight request can include the following headers:
 *
 * `Access-Control-Allow-Methods`
 * Indicates which methods are supported by the response’s URL for the purposes of the CORS protocol.
 *
 * The `Allow` header is not relevant for the purposes of the CORS protocol.
 *
 * `Access-Control-Allow-Headers`
 * Indicates which headers are supported by the response’s URL for the purposes of the CORS protocol.
 *
 * `Access-Control-Max-Age`
 * Indicates how long the information provided by the `Access-Control-Allow-Methods` and `Access-Control-Allow-Headers` headers can be cached.
 *
 * An HTTP response to a CORS request that is not a CORS-preflight request can also include the following header:
 *
 * `Access-Control-Expose-Headers`
 * Indicates which headers can be exposed as part of the response by listing their names.
 *
 */

type Config struct {
	Origin        []string
	Methods       []string
	Headers       []string
	ExposeHeaders []string
	Credentials   bool
	MaxAge        time.Duration
}

var _config = Config{
	Origin:      []string{"*"},
	Methods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "HEAD", "PATCH"},
	Headers:     []string{"Content-Type"},
	Credentials: false,
	MaxAge:      time.Hour,
}

/**
 * Merge user config with default
 */
func merge(source Config, target *Config) {
	if target.Origin == nil {
		target.Origin = source.Origin
	}
	if target.Methods == nil {
		target.Methods = source.Methods
	}
	if target.Headers == nil {
		target.Headers = source.Headers
	}
	if target.MaxAge == 0 {
		target.MaxAge = source.MaxAge
	}
}

/**
 * Search string in slice
 */
func hasMatch(data []string, str string) bool {
	for _, v := range data {
		if v == str {
			return true
		}
	}
	return false
}

/**
 * Value should be included
 */
func hasInclude(data []string, val []string) bool {
	out := make(map[string]bool)
	if len(data) < len(val) {
		return false
	}

	for _, d := range data {
		out[d] = true
	}

	for _, v := range val {
		if !out[v] {
			return false
		}
	}

	return true
}

/**
 * A CORS-preflight request is a CORS request that checks to see if the CORS protocol is understood. It uses `OPTIONS` as method and includes these headers:
 *
 * `Access-Control-Request-Method`
 * Indicates which method a future CORS request to the same resource might use.
 *
 * `Access-Control-Request-Headers`
 * Indicates which headers a future CORS request to the same resource might use.
 */
func corsPreFlightRequest(ctx *rest.Context, config Config) {
	method := ctx.Request.Header.Get("Access-Control-Request-Method")
	headers := ctx.Request.Header.Get("Access-Control-Request-Headers")

	if method != "" && !hasMatch(config.Methods, method) {
		ctx.Status(403).Throw(MethodNotAllowed)
		return
	}

	if headers != "" && !hasInclude(config.Headers, strings.Split(headers, ", ")) {
		ctx.Status(403).Throw(HeadersNotAllowed)
		return
	}

	if len(config.Methods) > 0 {
		ctx.SetHeader("Access-Control-Allow-Methods", strings.Join(config.Methods, ", "))
	}

	if len(config.Headers) > 0 {
		ctx.SetHeader("Access-Control-Allow-Headers", strings.Join(config.Headers, ", "))
	}

	if config.MaxAge > time.Duration(0) {
		ctx.SetHeader("Access-Control-Max-Age", strconv.FormatInt(int64(config.MaxAge/time.Second), 10))
	}

	ctx.Status(204).Text("")
	ctx.End()

}

/**
 * Cors request
 */
func Load(config Config) rest.Handler {
	merge(_config, &config)
	allowedAllOrigins := hasMatch(_config.Origin, "*")
	return func(ctx *rest.Context) {
		origin := ctx.Request.Header.Get("Origin")
		// STEP 1: check origin
		if origin == "" {
			return
		}

		// STEP 2: validate origin
		if !allowedAllOrigins && !hasMatch(config.Origin, origin) {
			ctx.Status(403)
			ctx.Throw(OriginNotAllowed)
			return
		}

		ctx.SetHeader("Access-Control-Allow-Origin", origin)

		//check: https://fetch.spec.whatwg.org/#cors-protocol-and-credentials
		if config.Credentials {
			ctx.SetHeader("Access-Control-Allow-Credentials", "true")
		}

		// STEP 3: check request method
		if ctx.Request.Method != "OPTIONS" {
			if len(config.ExposeHeaders) > 0 {
				ctx.SetHeader("Access-Control-Allow-Headers", strings.Join(config.ExposeHeaders, ", "))
			}
			return
		}

		corsPreFlightRequest(ctx, config)
	}
}
