package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/mixigroup/mixi2-application-sdk-go/auth"
	application_apiv1 "github.com/mixigroup/mixi2-application-sdk-go/gen/go/social/mixi/application/service/application_api/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type Quote struct {
	Text   string `json:"text"`
	Author string `json:"author"`
}

type State struct {
	PostedIndexes []int `json:"posted_indexes"`
}

const (
	stateFile  = "state.json"
	quotesFile = "quotes.json"
)

func loadQuotes() []Quote {
	data, err := os.ReadFile(quotesFile)
	if err != nil {
		log.Fatal("quotes.json読み込み失敗:", err)
	}

	var quotes []Quote
	if err := json.Unmarshal(data, &quotes); err != nil {
		log.Fatal("quotes.json解析失敗:", err)
	}

	return quotes
}

func loadState() State {
	data, err := os.ReadFile(stateFile)
	if err != nil {
		return State{PostedIndexes: []int{}}
	}

	var s State
	if err := json.Unmarshal(data, &s); err != nil {
		return State{PostedIndexes: []int{}}
	}

	return s
}

func saveState(s State) {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		log.Fatal("state.json変換失敗:", err)
	}

	if err := os.WriteFile(stateFile, data, 0644); err != nil {
		log.Fatal("state.json保存失敗:", err)
	}
}

func contains(list []int, target int) bool {
	for _, v := range list {
		if v == target {
			return true
		}
	}
	return false
}

func pickRandomQuoteIndex(quotes []Quote, state State) int {
	available := []int{}

	for i := range quotes {
		if !contains(state.PostedIndexes, i) {
			available = append(available, i)
		}
	}

	if len(available) == 0 {
		available = make([]int, len(quotes))
		for i := range quotes {
			available[i] = i
		}
	}

	rand.Seed(time.Now().UnixNano())

	return available[rand.Intn(len(available))]
}

func main() {
	quotes := loadQuotes()

	if len(quotes) == 0 {
		log.Fatal("名言が登録されていません")
	}

	state := loadState()

	if len(state.PostedIndexes) >= len(quotes) {
		state.PostedIndexes = []int{}
	}

	index := pickRandomQuoteIndex(quotes, state)
	q := quotes[index]

	authenticator, err := auth.NewAuthenticator(
		os.Getenv("CLIENT_ID"),
		os.Getenv("CLIENT_SECRET"),
		os.Getenv("TOKEN_URL"),
	)
	if err != nil {
		log.Fatal("認証設定作成失敗:", err)
	}

	authCtx, err := authenticator.AuthorizedContext(context.Background())
	if err != nil {
		log.Fatal("認証失敗:", err)
	}

	conn, err := grpc.Dial(
		os.Getenv("API_ADDRESS"),
		grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})),
	)
	if err != nil {
		log.Fatal("mixi2 API接続失敗:", err)
	}
	defer conn.Close()

	client := application_apiv1.NewApplicationServiceClient(conn)

	text := q.Text + "\n\n" + q.Author

	_, err = client.CreatePost(authCtx, &application_apiv1.CreatePostRequest{
		Text: text,
	})
	if err != nil {
		log.Fatal("投稿失敗:", err)
	}

	state.PostedIndexes = append(state.PostedIndexes, index)
	saveState(state)

	log.Println("投稿成功:", q.Text)
}
