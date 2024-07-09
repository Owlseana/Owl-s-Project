package model

import "time"

// type Article struct {
// 	ID      int    `json:"id"`
// 	Title   string `json:"title"`
// 	Content string `json:"content"`
// }

type Article struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement;comment: '文章编号'"`
	Title     string    `json:"title" gorm:"size:50;not null;default:'';comment: '文章标题'"`
	Author    string    `json:"author" gorm:"size:50;not null;default:'';comment: '作者'"`
	Content   string    `json:"content" gorm:"size:100;not null;default:'';comment: '内容'"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime;comment: '创作年份'"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime;comment: '更新时间'"`
}
