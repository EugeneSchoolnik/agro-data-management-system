package handler

import (
    "agro-data-management-system/internal/config"
    "net/http"
    "strings"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v5"
    "golang.org/x/crypto/bcrypt"
    "go.uber.org/zap"
)

var authCfg config.AuthConfig

// SetAuthConfig sets authentication configuration used by handlers/middleware.
func SetAuthConfig(c config.AuthConfig) {
    authCfg = c
    if authCfg.TokenTTL == 0 {
        authCfg.TokenTTL = 24 * time.Hour
    }
}

// login handles authentication and returns a JWT token.
func (h *Handler) login(c *gin.Context) {
    var req struct {
        Email    string `json:"email"`
        Password string `json:"password"`
    }

    if err := c.BindJSON(&req); err != nil {
        h.newErrorResponse(c, http.StatusBadRequest, "invalid request body")
        return
    }

    // Try to find user in DB first
    if h.services != nil && h.services.User != nil {
        user, err := h.services.User.GetByEmail(req.Email)
        if err == nil && user != nil {
            // Expecting bcrypt password hash in DB
            if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err == nil {
                // valid
                token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
                    "email": user.Email,
                    "role":  user.Role,
                    "exp":   time.Now().Add(authCfg.TokenTTL).Unix(),
                })

                secret := []byte(authCfg.JWTSecret)
                if len(secret) == 0 {
                    h.log.Warn("JWT secret is empty; using fallback (insecure)")
                    secret = []byte("insecure-development-secret")
                }

                tokenString, err := token.SignedString(secret)
                if err != nil {
                    h.log.Error("failed to sign token", zap.Error(err))
                    h.newErrorResponse(c, http.StatusInternalServerError, "failed to create token")
                    return
                }

                c.JSON(http.StatusOK, gin.H{"token": tokenString})
                return
            }
            // wrong password -> continue to fallback
        }
        // If error or not found, fallthrough to config fallback
    }

    // Fallback: Basic checks against configured admin credentials
    if authCfg.AdminEmail == "" || authCfg.AdminPassword == "" {
        h.log.Warn("Auth config not set; rejecting login attempt")
        h.newErrorResponse(c, http.StatusInternalServerError, "authentication not configured")
        return
    }

    if !strings.EqualFold(strings.TrimSpace(req.Email), strings.TrimSpace(authCfg.AdminEmail)) {
        h.newErrorResponse(c, http.StatusUnauthorized, "invalid credentials")
        return
    }

    // Support bcrypt hashed password in config or plain text
    passMatch := false
    if strings.HasPrefix(authCfg.AdminPassword, "$2a$") || strings.HasPrefix(authCfg.AdminPassword, "$2b$") || strings.HasPrefix(authCfg.AdminPassword, "$2y$") {
        if err := bcrypt.CompareHashAndPassword([]byte(authCfg.AdminPassword), []byte(req.Password)); err == nil {
            passMatch = true
        }
    } else {
        if subtleConstantTimeCompare(authCfg.AdminPassword, req.Password) {
            passMatch = true
        }
    }

    if !passMatch {
        h.newErrorResponse(c, http.StatusUnauthorized, "invalid credentials")
        return
    }

    // Create token
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "email": req.Email,
        "exp":   time.Now().Add(authCfg.TokenTTL).Unix(),
    })

    secret := []byte(authCfg.JWTSecret)
    if len(secret) == 0 {
        h.log.Warn("JWT secret is empty; using fallback (insecure)")
        secret = []byte("insecure-development-secret")
    }

    tokenString, err := token.SignedString(secret)
    if err != nil {
        h.log.Error("failed to sign token", zap.Error(err))
        h.newErrorResponse(c, http.StatusInternalServerError, "failed to create token")
        return
    }

    c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

// subtleConstantTimeCompare compares two strings in constant time.
func subtleConstantTimeCompare(a, b string) bool {
    if len(a) != len(b) {
        return false
    }
    // Use byte-wise comparison to avoid importing crypto/subtle directly
    var res byte = 0
    for i := 0; i < len(a); i++ {
        res |= a[i] ^ b[i]
    }
    return res == 0
}
