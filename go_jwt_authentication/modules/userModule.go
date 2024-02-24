package modules

import (
	databases "go_jwt_authenticaion/databases"
	"time"

	"go.mangodb.org./mango-driver/bson/primitive"
)

type User struct {
	Id            primitive.ObjectId `bson:"_id"`
	First_Name    *string            `json:"First_Name" validate:"required,min=2,max=100"`
	Last_Name     *string            `json:"Last_Name" validate:"required,min=2,max=100"`
	Password      *string            `json:"Password" validate:"required,min=6"`
	Email         *string            `json:"Email" validate:"required,min=2,max=100"`
	Phone         *string            `json:"Phone" validate:"required"`
	Token         *string            `json:"Token"`
	User_Type     *string            `json:"User_Type" validate:"required,eq=ADMIN|eq=USER"`
	Refresh_Token *string            `json:"Refresh_Token"`
	Created_at    time.Time          `json:"Created_at"`
	Updated_at    time.Time          `json:"Updated_at"`
	User_Id       *string            `json:"User_Id"`
}
