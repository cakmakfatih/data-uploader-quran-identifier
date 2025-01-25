package quran

import "encoding/xml"

type Quran struct {
	XMLName xml.Name `xml:"quran"`
	Suras   []Surah  `xml:"sura"`
}

type Surah struct {
	XMLName xml.Name `xml:"sura"`
	Index   int      `xml:"index,attr"`
	Name    string   `xml:"name,attr"`
	Ayas    []Aya    `xml:"aya"`
}

type Aya struct {
	XMLName xml.Name `xml:"aya" json:"aya"`
	Index   int      `xml:"index,attr" json:"index"`
	Text    string   `xml:"text,attr" json:"text"`
}

type QuranMetadata struct {
	XMLName xml.Name `xml:"quran"`
	Suras   Suras    `xml:"suras"`
	Juzes   Juzes    `xml:"juzs"`
	Hizbs   Hizbs    `xml:"hizbs"`
	Manzils Manzils  `xml:"manzils"`
	Rukus   Rukus    `xml:"rukus"`
	Pages   Pages    `xml:"pages"`
	Sajdas  Sajdas   `xml:"sajdas"`
}

type Suras struct {
	XMLName xml.Name `xml:"suras"`
	Sura    []Sura   `xml:"sura"`
}

type Sura struct {
	Index int    `xml:"index,attr"`
	Ayas  int    `xml:"ayas,attr"`
	Start int    `xml:"start,attr"`
	Name  string `xml:"name,attr"`
	Tname string `xml:"tname,attr"`
	Ename string `xml:"ename,attr"`
	Type  string `xml:"type,attr"`
	Order int    `xml:"order,attr"`
	Rukus int    `xml:"rukus,attr"`
}

type Juzes struct {
	XMLName xml.Name `xml:"juzs"`
	Juz     []Juz    `xml:"juz"`
}

type Juz struct {
	Index int    `xml:"index,attr"`
	Sura  string `xml:"sura,attr"`
	Aya   string `xml:"aya,attr"`
}

type Hizbs struct {
	XMLName xml.Name  `xml:"hizbs"`
	Alias   string    `xml:"alias,attr"`
	Quarter []Quarter `xml:"quarter"`
}

type Quarter struct {
	Index int    `xml:"index,attr"`
	Sura  string `xml:"sura,attr"`
	Aya   string `xml:"aya,attr"`
}

type Manzils struct {
	XMLName xml.Name `xml:"manzils"`
	Alias   string   `xml:"alias,attr"`
	Manzil  []Manzil `xml:"manzil"`
}

type Manzil struct {
	Index int    `xml:"index,attr"`
	Sura  string `xml:"sura,attr"`
	Aya   string `xml:"aya,attr"`
}

type Rukus struct {
	XMLName xml.Name `xml:"rukus"`
	Alias   string   `xml:"alias,attr"`
	Ruku    []Ruku   `xml:"ruku"`
}

type Ruku struct {
	Index int    `xml:"index,attr"`
	Sura  string `xml:"sura,attr"`
	Aya   string `xml:"aya,attr"`
}

type Pages struct {
	XMLName xml.Name `xml:"pages"`
	Page    []Page   `xml:"page"`
}

type Page struct {
	Index int    `xml:"index,attr"`
	Sura  string `xml:"sura,attr"`
	Aya   string `xml:"aya,attr"`
}

type Sajdas struct {
	XMLName xml.Name `xml:"sajdas"`
	Sajda   []Sajda  `xml:"sajda"`
}

type Sajda struct {
	Index int    `xml:"index,attr"`
	Sura  string `xml:"sura,attr"`
	Aya   string `xml:"aya,attr"`
	Type  string `xml:"type,attr"`
}
