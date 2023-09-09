package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/jayendramadaram/port-wardens/model"
	"github.com/jayendramadaram/port-wardens/server"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type Auth interface {
	server.Auth
}

type auth struct {
	dbUser *mongo.Database
}

type Claims struct {
	Id string `json:"id"`
	jwt.StandardClaims
}

func NewAuth(dbUser *mongo.Database) Auth {
	return &auth{dbUser: dbUser}
}

func (a *auth) SignUp(user model.CreateUser) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("error hashing password: %w", err)
	}

	User := model.User{
		Email:    user.Email,
		Username: user.Username,
		Password: string(hashedPassword),
	}

	_, err = a.dbUser.Collection("users").InsertOne(context.TODO(), User)
	if err != nil {
		return fmt.Errorf("error inserting user: %w", err)
	}
	return nil
}

func (a *auth) Login(user model.LoginUser) (*jwt.Token, error) {
	var UserRecord model.User
	if err := a.dbUser.Collection("users").FindOne(context.TODO(), bson.M{"email": user.Email}).Decode(&UserRecord); err != nil {
		return nil, fmt.Errorf("invalid email %w", err)
	}

	err := bcrypt.CompareHashAndPassword([]byte(UserRecord.Password), []byte(user.Password))
	if err != nil {
		return nil, fmt.Errorf("invalid password %w", err)
	}

	claims := &Claims{
		Id: UserRecord.ID.Hex(),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
			IssuedAt:  time.Now().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token, nil
}
