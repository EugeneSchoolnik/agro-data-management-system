package handler

import (
    "fmt"
    "net/http"
    "strings"

    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v5"
)

// authMiddleware validates JWT from Authorization header and aborts on failure.
func (h *Handler) authMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
            h.newErrorResponse(c, http.StatusUnauthorized, "missing or invalid authorization header")
            return
        }

        tokenString := strings.TrimPrefix(authHeader, "Bearer ")

        token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
            if t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
                return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
            }
            return []byte(authCfg.JWTSecret), nil
        })

        if err != nil || token == nil || !token.Valid {
            h.newErrorResponse(c, http.StatusUnauthorized, "invalid or expired token")
            return
        }

        // Optionally expose email claim
        if claims, ok := token.Claims.(jwt.MapClaims); ok {
            if email, ok := claims["email"].(string); ok {
                c.Set("userEmail", email)
            }
        }

        c.Next()
    }
}
