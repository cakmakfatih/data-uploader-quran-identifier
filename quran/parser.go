package quran

import (
	"encoding/xml"
	"os"
	"path/filepath"
)

func ReadMetadata() (*QuranMetadata, error) {
	workingDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	path := filepath.Join(workingDir, "static", "metadata.xml")
	byteValue, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var quranMetaData QuranMetadata
	err = xml.Unmarshal(byteValue, &quranMetaData)
	if err != nil {
		return nil, err
	}

	return &quranMetaData, nil
}

func ReadQuran(path string) (*Quran, error) {
	byteValue, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var quran Quran
	err = xml.Unmarshal(byteValue, &quran)
	if err != nil {
		return nil, err
	}

	return &quran, nil
}
