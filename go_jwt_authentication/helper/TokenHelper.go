package helper

import (
	"context"
	"fmt"
	"go_jwt_authentication/databases"
	"log"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SignedDetails struct {
	Email      string
	First_Name string
	Last_Name  string
	Uid        string
	User_Type  string
	jwt.StandardClaims
}

var userCollection *mongo.Collection = databases.OpenCollection(databases.Client, "user")
var SECRET_KEY = os.Getenv("SECRET_KEY")

func GenerateAllTokens(email string, firstName string, lastName string, userType string, uid string) (signedToken string, Signed_Refresh_Token string, err error) {

	claims := &SignedDetails{

		Email:      email,
		First_Name: firstName,
		Last_Name:  lastName,
		Uid:        uid,
		User_Type:  userType,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}

	refreshClaims := &SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(168)).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))

	refresh_token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(SECRET_KEY))

	if err != nil {
		log.Panic(err)
		return
	}
	return token, refresh_token, err
}

func ValidateToken(signedToken string) (claims *SignedDetails, msg string) {

	token, err := jwt.ParseWithClaims(

		signedToken,
		&SignedDetails{},
		func(t *jwt.Token) (interface{}, error) {
			return []byte(SECRET_KEY), nil
		},
	)

	if err != nil {
		msg = err.Error()
		return
	}
	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		msg = fmt.Sprintf("The Token Is Invalid")
		msg = err.Error()
		return

	}

	if claims.ExpiresAt < time.Now().Local().Unix() {

		msg = fmt.Sprintf("Token Is Expired")
		msg = err.Error()
		return
	}
	return claims, msg
}

func UpdateAllTokens(signedToken string, signedRefreshToken string, userId string) {

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

	var updateObj primitive.D

	updateObj = append(updateObj, bson.E{Key: "token", Value: signedToken})
	updateObj = append(updateObj, bson.E{Key: "refreshToken", Value: signedRefreshToken})

	Updated_at, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	updateObj = append(updateObj, bson.E{Key: "updated_at", Value: Updated_at})

	upsert := true
	filter := bson.M{"user_Id": userId}
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}

	_, err := userCollection.UpdateOne(
		ctx,
		filter,
		bson.D{
			{Key: "$set", Value: updateObj},
		},
		&opt,
	)
	defer cancel()
	if err != nil {
		log.Panic(err)
		return
	}
	return
}
