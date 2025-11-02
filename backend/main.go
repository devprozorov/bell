package main

import (
	kanban "bell-backend/board"
	"bell-backend/config"
	"bell-backend/routes"
	wikimodule "bell-backend/wiki"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Инициализация базы данных
	config.ConnectDB()

	// Создаём роутер
	r := gin.Default()
	// Разрешаем CORS для фронтенда (Nuxt на порту 3000)
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// === Роуты ===
	routes.UserRoutes(r)
	routes.RegisterAdminRoutes(r)
	routes.AuthRoutes(r)
	// === Роут kanban ===
	kanban.RegisterKanbanRoutes(r, config.ConnectDB())
	// === Роут вики ===
	wiki := wikimodule.NewWikiModule(config.DB.Client())
	wiki.RegisterRoutes(r)
	// Запуск сервера
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Ошибка запуска сервера:", err)
	}

	r.Use(func(c *gin.Context) {
		c.Next()
		for _, err := range c.Errors {
			log.Println("GIN ERROR:", err.Err)
		}
	})

}
