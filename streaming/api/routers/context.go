package routers

import (
	// "bytes"

	"home-media/streaming/core"

	expirable "github.com/hashicorp/golang-lru/v2/expirable"
)

// RouteContext comment
type RouteContext struct {
	Config          *core.Config
	SessionKeyVault *expirable.LRU[string, string]
}

// // GetDatabase comment
// func (ctx *RouteContext) GetDatabase(ginCtx *gin.Context) (db *gorm.DB, err error) {
// 	var organizationID string

// 	organizationID = ginCtx.Param("organizationID")
// 	if len(organizationID) == 0 {
// 		organizationID = ginCtx.Request.Header.Get("x-organization-id")
// 	}
// 	if len(organizationID) == 0 {
// 		err = errors.New("organization not found")
// 		return nil, err
// 	}

// 	customerDb := fmt.Sprintf("claimh_customer_%s", organizationID)
// 	return core.NewDatabase(ctx.Config, customerDb)
// }
