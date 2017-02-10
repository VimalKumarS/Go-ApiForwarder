package Utility

import (
	"database/sql"
	"errors"
	"gateway/model"
	"strings"

	_ "github.com/denisenkom/go-mssqldb"
)

type SQLDB struct {
	DB    *sql.DB
	Conf  *appsetting.AppSetting
	Error error
}

type UserInfo struct {
	TokenID  string
	User     string
	ExpireAt string
	Error    error
}

type QueryDB interface {
	InitDB()
	CloseDB()
	GetUser(string) UserInfo
	GetUserApiAccessInfo(string, string, string) (bool, error)
}

//InitDB : - Initialize Db
func (dbSeting *SQLDB) InitDB() {
	db, err := sql.Open("mssql", dbSeting.Conf.ConnectionStrings[3].ConnectionString)
	dbSeting.DB = db
	if err != nil {
		Log.Println(err.Error())
		dbSeting.Error = errors.New(err.Error())
	}
	err = dbSeting.DB.Ping()
	if err != nil {
		Log.Println(err.Error())
		dbSeting.Error = errors.New(err.Error())

	}

}

//CloseDB :  close the db connection
func (dbSeting *SQLDB) CloseDB() {
	if dbSeting.DB != nil {
		defer dbSeting.DB.Close()
	}
}

//GetUser : based on token
func (dbSeting *SQLDB) GetUser(token string) UserInfo {
	user := UserInfo{}
	rows, err := dbSeting.DB.Query("SELECT cast([Id] as char(36)),[User],ExpiresAt FROM [sec].[Token] where id=?1", token)
	if err != nil {
		user.Error = err
	}
	defer rows.Close()

	for rows.Next() {
		// Get the results of the query
		err := rows.Scan(&user.TokenID, &user.User, &user.ExpireAt)
		if err != nil {
			user.Error = err
		}
	}
	return user
}

//GetUserAPIAccessInfo : validate the user have access to
func (dbSeting *SQLDB) GetUserAPIAccessInfo(userID string, url string, httpMethod string) (bool, error) {
	var apiID, resourceURL, resourceRoot, method string
	var isValidURL = false
	rows, err := dbSeting.DB.Query(`select api.Id,api.ResourceUrl,api.ResourceRoot,api.Method from sec.Api api
									inner join sec.ApiResourceActionMap ara on api.Id = ara.ApiId
									inner join sec.UserGroupResourceActionMap ugra on ara.ResourceActionId = ugra.ResourceActionID 
									inner join  sec.UserGroupMap ugm on ugm.UserGroupID=ugra.UserGroupID
									where ugm.UserID=?1 and api.Enabled=1`, userID)

	if err != nil {
		return false, err
	}
	defer rows.Close()

	for rows.Next() {
		// Get the results of the query
		err := rows.Scan(&apiID, &resourceURL, &resourceRoot, &method)
		if err != nil {
			return false, err
		}
		if strings.ToLower(url) == strings.ToLower(resourceURL) && strings.ToLower(method) == strings.ToLower(httpMethod) {
			isValidURL = true
			break
		}
	}
	if isValidURL {
		return true, nil
	}

	return false, nil

}
