package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"mangahub/config"
	"mangahub/models"
)

type MangaDexResponse struct {
	Data []MangaDexManga `json:"data"`
}
type MangaDexManga struct {
	ID         string             `json:"id"`
	Attributes MangaDexAttributes `json:"attributes"`
}
type MangaDexAttributes struct {
	Title       map[string]string `json:"title"`
	Description map[string]string `json:"description"`
	Status      string            `json:"status"`
	Tags        []MangaDexTag     `json:"tags"`
}
type MangaDexTag struct {
	Attributes struct {
		Name map[string]string `json:"name"`
	} `json:"attributes"`
}
type MangaDexCoverResponse struct {
	Data []struct {
		Attributes struct {
			FileName string `json:"fileName"`
		} `json:"attributes"`
	} `json:"data"`
}

func fetchCoverURL(mangaID string) string {
	resp, err := http.Get(fmt.Sprintf("https://api.mangadex.org/cover?manga[]=%s&limit=1", mangaID))
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var cr MangaDexCoverResponse
	if json.Unmarshal(body, &cr) != nil || len(cr.Data) == 0 {
		return ""
	}
	return fmt.Sprintf("https://uploads.mangadex.org/covers/%s/%s.256.jpg", mangaID, cr.Data[0].Attributes.FileName)
}

func fetchMangaDex(limit, offset int) ([]models.Manga, error) {
	resp, err := http.Get(fmt.Sprintf(
		"https://api.mangadex.org/manga?limit=%d&offset=%d&contentRating[]=safe&contentRating[]=suggestive",
		limit, offset,
	))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	var mdResp MangaDexResponse
	if err := json.Unmarshal(body, &mdResp); err != nil {
		return nil, err
	}

	var mangas []models.Manga
	for _, item := range mdResp.Data {
		attr := item.Attributes
		title := attr.Title["en"]
		if title == "" {
			title = attr.Title["ja-ro"]
		}
		if title == "" {
			for _, v := range attr.Title {
				title = v
				break
			}
		}
		if title == "" {
			continue
		}
		desc := attr.Description["en"]
		if len(desc) > 500 {
			desc = desc[:500] + "..."
		}
		var genres []string
		for _, tag := range attr.Tags {
			if name := tag.Attributes.Name["en"]; name != "" {
				genres = append(genres, name)
			}
		}
		mangas = append(mangas, models.Manga{
			MangaDexID:  item.ID,
			Title:       title,
			Description: desc,
			Author:      "MangaDex",
			Genre:       strings.Join(genres, ", "),
			Status:      attr.Status,
			CoverImage:  fetchCoverURL(item.ID),
		})
	}
	return mangas, nil
}

func SeedMangaFromMangaDex() {
	fmt.Println("Seeding manga from MangaDex...")
	total := 0
	for offset := 0; offset < 200; offset += 100 {
		mangas, err := fetchMangaDex(100, offset)
		if err != nil {
			fmt.Printf("  Error offset %d: %v\n", offset, err)
			continue
		}
		for _, m := range mangas {
			var existing models.Manga
			if config.DB.Where("manga_dex_id = ?", m.MangaDexID).First(&existing).Error != nil {
				config.DB.Create(&m)
				total++
			}
		}
		fmt.Printf("  Offset %d done, total: %d\n", offset, total)
	}
	fmt.Printf("Seed done. Total: %d\n", total)
}
