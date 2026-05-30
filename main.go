package main

import (
	"context"
	"encoding/json"
	"log"
	"math/rand/v2"
	"os"
	"strings"

	"github.com/mixigroup/mixi2-application-sdk-go/auth"
	application_apiv1 "github.com/mixigroup/mixi2-application-sdk-go/gen/go/social/mixi/application/service/application_api/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	stateFile  = "state.json"
	quotesFile = "quotes.json"
)

type Quote struct {
	Text   string `json:"text"`
	Author string `json:"author"`
}

type State struct {
	PostedIndexes []int `json:"posted_indexes"`
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

	posted := buildPostedMap(state)

	index := pickRandomQuoteIndex(quotes, posted)
	q := quotes[index]

	quoteText := normalizeText(q.Text)
	text := trimPostText(quoteText + "\n\n" + q.Author)

	if os.Getenv("PREVIEW") == "1" {
		log.Println("プレビュー:")
		log.Println(text)
		return
	}

	authenticator, err := auth.NewAuthenticator(
		requireEnv("CLIENT_ID"),
		requireEnv("CLIENT_SECRET"),
		requireEnv("TOKEN_URL"),
	)
	if err != nil {
		log.Fatal("認証設定作成失敗:", err)
	}

	authCtx, err := authenticator.AuthorizedContext(context.Background())
	if err != nil {
		log.Fatal("認証失敗:", err)
	}

	conn, err := grpc.NewClient(
		requireEnv("API_ADDRESS"),
		grpc.WithTransportCredentials(
			credentials.NewClientTLSFromCert(nil, ""),
		),
	)
	if err != nil {
		log.Fatal("mixi2 API接続失敗:", err)
	}
	defer conn.Close()

	client := application_apiv1.NewApplicationServiceClient(conn)

	_, err = client.CreatePost(authCtx, &application_apiv1.CreatePostRequest{
		Text: text,
	})
	if err != nil {
		log.Fatal("投稿失敗:", err)
	}

	state.PostedIndexes = append(state.PostedIndexes, index)
	saveState(state)

	log.Println("投稿成功:", text)
}

func requireEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatal(key + " missing value")
	}

	return value
}

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

func buildPostedMap(state State) map[int]bool {
	posted := make(map[int]bool, len(state.PostedIndexes))

	for _, index := range state.PostedIndexes {
		posted[index] = true
	}

	return posted
}

func pickRandomQuoteIndex(quotes []Quote, posted map[int]bool) int {
	available := []int{}

	for i := range quotes {
		if !posted[i] {
			available = append(available, i)
		}
	}

	if len(available) == 0 {
		available = make([]int, len(quotes))
		for i := range quotes {
			available[i] = i
		}
	}

	return available[rand.IntN(len(available))]
}

func trimPostText(text string) string {
	// mixi2投稿本文の上限に収まるようにするための最大文字数
	const maxLen = 147

	runes := []rune(text)
	if len(runes) <= maxLen {
		return text
	}

	return string(runes[:maxLen-1]) + "…"
}

func normalizeText(text string) string {
	return strings.ReplaceAll(text, "\\n", "\n")
}
