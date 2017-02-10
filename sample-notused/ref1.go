//DataRepository :- make sql connection
func DataRepository() error {
	query := url.Values{}
	u := &url.URL{
		Scheme: "sqlserver",
		User:   url.User("teslamotors\\vimkumar"),
		//Host:   fmt.Sprintf("%s:%d", hostname, port),
		Path:     "sjc04d1wrpdb01.sdlc.teslamotors.com", // if connecting to an instance instead of a port
		RawQuery: query.Encode(),
	}
	connectionString := u.String()
	db, err := sql.Open("mssql", "server=sjc04d1wrpdb01.sdlc.teslamotors.com;port=1433;user id=teslamotors\\vimkumar;password=NewLogin61!;database=WarpUser;log=63")
	println(connectionString)
	if err != nil {
		Log.Println(err.Error())
		return errors.New(err.Error())
	}
	defer db.Close()
	rows, err := db.Query("SELECT User  FROM WarpUser.Sec.Token")
	if err != nil {
		Log.Println(err.Error())
	}
	defer rows.Close()

	err = db.Ping()
	if err != nil {
		Log.Println(err.Error())
		return errors.New(err.Error())
		//panic(err.Error())
	}

	return nil
}

//Authenticate - user using token
func (val *AuthenticateModel) Authenticate() error {
	//Todo: Call Sql get the authneticate user

	var TokenID string
	var User string
	var ExpireAt string

	db, err := InitDB(val.Conf)
	if err != nil {
		return errors.New(err.Error()) // db failure
	}
	defer db.Close()
	rows, err := db.Query("SELECT cast([Id] as char(36)),[User],ExpiresAt FROM [sec].[Token] where id=?1", val.Token)

	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {

		// Get the results of the query
		err := rows.Scan(&TokenID, &User, &ExpireAt)
		if err != nil {
			return err
		}
	}
	val.User = User // set the user

	return nil
	//errors.New("Not Authorized")
}


func testSqlDB() {
	p := fmt.Println
	t := time.Now()

	p(t.Format(time.RFC3339))
	RFC3339local := "2006-01-02T15:04:05Z"
	fmt.Println(time.Parse(RFC3339local, "2018-01-16T18:55:00.05Z"))
	fmt.Println(t.Format("2018-01-16T18:55:00.05Z"))
	file, e := ioutil.ReadFile("./config/appsetting.json")
	if e != nil {
		fmt.Printf("File read error: %v", e)
		os.Exit(1)
	}
	var configuration appsetting.AppSetting
	json.Unmarshal(file, &configuration)
	authMode := Utility.AuthenticateModel{Conf: &configuration, Token: "6A253656-F2BE-4AC5-8BEA-2B4CBCE8F836", HTTPMethod: "GET", URL: "/a/b"}
	authMode.Authenticate()
}

//2017-02-05 15:47:14.5627001 -0800 PST
//2017-02-05T15:47:14-08:00
//http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

//func ExecuteApiHandler(conf appsetting.AppSetting) http.Handler {
//	return http.HandlerFunc(ApiHandlerGateway)
//}

//Handle the incomming request
//func (Conf *CallingGatewayHandler) ApiHandlerGateway(w http.ResponseWriter, r *http.Request) (int, error) {
//http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
//	io.WriteString(w, "Hello World")
//	return http.StatusOK, nil
//}
