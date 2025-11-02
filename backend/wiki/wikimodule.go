package modules

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Модель страницы Wiki
type WikiPage struct {
	ID        primitive.ObjectID  `bson:"_id,omitempty" json:"_id"`
	ParentID  *primitive.ObjectID `bson:"parentId,omitempty" json:"parentId,omitempty"`
	Title     string              `bson:"title" json:"title"`
	Content   string              `bson:"content" json:"content"`
	Emoji     string              `bson:"emoji" json:"emoji"`
	CreatedAt time.Time           `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time           `bson:"updatedAt" json:"updatedAt"`
}

// Структура модуля
type WikiModule struct {
	Collection *mongo.Collection
}

// Инициализация модуля
func NewWikiModule(client *mongo.Client) *WikiModule {
	return &WikiModule{
		Collection: client.Database("bell").Collection("wiki_pages"),
	}
}

// Регистрация маршрутов
func (m *WikiModule) RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api/wiki")
	{
		api.GET("/pages", m.GetAllPages)
		api.GET("/pages/:id", m.GetPage)
		api.POST("/pages", m.CreatePage)
		api.PUT("/pages/:id", m.UpdatePage)
		api.DELETE("/pages/:id", m.DeletePage)
	}
}

// Получение всех страниц
func (m *WikiModule) GetAllPages(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := m.Collection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var pages []WikiPage
	if cursor != nil {
		if err := cursor.All(ctx, &pages); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	if pages == nil {
		pages = []WikiPage{}
	}

	c.JSON(http.StatusOK, pages)
}

// Получение страницы по ID
func (m *WikiModule) GetPage(c *gin.Context) {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный ID"})
		return
	}

	var page WikiPage
	err = m.Collection.FindOne(context.Background(), bson.M{"_id": objID}).Decode(&page)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Страница не найдена"})
		return
	}

	c.JSON(http.StatusOK, page)
}

// Создание новой страницы
func (m *WikiModule) CreatePage(c *gin.Context) {
	var input WikiPage
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input.ID = primitive.NewObjectID()
	input.CreatedAt = time.Now()
	input.UpdatedAt = time.Now()

	_, err := m.Collection.InsertOne(context.Background(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, input)
}

// Обновление страницы
func (m *WikiModule) UpdatePage(c *gin.Context) {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный ID"})
		return
	}

	var body struct {
		Title   string `json:"title"`
		Content string `json:"content"`
		Emoji   string `json:"emoji"`
	}
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	update := bson.M{
		"$set": bson.M{
			"title":     body.Title,
			"content":   body.Content,
			"emoji":     body.Emoji,
			"updatedAt": time.Now(),
		},
	}

	res, err := m.Collection.UpdateOne(context.Background(), bson.M{"_id": objID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if res.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Страница не найдена"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Страница обновлена"})
}

// Удаление страницы
func (m *WikiModule) DeletePage(c *gin.Context) {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный ID"})
		return
	}

	// Удаляем страницу
	_, err = m.Collection.DeleteOne(context.Background(), bson.M{"_id": objID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Удаляем дочерние
	_, _ = m.Collection.DeleteMany(context.Background(), bson.M{"parentId": objID})

	c.JSON(http.StatusOK, gin.H{"message": "Страница и дочерние удалены"})
}
