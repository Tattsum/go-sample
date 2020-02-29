package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/joho/godotenv"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/facebook"
	"github.com/stretchr/gomniauth/providers/github"
	"github.com/stretchr/gomniauth/providers/google"
	"github.com/stretchr/objx"
	"github.com/tattsum/go-sample/trace"
)

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

// ServeHTTPはTHHPリクエストを処理します．
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	data := map[string]interface{}{
		"Host": r.Host,
	}
	if authCookie, err := r.Cookie("auth"); err == nil {
		data["UserData"] = objx.MustFromBase64(authCookie.Value)
	}

	t.templ.Execute(w, r)
}

func main() {
	var addr = flag.String("addr", ":8080", "アプリケーションのアドレス")
	flag.Parse() // フラグを解釈する．

	// 環境変数のインポート
	err := godotenv.Load(fmt.Sprintf("envfiles/develop.env"))
	if err == nil {
		log.Fatal("Error Loading .env file")
	}
	facebook_id := os.Getenv("FACEBOOK_CLIENT_ID")
	facebook_secret := os.Getenv("FACEBOOK_SERCRET_ID")
	github_id := os.Getenv("GITHUB_CLIENT_ID")
	github_secret := os.Getenv("GITHUB_SECRET_ID")
	google_id := os.Getenv("GOOGLE_CLIENT_ID")
	google_secret := os.Getenv("GOOGLE_SECRET_ID")

	//Gomniauthのセットアップ
	gomniauth.SetSecurityKey("セキュリティキー")
	gomniauth.WithProviders(
		// provider.New("クライアントID", "秘密の値", "callbackURL")
		facebook.New(facebook_id, facebook_secret, "http://localhost:8080/auth/callback/facebook"),
		github.New(github_id, github_secret, "http://localhost:8080/auth/callback/github"),
		google.New(google_id, google_secret, "http://localhost:8080/auth/callback/google"),
	)

	r := newRoom(UseGravatar)
	r.tracer = trace.New(os.Stdout)
	// ルート
	http.Handle("/", MustAuth(&templateHandler{filename: "/chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandler)
	http.Handle("/room", r)
	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:   "auth",
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		})
		w.Header()["Location"] = []string{"/chat"}
		w.WriteHeader(http.StatusTemporaryRedirect)
	})

	// チャットルームを開始します．
	go r.run()

	// Webサーバを開始
	log.Println("Webサーバを開始します．port: ", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe", err)
	}
}
