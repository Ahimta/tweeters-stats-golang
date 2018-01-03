package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/Ahimta/tweeters-stats-golang/auth"
	"github.com/Ahimta/tweeters-stats-golang/config"
	"github.com/Ahimta/tweeters-stats-golang/handlers"
)

func hello(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello world!")
}

func main() {
	c := config.GetConfig()
	oauthClient, err := auth.NewOauth1Client(c)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", hello)
	mux.HandleFunc("/login/twitter", handlers.LoginHandlerFactory(c, oauthClient))
	mux.HandleFunc("/oauth/twitter/callback", handlers.OauthTwitterHandlerFactory(c, oauthClient))
	mux.HandleFunc("/tweeters-stats", handlers.GetTweetersStatsHandlerFactory(c, oauthClient))

	fmt.Printf("Server running on http://localhost:%s\n", c.Port)
	http.ListenAndServe(fmt.Sprintf(":%s", c.Port), mux)
}
