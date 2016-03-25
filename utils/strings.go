package utils

import (
	"encoding/hex"
	"io/ioutil"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

var (
	spaceReplacer = strings.NewReplacer(
		"　", " ",
		// "  ", " ",
	)

	zenHanNumberReplacer = strings.NewReplacer(
		"０", "0", "１", "1", "２", "2",
		"３", "3", "４", "4", "５", "5",
		"６", "6", "７", "7", "８", "8",
		"９", "9",
	)

	zenHanAlphabetReplacer = strings.NewReplacer(
		"ａ", "a", "ｂ", "b", "ｃ", "c", "ｄ", "d",
		"ｅ", "e", "ｆ", "f", "ｇ", "g", "ｈ", "h",
		"ｉ", "i", "ｊ", "j", "ｋ", "k", "ｌ", "l",
		"ｍ", "m", "ｎ", "n", "ｏ", "o", "ｐ", "p",
		"ｑ", "q", "ｒ", "r", "ｓ", "s", "ｔ", "t",
		"ｕ", "u", "ｖ", "v", "ｗ", "w", "ｘ", "x",
		"ｙ", "y", "ｚ", "z",

		"Ａ", "A", "Ｂ", "B", "Ｃ", "C", "Ｄ", "D",
		"Ｅ", "E", "Ｆ", "F", "Ｇ", "G", "Ｈ", "H",
		"Ｉ", "I", "Ｊ", "J", "Ｋ", "K", "Ｌ", "L",
		"Ｍ", "M", "Ｎ", "N", "Ｏ", "O", "Ｐ", "P",
		"Ｑ", "Q", "Ｒ", "R", "Ｓ", "S", "Ｔ", "T",
		"Ｕ", "U", "Ｖ", "V", "Ｗ", "W", "Ｘ", "X",
		"Ｙ", "Y", "Ｚ", "Z",
	)

	hanZenKanaReplacer = strings.NewReplacer(
		"ｳﾞ", "ヴ",
		"ｶﾞ", "ガ", "ｷﾞ", "ギ", "ｸﾞ", "グ", "ｹﾞ", "ゲ", "ｺﾞ", "ゴ",
		"ｻﾞ", "ザ", "ｼﾞ", "ジ", "ｽﾞ", "ズ", "ｾﾞ", "ゼ", "ｿﾞ", "ゾ",
		"ﾀﾞ", "ダ", "ﾁﾞ", "ヂ", "ﾂﾞ", "ヅ", "ﾃﾞ", "デ", "ﾄﾞ", "ド",
		"ﾊﾞ", "バ", "ﾋﾞ", "ビ", "ﾌﾞ", "ブ", "ﾍﾞ", "ベ", "ﾎﾞ", "ボ",
		"ﾊﾟ", "パ", "ﾋﾟ", "ピ", "ﾌﾟ", "プ", "ﾍﾟ", "ペ", "ﾎﾟ", "ポ",

		"ｱ", "ア", "ｲ", "イ", "ｳ", "ウ", "ｴ", "エ", "ｵ", "オ",
		"ｶ", "カ", "ｷ", "キ", "ｸ", "ク", "ｹ", "ケ", "ｺ", "コ",
		"ｻ", "サ", "ｼ", "シ", "ｽ", "ス", "ｾ", "セ", "ｿ", "ソ",
		"ﾀ", "タ", "ﾁ", "チ", "ﾂ", "ツ", "ﾃ", "テ", "ﾄ", "ト",
		"ﾅ", "ナ", "ﾆ", "ニ", "ﾇ", "ヌ", "ﾈ", "ネ", "ﾉ", "ノ",
		"ﾊ", "ハ", "ﾋ", "ヒ", "ﾌ", "フ", "ﾍ", "ヘ", "ﾎ", "ホ",
		"ﾏ", "マ", "ﾐ", "ミ", "ﾑ", "ム", "ﾒ", "メ", "ﾓ", "モ",
		"ﾔ", "ヤ", "ﾕ", "ユ", "ﾖ", "ヨ",
		"ﾗ", "ラ", "ﾘ", "リ", "ﾙ", "ル", "ﾚ", "レ", "ﾛ", "ロ",
		"ﾜ", "ワ", "ｦ", "ヲ", "ﾝ", "ン",
		"ｧ", "ァ", "ｨ", "ィ", "ｩ", "ゥ", "ｪ", "ェ", "ｫ", "ォ",
		"ｬ", "ャ", "ｭ", "ュ", "ｮ", "ョ",
		"ｯ", "ッ",

		"ｰ", "ー",
		"ﾞ", "＂",
		"ﾟ", "",
		"｡", "。",
		"｢", "「",
		"｣", "」",
		"､", "、",
		"･", "・",
	)

	zenHanSymbolReplacer = strings.NewReplacer(
		"！", "!",
		"＂", "\"",
		"＃", "#",
		"＄", "$",
		"％", "%",
		"＆", "&",
		"＇", "'",
		"（", "(",
		"）", ")",
		"＊", "*",
		"＋", "+",
		"，", ",",
		"－", "-",
		"．", ".",
		"／", "/",
		"：", ":",
		"；", ";",
		"＜", "<",
		"＝", "=",
		"＞", ">",
		"？", "?",
		"＠", "@",
		"［", "[",
		"＼", "\\",
		"］", "]",
		"＾", "^",
		"＿", "_",
		"｀", "`",
		"｛", "{",
		"｜", "|",
		"｝", "}",
		"～", "~",
	)

	hiraganaKatakanaReplacer = strings.NewReplacer(
		"ゔ", "ヴ",
		"が", "ガ", "ぎ", "ギ", "ぐ", "グ", "げ", "ゲ", "ご", "ゴ",
		"ざ", "ザ", "じ", "ジ", "ず", "ズ", "ぜ", "ゼ", "ぞ", "ゾ",
		"だ", "ダ", "ぢ", "ヂ", "づ", "ヅ", "で", "デ", "ど", "ド",
		"ば", "バ", "び", "ビ", "ぶ", "ブ", "べ", "ベ", "ぼ", "ボ",
		"ぱ", "パ", "ぴ", "ピ", "ぷ", "プ", "ぺ", "ペ", "ぽ", "ポ",

		"あ", "ア", "い", "イ", "う", "ウ", "え", "エ", "お", "オ",
		"か", "カ", "き", "キ", "く", "ク", "け", "ケ", "こ", "コ",
		"さ", "サ", "し", "シ", "す", "ス", "せ", "セ", "そ", "ソ",
		"た", "タ", "ち", "チ", "つ", "ツ", "て", "テ", "と", "ト",
		"な", "ナ", "に", "ニ", "ぬ", "ヌ", "ね", "ネ", "の", "ノ",
		"は", "ハ", "ひ", "ヒ", "ふ", "フ", "へ", "ヘ", "ほ", "ホ",
		"ま", "マ", "み", "ミ", "む", "ム", "め", "メ", "も", "モ",
		"や", "ヤ", "ゆ", "ユ", "よ", "ヨ",
		"ら", "ラ", "り", "リ", "る", "ル", "れ", "レ", "ろ", "ロ",
		"わ", "ワ", "を", "ヲ", "ん", "ン",
		"ぁ", "ァ", "ぃ", "ィ", "ぅ", "ゥ", "ぇ", "ェ", "ぉ", "ォ",
		"ゃ", "ャ", "ゅ", "ュ", "ょ", "ョ",
		"っ", "ッ",
	)

	katakanaHiraganaReplacer = strings.NewReplacer(
		"ヴ", "ゔ",
		"ガ", "が", "ギ", "ぎ", "グ", "ぐ", "ゲ", "げ", "ゴ", "ご",
		"ザ", "ざ", "ジ", "じ", "ズ", "ず", "ゼ", "ぜ", "ゾ", "ぞ",
		"ダ", "だ", "ヂ", "ぢ", "ヅ", "づ", "デ", "で", "ド", "ど",
		"バ", "ば", "ビ", "び", "ブ", "ぶ", "ベ", "べ", "ボ", "ぼ",
		"パ", "ぱ", "ピ", "ぴ", "プ", "ぷ", "ペ", "ぺ", "ポ", "ぽ",

		"ア", "あ", "イ", "い", "ウ", "う", "エ", "え", "オ", "お",
		"カ", "か", "キ", "き", "ク", "く", "ケ", "け", "コ", "こ",
		"サ", "さ", "シ", "し", "ス", "す", "セ", "せ", "ソ", "そ",
		"タ", "た", "チ", "ち", "ツ", "つ", "テ", "て", "ト", "と",
		"ナ", "な", "ニ", "に", "ヌ", "ぬ", "ネ", "ね", "ノ", "の",
		"ハ", "は", "ヒ", "ひ", "フ", "ふ", "ヘ", "へ", "ホ", "ほ",
		"マ", "ま", "ミ", "み", "ム", "む", "メ", "め", "モ", "も",
		"ヤ", "や", "ユ", "ゆ", "ヨ", "よ",
		"ラ", "ら", "リ", "り", "ル", "る", "レ", "れ", "ロ", "ろ",
		"ワ", "わ", "ヲ", "を", "ン", "ん",
		"ァ", "ぁ", "ィ", "ぃ", "ゥ", "ぅ", "ェ", "ぇ", "ォ", "ぉ",
		"ャ", "ゃ", "ュ", "ゅ", "ョ", "ょ",
		"ッ", "っ",
	)
)

func NormalizeString(str string) string {
	str = NormalizeSpace(str)
	str = ZenToHanNumer(str)
	str = ZenToHanAlphabet(str)
	str = HanToZenKana(str)
	str = ZenToHanSymbol(str)

	return str
}

func NormalizeSpace(str string) string {
	str = strings.TrimSpace(str)
	str = spaceReplacer.Replace(str)

	return str
}

func ZenToHanNumer(str string) string {
	return zenHanNumberReplacer.Replace(str)
}

func ZenToHanAlphabet(str string) string {
	return zenHanAlphabetReplacer.Replace(str)
}

func HanToZenKana(str string) string {
	return hanZenKanaReplacer.Replace(str)
}

func ZenToHanSymbol(str string) string {
	return zenHanSymbolReplacer.Replace(str)
}

func SjisToUtf8(str string) string {
	ret, _ := ioutil.ReadAll(transform.NewReader(strings.NewReader(str), japanese.ShiftJIS.NewDecoder()))

	return string(ret)
}

func ToHash(s string) (string, error) {
	converted, err := bcrypt.GenerateFromPassword([]byte(s), 10)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(converted[:]), nil
}

func CheckHash(shash string, s string) error {
	_shash, _ := hex.DecodeString(shash)

	err := bcrypt.CompareHashAndPassword(_shash, []byte(s))
	if err != nil {
		return err
	}

	return nil
}

func HiraganaToKatakana(str string) string {
	return hiraganaKatakanaReplacer.Replace(str)
}

func KatakanaToHiragana(str string) string {
	return katakanaHiraganaReplacer.Replace(str)
}
