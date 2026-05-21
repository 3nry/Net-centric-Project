package services

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func FetchChapters(mangaDexID string) ([]gin.H, error) {

	url := fmt.Sprintf(
		"https://api.mangadex.org/manga/%s/feed?translatedLanguage[]=en&order[chapter]=asc&limit=50",
		mangaDexID,
	)

	resp, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var result map[string]interface{}

	json.NewDecoder(resp.Body).Decode(&result)

	chapters := []gin.H{}

	if data, ok := result["data"].([]interface{}); ok {

		for _, item := range data {

			ch := item.(map[string]interface{})

			attrs := ch["attributes"].(map[string]interface{})

			chapters = append(chapters, gin.H{
				"id":      ch["id"],
				"chapter": attrs["chapter"],
				"title":   attrs["title"],
				"pages":   attrs["pages"],
			})
		}
	}

	return chapters, nil
}

func FetchPages(chapterID string) ([]string, error) {

	url := fmt.Sprintf(
		"https://api.mangadex.org/at-home/server/%s",
		chapterID,
	)

	resp, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var result map[string]interface{}

	json.NewDecoder(resp.Body).Decode(&result)

	baseURL := result["baseUrl"].(string)

	chapter := result["chapter"].(map[string]interface{})

	hash := chapter["hash"].(string)

	dataArr := chapter["data"].([]interface{})

	pages := []string{}

	for _, p := range dataArr {

		filename := p.(string)

		pages = append(
			pages,
			fmt.Sprintf("%s/data/%s/%s", baseURL, hash, filename),
		)
	}

	return pages, nil
}
