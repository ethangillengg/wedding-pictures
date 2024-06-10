package services

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

const (
	key            = "ZoAgQHNNxcrgBxCewXuRQZPjaqgtfSyX"
	MaxAge         = 86400 * 30 // 30 days
	IsProd         = false
	UserCookieKey  = "user-session"
	UserSessionKey = "user"
	HttpOnly       = true
)

type AuthService struct{}

func NewAuthService() *AuthService {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	googleClientId := os.Getenv("GOOGLE_CLIENT_ID")
	googleClientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")

	store := sessions.NewCookieStore([]byte(key))
	store.MaxAge(MaxAge)

	store.Options.Path = "/"
	store.Options.HttpOnly = HttpOnly
	store.Options.Secure = IsProd

	gothic.Store = store

	goth.UseProviders(google.New(googleClientId, googleClientSecret, "http://localhost:7331/auth/google"))

	return &AuthService{}
}

func (as *AuthService) GetUserSession(r *http.Request) (goth.User, error) {
	session, err := gothic.Store.Get(r, UserCookieKey)
	if err != nil {
		return goth.User{}, err
	}

	u := session.Values[UserSessionKey]
	if u == nil {
		return goth.User{}, fmt.Errorf("user is not authenticated! %v", u)
	}

	return u.(goth.User), nil
}

func (as *AuthService) StoreUserSession(w http.ResponseWriter, r *http.Request, user goth.User) error {
	session, _ := gothic.Store.Get(r, UserCookieKey)

	session.Values[UserSessionKey] = user

	err := session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	return nil
}

func (as *AuthService) ClearUserSession(w http.ResponseWriter) {
	c := &http.Cookie{
		Name:     UserCookieKey,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: HttpOnly,
	}

	http.SetCookie(w, c)
}
