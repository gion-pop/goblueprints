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
)

// templ は 1 つのテンプレートを表す
type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

// ServeHTTP は HTTP リクエストを処理する
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(
			template.ParseFiles(filepath.Join("templates", t.filename)),
		)
	})
	t.templ.Execute(w, r)
}

func main() {
	addr := flag.String("addr", ":8080", "アプリケーションのアドレス")
	flag.Parse() // 引数を解釈

	// gomniauth のセットアップ
	gomniauth.SetSecurityKey(os.Getenv("_CHAT_SECURITY_KEY"))
	gomniauth.WithProviders(
		google.New(
			os.Getenv("_CHAT_GOOGLE_ID"),
			os.Getenv("_CHAT_GOOGLE_SECRET"),
			"http://localhost:8080/auth/callback/google",
		),
		facebook.New(
			os.Getenv("_CHAT_FACEBOOK_ID"),
			os.Getenv("_CHAT_FACEBOOK_SECRET"),
			"http://localhost:8080/auth/callback/facebook",
		),
		github.New(
			os.Getenv("_CHAT_GITHUB_CK"),
			os.Getenv("_CHAT_GITHUB_CS"),
			"http://localhost:8080/auth/callback/github",
		),
	)

	r := NewRoom()
	http.Handle("/", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandler)
	http.Handle("/room", r)
	// チャットルームを開始
	go r.run()
	// Web サーバーを起動
	log.Println("Web サーバーを起動します．ポート：", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
