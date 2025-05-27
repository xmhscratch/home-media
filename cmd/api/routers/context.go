package routers

import (
	// "bytes"

	"home-media/sys"

	expirable "github.com/hashicorp/golang-lru/v2/expirable"
)

// RouteContext comment
type RouteContext struct {
	Config          *sys.Config
	SessionKeyVault *expirable.LRU[string, string]
}
