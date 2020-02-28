package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/facebook"
	"github.com/stretchr/gomniauth/providers/github"
	"github.com/stretchr/gomniauth/providers/google"
	"github.com/tattsum/go-project/trace"
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
	t.templ.Execute(w, r)
}

func main() {
	var addr = flag.String("addr", ":8080", "アプリケーションのアドレス")
	flag.Parse() // フラグを解釈する．

	//Gomniauthのセットアップ
	gomniauth.SetSecurityKey("セキュリティキー")
	gomniauth.WithProviders(
		// provider.New("クライアントID", "秘密の値", "callbackURL")
		facebook.New("FACEBOOK_CLIENT_ID", "FACEBOOK_SERCRET_ID", "http://localhost:8080/auth/callback/facebook"),
		github.New("GITHUB_CLIENT_ID", "GITHUB_SECRET_ID", "http://localhost:8080/auth/callback/github"),
		google.New("GOOGLE_CLIENT_ID", "GOOGLE_SERCET_ID", "http://localhost:8080/auth/callback/google"),
	)

	r := newRoom()
	r.tracer = trace.New(os.Stdout)
	// ルート
	http.Handle("/", MustAuth(&templateHandler{filename: "/chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandler)
	http.Handle("/room", r)

	// チャットルームを開始します．
	go r.run()

	// Webサーバを開始
	log.Println("Webサーバを開始します．port: ", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe", err)
	}
}
