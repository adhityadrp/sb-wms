package routes

import (
	"sb-wms/controllers"
	"sb-wms/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	api := r.Group("/api")
	{
		api.POST("/users/register", controllers.Register)
		api.POST("/users/login", controllers.Login)

		// protected routes (JWT)
		protected := api.Group("/")
		protected.Use(middlewares.JWTAuthMiddleware())
		{
			// categories
			protected.GET("/categories", controllers.GetCategories)
			protected.POST("/categories", controllers.CreateCategory)
			protected.GET("/categories/:id", controllers.GetCategoryByID)
			protected.PUT("/categories/:id", controllers.UpdateCategory)
			protected.DELETE("/categories/:id", controllers.DeleteCategory)

			// items
			protected.GET("/items", controllers.GetItems)
			protected.POST("/items", controllers.CreateItem)
			protected.GET("/items/:id", controllers.GetItemByID)
			protected.PUT("/items/:id", controllers.UpdateItem)
			protected.DELETE("/items/:id", controllers.DeleteItem)

			// transactions
			protected.GET("/transactions", controllers.GetTransactions)
			protected.POST("/transactions", controllers.CreateTransaction)
		}

	}

	return r
}
