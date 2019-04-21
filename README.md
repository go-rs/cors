# Middleware: CORS Request

## Reference 
- https://fetch.spec.whatwg.org/#http-cors-protocol
- https://en.m.wikipedia.org/wiki/Cross-origin_resource_sharing

## Config
```
type Config struct {
	Origin        []string
	Methods       []string
	Headers       []string
	ExposeHeaders []string
	Credentials   bool
	MaxAge        time.Duration
}

// Default
var _config = Config{
	Origin:      []string{"*"},
	Methods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "HEAD", "PATCH"},
	Headers:     []string{"Content-Type"},
	Credentials: false,
	MaxAge:      time.Hour,
}
```

## How to use?

```
var api rest.API

config := cors.Config{
    Methods: []string{"GET", "POST"},
    Credentials: true,
    MaxAge: 6 * time.Hour,
}

api.Use(cors.Load(config))

```