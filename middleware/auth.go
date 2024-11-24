package middleware

import (
	"net/http"
	"strings"

	"a21hc3NpZ25tZW50/model" // Sesuaikan path sesuai struktur project Anda

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Ambil cookie bernama session_token
		tokenString, err := ctx.Cookie("session_token")
		if err != nil {
			// Jika token tidak ada
			if strings.Contains(ctx.GetHeader("Content-Type"), "application/json") {
				ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			} else {
				ctx.Redirect(http.StatusSeeOther, "/login")
			}
			ctx.Abort()
			return
		}

		// Parse token
		claims := &model.Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return model.JwtKey, nil // Gunakan secret key dari model.JwtKey
		})

		if err != nil || !token.Valid {
			// Jika token tidak valid
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid token"})
			ctx.Abort()
			return
		}

		// Simpan UserID dari claims ke dalam context
		ctx.Set("email", claims.Email)

		// Lanjutkan ke middleware atau handler berikutnya
		ctx.Next()
	}
}
