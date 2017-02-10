package Utility

import (
	"encoding/json"
	"errors"
	"gateway/model"
	"io/ioutil"
	"net/http"
	"strings"
)

// CallingGatewayHandler struct type
type CallingGatewayHandler struct {
	Conf appsetting.AppSetting
	H    func(*appsetting.AppSetting, http.ResponseWriter, *http.Request) (int, error)
}

//APIHandlerGateway Handle the incomming  request at gateway
//Return stat code and error if any
func APIHandlerGateway(conf *appsetting.AppSetting, w http.ResponseWriter, req *http.Request) (int, error) {
	Log.Println(req.URL.Path)
	Log.Println(req.Method)

	//Todo: check of Token
	if req.Header.Get(conf.GatewaySetting.Headers.AuthToken) == "" {
		return http.StatusUnauthorized, errors.New("Unauthourized")
	}
	authToken := req.Header.Get(conf.GatewaySetting.Headers.AuthToken)

	authMode := AuthenticateModel{Conf: conf, Token: authToken, HTTPMethod: req.Method, URL: req.URL.Path}
	if err := authMode.Authenticate(); err != nil { //could be multiple reason
		return http.StatusUnauthorized, errors.New("Unauthourized")
	}

	auth, err := authMode.Authorize()
	if err != nil || auth == nil { //could be multiple reason
		return http.StatusUnauthorized, errors.New("Unauthourized")
	}

	url := req.URL.Path //strings.Split(r.URL.String(), "/")
	service := ServiceWebClient{URL: conf.GatewaySetting.URL}

	req.Header.Set("auth-system", auth.System)
	req.Header.Set("auth-user", auth.User)
	req.Header.Set("auth-user-roles", strings.Join(auth.Roles, ",")) // Todo: Role part is not implemented

	response := service.SendCommand(req.Method, url, req.Body, req.Header)
	if response.StatusCode == 404 {
		return http.StatusNotFound, errors.New(response.Status)
	}
	payload, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	var res interface{}
	json.Unmarshal(payload, &res)

	w.Header().Set("Content-Type", response.Header.Get("Content-Type"))
	if _, err = w.Write(payload); err != nil {
		return http.StatusInternalServerError, err
	}
	//formatter.JSON(w, response.StatusCode, res)
	//io.WriteString(w, "Hello World")
	return response.StatusCode, nil

}

func (ah CallingGatewayHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Updated to pass ah.appContext as a parameter to our handler type.
	status, err := ah.H(&ah.Conf, w, r)
	if err != nil {
		Log.Printf("HTTP %d: %q", status, err)
		Log.Println("HTTP %d: %q", status, err)
		switch status {
		case http.StatusNotFound:
			http.NotFound(w, r)
		case http.StatusInternalServerError:
			http.Error(w, err.Error(), status)
		default:
			http.Error(w, http.StatusText(status), status)
		}
	}
}
