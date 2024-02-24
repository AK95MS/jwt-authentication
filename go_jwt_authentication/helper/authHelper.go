package helper

import (
	"errors"

	"github.com/gin-gonic/gin"
)

func CheckUserType(c *gin.Context, role string) (err error) {
	userType := c.GetString("UserType")
	err = nil

	if userType != role {
		err = errors.New("Unauthorized to access this resource")
	}
	return err

}

func MatchUserTypeUid(c *gin.Context, user_id string) (err error) {

	userType := c.GetString("user_type")

	userId := c.GetString("uid")

	err = nil

	if userType == "USER" && userId != user_id {
		err = errors.New("Unauthorized to access this resource")
		return err
	}

	err = CheckUserType(c, userType)
	return err

}
