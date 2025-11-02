package kanban

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// === Models ===
// Простые структуры: ID — string (unix nano), чтобы фронт мог работать с ними легко.
type Board struct {
	ID        string    `bson:"_id,omitempty" json:"_id"`
	Name      string    `bson:"name" json:"name"`
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt" json:"updatedAt"`
}

type Status struct {
	ID        string    `bson:"_id,omitempty" json:"_id"`
	Name      string    `bson:"name" json:"name"`
	BoardID   string    `bson:"boardId" json:"boardId"`
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt" json:"updatedAt"`
}

type Card struct {
	ID          string    `bson:"_id,omitempty" json:"_id"`
	BoardID     string    `bson:"boardId" json:"boardId"`
	StatusID    string    `bson:"statusId" json:"statusId"`
	Title       string    `bson:"title" json:"title"`
	Description string    `bson:"description" json:"description"`
	Color       string    `bson:"color" json:"color"`
	Image       string    `bson:"image,omitempty" json:"image,omitempty"` // "/uploads/..."
	Tags        []string  `bson:"tags,omitempty" json:"tags,omitempty"`
	DueDate     string    `bson:"dueDate,omitempty" json:"dueDate,omitempty"` // RFC3339 string or empty
	CreatedAt   time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time `bson:"updatedAt" json:"updatedAt"`
}

// RegisterKanbanRoutes registers all endpoints under /api/*
// router: *gin.Engine, db: *mongo.Database
func RegisterKanbanRoutes(router *gin.Engine, db *mongo.Database) {
	boardColl := db.Collection("boards")
	statusColl := db.Collection("statuses")
	cardColl := db.Collection("cards")

	// Static uploads
	uploadDir := "./uploads"
	_ = os.MkdirAll(uploadDir, os.ModePerm)
	router.Static("/uploads", uploadDir)

	// Simple CORS middleware (safe to include; if you already have CORS, this is harmless)
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})

	// --- BOARDS ---
	router.GET("/api/boards", func(c *gin.Context) {
		cur, err := boardColl.Find(context.TODO(), bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		var out []Board
		if err := cur.All(context.TODO(), &out); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, out)
	})

	router.POST("/api/boards", func(c *gin.Context) {
		var b Board
		if err := c.BindJSON(&b); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		now := time.Now()
		b.ID = fmt.Sprintf("%d", now.UnixNano())
		b.CreatedAt = now
		b.UpdatedAt = now
		if _, err := boardColl.InsertOne(context.TODO(), b); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, b)
	})

	router.DELETE("/api/boards/:id", func(c *gin.Context) {
		id := c.Param("id")
		// optionally: delete statuses/cards referencing this board
		if _, err := boardColl.DeleteOne(context.TODO(), bson.M{"_id": id}); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		// cascade: delete statuses and cards for that board
		_, _ = statusColl.DeleteMany(context.TODO(), bson.M{"boardId": id})
		_, _ = cardColl.DeleteMany(context.TODO(), bson.M{"boardId": id})
		c.JSON(http.StatusOK, gin.H{"status": "deleted"})
	})

	// --- STATUSES (columns) ---
	router.GET("/api/statuses", func(c *gin.Context) {
		boardId := c.Query("boardId")
		filter := bson.M{}
		if boardId != "" {
			filter["boardId"] = boardId
		}
		cur, err := statusColl.Find(context.TODO(), filter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		var out []Status
		if err := cur.All(context.TODO(), &out); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, out)
	})

	router.POST("/api/statuses", func(c *gin.Context) {
		var s Status
		if err := c.BindJSON(&s); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		now := time.Now()
		s.ID = fmt.Sprintf("%d", now.UnixNano())
		s.CreatedAt = now
		s.UpdatedAt = now
		if _, err := statusColl.InsertOne(context.TODO(), s); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, s)
	})

	router.DELETE("/api/statuses/:id", func(c *gin.Context) {
		id := c.Param("id")
		if _, err := statusColl.DeleteOne(context.TODO(), bson.M{"_id": id}); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		// optional: set cards with this status to empty statusId or delete them — here we set statusId=""
		_, _ = cardColl.UpdateMany(context.TODO(), bson.M{"statusId": id}, bson.M{"$set": bson.M{"statusId": ""}})
		c.JSON(http.StatusOK, gin.H{"status": "deleted"})
	})

	// --- CARDS ---
	router.GET("/api/cards", func(c *gin.Context) {
		boardId := c.Query("boardId")
		tag := c.Query("tag") // optional
		filter := bson.M{}
		if boardId != "" {
			filter["boardId"] = boardId
		}
		if tag != "" {
			// tag stored without leading '#', but we'll match either way
			filter["tags"] = tag
		}
		cur, err := cardColl.Find(context.TODO(), filter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		var out []Card
		if err := cur.All(context.TODO(), &out); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, out)
	})

	router.POST("/api/cards", func(c *gin.Context) {
		var card Card
		if err := c.BindJSON(&card); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		now := time.Now()
		card.ID = fmt.Sprintf("%d", now.UnixNano())
		card.CreatedAt = now
		card.UpdatedAt = now
		if _, err := cardColl.InsertOne(context.TODO(), card); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, card)
	})

	router.PUT("/api/cards/:id", func(c *gin.Context) {
		id := c.Param("id")
		var update Card
		if err := c.BindJSON(&update); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		update.UpdatedAt = time.Now()
		// Don't allow changing ID
		_, err := cardColl.UpdateOne(context.TODO(), bson.M{"_id": id}, bson.M{"$set": update})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		// Return updated document
		var out Card
		_ = cardColl.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&out)
		c.JSON(http.StatusOK, out)
	})

	router.DELETE("/api/cards/:id", func(c *gin.Context) {
		id := c.Param("id")
		res, err := cardColl.DeleteOne(context.TODO(), bson.M{"_id": id})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if res.DeletedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "deleted"})
	})

	// --- UPLOAD ---
	router.POST("/api/upload", func(c *gin.Context) {
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// sanitize filename (use only base)
		filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), filepath.Base(file.Filename))
		savePath := filepath.Join(uploadDir, filename)
		if err := c.SaveUploadedFile(file, savePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		// return relative path that frontend can prefix with server origin
		c.JSON(http.StatusOK, gin.H{"url": "/uploads/" + filename})
	})

	// --- CLEANUP ---
	router.GET("/api/cleanup", func(c *gin.Context) {
		removed := cleanupUploads(db)
		c.JSON(http.StatusOK, gin.H{"removed": removed})
	})

	// background cleanup ticker (non-blocking)
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()
		for {
			<-ticker.C
			removed := cleanupUploads(db)
			if removed > 0 {
				log.Println("kanban cleanup removed files:", removed)
			}
		}
	}()
}

// cleanupUploads deletes files in ./uploads that are not referenced by any card.image
func cleanupUploads(db *mongo.Database) int {
	uploadDir := "./uploads"
	files, err := os.ReadDir(uploadDir)
	if err != nil {
		return 0
	}
	cardColl := db.Collection("cards")
	removed := 0
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		path := "/uploads/" + f.Name()
		cnt, err := cardColl.CountDocuments(context.TODO(), bson.M{"image": path})
		if err != nil {
			continue
		}
		if cnt == 0 {
			_ = os.Remove(filepath.Join(uploadDir, f.Name()))
			removed++
		}
	}
	return removed
}
