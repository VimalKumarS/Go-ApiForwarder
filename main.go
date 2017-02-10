package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"gateway/gatewayHandler"
	"gateway/model"
	"io/ioutil"
	"net/http"
	"os"
)

var (
	logpath = flag.String("logpath", "./log/gateway.log", "Log Path")
)

func main() {

	fmt.Println("Api Gateway Started")
	file, e := ioutil.ReadFile("./config/appsetting.json")
	if e != nil {
		fmt.Printf("File read error: %v", e)
		os.Exit(1)
	}

	var configuration appsetting.AppSetting
	json.Unmarshal(file, &configuration)
	//fmt.Printf("%v", configuration.ConnectionStrings[0].ConnectionString)
	//fmt.Println(configuration.ConnectionStrings[1].ConnectionString)
	conf := &Utility.CallingGatewayHandler{Conf: configuration, H: Utility.APIHandlerGateway}

	mux := http.NewServeMux()
	http.HandleFunc("/favicon.ico", nil)
	mux.HandleFunc("/", conf.ServeHTTP)
	fmt.Println("Api Gateway Listening at http://localhost:8000")
	fmt.Println("To exit press : - Crtl+c")

	flag.Parse()
	Utility.NewLog(*logpath)
	Utility.Log.Println("Api Gateway Listening at http://localhost:8000")
	http.ListenAndServe(":8000", mux)

}
