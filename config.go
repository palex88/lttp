package main

//var (
//	Database   *sql.DB
//	ConfigFile = "config.json"
	//OAuthConfig  *oauth2.Config
	//SessionStore sessions.Store
//)

//func init() {
	//OAuthConfig = configureOAuthClient("clientid", "clientsecret")

	//cookieStore := sessions.NewCookieStore([]byte("something-very-secret"))
	//cookieStore.Options = &sessions.Options{
	//	HttpOnly: true,
	//}
	//SessionStore = cookieStore
//}



//func configureOAuthClient(clientID, clientSecret string) *oauth2.Config {
//	redirectURL := os.Getenv("OAUTH2_CALLBACK")
//	if redirectURL == "" {
//		redirectURL = "http://localhost:8080/oauth2callback"
//	}
//	return &oauth2.Config{
//		ClientID:     clientID,
//		ClientSecret: clientSecret,
//		RedirectURL:  redirectURL,
//		Scopes:       []string{"email", "profile"},
//		Endpoint:     google.Endpoint,
//	}
//}
