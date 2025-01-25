package models

type Manzil struct {
	Index  int      `json:"index"`
	Verses []*Verse `json:"verses"`
}
type Quarter struct {
	Index  int      `json:"index"`
	Verses []*Verse `json:"verses"`
}
type Transliteration struct {
	Verse *Verse `json:"verse"`
	Lang  string `json:"lang"`
	Text  string `json:"text"`
}
type Translation struct {
	Verse *Verse `json:"verse"`
	Lang  string `json:"lang"`
	Text  string `json:"text"`
	By    string `json:"by"`
}
type Sajda struct {
	Index int    `json:"index"`
	Type  string `json:"type"`
	Verse *Verse `json:"verse"`
}

type Verse struct {
	Chapter          *Chapter           `json:"chapter"`
	Reference        string             `json:"reference"`
	ChapterIndex     string             `json:"chapter_index"`
	VerseIndex       string             `json:"verse_index"`
	Text             string             `json:"text"`
	Ruku             *Ruku              `json:"ruku,omitempty"`
	Transliterations []*Transliteration `json:"transliterations,omitempty"`
	Translations     []*Translation     `json:"translations,omitempty"`
	Sajda            *Sajda             `json:"sajda,omitempty"`
}

type Ruku struct {
	Index int    `json:"index"`
	Verse *Verse `json:"verse"`
}

type Chapter struct {
	Uid   string `json:"uid"`
	Index int    `json:"index"`
	Start int    `json:"start"`
	Tname string `json:"tname"`
	Name  string `json:"name"`
	Type  string `json:"type"`
	Order int    `json:"order"`
}
