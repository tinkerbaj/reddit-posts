package router

import (
	"github.com/gin-gonic/gin"
	"github.com/tinkerbaj/redditjobs/handlers"

)

//Separated router just to main.go looks cleaner
func SetRouter(app *gin.Engine) {


	app.GET("/", handlers.Home)
	app.POST("/add-subreddit", handlers.AddSubreddit)
	app.PATCH("/edit-subreddit/:id", handlers.EditSubredditPost)
	app.GET("/sidebar", handlers.Sidebar)
	app.GET("/sub/:name/:flair", handlers.Subreddit)
	app.GET("/check/:sub/:flair", handlers.Check)


}
