package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"opserver/common"
	"opserver/config"
)

func VerifyTokenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		Authorization := c.Request.Header.Get("Authorization")
		if Authorization == "" {
			common.ResponseHandle(http.StatusUnauthorized, "no token", c)
			return
		}
		t, err := jwt.Parse(Authorization, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.GetValue("token", "secret")), nil
		})

		if err != nil {
			common.ResponseHandle(http.StatusUnauthorized, "token 不合法", c)
			return
		}

		if t.Valid {
			c.Next()
		} else if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				common.ResponseHandle(http.StatusUnauthorized, "无效的token", c)
				return
			} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
				// Token is either expired or not active yet
				common.ResponseHandle(http.StatusUnauthorized, "token 已过期", c)
				return
			} else {
				common.ResponseHandle(http.StatusUnauthorized, "token 不合法", c)
				return
			}
		} else {
			common.ResponseHandle(http.StatusUnauthorized, "token 不合法", c)
			return
		}
	}
}