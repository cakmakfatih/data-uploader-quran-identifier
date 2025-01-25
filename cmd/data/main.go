package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"quranarango/core"
	"quranarango/quran"
	"strconv"
	"strings"
	"sync"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func getPageOfVerse(chapterIndex int, verseIndex int, qMetaData quran.QuranMetadata) int {
	for _, p := range qMetaData.Pages.Page {
		suraIdx, _ := strconv.Atoi(p.Sura)
		ayaIdx, _ := strconv.Atoi(p.Aya)

		if suraIdx > chapterIndex {
			return p.Index - 1
		} else if suraIdx == chapterIndex && ayaIdx > verseIndex {
			return p.Index - 1
		} else if suraIdx == chapterIndex && ayaIdx == verseIndex {
			return p.Index
		}
	}

	return qMetaData.Pages.Page[len(qMetaData.Pages.Page)-1].Index
}

func getQuarterOfVerse(chapterIndex, verseIndex int, qMetaData quran.QuranMetadata) int {
	for _, q := range qMetaData.Hizbs.Quarter {
		suraIdx, _ := strconv.Atoi(q.Sura)
		ayaIdx, _ := strconv.Atoi(q.Aya)

		if suraIdx > chapterIndex {
			return q.Index - 1
		} else if suraIdx == chapterIndex && ayaIdx > verseIndex {
			return q.Index - 1
		} else if suraIdx == chapterIndex && ayaIdx == verseIndex {
			return q.Index
		}
	}

	return qMetaData.Hizbs.Quarter[len(qMetaData.Hizbs.Quarter)-1].Index
}

func getJuzOfVerse(chapterIndex, verseIndex int, qMetaData quran.QuranMetadata) int {
	for _, j := range qMetaData.Juzes.Juz {
		suraIdx, _ := strconv.Atoi(j.Sura)
		ayaIdx, _ := strconv.Atoi(j.Aya)

		if suraIdx > chapterIndex {
			return j.Index - 1
		} else if suraIdx == chapterIndex && ayaIdx > verseIndex {
			return j.Index - 1
		} else if suraIdx == chapterIndex && ayaIdx == verseIndex {
			return j.Index
		}
	}

	return qMetaData.Juzes.Juz[len(qMetaData.Juzes.Juz)-1].Index
}

func getIsRuku(chapterIdx int, verseIdx int, qMetaData quran.QuranMetadata) bool {
	for _, ruku := range qMetaData.Rukus.Ruku {
		chapter, _ := strconv.Atoi(ruku.Sura)
		verse, _ := strconv.Atoi(ruku.Aya)

		if chapter == chapterIdx && verse == verseIdx {
			return true
		}
	}

	return false
}

func getSajda(chapterIndex int, verseIndex int, qMetaData quran.QuranMetadata) string {
	for _, sajda := range qMetaData.Sajdas.Sajda {
		if fmt.Sprintf("%d", chapterIndex) == sajda.Sura && fmt.Sprintf("%d", verseIndex) == sajda.Aya {
			return sajda.Type
		}
	}

	return "none"
}

func findSuraByIndex(index int, qMetaData quran.QuranMetadata) (*quran.Sura, bool) {
	for _, sura := range qMetaData.Suras.Sura {
		if sura.Index == index {
			return &sura, true
		}
	}

	return nil, false
}

func initializeData(driver neo4j.DriverWithContext) {
	ctx := context.Background()
	session := driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: "neo4j"})
	defer session.Close(ctx)

	wd, _ := os.Getwd()

	xmlPathQuran := filepath.Join(wd, "static", "quran.xml")
	xmlPathQuranMetadata := filepath.Join(wd, "static", "metadata.xml")
	qXmlData, _ := os.ReadFile(xmlPathQuran)
	qmXmlData, _ := os.ReadFile(xmlPathQuranMetadata)

	translationsPath := filepath.Join(wd, "static", "translations")
	transliterationsPath := filepath.Join(wd, "static", "transliterations")
	translationPaths, _ := os.ReadDir(translationsPath)
	transliterationPaths, _ := os.ReadDir(transliterationsPath)

	xmlPaths := []string{}

	for _, p := range translationPaths {
		xmlPaths = append(xmlPaths, filepath.Join(translationsPath, p.Name()))
	}

	for _, p := range transliterationPaths {
		xmlPaths = append(xmlPaths, filepath.Join(transliterationsPath, p.Name()))
	}

	var q quran.Quran
	var qm quran.QuranMetadata

	quranTrData := map[string]*quran.Quran{}

	for _, p := range xmlPaths {
		var d quran.Quran
		xmlData, _ := os.ReadFile(p)
		xml.Unmarshal(xmlData, &d)

		quranTrData[p] = &d
	}

	xml.Unmarshal(qXmlData, &q)
	xml.Unmarshal(qmXmlData, &qm)

	for _, juz := range qm.Juzes.Juz {
		session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
			startChapter, _ := strconv.Atoi(juz.Sura)
			startVerse, _ := strconv.Atoi(juz.Aya)

			juzId, err := saveJuz(tx, &CreateJuzTransaction{
				Index:        juz.Index,
				StartChapter: startChapter,
				StartVerse:   startVerse,
			})

			if err != nil {
				panic(err)
			}

			return juzId, nil
		})
	}

	for _, quarter := range qm.Hizbs.Quarter {
		session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
			startChapter, _ := strconv.Atoi(quarter.Sura)
			startVerse, _ := strconv.Atoi(quarter.Aya)

			quarterId, err := saveQuarter(tx, &CreateQuarterTransaction{
				Index:        quarter.Index,
				StartChapter: startChapter,
				StartVerse:   startVerse,
			})

			if err != nil {
				panic(err)
			}

			return quarterId, nil
		})
	}

	for _, page := range qm.Pages.Page {
		session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
			startChapterIdx, _ := strconv.Atoi(page.Sura)
			startVerseIdx, _ := strconv.Atoi(page.Aya)

			pageId, err := savePage(tx, &CreatePageTransaction{
				PageIdx:      page.Index,
				StartChapter: startChapterIdx,
				StartVerse:   startVerseIdx,
			})

			if err != nil {
				panic(err)
			}

			return pageId, err
		})
	}

	var wg sync.WaitGroup
	for chapterIdx, chapter := range q.Suras {
		chapterMetaData, exists := findSuraByIndex(chapter.Index, qm)
		if !exists {
			log.Panic("could not find chapter meta data")
		}

		chapterId, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
			chapterId, err := saveChapter(tx, &ChapterTransaction{
				Index: chapter.Index,
				Start: chapterMetaData.Start,
				Name:  chapterMetaData.Name,
				Tname: chapterMetaData.Tname,
				Ename: chapterMetaData.Ename,
				Type:  chapterMetaData.Type,
				Order: chapterMetaData.Order,
			})

			if err != nil {
				return nil, err
			}

			return chapterId, nil
		})

		if err != nil {
			log.Panic(err)
		}

		log.Printf("added chapter: %s", chapterMetaData.Tname)

		wg.Add(1)

		go func() {
			for verseIdx, verse := range chapter.Ayas {
				reference := fmt.Sprintf("%v:%v", chapter.Index, verse.Index)
				sajda := getSajda(chapter.Index, verse.Index, qm)
				isRuku := getIsRuku(chapter.Index, verse.Index, qm)

				verseId, err := session.ExecuteWrite(ctx,
					func(tx neo4j.ManagedTransaction) (any, error) {
						verseId, err := saveVerse(tx, &VerseTransaction{
							ChapterId: chapterId.(*int64),
							Reference: reference,
							Chapter:   chapter.Index,
							Index:     verse.Index,
							Text:      verse.Text,
							Sajda:     sajda,
							IsRuku:    isRuku,
						})

						if err != nil {
							return nil, err
						}

						for xmlPath, quranData := range quranTrData {
							saveQuranTrData(tx, &QuranTrDataTransaction{
								VerseId:    verseId,
								Reference:  reference,
								ChapterIdx: chapterIdx,
								VerseIdx:   verseIdx,
								XmlPath:    xmlPath,
								QuranData:  quranData,
							})
						}

						pageIdx := getPageOfVerse(chapter.Index, verse.Index, qm)
						juzIdx := getJuzOfVerse(chapter.Index, verse.Index, qm)
						quarterIdx := getQuarterOfVerse(chapter.Index, verse.Index, qm)

						linkVerseToJuz(tx, &LinkVerseToJuzTransaction{
							VerseId: verseId,
							JuzIdx:  juzIdx,
						})
						linkVerseToPage(tx, &LinkVerseToPageTransaction{
							PageIdx: pageIdx,
							VerseId: verseId,
						})
						linkVerseToQuarter(tx, &LinkVerseToQuarterTransaction{
							QuarterIdx: quarterIdx,
							VerseId:    verseId,
						})

						if isRuku {
							linkRukuToVerse(tx, &LinkRukuToVerseTransaction{
								VerseId:    verseId,
								ChapterIdx: chapter.Index,
								VerseIdx:   verse.Index,
							})
						}

						return verseId, nil
					})

				if err != nil {
					log.Println(err)
				}
				log.Printf("saved verse with id %v", verseId)
			}

			wg.Done()
		}()
	}
	wg.Wait()
}

type QuranTrDataTransaction struct {
	VerseId    *int64
	Reference  string
	ChapterIdx int
	VerseIdx   int
	XmlPath    string
	QuranData  *quran.Quran
}

type LinkVerseToPageTransaction struct {
	PageIdx int
	VerseId *int64
}

type LinkVerseToJuzTransaction struct {
	VerseId *int64
	JuzIdx  int
}

type LinkRukuToVerseTransaction struct {
	VerseId    *int64
	ChapterIdx int
	VerseIdx   int
}

func linkVerseToJuz(tx neo4j.ManagedTransaction, ts *LinkVerseToJuzTransaction) (*int64, error) {
	ctx := context.Background()
	res, err := tx.Run(ctx, `
	MATCH(j: Juz)
	MATCH(v: Verse)
	WHERE id(v) = $verseId AND
	j.index = $juzIdx
	CREATE (v) <- [rel:HAS] - (j)
	RETURN id(rel) AS nodeId`,
		map[string]any{
			"verseId": ts.VerseId,
			"juzIdx":  ts.JuzIdx,
		})

	if err != nil {
		return nil, err
	}

	resSingle, err := res.Single(ctx)
	if err != nil {
		return nil, err
	}

	relId := resSingle.AsMap()["nodeId"].(int64)
	return &relId, nil
}

func linkRukuToVerse(tx neo4j.ManagedTransaction, ts *LinkRukuToVerseTransaction) (*int64, error) {
	ctx := context.Background()
	res, err := tx.Run(ctx, `
	MATCH(r: Ruku)
	MATCH(v: Verse) 
	WHERE id(v) = $verseId AND
	r.chapter = $chapterIdx AND r.verse = $verseIdx
	CREATE (v) - [rel:IS] -> (r)
	RETURN id(rel) as nodeId`,
		map[string]any{
			"verseId":    ts.VerseId,
			"verseIdx":   ts.VerseIdx,
			"chapterIdx": ts.ChapterIdx,
		})

	if err != nil {
		return nil, err
	}

	resSingle, err := res.Single(ctx)
	if err != nil {
		return nil, err
	}

	relationId := resSingle.AsMap()["nodeId"].(int64)
	return &relationId, nil
}

type LinkVerseToQuarterTransaction struct {
	QuarterIdx int
	VerseId    *int64
}

func linkVerseToQuarter(tx neo4j.ManagedTransaction, ts *LinkVerseToQuarterTransaction) (*int64, error) {
	ctx := context.Background()
	res, err := tx.Run(ctx, `
	MATCH(v: Verse) WHERE id(v) = $verseId
	MATCH(q: Quarter) WHERE q.index = $quarterIdx
	CREATE (v) <- [r:HAS] - (q)
	RETURN id(r) as nodeId`,
		map[string]any{
			"verseId":    ts.VerseId,
			"quarterIdx": ts.QuarterIdx,
		})

	if err != nil {
		return nil, err
	}

	resSingle, err := res.Single(ctx)
	if err != nil {
		return nil, err
	}

	relId := resSingle.AsMap()["nodeId"].(int64)

	return &relId, nil
}

func linkVerseToPage(tx neo4j.ManagedTransaction, ts *LinkVerseToPageTransaction) (*int64, error) {
	ctx := context.Background()
	res, err := tx.Run(ctx, `
			MATCH(v: Verse) WHERE id(v) = $verseId
			MATCH(p: Page) WHERE p.index = $pageIdx
			CREATE (v) <- [r:HAS] - (p)
			RETURN id(r) as nodeId`,
		map[string]any{
			"verseId": ts.VerseId,
			"pageIdx": ts.PageIdx,
		})

	if err != nil {
		return nil, err
	}

	resSingle, err := res.Single(ctx)
	if err != nil {
		return nil, err
	}

	relationId := resSingle.AsMap()["nodeId"].(int64)
	return &relationId, nil
}

func saveQuranTrData(tx neo4j.ManagedTransaction, trDataTransaction *QuranTrDataTransaction) {
	xmlPathSplitted := strings.Split(trDataTransaction.XmlPath, "\\")
	fileName := xmlPathSplitted[len(xmlPathSplitted)-1]
	fileNameSplitted := strings.Split(fileName, ".")

	language := fileNameSplitted[0]

	if strings.Contains(trDataTransaction.XmlPath, "transliterations") {
		text := trDataTransaction.QuranData.Suras[trDataTransaction.ChapterIdx].Ayas[trDataTransaction.VerseIdx].Text

		saveTransliteration(tx, &TransliterationTransaction{
			VerseId:   trDataTransaction.VerseId,
			Lang:      language,
			Reference: trDataTransaction.Reference,
			Text:      text,
		})
	} else if strings.Contains(trDataTransaction.XmlPath, "translations") {
		by := fileNameSplitted[1]

		text := trDataTransaction.QuranData.Suras[trDataTransaction.ChapterIdx].Ayas[trDataTransaction.VerseIdx].Text

		saveTranslation(tx, &TranslationTransaction{
			VerseId:   trDataTransaction.VerseId,
			Lang:      language,
			Reference: trDataTransaction.Reference,
			Text:      text,
			By:        by,
		})
	}
}

type TranslationTransaction struct {
	VerseId   *int64
	Reference string
	Text      string
	By        string
	Lang      string
}

type VerseTransaction struct {
	ChapterId *int64
	Reference string
	Chapter   int
	Index     int
	Text      string
	Sajda     string
	IsRuku    bool
}

type ChapterTransaction struct {
	Index int
	Start int
	Name  string
	Tname string
	Ename string
	Type  string
	Order int
}

type TransliterationTransaction struct {
	VerseId   *int64
	Reference string
	Text      string
	Lang      string
}

func saveTransliteration(tx neo4j.ManagedTransaction, transliteration *TransliterationTransaction) (*int64, error) {
	ctx := context.Background()
	res, err := tx.Run(ctx, `
	MATCH(v: Verse) WHERE id(v) = $verseId
	CREATE(t: Transliteration {
		reference: $reference,
		text: $text,
		lang: $lang
	}) - [:TRANSLITERATE] -> (v) 
		RETURN id(t) AS nodeId`, map[string]any{
		"verseId":   transliteration.VerseId,
		"reference": transliteration.Reference,
		"text":      transliteration.Text,
		"lang":      transliteration.Lang,
	})

	if err != nil {
		return nil, err
	}

	resSingle, err := res.Single(ctx)
	if err != nil {
		return nil, err
	}

	transliterationId := resSingle.AsMap()["nodeId"].(int64)

	return &transliterationId, nil
}

func saveVerse(tx neo4j.ManagedTransaction, verse *VerseTransaction) (*int64, error) {
	ctx := context.Background()
	res, err := tx.Run(ctx, `
		MATCH(c: Chapter) WHERE id(c) = $chapterId
		CREATE(v: Verse {
			reference: $reference,
			chapter: $chapter,
			index: $index,
			text: $text
		}) - [:OWNED_BY] -> (c)
		RETURN id(v) AS nodeId`,
		map[string]any{
			"chapterId": verse.ChapterId,
			"reference": verse.Reference,
			"chapter":   verse.Chapter,
			"index":     verse.Index,
			"text":      verse.Text,
		})

	if err != nil {
		log.Panic(err)
		return nil, err
	}

	resSingle, err := res.Single(ctx)

	if err != nil {
		return nil, err
	}

	savedId := resSingle.AsMap()["nodeId"].(int64)
	if verse.IsRuku {
		_, err = tx.Run(ctx, `
		MATCH(v: Verse) WHERE id(v) = $verseId
		SET v:Ruku
		RETURN id(v) as nodeId`, map[string]any{
			"verseId": savedId,
		})

		if err != nil {
			log.Panic(err)
		}
	}

	if verse.Sajda != "none" {
		_, err = tx.Run(ctx, `
		MATCH(v: Verse) WHERE id(v) = $verseId
		SET v.sajda = $sajda
		SET v:Sajda
		RETURN id(v) as nodeId`, map[string]any{
			"verseId": savedId,
			"sajda":   verse.Sajda,
		})

		if err != nil {
			log.Panic(err)
		}
	}

	return &savedId, nil
}

func saveTranslation(tx neo4j.ManagedTransaction, translation *TranslationTransaction) (*int64, error) {
	ctx := context.Background()
	res, err := tx.Run(ctx, `
		MATCH(v: Verse) WHERE id(v) = $verseId
		CREATE(t: Translation {
			text: $text,
			reference: $reference,
			by: $by,
			lang: $lang
		}) - [:TRANSLATE] -> (v)
		RETURN id(t) AS nodeId`,
		map[string]any{
			"verseId":   translation.VerseId,
			"text":      translation.Text,
			"reference": translation.Reference,
			"by":        translation.By,
			"lang":      translation.Lang,
		})

	if err != nil {
		return nil, err
	}

	resSingle, err := res.Single(ctx)

	if err != nil {
		return nil, err
	}

	savedId := resSingle.AsMap()["nodeId"].(int64)

	return &savedId, nil
}

func saveChapter(tx neo4j.ManagedTransaction, chapter *ChapterTransaction) (*int64, error) {
	ctx := context.Background()
	res, err := tx.Run(
		ctx, `
			CREATE(c: Chapter {
				index: $index,
				start: $start,
				name: $name,
				tname: $tname,
				ename: $ename,
				type: $type,
				order: $order
			})
			RETURN id(c) AS nodeId`,
		map[string]any{
			"index": chapter.Index,
			"start": chapter.Start,
			"name":  chapter.Name,
			"tname": chapter.Tname,
			"ename": chapter.Ename,
			"type":  chapter.Type,
			"order": chapter.Order,
		},
	)

	if err != nil {
		return nil, err
	}

	resSingle, err := res.Single(ctx)

	if err != nil {
		return nil, err
	}

	savedId := resSingle.AsMap()["nodeId"].(int64)

	return &savedId, nil
}

type CreatePageTransaction struct {
	PageIdx      int
	StartChapter int
	StartVerse   int
}

type CreateJuzTransaction struct {
	Index        int
	StartChapter int
	StartVerse   int
}

type CreateQuarterTransaction struct {
	Index        int
	StartChapter int
	StartVerse   int
}

func saveQuarter(tx neo4j.ManagedTransaction, tr *CreateQuarterTransaction) (*int64, error) {
	ctx := context.Background()
	res, err := tx.Run(ctx, `
	CREATE(q: Quarter{
		index: $index,
		start_chapter: $startChapter,
		start_verse: $startVerse
	})
		RETURN id(q) AS nodeId`,
		map[string]any{
			"index":        tr.Index,
			"startChapter": tr.StartChapter,
			"startVerse":   tr.StartVerse,
		})

	if err != nil {
		return nil, err
	}

	resSingle, err := res.Single(ctx)
	if err != nil {
		return nil, err
	}

	quarterId := resSingle.AsMap()["nodeId"].(int64)
	return &quarterId, nil
}

func saveJuz(tx neo4j.ManagedTransaction, tr *CreateJuzTransaction) (*int64, error) {
	ctx := context.Background()
	res, err := tx.Run(ctx, `
	CREATE(j: Juz{
		index: $index,
		start_chapter: $startChapter,
		start_verse: $startVerse
	})
		RETURN id(j) as nodeId`,
		map[string]any{
			"index":        tr.Index,
			"startChapter": tr.StartChapter,
			"startVerse":   tr.StartVerse,
		})

	if err != nil {
		return nil, err
	}

	resSingle, err := res.Single(ctx)
	if err != nil {
		return nil, err
	}

	juzId := resSingle.AsMap()["nodeId"].(int64)
	return &juzId, nil
}

func savePage(tx neo4j.ManagedTransaction, transaction *CreatePageTransaction) (*int64, error) {
	ctx := context.Background()
	res, err := tx.Run(ctx, `
	CREATE (p: Page{
		index: $index,
		start_chapter: $startChapterIdx,
		start_verse: $startVerseIdx
	})
		RETURN id(p) AS nodeId`, map[string]any{
		"index":           transaction.PageIdx,
		"startChapterIdx": transaction.StartChapter,
		"startVerseIdx":   transaction.StartVerse,
	})

	if err != nil {
		return nil, err
	}

	resSingle, err := res.Single(ctx)
	if err != nil {
		return nil, err
	}

	pageId := resSingle.AsMap()["nodeId"].(int64)

	return &pageId, nil
}

func main() {
	core.LoadEnv()
	driver := core.ConnToNeoDb()

	initializeData(driver)
}
