package auth

import (
	"errors"
	"net/http"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"

	"github.com/xanderflood/fruit-pi-server/lib/tools"
	"github.com/xanderflood/fruit-pi-server/pkg/db"
)

//AuthorizationContextKey is the key used to store the Authorization
//in the context
const AuthorizationContextKey = "FRUIT_PI_SERVER_PUBLIC_API_AUTHORIZATION"

//Getter is a helper for grabbing the Authorization
//that the middleware stores in the context.
//go:generate counterfeiter . Getter
type Getter func(c *gin.Context) (Authorization, bool)

//GetAuthorizationFromContext is the default Getter
func GetAuthorizationFromContext(c *gin.Context) (Authorization, bool) {
	authIface := c.Value(AuthorizationContextKey)
	auth, ok := authIface.(Authorization)
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization object was found"})
	}
	return auth, ok
}

//Authorizer represents the needed interactions with jwt.Parser
//go:generate counterfeiter . Authorizer
type Authorizer interface {
	ParseWithClaims(tokenString string, claims jwt.Claims, keyFunc jwt.Keyfunc) (*jwt.Token, error)
}

//AuthorizationManager exposes middleware functionality for authorization
type AuthorizationManager interface {
	BackendMiddleware() gin.HandlerFunc
}

//JWTAuthorizationManager provides a JWT-based implementation of AuthorizationManager
type JWTAuthorizationManager struct {
	logger        tools.Logger
	signingSecret string
	authorizer    Authorizer
	db            db.DB
}

//NewAuthorizationManager creates a new JWTAuthorizationManager
func NewAuthorizationManager(
	logger tools.Logger,
	signingSecret string,
	authorizer Authorizer,
	db db.DB,
) JWTAuthorizationManager {
	return JWTAuthorizationManager{
		logger:        logger,
		signingSecret: signingSecret,
		authorizer:    authorizer,
		db:            db,
	}
}

func (a JWTAuthorizationManager) getAuthorizationFromString(tokenString string) (Authorization, error) {
	var auth Authorization
	_, err := a.authorizer.ParseWithClaims(tokenString, &auth, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("must use HMAC signing")
		}

		return []byte(a.signingSecret), nil
	})
	if err != nil {
		return Authorization{}, err
	}

	return auth, nil
}

//BackendMiddleware checks for a JWT in a bearer token on the request
//and converts it into an Authorzation struct, which is stored in
//the context.
func (a JWTAuthorizationManager) BackendMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "no authorization provided"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		auth, err := a.getAuthorizationFromString(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		c.Set(AuthorizationContextKey, auth)
		c.Next()
	}
}
