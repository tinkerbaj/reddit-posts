package reddit

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
	"time"

	"github.com/tinkerbaj/redditjobs/db"
	"github.com/tinkerbaj/redditjobs/models"
	"github.com/tinkerbaj/redditjobs/utils"
)

//Here is logic for fetching posts from reddit
const (
	userAgent = "Mozilla/5.0 (X11; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/111.0"
)

//RedditPost and RedditResponse I generated online based on json response I get from the reddit
type RedditPost struct {
	Data struct {
		Title     string  `json:"title"`
		Author    string  `json:"author"`
		Created   float64 `json:"created_utc"`
		Permalink string  `json:"permalink"`
		ID        string  `json:"id"`
	} `json:"data"`
}

type RedditResponse struct {
	Data struct {
		Children []RedditPost `json:"children"`
	} `json:"data"`
}

func CheckForNewPosts(subreddit, flair string) ([]RedditPost, error) {
	log.Println("Checking for new posts...")

	var url string
	//here I need only first url but I was testing some stuff I leave url for reference 😅
	if flair != "" {
		url = fmt.Sprintf("https://www.reddit.com/r/%s/search.json?q=flair:\"%s\"&restrict_sr=1&sort=new", subreddit, flair)
	} else {
		url = fmt.Sprintf("https://www.reddit.com/r/%s/new.json?sort=new", subreddit)
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		
		utils.ColorError("Error creating request: ", err)
		return nil, err
	}

	req.Header.Set("User-Agent", userAgent)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		utils.ColorError("Error making request:", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Received non-OK status: %d\n", resp.StatusCode)
		return nil, fmt.Errorf("status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		utils.ColorError("Error reading response body: ", err)
		return nil, err
	}

	var redditResponse RedditResponse
	err = json.Unmarshal(body, &redditResponse)
	if err != nil {
		utils.ColorError("Error unmarshaling JSON: ", err)
		return nil, err
	}

	if len(redditResponse.Data.Children) == 0 {
		utils.ColorError("No posts found")
		return nil, nil
	}

	// Sort posts by creation time (newest first)
	sort.Slice(redditResponse.Data.Children, func(i, j int) bool {
		return redditResponse.Data.Children[i].Data.Created > redditResponse.Data.Children[j].Data.Created
	})
	subredditID, _ := getSubredditID(subreddit)
	flairID, _ := getFlairID(flair, subredditID) // or nil
	_ = saveNewPostsToDB(redditResponse.Data.Children, subredditID, flairID)

	// Return all posts as new for now
	return redditResponse.Data.Children, nil
}


func getSubredditID(name string) (int64, error) {
	var subreddit models.Subreddit
	db.DB.Where("name = ?", name).Find(&subreddit)
	return int64(subreddit.ID), nil
}

func getFlairID(name string, subredditID int64) (int64, error) {
	var flair models.Flair
	db.DB.Where("name = ? AND subreddit_id = ?", name, subredditID).Find(&flair)
	return int64(flair.ID), nil
}

func GetFilteredPosts(subredditName, flairName string) ([]models.Post, error) {

	var posts []models.Post
		subredditID, _ := getSubredditID(subredditName)
	flairID, _ := getFlairID(flairName, subredditID)
	db.DB.Where("subreddit_id = ? AND flair_id = ?", subredditID, flairID).Preload("Subreddit").Preload("Flair").Find(&posts)

	return posts, nil
}

//I started with RedditPost struct later I added models than I feel lazy to convert rest of the code
func saveNewPostsToDB(posts []RedditPost, subredditID int64, flairID int64) error {

	for _, post := range posts {

		var newpost models.Post
		newpost.Author = post.Data.Author
		newpost.Title = post.Data.Title
		newpost.Permalink = post.Data.Permalink
		k := int64(post.Data.Created)
		newpost.CreatedAt = time.Unix(k, 0)
		i := int(flairID)
		newpost.FlairID = &i
		newpost.SubredditID = int(subredditID)
		db.DB.Save(&newpost)
	}
	return nil
}

