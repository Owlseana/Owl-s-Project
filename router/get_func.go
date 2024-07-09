package router

import (
	"MyFirstProgram/dao"
	"MyFirstProgram/model"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// GetArticleByID 从数据库中获取文章
// func (h *ArticleHandler)  (id string) (model.Article, error) {
// 	var article model.Article
// 	query := "SELECT id, title, author, content, created_at, updated_at FROM articles WHERE id=?"
// 	err := h.DB.QueryRow(query, id).Scan(&article.ID, &article.Title, &article.Author, &article.Content, &article.CreatedAt, &article.UpdatedAt)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return article, errors.New("article not found")
// 		}
// 		return article, err
// 	}
// 	return article, nil
// }

// 通过id查找文章
func getArticleByID(c *gin.Context) {
	var article model.Article
	id := c.Param("id")
	query := "SELECT id, title, author, content, created_at, updated_at FROM articles WHERE id=?"
	err := dao.DB.QueryRow(query, id).Scan(&article.ID, &article.Title, &article.Author, &article.Content, &article.CreatedAt, &article.UpdatedAt)
	if err != nil {
		c.JSON(404, gin.H{"message": "Article not found"})
		return
	}
	c.JSON(200, article)
}

// 通过title查找文章
func getArticleByTitle(c *gin.Context) {
	var article model.Article
	title := c.Param("title")
	query := "SELECT id, title, author, content, created_at, updated_at FROM articles WHERE title=?"
	err := dao.DB.QueryRow(query, title).Scan(&article.ID, &article.Title, &article.Author, &article.Content, &article.CreatedAt, &article.UpdatedAt)
	if err != nil {
		c.JSON(404, gin.H{"message": "Article not found"})
		return
	}
	c.JSON(200, article)
}

// 通过创作时间区间查找文章
func getArticleByCreatedTime(c *gin.Context) {
	// 获取查询参数
	startTime := c.Query("start_time")
	endTime := c.Query("end_time")

	start, err := time.Parse("2006-01-02 15:04:05", startTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid start_time format"})
		return
	}

	end, err := time.Parse("2006-01-02 15:04:05", endTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid end_time format"})
		return
	}

	query := "SELECT id, title, author, content, created_at, updated_at FROM articles WHERE created_at BETWEEN ? AND ?"
	rows, err := dao.DB.Query(query, start, end)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Database query failed"})
		return
	}
	defer rows.Close()

	var articles []model.Article

	for rows.Next() {
		var article model.Article
		if err := rows.Scan(&article.ID, &article.Title, &article.Author, &article.Content, &article.CreatedAt, &article.UpdatedAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to scan article"})
			return
		}
		articles = append(articles, article)
	}

	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error occurred during rows iteration"})
		return
	}

	// 检查是否查询到结果
	if len(articles) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "No articles found"})
		return
	}

	// 返回结果
	c.JSON(http.StatusOK, articles)
}

// 增加功能
func createArticle(c *gin.Context) {
	var article model.Article
	err := c.BindJSON(&article)
	if err != nil {
		c.JSON(400, gin.H{"message": "Invalid request body"})
		return
	}
	result, err := dao.DB.Exec("INSERT INTO articles (title, content, author) VALUES (?, ?, ?)", article.Title, article.Content, article.Author)
	if err != nil {
		c.JSON(500, gin.H{"message": "Failed to create article"})
		return
	}
	id, _ := result.LastInsertId()
	article.ID = uint(id)
	var inserted model.Article
	query := "SELECT id, title, author, content, created_at, updated_at FROM articles WHERE id=?"
	dao.DB.QueryRow(query, id).Scan(&inserted.ID, &inserted.Title, &inserted.Author, &inserted.Content, &inserted.CreatedAt, &inserted.UpdatedAt)
	c.JSON(201, inserted)
	fmt.Println("Inserted Successfully!")
}

// 测试功能
func health(c *gin.Context) {
	c.JSON(200, gin.H{"status": "ok"})
}

// 修改功能(只能通过ID更改)
func updateArticle(c *gin.Context) {
	var article model.Article
	id := c.Param("id")
	err := c.BindJSON(&article)
	if err != nil {
		c.JSON(400, gin.H{"message": "Invalid request body"})
		return
	}
	var currentTimestamp string
	dao.DB.QueryRow("SELECT CURRENT_TIMESTAMP").Scan(&currentTimestamp)

	_, err = dao.DB.Exec("UPDATE articles SET title=?, author=?, content=?, updated_at=(SELECT CURRENT_TIMESTAMP) WHERE id=?", article.Title, article.Author, article.Content, id)
	if err != nil {
		c.JSON(500, gin.H{"message": "Failed to update article"})
		return
	}
	var updated model.Article
	query := "SELECT id, title, author, content, created_at, updated_at FROM articles WHERE id=?"
	dao.DB.QueryRow(query, id).Scan(&updated.ID, &updated.Title, &updated.Author, &updated.Content, &updated.CreatedAt, &updated.UpdatedAt)
	//updated.UpdatedAt = service.TimeinBeiJing()
	c.JSON(200, updated)
	fmt.Println("Updated Successfully!")
}

// 删除功能
func deleteArticle(c *gin.Context) {
	id := c.Param("id")
	// 查询文章
	var article model.Article
	err := dao.DB.QueryRow("SELECT * FROM articles WHERE id=?", id).Scan(&article.ID, &article.Title, &article.Author, &article.Content, &article.CreatedAt, &article.UpdatedAt)
	if err != nil {
		c.JSON(404, gin.H{"message": "Article not found"})
		return
	}
	c.JSON(200, gin.H{
		"message": "Are you sure you want to delete this article?",
		"article": article,
		"confirm": "Input 1 to confirm, 0 to cancel",
	})

	// 读取用户输入
	var input int
	if err := c.BindJSON(&input); err != nil {
		c.JSON(400, gin.H{"message": "Invalid input"})
		return
	}

	if input == 1 {
		_, err := dao.DB.Exec("DELETE FROM articles WHERE id=?", id)
		if err != nil {
			c.JSON(500, gin.H{"message": "Failed to delete article"})
			return
		}
		c.JSON(200, gin.H{"message": "Article deleted"})
	} else {
		c.JSON(200, gin.H{"message": "Article deletion cancelled"})
	}
}
