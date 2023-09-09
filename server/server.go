package server

import (
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/jayendramadaram/port-wardens/model"
)

type Server struct {
	router *gin.Engine
	store  Store
	auth   Auth
	secret string
}

func NewServer(store Store, auth Auth, secret string) *Server {
	return &Server{
		router: gin.Default(),
		store:  store,
		secret: secret,
		auth:   auth,
	}
}

func (s *Server) Run() {
	authRoutes := s.router.Group("/")
	authRoutes.Use(s.authenticateJWT)

	s.router.POST("/signup", s.Signup()) // Manual Test done
	s.router.POST("/login", s.Login())   // Manual Test done
	s.router.GET("/health", s.HealthCheck())

	{
		// authRoutes.GET("/dashboard", c)
	}

	s.router.Run(":8080")
}

type Store interface {
	HealthCheck() error
}

type Auth interface {
	SignUp(user model.CreateUser) error
	Login(user model.LoginUser) (*jwt.Token, error)
}

func (s *Server) Signup() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user := model.CreateUser{}
		if err := ctx.ShouldBindJSON(&user); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := s.auth.SignUp(user); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
	}
}

func (s *Server) Login() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user := model.LoginUser{}
		if err := ctx.ShouldBindJSON(&user); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		token, err := s.auth.Login(user)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		tokenString, err := token.SignedString([]byte(s.secret))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"token": tokenString})
	}
}

func (s *Server) authenticateJWT(ctx *gin.Context) {
	tokenString := ctx.GetHeader("Authorization")
	if tokenString == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Missing authorization token"})
		ctx.Abort()
		return
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Invalid signing method")
		}

		return []byte(s.secret), nil
	})

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		ctx.Abort()
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if idClaim, exists := claims["id"]; exists {
			// fmt.Println("User ID: ", idClaim.(float64), "Authorized 1")
			ctx.Set("userID", idClaim)

		} else {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			ctx.Abort()
			return
		}
	} else {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
		ctx.Abort()
		return
	}

	ctx.Next()
}

func (s *Server) HealthCheck() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if err := s.store.HealthCheck(); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"status": "online"})
	}
}
