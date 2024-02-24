package controller

import (
	"context"
	"fmt"
	databases "go_jwt_authentication/databases"
	helper "go_jwt_authentication/helper"
	"go_jwt_authentication/modules"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	//db "go_jwt_authentication/databesConnection"
)

var userCollection *mongo.Collection = databases.OpenCollection(databases.Client, "user")

var validate = validator.New()

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)

	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}

func VerifyPassword(userPassword string, providedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))

	check := true
	msg := ""

	if err != nil {
		msg = "email or password is incorrect"
		check = false
		return check, msg
	}
	return check, msg

}

func Signup() gin.HandlerFunc {
	return func(c *gin.Context) {

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

		var user modules.User

		if err := c.Bind(&user); err != nil {

			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

			return
		}

		validationErr := validate.Struct(user)

		if validationErr != nil {
			log.Fatal("Error While Validating User Request")
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		Email_count, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		defer cancel()
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Error Occured While Checking For The Email"})
		}

		password := HashPassword(*user.Password)

		user.Password = &password

		Phone_count, err := userCollection.CountDocuments(ctx, bson.M{"Phone_Number": user.Phone})
		defer cancel()
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusBadRequest, gin.H{"Error": "Error Occured While Chceking For The Phone Number"})
		}

		if Email_count > 0 || Phone_count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "This Email and Phone Number Is Already Present"})
		}

		user.Created_at, _ = time.Parse(time.RFC1123, time.Now().Format(time.RFC3339))

		user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		user.Id = primitive.NewObjectID()

		user.User_Id = user.Id.Hex()
		token, refresh_token, _ := helper.GenerateAllTokens(*user.Email, *user.First_Name, *user.Last_Name, *user.User_Type, *user.User_Id)

		user.Token = &token
		user.Refresh_Token = &refresh_token

		resultInsertionNumber, InsertErr := userCollection.InsertOne(ctx, user)

		if InsertErr != nil {
			msg := fmt.Sprint("User Item Was Not Created")
			c.JSON(http.StatusInternalServerError, gin.H{"Error": msg})
		}
		defer cancel()
		c.JSON(http.StatusOK, resultInsertionNumber)
	}
}
func Login() gin.HandlerFunc {
	return func(c *gin.Context) {

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

		var user modules.User

		var foundUser modules.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, g.H{"error": "Email and Password not match"})
			return
		}

		passwordIsValid, msg := VerifyPassword(*user.Password, *foundUser.Password)

		defer cancel()

		if passwordIsValid != true {
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		}

		if foundUser.Email == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "User Not Found"})
		}

		token, refreshToken, _ := helper.GenerateAllTokens(*foundUser.Email, *foundUser.First_Name, *foundUser.Last_Name, *foundUser.User_Type, *foundUser.User_Id)
		helper.UpdateAllTokens(token, refreshToken, foundUser.User_Id)

		err = userCollection.FindOne(ctx, bson.M{"user_id": foundUser.User_Id}).Decode(&foundUser)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, foundUser)
	}

}

func GetUsers() {

}

func GetUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId := ctx.Param("user_id")

		if err := helper.MatchUserTypeUid(ctx, userId); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

		var user modules.User

		err := userCollection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&user)
		defer cancel()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, user)
	}
}
