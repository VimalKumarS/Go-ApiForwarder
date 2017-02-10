package Utility

import (
	"errors"
	"gateway/model"
	"time"
)

//AuthenticateModel -struct represent authenticate model
type AuthenticateModel struct {
	Conf       *appsetting.AppSetting
	Token      string
	URL        string
	HTTPMethod string
	User       string
}

//Auth -struct represent return from authorize
type Auth struct {
	User   string
	System string
	Roles  []string
}

//GateWayAuthorization interface
type GateWayAuthorization interface {
	Authenticate() error
	Authorize() (*Auth, error)
}

//Authenticate - user using token
func (val *AuthenticateModel) Authenticate() error {
	//Todo: Call Sql get the authneticate user
	sqldb := SQLDB{Conf: val.Conf}
	sqldb.InitDB()

	if sqldb.Error != nil {
		return errors.New(sqldb.Error.Error()) // db failure
	}

	//defer close db connection
	defer sqldb.CloseDB()

	userInfo := sqldb.GetUser(val.Token)
	if userInfo.Error != nil {
		return errors.New(userInfo.Error.Error()) // Fetching row error
	}
	//validate Expire date time
	if !ValidateTokenExpireDate(userInfo.ExpireAt) {
		return errors.New("Invalid Token, UnAuthorized") // Fetching row error
	}
	//Authorized user
	val.User = userInfo.User

	return nil
	//errors.New("Not Authorized")
}

//Authorize - user have access to api
func (val *AuthenticateModel) Authorize() (*Auth, error) {

	//Todo: Call Sql get the authorize role for user
	sqldb := SQLDB{Conf: val.Conf}
	sqldb.InitDB()

	if sqldb.Error != nil {
		return nil, errors.New(sqldb.Error.Error()) // db failure
	}

	//defer close db connection
	defer sqldb.CloseDB()

	IsValid, err := sqldb.GetUserAPIAccessInfo(val.User, val.URL, val.HTTPMethod)
	if err != nil {
		return nil, errors.New(sqldb.Error.Error())
	}

	if IsValid {
		auth := Auth{System: "Warp", User: val.User, Roles: []string{"role1", "role2"}}
		return &auth, nil
	}
	return nil, errors.New("Invalid Token, UnAuthorized")
}

//ValidateTokenExpireDate : validate the token expire time
func ValidateTokenExpireDate(expireAt string) bool {
	RFC3339local := "2006-01-02T15:04:05Z" // time reference for format
	t, _ := time.Parse(RFC3339local, expireAt)
	now := time.Now()
	if t.After(now) {
		return true // Valid token expire time
	}
	return false // invalid token expire time
}
