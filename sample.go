package main

import (
	"net/http"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/memcache"
	"log"
	"html/template"
)

type Sample struct {
	// テスト用に整数を格納する
	No int
}
// ハンドラ
// 「/record」のリクエストをクライアントから受けたら、ここで一括で受けて、リクエストメソッドで処理を分岐する
func HandleSample(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		doGet(w, r)
	} else if r.Method == "POST" {
		doPost(w, r)
	}
}

// 「/record」でPostメソッドが飛んできた時に発動する
// 構造体SampleをDatastoreに保存する
func doPost(w http.ResponseWriter, r *http.Request) {
	// コンテキスト生成
	c := appengine.NewContext(r)

	i := 0
	// 構造体Sampleを100個生成するためにfor文をする
	// 構造体SampleのNoフィールドを各々0~99に設定する
	for {
		key := datastore.NewIncompleteKey(c, "Sample", nil)
		// 構造体Sampleをインスタス化
		sample := &Sample{No:i}
		i++
		_, err := datastore.Put(c, key, sample)

		if err != nil {
			// エラーハンドル
		}

		// 0~999まで処理を行うため、iが100になったら、ループを抜ける
		if i == 1000 {
			break
		}
	}

	// post.htmlを表示
	tmpl := template.Must(template.ParseFiles("./template/post.html"))
	tmpl.Execute(w, nil)
}

// 「/record」でGetメソッドが飛んできた時に発動する
// クエリカーソルでSampleカインドに登録してあるデータを取得
func doGet(w http.ResponseWriter, r *http.Request) {
	// SampleカインドでQueryを発行、Noで昇順に指定する
	q := datastore.NewQuery("Sample").Limit(20).Order("No")

	// コンテキスト生成
	c := appengine.NewContext(r)

	// memcacheから、key "sample"でitemをGet
	item, err := memcache.Get(c, "sample")
	// エラーがnilなら、カーソルを使用することができるので、カーソルでの処理を行う
	if err == nil {
		// DecodeCursorはbase-64からcursorにデコードする
		cursor, err := datastore.DecodeCursor(string(item.Value))
		if err == nil {
			// Start() : 開始点が指定されたQueryを返す
			q = q.Start(cursor)
		}
	}

	// Datastoreから取得したデータを格納するためのスライスを作成
	s := make([]Sample, 1100)
	// Queryを実行
	t := q.Run(c)
	i := 0
	for {
		// t.Nextで取得したEntityを格納するために変数を宣言
		var sample Sample
		// Datastoreから取得したEntityを引数で与えた&sampleに格納
		_, err := t.Next(&sample)
		// イテレータが最後まで進み、これ以上Datastoreにデータ(Entity)が存在しない場合は、t.Nextはdatastore.Doneを返す
		if err == datastore.Done {
			log.Println("==datastore.Done==")
			break
		}
		if err != nil {
			// エラーハンドル
			break
		}
		// スライスsのi番目に取得したEntityを格納
		s[i] = sample
		i++
	}

	// t.Cursor() : イテレータの現在の位置を表すcursorを返す
	if cursor, err := t.Cursor(); err == nil {
		// memcache.Itemを設定する
		item := &memcache.Item{
			Key:   "sample",
			// cursor.String()は、cursorのbase-64を返す
			// memcache.Valueにcursorを保存する
			Value: []byte(cursor.String()),
		}

		// itemをmemcacheに格納
		err := memcache.Set(c, item)

		if err != nil {
			// エラーハンドル
		}

	}

	s2 := make([]Sample, i)

	copy(s2, s)

	tmpl := template.Must(template.ParseFiles("./template/record.html"))
	tmpl.Execute(w, s2)

}


