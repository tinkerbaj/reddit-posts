package handlers

import (
	"net/http"

	// "os/signal"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tinkerbaj/redditjobs/db"
	"github.com/tinkerbaj/redditjobs/models"
	"github.com/tinkerbaj/redditjobs/reddit"
	"github.com/tinkerbaj/redditjobs/requests"
	"github.com/tinkerbaj/redditjobs/templates"
	"github.com/tinkerbaj/redditjobs/templates/components"

	// "github.com/tinkerbaj/redditjobs/templates/components"
	"github.com/tinkerbaj/redditjobs/utils"

	datastar "github.com/starfederation/datastar/sdk/go"
)

// Showing all posts
func Home(c *gin.Context) {
	var posts []models.Post
	db.DB.Preload("Subreddit").Preload("Flair").Order("created_at DESC").Find(&posts)
	hom := templates.Home(posts)
	lay := templates.Layout(hom)
	utils.Render(lay, c)
}

// Checking for new posts it is on click to prevent ban from reddit
func Check(c *gin.Context) {
	subname := c.Param("sub")
	flair := c.Param("flair")
	reddit.CheckForNewPosts(subname, flair)
}

// Showing subreddit / flair page
func Subreddit(c *gin.Context) {
	subname := c.Param("name")
	flair := c.Param("flair")

	posts, err := reddit.GetFilteredPosts(subname, flair)
	utils.ColorError("Posts not found", err)
	sub := templates.SubRedditPage(subname, flair, posts)
	lay := templates.Layout(sub)
	utils.Render(lay, c)
}

// Sidebar with link to the subreddits and there flairs
func Sidebar(c *gin.Context) {
	sse := datastar.NewSSE(c.Writer, c.Request)

	var subreddits []models.Subreddit
	db.DB.Preload("Flairs").Find(&subreddits)

	sse.MergeFragmentTempl(components.Sidebar(subreddits), datastar.WithMergeAppend(), datastar.WithSelectorID("sidebar"))
}


//Form add subreddit endpoint
func AddSubreddit(c *gin.Context) {
	var subPostReq requests.SubredditPostRequest
	if err := c.ShouldBind(&subPostReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data" + err.Error()})
		return
	}

	var subreddit models.Subreddit

	subreddit.Name = subPostReq.Name
	err := db.DB.Save(&subreddit).Error
	utils.ColorError("Cant create subreddit", err)

	//Because flairs is separated by , split make array and insert them
	flairNames := strings.Split(subPostReq.Flairs, ",")

	for _, flairName := range flairNames {
		flairName = strings.TrimSpace(flairName)
		if flairName == "" {
			continue
		}

		var flair models.Flair

		flair.Name = flairName
		flair.SubredditID = int(subreddit.ID)

		err = db.DB.Save(&flair).Error
		utils.ColorError("Cant create flair", err)
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":    "Subreddit and flairs added successfully",
		"subreddit":  subPostReq.Name,
		"flairCount": len(flairNames),
	})
}


func EditSubredditPost(c *gin.Context) {
	id := c.Param("id")
	var form requests.SubredditEditPostRequest
	if err := c.Bind(&form); err != nil {
		c.String(http.StatusBadRequest, "Invalid form data")
		return
	}

	var subreddit models.Subreddit
	if err := db.DB.First(&subreddit, id).Error; err != nil {
		c.String(http.StatusNotFound, "Subreddit not found")
		return
	}
	subreddit.Name = form.Name
	db.DB.Save(&subreddit)

	// Update existing flairs
	for i := range form.FlairIDs {
		if i < len(form.FlairNames) {
			var flair models.Flair
			if err := db.DB.First(&flair, form.FlairIDs[i]).Error; err == nil {
				flair.Name = form.FlairNames[i]
				db.DB.Save(&flair)
			}
		}
	}

	// Add new flairs
	for _, name := range strings.Split(form.NewFlairs, ",") {
		name = strings.TrimSpace(name)
		if name != "" {
			db.DB.Create(&models.Flair{Name: name, SubredditID: int(subreddit.ID)})
		}
	}


}