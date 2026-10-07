package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type userIDKey struct{}

func UserID(ctx context.Context) uuid.UUID {
	id, _ := ctx.Value(userIDKey{}).(uuid.UUID)
	return id
}

func RequireJWT(api huma.API, secret []byte) func(huma.Context, func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		parts := strings.Fields(ctx.Header("Authorization"))
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			_ = huma.WriteErr(api, ctx, http.StatusUnauthorized, "missing or invalid bearer token")
			return
		}
		token, err := jwt.Parse(parts[1], func(*jwt.Token) (any, error) {
			return secret, nil
		}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}), jwt.WithExpirationRequired())
		if err != nil || !token.Valid {
			_ = huma.WriteErr(api, ctx, http.StatusUnauthorized, "invalid token")
			return
		}
		sub, err := token.Claims.GetSubject()
		if err != nil {
			_ = huma.WriteErr(api, ctx, http.StatusUnauthorized, "invalid token subject")
			return
		}
		id, err := uuid.Parse(sub)
		if err != nil {
			_ = huma.WriteErr(api, ctx, http.StatusUnauthorized, "invalid token subject")
			return
		}
		next(huma.WithContext(ctx, context.WithValue(ctx.Context(), userIDKey{}, id)))
	}
}
