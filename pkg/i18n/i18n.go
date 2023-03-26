package i18n

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type Interface interface {
	Translate(key, language string, data map[string]string) string
}

type I18n struct {
	translations map[string]map[string]any
}

func NewI18n(localesPath string) (Interface, error) {
	entries, err := os.ReadDir(localesPath)
	if err != nil {
		return nil, err
	}

	translations := make(map[string]map[string]any)

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := strings.Replace(entry.Name(), ".json", "", 1)

		fileContent := make(map[string]any)
		f, err := os.ReadFile(filepath.Join(localesPath, entry.Name()))
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(f, &fileContent)
		translations[name] = fileContent
	}

	for langCode, _ := range translations {
		langNestedDirs, err := os.ReadDir(filepath.Join(localesPath, langCode))
		if err != nil {
			return nil, err
		}

		for _, dir := range langNestedDirs {
			files, err := os.ReadDir(filepath.Join(localesPath, langCode, dir.Name()))
			if err != nil {
				return nil, err
			}

			translations[langCode][dir.Name()] = make(map[string]any)

			for _, file := range files {
				fileContent := make(map[string]any)
				f, err := os.ReadFile(filepath.Join(localesPath, langCode, dir.Name(), file.Name()))
				if err != nil {
					return nil, err
				}
				err = json.Unmarshal(f, &fileContent)
				translations[langCode][dir.Name()].(map[string]any)[strings.Replace(file.Name(), ".json", "", 1)] = fileContent
			}
		}

	}

	return &I18n{
		translations: translations,
	}, nil
}

func (i *I18n) Translate(key, language string, data map[string]string) string {
	str, _ := GetNested[string](i.translations[language], strings.Split(key, ".")...)

	str = strings.ReplaceAll(str, "{{ ", "{{.")

	tmpl, err := template.New("t").Parse(str)
	if err != nil {
		return str
	}

	res := &strings.Builder{}
	err = tmpl.Execute(res, data)
	if err != nil {
		log.Fatal(err)
	}

	return res.String()
}
