package models

//FailedPage 正常に取得できなかったページ
type FailedPage struct {
	ID         uint
	PageNumber uint
	Text       string
}
