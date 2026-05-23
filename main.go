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
	Text   string
	Author string
}

type State struct {
	PostedIndexes []int `json:"posted_indexes"`
}

const stateFile = "state.json"

var quotes = []Quote{
	{"（落葉が雪にを）六畳一間の高田馬場のアパートのベットに座ってギターを弾きながら作ったとは思えないでしょ。ええそうなんです。酒飲みな、酒飲みながら。毎日一曲っていうのを作ってた中の一曲なんです。", "2025-03-20 ラジオ日本 きのうの続きのつづき"},
	{"いやぁ、しかしね（君は薔薇より美しいを）歌って気持ちいって言われるとホントにもぅ悔しいよねぇ。歌ってこんなに辛いもんだっての分からせたいよ君らに。へへへへｗ", "2015-03-22 東海ラジオ 源石和輝音楽博覧会"},
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
		log.Fatal(err)
	}

	if err := os.WriteFile(stateFile, data, 0644); err != nil {
		log.Fatal(err)
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

func pickRandomQuoteIndex(state State) int {
	if len(state.PostedIndexes) >= len(quotes) {
		state.PostedIndexes = []int{}
	}

	available := []int{}

	for i := range quotes {
		if !contains(state.PostedIndexes, i) {
			available = append(available, i)
		}
	}

	rand.Seed(time.Now().UnixNano())

	return available[rand.Intn(len(available))]
}

func main() {
	if len(quotes) == 0 {
		log.Fatal("名言が登録されていません")
	}

	state := loadState()

	if len(state.PostedIndexes) >= len(quotes) {
		state.PostedIndexes = []int{}
	}

	index := pickRandomQuoteIndex(state)
	q := quotes[index]

	authenticator, err := auth.NewAuthenticator(
		os.Getenv("CLIENT_ID"),
		os.Getenv("CLIENT_SECRET"),
		os.Getenv("TOKEN_URL"),
	)
	if err != nil {
		log.Fatal(err)
	}

	authCtx, err := authenticator.AuthorizedContext(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	conn, err := grpc.Dial(
		os.Getenv("API_ADDRESS"),
		grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := application_apiv1.NewApplicationServiceClient(conn)

	text := q.Text + "\n\n" + q.Author

	_, err = client.CreatePost(authCtx, &application_apiv1.CreatePostRequest{
		Text: text,
	})
	if err != nil {
		log.Fatal(err)
	}

	state.PostedIndexes = append(state.PostedIndexes, index)
	saveState(state)

	log.Println("投稿成功:", q.Text)
}
