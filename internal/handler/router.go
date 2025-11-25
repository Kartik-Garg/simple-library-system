package handler

import "github.com/gin-gonic/gin"

// is this controller layer?
// setting up the gin server

func SetUpRouter(BookHandler *BookHandler) *gin.Engine {
	router := gin.Default()

	// adding health check
	// how does this health endpoint work? and why we need it?
	router.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"status":  "healthy",
			"service": "library-api",
		})
	})

	// API routes. Adding versioning (making more extensible)

	// why do we add versioning? explain in detail

	// what does group do?
	v1 := router.Group("/api/v1")
	{
		// what does this Group does
		books := v1.Group("/books")
		{
			books.GET("", BookHandler.GetAllBooks)
			books.GET("/:id", BookHandler.GetBook)
			books.POST("", BookHandler.CreateBook)
			books.PUT("/:id", BookHandler.UpdateBook)
			books.DELETE("/:id", BookHandler.DeleteBook)
		}
	}
	return router
}
