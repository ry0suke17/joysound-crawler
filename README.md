# joysound-crawler

JOYSOUNDの曲を収集するクローラー

##　起動

```
docker-compose up
docker exec exsongs_elasticsearch_1 /usr/share/elasticsearch/bin/plugin install analysis-kuromoji
go run exsongs.go --mode=crawl
```

## 曲情報保存までの流れ

1. [spalsh](https://github.com/scrapinghub/splash)へ曲ページリクエストしjavascriptをレンダリングしたhtmlを取得
2. elasticsearchの[kuromoji](https://www.atilika.com/ja/kuromoji/)で漢字の読み仮名を取得
3. 情報保存
