package controllers

import (
	"encoding/json"
	"entale-test/models"
	"fmt"
	"io"
	"log"
	"net/http"

	"gorm.io/gorm"
)

type Article struct {
	ID          uint       `json:"id"`
	Title       string     `json:"title"`
	Body        string     `json:"body"`
	Medias      []Media    `json:"medias"`
	PublishedAt string     `json:"publishedAt"`
}

type Media struct {
	ID          uint   `json:"id"`
	ArticleID uint `json:"article_id"`
	ContentUrl  string `json:"contentUrl"`
	ContentType string `json:"contentType"`
}

func SaveArticles(w http.ResponseWriter, r *http.Request) {
	log.Println("Saving articles to database")
	url := "https://gist.githubusercontent.com/gotokatsuya/cc78c04d3af15ebe43afe5ad970bc334/raw/dc39bacb834105c81497ba08940be5432ed69848/articles.json"
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var articles []Article
	err = json.Unmarshal(body, &articles)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)
	}

	saveArticlesToDB(articles)
	fmt.Printf("Articles: %v", articles)
}

func saveArticlesToDB(articles []Article) {
	for _, article := range articles {
		err := models.DB.Transaction(func(tx *gorm.DB) error {
			var existingArticle models.Article
			if err := tx.Where("id = ?", article.ID).First(&existingArticle).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					if err := tx.Create(&article).Error; err != nil {
						return err
					}
				} else {
					return err
				}
			}

			for _, media := range article.Medias {
				media.ArticleID = article.ID
				var existingMedia models.Media
				if err := tx.Where("id = ?", media.ID).First(&existingMedia).Error; err != nil {
					if err == gorm.ErrRecordNotFound {
						if err := tx.Create(&media).Error; err != nil {
							return err
						}
					} else {
						return err
					}
				}
			}
			return nil
		})
		if err != nil {
			log.Printf("Error saving article %d: %v", article.ID, err)
		}
	}
}

func GetArticles(w http.ResponseWriter, r *http.Request){
	var articles []models.Article
	if err := models.DB.Preload("Medias").Find(&articles).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(articles); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
