package main

import (
	"context"
	"crypto/tls"
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

var quotes = []Quote{
	{"僕もね、一応ハーフなんです。ええ。おやじが滋賀県で、おふくろが岩手県。", "2015-02-26 関西テレビ よ～いドン！"},
	{"僕、森川さん好きですよ、大好き。", "2011"},
	{"（落葉が雪にを）六畳一間の高田馬場のアパートのベットに座ってギターを弾きながら作ったとは思えないでしょ。ええそうなんです。酒飲みな、酒飲みながら。毎日一曲っていうのを作ってた中の一曲なんです。", "2025-03-20 ラジオ日本 きのうの続きのつづき"},
	{"いやぁ、しかしね（君は薔薇より美しいを）歌って気持ちいって言われるとホントにもぅ悔しいよねぇ。歌ってこんなに辛いもんだっての分からせたいよ君らに。へへへへｗ", "2015-03-22 東海ラジオ 源石和輝音楽博覧会"},
}

func main() {
	rand.Seed(time.Now().UnixNano())

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

	client := application_apiv1.NewApplicationAPIServiceClient(conn)

	q := quotes[rand.Intn(len(quotes))]

	text := "今日の名言\n\n" +
		"「" + q.Text + "」\n" +
		"— " + q.Author

	_, err = client.CreatePost(authCtx, &application_apiv1.CreatePostRequest{
		Content: text,
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Println("投稿成功:", q.Text)
}