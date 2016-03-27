package models

//Song 曲モデル
type Song struct {
	ID               uint
	PageNumber       uint
	ArtistName       string
	ArtistNameR      string
	ArtistNameK      string
	LyricWriterName  string
	LyricWriterNameK string
	LyricWriterNameR string
	SongWriterName   string
	SongWriterNameK  string
	SongWriterNameR  string
	Name             string
	NameK            string
	NameR            string
	Number           string
	OriginalKey      string
	DeliveryStatus   string
	ModelNames       string
	Lyric            string `gorm:"type:text"`
	DeliveryTerm     string
}

//CanCreate 登録できるかチェック
func (m *Song) CanCreate() bool {
	if m.PageNumber <= 0 {
		return false
	}

	if m.ArtistName == "" {
		return false
	}

	if m.ArtistNameK == "" {
		return false
	}

	if m.ArtistNameR == "" {
		return false
	}

	if m.Name == "" {
		return false
	}

	if m.NameK == "" {
		return false
	}

	if m.NameR == "" {
		return false
	}

	if m.Number == "" {
		return false
	}

	return true
}
