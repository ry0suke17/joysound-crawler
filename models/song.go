package models

//Song 曲モデル
type Song struct {
	ID              uint
	PageNumber      uint
	ArtistName      string
	LyricWriterName string
	SongWriterName  string
	Name            string
	Number          string
	OriginalKey     string
	DeliveryStatus  string
	ModelNames      string
	Lyric           string `gorm:"type:text"`
	DeliveryTerm    string
}

//CanCreate 登録できるかチェック
func (m *Song) CanCreate() bool {
	if m.PageNumber <= 0 {
		return false
	}

	if m.ArtistName == "" {
		return false
	}

	if m.Name == "" {
		return false
	}

	if m.Number == "" {
		return false
	}

	return true
}
