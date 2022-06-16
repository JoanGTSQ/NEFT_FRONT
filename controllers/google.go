package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	context2 "jgt.solutions/context"
	"jgt.solutions/errorController"
	"jgt.solutions/models"
	"jgt.solutions/views"
)

var (
	oauthConfGl = &oauth2.Config{
		ClientID:     "",
		ClientSecret: "",
		RedirectURL:  "",
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.profile", "https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/user.birthday.read"},
		Endpoint:     google.Endpoint,
	}
	oauthConfGP = &oauth2.Config{
		ClientID:     "",
		ClientSecret: "",
		RedirectURL:  "",
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.profile", "https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/user.birthday.read"},
		Endpoint:     google.Endpoint,
	}
	oauthStateStringGl = ""
)

func InitializeOAuthGoogle() {
	oauthConfGl.ClientID = viper.GetString("google.clientID")
	oauthConfGP.ClientID = viper.GetString("google.clientID")
	oauthConfGl.ClientSecret = viper.GetString("google.clientSecret")
	oauthConfGP.ClientSecret = viper.GetString("google.clientSecret")
	oauthConfGl.RedirectURL = "https://" + viper.GetString("urlRedirect") + "/registergoogle"
	oauthConfGP.RedirectURL = "https://" + viper.GetString("urlRedirect") + "/logingl"
	oauthStateStringGl = viper.GetString("oauthStateString")
}
func HandleLogin(w http.ResponseWriter, r *http.Request, oauthConf *oauth2.Config, oauthStateString string) {
	URL, err := url.Parse(oauthConf.Endpoint.AuthURL)
	if err != nil {
		fmt.Println("Parse: " + err.Error())
	}
	//Log.Info(URL.String())
	parameters := url.Values{}
	parameters.Add("client_id", oauthConf.ClientID)
	parameters.Add("scope", strings.Join(oauthConf.Scopes, " "))
	parameters.Add("redirect_uri", oauthConf.RedirectURL)
	parameters.Add("response_type", "code")
	parameters.Add("state", oauthStateString)
	URL.RawQuery = parameters.Encode()
	urlS := URL.String()
	//Log.Info(url)
	http.Redirect(w, r, urlS, http.StatusTemporaryRedirect)
}

/*
HandleGoogleLogin Function
*/
func HandleGoogleRegister(w http.ResponseWriter, r *http.Request) {
	HandleLogin(w, r, oauthConfGl, oauthStateStringGl)
}
func HandleGoogleLogin(w http.ResponseWriter, r *http.Request) {
	HandleLogin(w, r, oauthConfGP, oauthStateStringGl)
}

/*
CallBackFromGoogle Function
*/
func (u *Users) CallBackFromGoogle(w http.ResponseWriter, r *http.Request) {

	state := r.FormValue("state")
	if state != oauthStateStringGl {
		fmt.Println("invalid oauth state, expected " + oauthStateStringGl + ", got " + state + "\n")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	code := r.FormValue("code")

	if code == "" {
		fmt.Println("code not found")
		w.Write([]byte("Code Not Found to provide AccessToken..\n"))
		reason := r.FormValue("error_reason")
		if reason == "user_denied" {
			w.Write([]byte("User has denied Permission.."))
		}
		// User has denied access..
		// http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	} else {
		token, err := oauthConfGl.Exchange(context.Background(), code)
		if err != nil {
			return
		}
		resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + url.QueryEscape(token.AccessToken))
		if err != nil {
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}
		defer resp.Body.Close()

		response, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}

		type userStruct struct {
			Id      string
			Name    string
			Email   string
			Picture string
		}
		var userGoogle userStruct
		if err := json.Unmarshal(response, &userGoogle); err != nil { // Parse []byte to go struct pointer
			fmt.Println("Can not unmarshal JSON")
		}
		type DOB struct {
			Year  int
			Month int
			Day   int
		}
		type SecondPas struct {
			Metadata interface{}
			Date     DOB
		}
		type ReceiverDOB struct {
			Birthdays []SecondPas `json:"birthdays"`
		}

		var receiver ReceiverDOB
		resp2, err := http.Get(fmt.Sprintf("https://people.googleapis.com/v1/people/%s?personFields=birthdays&sources=READ_SOURCE_TYPE_PROFILE&access_token=%s", userGoogle.Id, url.QueryEscape(token.AccessToken)))
		if err != nil {
			fmt.Println("Get: " + err.Error() + "\n")
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}
		defer resp2.Body.Close()
		response2, _ := ioutil.ReadAll(resp2.Body)
		if err := json.Unmarshal(response2, &receiver); err != nil { // Parse []byte to go struct pointer
			fmt.Println("Can not unmarshal JSON")
		}

		dateElem := []int{
			receiver.Birthdays[0].Date.Year,
			receiver.Birthdays[0].Date.Day,
			receiver.Birthdays[0].Date.Month,
		}

		date, err := time.Parse("2006-2-1", strings.Trim(strings.Join(strings.Fields(fmt.Sprint(dateElem)), "-"), "[]"))
		if err != nil {
			fmt.Println(err)
		}
		user := models.User{
			Name:     userGoogle.Name,
			Email:    userGoogle.Email,
			Password: userGoogle.Id + userGoogle.Name + userGoogle.Email,
			Photo:    userGoogle.Picture,
			DOB:      date,
		}
		var vd views.Data
		err = u.us.Create(&user)
		if err != nil {
			vd.Alert = &views.Alert{
				Level:   views.AlertLvlError,
				Message: err.Error(),
			}
			vd.Yield = &userGoogle
			u.NewView.Render(w, r, &vd)
			return
		}
		ctx := r.Context()
		ctx = context2.WithUser(ctx, &user)
		r = r.WithContext(ctx)
		err = u.signIn(w, &user)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
}

func (u *Users) CallBackLoginFromGoogle(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	state := r.FormValue("state")
	if state != oauthStateStringGl {
		errorController.WD.Content = "invalid oauth state, expected " + oauthStateStringGl + ", got " + state
		errorController.WD.Site = "login google"
		errorController.WD.SendErrorWHWeb()
		vd.Alert = &views.Alert{
			Level:   views.AlertLvlError,
			Message: views.AlertMsgGeneric,
		}
		u.LoginView.Render(w, r, &vd)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	code := r.FormValue("code")

	if code == "" {

		fmt.Println("code not found")
		w.Write([]byte("Code Not Found to provide AccessToken..\n"))
		reason := r.FormValue("error_reason")
		if reason == "user_denied" {
			vd.Alert = &views.Alert{
				Level:   views.AlertLvlError,
				Message: "User don't let enough permission",
			}
			u.LoginView.Render(w, r, &vd)
		}
		// User has denied access..
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	} else {
		token, err := oauthConfGP.Exchange(context.Background(), code)
		if err != nil {
			return
		}
		resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + url.QueryEscape(token.AccessToken))
		if err != nil {
			fmt.Println("Get: " + err.Error() + "\n")
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}
		defer resp.Body.Close()

		response, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("ReadAll: " + err.Error() + "\n")
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}

		if err != nil {
			fmt.Println("ReadAll: " + err.Error() + "\n")
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}
		type userStruct struct {
			Id      string
			Name    string
			Email   string
			Picture string
		}

		var userGoogle userStruct

		if err := json.Unmarshal(response, &userGoogle); err != nil { // Parse []byte to go struct pointer
			fmt.Println("Can not unmarshal JSON")
		}
		user, err := u.us.Authenticate(userGoogle.Email, userGoogle.Id+userGoogle.Name+userGoogle.Email)
		if err != nil {
			vd.Alert = &views.Alert{
				Level:   views.AlertLvlError,
				Message: err.Error(),
			}
			vd.Yield = &userGoogle
			u.LoginView.Render(w, r, &vd)
			return
		}
		if user.Photo != userGoogle.Picture {
			user.Photo = userGoogle.Picture
			u.us.Update(user)
		}
		ctx := r.Context()
		ctx = context2.WithUser(ctx, user)
		r = r.WithContext(ctx)
		err = u.signIn(w, user)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
}
