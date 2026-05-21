package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"log"
	"os"

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
	Index int `json:"index"`
}

const stateFile = "state.json"

var quotes = []Quote{
	{"テスト", "1999-01-01 テスト"},
	{"僕もね、一応ハーフなんです。ええ。おやじが滋賀県で、おふくろが岩手県。", "2015-02-26 関西テレビ よ～いドン！"},
	{"僕、森川さん好きですよ、大好き。", "2007-10-23 笑っていいとも！"},
	{"（落葉が雪にを）六畳一間の高田馬場のアパートのベットに座ってギターを弾きながら作ったとは思えないでしょ。ええそうなんです。酒飲みな、酒飲みながら。毎日一曲っていうのを作ってた中の一曲なんです。", "2025-03-20 ラジオ日本 きのうの続きのつづき"},
	{"いやぁ、しかしね（君は薔薇より美しいを）歌って気持ちいって言われるとホントにもぅ悔しいよねぇ。歌ってこんなに辛いもんだっての分からせたいよ君らに。へへへへｗ", "2015-03-22 東海ラジオ 源石和輝音楽博覧会"},
}

func loadState() State {
	data, err := os.ReadFile(stateFile)
	if err != nil {
		return State{Index: 0}
	}

	var s State
	if err := json.Unmarshal(data, &s); err != nil {
		return State{Index: 0}
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

func main() {
	if len(quotes) == 0 {
		log.Fatal("名言が登録されていません")
	}

	state := loadState()

	if state.Index < 0 || state.Index >= len(quotes) {
		state.Index = 0
	}

	q := quotes[state.Index]

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

	text := "今日の一言“φ(･ω･｡)\n\n" +
		"「" + q.Text + "」\n" +
		"— " + q.Author

	_, err = client.CreatePost(authCtx, &application_apiv1.CreatePostRequest{
		Text: text,
	})
	if err != nil {
		log.Fatal(err)
	}

	state.Index++
	if state.Index >= len(quotes) {
		state.Index = 0
	}

	saveState(state)

	log.Println("投稿成功:", q.Text)
}
