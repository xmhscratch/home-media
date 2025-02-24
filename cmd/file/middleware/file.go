package middleware

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	// "github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/utils/v2"
	"github.com/valyala/fasthttp"
)

// Config defines the config for middleware.
type StorageConfig struct {
	// FS is the file system to serve the static files from.
	// You can use interfaces compatible with fs.FS like embed.FS, os.DirFS etc.
	//
	// Optional. Default: nil
	FS fs.FS

	// Next defines a function to skip this middleware when returned true.
	//
	// Optional. Default: nil
	Next func(c *fiber.Ctx) bool

	// ModifyResponse defines a function that allows you to alter the response.
	//
	// Optional. Default: nil
	ModifyResponse fiber.Handler

	// NotFoundHandler defines a function to handle when the path is not found.
	//
	// Optional. Default: nil
	NotFoundHandler fiber.Handler

	// The names of the index files for serving a directory.
	//
	// Optional. Default: []string{"index.html"}.
	IndexNames []string `json:"index"`

	// Expiration duration for inactive file handlers.
	// Use a negative time.Duration to disable it.
	//
	// Optional. Default: 10 * time.Second.
	CacheDuration time.Duration `json:"cache_duration"`

	// The value for the Cache-Control HTTP-header
	// that is set on the file response. MaxAge is defined in seconds.
	//
	// Optional. Default: 0.
	MaxAge int `json:"max_age"`

	// When set to true, the server tries minimizing CPU usage by caching compressed files.
	// This works differently than the github.com/gofiber/compression middleware.
	//
	// Optional. Default: false
	Compress bool `json:"compress"`

	// When set to true, enables byte range requests.
	//
	// Optional. Default: false
	ByteRange bool `json:"byte_range"`

	// When set to true, enables directory browsing.
	//
	// Optional. Default: false.
	Browse bool `json:"browse"`

	// When set to true, enables direct download.
	//
	// Optional. Default: false.
	Download bool `json:"download"`
}

// ConfigDefault is the default config
var StorageConfigDefault = StorageConfig{
	IndexNames:    []string{"index.html"},
	CacheDuration: 10 * time.Second,
}

// Helper function to set default values
func configDefault(config ...StorageConfig) StorageConfig {
	// Return default config if nothing provided
	if len(config) < 1 {
		return StorageConfigDefault
	}

	// Override default config
	cfg := config[0]

	// Set default values
	if len(cfg.IndexNames) == 0 {
		cfg.IndexNames = StorageConfigDefault.IndexNames
	}

	if cfg.CacheDuration == 0 {
		cfg.CacheDuration = StorageConfigDefault.CacheDuration
	}

	return cfg
}

func NewStorageHandler(root string, cfg ...StorageConfig) fiber.Handler {
	config := configDefault(cfg...)

	var createFS sync.Once
	var fileHandler fasthttp.RequestHandler
	var cacheControlValue string

	// adjustments for io/fs compatibility
	if config.FS != nil && root == "" {
		root = "."
	}

	return func(c *fiber.Ctx) error {
		// Don't execute middleware if Next returns true
		if config.Next != nil && config.Next(c) {
			return c.Next()
		}

		// We only serve static assets on GET or HEAD methods
		method := c.Method()
		if method != fiber.MethodGet && method != fiber.MethodHead {
			return c.Next()
		}

		// Initialize FS
		createFS.Do(func() {
			prefix := c.Route().Path

			// Is prefix a partial wildcard?
			if strings.Contains(prefix, "*") {
				// /john* -> /john
				prefix = strings.Split(prefix, "*")[0]
			}

			prefixLen := len(prefix)
			if prefixLen > 1 && prefix[prefixLen-1:] == "/" {
				// /john/ -> /john
				prefixLen--
			}

			fs := &fasthttp.FS{
				Root:               root,
				FS:                 config.FS,
				AllowEmptyRoot:     true,
				GenerateIndexPages: config.Browse,
				AcceptByteRange:    config.ByteRange,
				Compress:           config.Compress,
				CompressBrotli:     config.Compress, // Brotli compression won't work without this
				CompressedFileSuffixes: map[string]string{
					"gzip": ".fiber.gz",
					"br":   ".fiber.br",
					"zstd": ".fiber.zst",
				},
				CacheDuration: config.CacheDuration,
				SkipCache:     config.CacheDuration < 0,
				IndexNames:    config.IndexNames,
				PathNotFound: func(fctx *fasthttp.RequestCtx) {
					fctx.Response.SetStatusCode(fiber.StatusNotFound)
				},
			}

			fs.PathRewrite = func(fctx *fasthttp.RequestCtx) []byte {
				path := fctx.Path()

				if len(path) >= prefixLen {
					checkFile, err := isFile(root, fs.FS)
					if err != nil {
						return path
					}

					// If the root is a file, we need to reset the path to "/" always.
					switch {
					case checkFile && fs.FS == nil:
						path = []byte("/")
					case checkFile && fs.FS != nil:
						path = utils.UnsafeBytes(root)
					default:
						path = path[prefixLen:]
						if len(path) == 0 || path[len(path)-1] != '/' {
							path = append(path, '/')
						}
					}
				}

				if len(path) > 0 && path[0] != '/' {
					path = append([]byte("/"), path...)
				}

				return path
			}

			maxAge := config.MaxAge
			if maxAge > 0 {
				cacheControlValue = "public, max-age=" + strconv.Itoa(maxAge)
			}

			fileHandler = fs.NewRequestHandler()
		})

		// Serve file
		fileHandler(c.Context())

		// Sets the response Content-Disposition header to attachment if the Download option is true
		if config.Download {
			c.Attachment()
		}

		// Return request if found and not forbidden
		status := c.Context().Response.StatusCode()

		if status != fiber.StatusNotFound && status != fiber.StatusForbidden {
			if len(cacheControlValue) > 0 {
				c.Context().Response.Header.Set(fiber.HeaderCacheControl, cacheControlValue)
			}

			if config.ModifyResponse != nil {
				return config.ModifyResponse(c)
			}

			return nil
		}

		// Return custom 404 handler if provided.
		if config.NotFoundHandler != nil {
			return config.NotFoundHandler(c)
		}

		// Reset response to default
		c.Context().SetContentType("") // Issue #420
		c.Context().Response.SetStatusCode(fiber.StatusOK)
		c.Context().Response.SetBodyString("")

		// Next middleware
		return c.Next()
	}
}

// isFile checks if the root is a file.
func isFile(root string, filesystem fs.FS) (bool, error) {
	var file fs.File
	var err error

	if filesystem != nil {
		file, err = filesystem.Open(root)
		if err != nil {
			return false, fmt.Errorf("static: %w", err)
		}
	} else {
		file, err = os.Open(filepath.Clean(root))
		if err != nil {
			return false, fmt.Errorf("static: %w", err)
		}
	}

	stat, err := file.Stat()
	if err != nil {
		return false, fmt.Errorf("static: %w", err)
	}

	return stat.Mode().IsRegular(), nil
}
