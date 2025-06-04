package requests

//Form request for adding new subreddit
type SubredditPostRequest struct {
	Name   string `json:"name" form:"name"`
	Flairs string `json:"flairs" form:"flairs"`
}

//Edit form for subreddit and flairs for it
type SubredditEditPostRequest struct {
	Name       string   `form:"name"`
	FlairIDs   []uint   `form:"flair_id[]"`   // Existing flair IDs
	FlairNames []string `form:"flairs[]"`     // Existing flair names
	NewFlairs  string   `form:"newflairs"`     // Comma-separated new flairs
}