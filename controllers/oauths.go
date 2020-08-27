package controllers

import (
	"context"
	"fmt"
	llctx "gallery/context"
	dbx "gallery/dropbox"
	"gallery/models"
	"net/http"

	"golang.org/x/oauth2"
)

//NewOAuths Creates new Oauths controller
func NewOAuths(os models.OAuthService, dbxConfid *oauth2.Config) *OAuths {
	return &OAuths{
		os:       os,
		dbxOauth: dbxConfid,
	}
}

//OAuths is custom sturct of controller
type OAuths struct {
	os       models.OAuthService
	dbxOauth *oauth2.Config
}

func (o *OAuths) DropboxConnect(w http.ResponseWriter, r *http.Request) {
	state := "random-state"
	url := o.dbxOauth.AuthCodeURL(state)
	http.Redirect(w, r, url, http.StatusFound)
}

func (o *OAuths) DropboxCallback(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	state := r.FormValue("state")
	if state != "random-state" {
		http.Error(w, "Invalid state", http.StatusBadRequest)
		return
	}
	code := r.FormValue("code")
	token, err := o.dbxOauth.Exchange(context.TODO(), code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	user := llctx.User(r.Context())
	existing, err := o.os.Find(user.ID, models.OAuthDropbox)
	if err == models.ErrNotFound {

	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		o.os.Delete(existing.ID)
	}
	userOAuth := models.OAuth{
		UserID:  user.ID,
		Token:   *token,
		Service: models.OAuthDropbox,
	}
	err = o.os.Create(&userOAuth)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/galleries", http.StatusMovedPermanently)
}

func (o *OAuths) DropboxTest(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	path := r.FormValue("path")

	user := llctx.User(r.Context())
	userOAuth, err := o.os.Find(user.ID, models.OAuthDropbox)
	if err != nil {
		panic(err)
	}
	token := userOAuth.Token

	folders, files, err := dbx.List(token.AccessToken, path)
	if err != nil {
		panic(err)
	}
	fmt.Fprintln(w, "Folders=", folders)
	fmt.Fprintln(w, "Files=", files)
}
