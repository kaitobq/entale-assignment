package models

type Article struct {
	ID          uint    `gorm:"primarykey" json:"id"`
	Title       string  `gorm:"type:varchar(255)" json:"title"`
	Body        string  `gorm:"type:text" json:"body"`
	Medias      []Media `gorm:"foreignKey:ArticleID" json:"medias"`
	PublishedAt string  `gorm:"type:varchar(50)"`
}

type Media struct {
	ID          uint   `gorm:"primarykey" json:"id"`
	ArticleID   uint   `gorm:"index" json:"article_id"`
	ContentUrl  string `gorm:"type:varchar(255)" json:"contentUrl"`
	ContentType string `gorm:"type:varchar(255)" json:"contentType"`
}
