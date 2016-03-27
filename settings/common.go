package settings

const (
	//DbInfo DbInfo
	DbInfo = "root:rooT555_@([mysql]:3306)/exsongs?charset=utf8&parseTime=true&loc=Asia%2FTokyo"

	//SplashInfo SplashInfo
	SplashInfo = "http://splash:8050/render.html"

	//ElasticsearchInfo ElasticsearchInfo
	ElasticsearchInfo = "http://elasticsearch:9200"

	//SongsInfo SongsInfo
	SongsInfo = "https://www.joysound.com/web/search/song/"

	//ElasticsearchSettings ElasticsearchSettings
	ElasticsearchSettings = `
{
  "settings": {
    "index": {
      "analysis": {
        "analyzer": {
          "romaji_analyzer": {
            "tokenizer": "kuromoji_tokenizer",
            "filter": [
              "romaji_readingform"
            ]
          },
          "katakana_analyzer": {
            "tokenizer": "kuromoji_tokenizer",
            "filter": [
              "katakana_readingform"
            ]
          }
        },
        "filter": {
          "romaji_readingform": {
            "type": "kuromoji_readingform",
            "use_romaji": true
          },
          "katakana_readingform": {
            "type": "kuromoji_readingform",
            "use_romaji": false
          }
        }
      }
    }
  }
}
`
)
