package main

import (
	"flag"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/yneee/exsongs/models"
	"github.com/yneee/exsongs/settings"
	"github.com/yneee/exsongs/utils"
)

var (
	baseURL = "https://www.joysound.com/web/search/song/"
	db      *gorm.DB
	mode    = flag.String("mode", "crawl", "mode. select 'crawl' or 'crawl-failed-page'.")
)

func init() {
	_db, err := gorm.Open("mysql", settings.DbInfo)
	if err != nil {
		log.Fatalln(err)
	}

	db = _db
}

func main() {
	flag.Parse()

	migrate()

	switch *mode {
	case "crawl":
		crawl()
	case "crawl-failed-page":
		crawlFailedPage()
	}

}

func migrate() {
	db.Set("gorm:table_options", "ENGINE=InnoDB")

	db.AutoMigrate(&models.Log{})
	db.AutoMigrate(&models.Song{})
	db.AutoMigrate(&models.FailedPage{})
}

func crawl() {
	lastLog := models.Log{}
	db.Last(&lastLog)

	//pagenumberセット
	pageNumber := uint(1)
	if 0 < lastLog.PageNumber {
		pageNumber = lastLog.PageNumber
	}

	for {
		//html取得
		doc, err := getDoc(pageNumber)
		if err != nil {
			log.Fatalln("get doc error. err: ", err)
		}

		//曲取得
		songs := getSongs(doc)

		//曲登録
		enable := true
		for i := 0; i < len(songs); i++ {
			existSong := models.Song{}
			db.Where(&models.Song{Number: songs[i].Number}).First(&existSong)

			//すでに登録されている
			if 0 < existSong.ID {
				log.Println("already crate. song: ", existSong)
				continue
			}

			//ページ数セット
			songs[i].PageNumber = pageNumber

			//1つでも登録正常でない場合はログに残す
			if !songs[i].CanCreate() {
				enable = false
				continue
			}

			//曲登録
			db.Create(&songs[i])
		}

		logText := "create"

		//そもそもページが存在するか
		notFundText := trim(doc.Find("#jp-cmp-main > .jp-cmp-box-005 > .jp-cmp-h1-error > span").Text())
		if notFundText == "このページは存在しません。" {
			//存在しない
			logText = "not_found_page"
		} else {
			//存在したけど曲なかった
			if len(songs) <= 0 {
				logText = "none_songs"
			}

			//存在したけど曲のどれかがデータ歯抜け
			if !enable {
				logText = "get_songs_failed"
			}
		}

		//失敗したページ登録
		if logText == "none_songs" || logText == "get_songs_failed" {
			db.Create(&models.FailedPage{
				PageNumber: pageNumber,
				Text:       logText,
			})
		}

		//ログ残す
		db.Create(&models.Log{
			PageNumber: pageNumber,
			Text:       logText,
		})

		log.Println(logText, "pageNumber:", pageNumber)

		//とりあえずこの番号まで
		if 600000 <= pageNumber {
			break
		}

		time.Sleep(700 * time.Millisecond)

		pageNumber++
	}
}

func crawlFailedPage() {

	startID := uint(1)
	endID := uint(10000)
	rangeID := uint(10000)

	for {
		failedPages := []models.FailedPage{}
		db.Where("? <= id and id < ?", startID, endID).Find(&failedPages)

		//なくなったら終わり
		if len(failedPages) <= 0 {
			break
		}

		//html取得
		for _, failedPage := range failedPages {
			doc, err := getDoc(failedPage.PageNumber)
			if err != nil {
				log.Fatalln("get doc error. err: ", err)
			}

			//曲取得
			songs := getSongs(doc)

			//曲登録
			enable := true
			for i := 0; i < len(songs); i++ {
				existSong := models.Song{}
				db.Where(&models.Song{Number: songs[i].Number}).First(&existSong)

				//すでに登録されている
				if 0 < existSong.ID {
					log.Println("already crate. song: ", existSong)
					continue
				}

				//ページ数セット
				songs[i].PageNumber = failedPage.PageNumber

				//1つでも登録正常でない場合はログに残す
				if !songs[i].CanCreate() {
					enable = false
					continue
				}

				//曲登録
				db.Create(&songs[i])
			}

			//曲ない
			if len(songs) <= 0 {
				log.Println("none songs. pageNumber:", failedPage.PageNumber)
			} else {
				if !enable {
					//また失敗
					log.Println("failed. pageNumber:", failedPage.PageNumber)
				} else {
					//okな場合はfailedPageから削除してあげる
					log.Println("create! pageNumber:", failedPage.PageNumber)
					db.Delete(&failedPage)
				}
			}

			time.Sleep(700 * time.Millisecond)
		}

		log.Println("end id:", endID)

		//id更新
		startID = endID
		endID = endID + rangeID
	}
}

func getSongs(doc *goquery.Document) []models.Song {
	songs := []models.Song{}

	//アーティスト名等取得（ない場合は、存在しないか、javascriptレンダリングがうまくできていないか）
	artistName, lyricWriterName, songWriterName := getSongRelationName(doc)
	lyric := getLyric(doc)

	doc.Find(".jp-cmp-karaoke-list-001 > ul > li").Each(func(i int, s *goquery.Selection) {
		song := models.Song{}

		song.ArtistName = artistName
		song.LyricWriterName = lyricWriterName
		song.SongWriterName = songWriterName
		song.Lyric = lyric

		//曲名
		song.Name = trim(s.Find(".jp-cmp-karaoke-details > h4").Text())

		//曲番号等
		s.Find(".jp-cmp-karaoke-details > .jp-cmp-movie-status-001 > dl").Children().Each(func(i int, _s *goquery.Selection) {
			switch trim(_s.Text()) {
			case "曲番号:":
				song.Number = trim(_s.Next().Text())

			case "キー":
				song.OriginalKey = trim(_s.Next().Text())

			case "配信予定:":
				song.DeliveryStatus = trim(_s.Next().Text())

			case "配信期間:":
				song.DeliveryTerm = trim(_s.Next().Text())
			}
		})

		//モデル
		modelNames := []string{}
		s.Find(".jp-cmp-karaoke-platform > ul > li").Each(func(i int, _s *goquery.Selection) {
			if modelName, exist := _s.Find("img").Attr("alt"); exist {
				modelNames = append(modelNames, modelName)
			}
		})

		song.ModelNames = strings.Join(modelNames, ", ")

		songs = append(songs, song)
	})

	return songs
}

func getLyric(doc *goquery.Document) string {
	str := doc.Find("#lyrics > .jp-cmp-song-words-contents > .jp-cmp-song-words-details p").Text()

	reg := regexp.MustCompile(`\n`)
	str = reg.ReplaceAllString(str, " ")

	utils.NormalizeString(str)

	return trim(str)
}

func getSongRelationName(doc *goquery.Document) (string, string, string) {
	var artistName, lyricWriterName, songWriterName string

	selection := doc.Find(".jp-cmp-song-block-001 .jp-cmp-song-visual .jp-cmp-song-table-001 tr")

	selection.Each(func(i int, s *goquery.Selection) {
		if trim(s.Find("th").Text()) == "歌手名" {
			artistName = trim(s.Find("td a").Text())
		}

		if trim(s.Find("th").Text()) == "作詞" {
			lyricWriterName = trim(s.Find("td span").Text())
		}

		if trim(s.Find("th").Text()) == "作曲" {
			songWriterName = trim(s.Find("td span").Text())
		}
	})

	return artistName, lyricWriterName, songWriterName
}

func getDoc(pageNum uint) (*goquery.Document, error) {
	url := []string{
		settings.SplashInfo,
		"?url=",
		"https://www.joysound.com/web/search/song/",
		strconv.Itoa(int(pageNum)),
		"&images=0",
	}

	doc, err := goquery.NewDocument(strings.Join(url, ""))
	if err != nil {
		return &goquery.Document{}, err
	}

	return doc, nil
}

func trim(str string) string {
	reg := regexp.MustCompile(`\n`)
	str = reg.ReplaceAllString(str, "")

	reg = regexp.MustCompile(`^(\s+?)(\S.*?)(\s+?)$`)
	str = reg.ReplaceAllString(str, "$2")

	return utils.NormalizeString(str)
}
