package handlers

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/api/idtoken"
)

const header = "Authorization"
const prefix = "Bearer "
const testToken = "TOKEN"
const testUserProvider = "Google"
const testUserId = "1"

func AuthenticateGoogle() gin.HandlerFunc {
  return func (c *gin.Context) {
		if strings.ToLower(os.Getenv("AUTH_ENABLED")) == "false" {
			log.Printf("Auth is not enabled")
			c.JSON(http.StatusOK, gin.H{
				"sessionToken": testToken,
			})
			return
		}

		var body struct {
			Token string `json:"token"`
		}
		if err := c.BindJSON(&body); err != nil || body.Token == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing or invalid token in request body"})
			return
		}

		audience := os.Getenv("AUTH_GOOGLE_CLIENT_ID")
		payload, err := idtoken.Validate(context.Background(), body.Token, audience)
		if err != nil {
			log.Printf("Token validation error: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid ID token"})
			return
		}

		// Extract claims
		subVal, ok := payload.Claims["sub"].(string)
		if !ok {
			log.Printf("sub claim is not a string")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid ID token"})
			return
		}

		// Generate session token
		sessionToken, err := generateSessionToken("Google", subVal, time.Hour, os.Getenv("AUTH_SECRET_KEY"))
		if err != nil {
			log.Printf("Session token generation error: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid ID token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"sessionToken": sessionToken,
		})
	}
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if strings.ToLower(os.Getenv("AUTH_ENABLED")) == "false" {
			log.Printf("Auth is not enabled")
			c.Set("user_provider", testUserProvider)
			c.Set("user_id", testUserId)
			c.Next()
			return
		}

		// Extract token from header
		token := extractSessionToken(c)

		// Extract claims 
		claims, err := verifySessionToken(token, os.Getenv("AUTH_SECRET_KEY"))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid session token"})
			return
		}

		// Add claims to context
		provider, ok := claims["provider"].(string)
		if !ok {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "provider claim is not a string"})
				return
		}
		
		sub, ok := claims["sub"].(string)
		if !ok {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "sub claim is not a string"})
				return
		}

		c.Set("user_provider", provider)
		c.Set("user_id", sub)
		c.Next()
	}
}

func generateSessionToken(provider string, sub string, expiry time.Duration, secretBase64 string) (string, error) {
	claims := jwt.MapClaims{
		"provider": provider,
		"sub":      sub,
		"iat":      time.Now().Unix(),
		"exp":      time.Now().Add(expiry).Unix(),
		"scope":    "session",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secretKey, err := getSecretKey(secretBase64)
	if err != nil {
		return "", err
	}

	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return signedToken, nil
}

func extractSessionToken(c *gin.Context) string {
	// Extract from Authorization header (Bearer ...)
	authHeader := c.GetHeader(header)
	if authHeader == "" {
		return ""
	}

	if len(authHeader) > len(prefix) && authHeader[:len(prefix)] == prefix {
		return authHeader[len(prefix):]
	}

	return ""
}

func verifySessionToken(token string, secretBase64 string) (map[string]interface{}, error) {
	// Get secret key
	secretKey, err := getSecretKey(secretBase64)
	if err != nil {
		return nil, err
	}

	// Extract jwt from token
	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	// Verify the claims
	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		if scope, ok := claims["scope"].(string); !ok || scope != "session" {
			return nil, fmt.Errorf("invalid scope claim")
		}

		return claims, nil
	}

	return nil, fmt.Errorf("invalid token or claims")
}

func getSecretKey(secretBase64 string) ([]byte, error) {
	if secretBase64 == "" {
		return nil, fmt.Errorf("AUTH_SECRET_KEY not set in environment")
	}

	secretKey, err := base64.StdEncoding.DecodeString(secretBase64)
	if err != nil {
		return nil, fmt.Errorf("failed to decode secret key: %w", err)
	}

	return secretKey, nil
}
