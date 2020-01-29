package server

import (
	"github.com/gin-gonic/gin"

	"github.com/xanderflood/fruit-pi-server/cmd/api/server/auth"
	"github.com/xanderflood/fruit-pi-server/internal/pkg/db"
	"github.com/xanderflood/fruit-pi-server/lib/tools"
)

//Server is the gin server interface for the public API
//go:generate counterfeiter . Server
type Server interface {
	// iot api
	GetDeviceConfig(c *gin.Context)
	InsertReading(c *gin.Context)

	// admin api
	RegisterDevice(c *gin.Context)
	ConfigureDevice(c *gin.Context)
	GetDeviceTokenFor(c *gin.Context)

	// authorization code
	BackendAuthorizationMiddleware(c *gin.Context)
}

//ServerAgent implements Server
type ServerAgent struct {
	logger tools.Logger

	authorize        auth.Getter
	dbClient         db.DB
	jwtSigningSecret string

	backendJWTMiddleware gin.HandlerFunc
}

//AddRoutes accepts a *gin.Engine and adds all the
//necessary routes to it for this API.
func AddRoutes(e *gin.Engine, a Server) {
	backend := e.Group("/api/v1", a.BackendAuthorizationMiddleware)
	backend.POST("/register-device", a.RegisterDevice)
	backend.POST("/configure-device", a.ConfigureDevice)
	backend.POST("/insert-reading", a.InsertReading)
	backend.GET("/get-device-config", a.GetDeviceConfig)
	backend.GET("/get-device-config/:uuid", a.GetDeviceConfig)
	backend.POST("/get-device-token", a.GetDeviceTokenFor)
}

//NewServer creates a new Server.
func NewServer(
	logger tools.Logger,

	authMgr auth.AuthorizationManager,
	authorize auth.Getter,
	dbClient db.DB,
	jwtSigningSecret string,
) ServerAgent {
	return ServerAgent{
		logger: logger,

		authorize:        authorize,
		dbClient:         dbClient,
		jwtSigningSecret: jwtSigningSecret,

		backendJWTMiddleware: authMgr.BackendMiddleware(),
	}
}

//BackendAuthorizationMiddleware callback for the backend authorization middleware
func (a ServerAgent) BackendAuthorizationMiddleware(c *gin.Context) {
	a.backendJWTMiddleware(c)
}
