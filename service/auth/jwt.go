package auth

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/KristianKjerstad/go-e-commerce-api/config"
	"github.com/KristianKjerstad/go-e-commerce-api/types"
	"github.com/KristianKjerstad/go-e-commerce-api/utils"
	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const UserKey contextKey = "usedID"

func CreateJWT(secret []byte, userID int) (string, error) {
	expiration := time.Second * time.Duration(config.EnvConfig.JWTExpirationInSeconds)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":    strconv.Itoa(userID),
		"expiredAt": time.Now().Add(expiration).Unix(),
	})

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func WithJWTAuth(handlerFunc http.HandlerFunc, store types.UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get token from the user request
		tokenString := getTokenFromRequest(r)
		// validate JWT
		token, err := validateToken(tokenString)
		if err != nil {
			log.Printf("Failed to validate token %v", err)
			PermissionDenied(w)
			return
		}

		if !token.Valid {
			log.Println("invalid token")
			PermissionDenied(w)
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		str := claims["userID"].(string)

		userID, err := strconv.Atoi(str)
		u, err := store.GetUserById(userID)
		if err != nil {
			log.Printf("failed to get user by id %v", err)
			PermissionDenied(w)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, UserKey, u.ID)
		r = r.WithContext(ctx)

		handlerFunc(w, r)

	}
}

func getTokenFromRequest(r *http.Request) string {
	tokenAuth := r.Header.Get("Authorization")
	const prefix = "Bearer "
	if !strings.HasPrefix(tokenAuth, prefix) {
		return "invalid token format"
	}
	tokenAuth = strings.TrimPrefix(tokenAuth, prefix)
	if tokenAuth != "" {
		return tokenAuth
	}
	return ""
}

func validateToken(t string) (*jwt.Token, error) {
	return jwt.Parse(t, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", t.Header["alg"])
		}
		return []byte(config.EnvConfig.JWTSecret), nil
	})
}

func PermissionDenied(w http.ResponseWriter) {
	utils.WriteError(w, http.StatusForbidden, fmt.Errorf("permission denied"))
}

func GetUserIDFromContext(ctx context.Context) int {
	userID, ok := ctx.Value(UserKey).(int)
	if !ok {
		return -1
	}

	return userID
}
