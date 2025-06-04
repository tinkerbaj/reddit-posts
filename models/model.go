package models

import "gorm.io/gorm"


//Models for subreddit, flairs and post (gorm actually save me like 100 extra lines of code)
type Subreddit struct {
    gorm.Model
    Name   string  `json:"name"`
    Flairs []Flair `json:"flairs" gorm:"foreignKey:SubredditID"`
    Posts  []Post  `json:"posts" gorm:"foreignKey:SubredditID"`
}

type Flair struct {
    gorm.Model
    Name        string     `json:"name"`
    SubredditID int        `json:"subreddit_id"`
    Subreddit   *Subreddit `json:"subreddit" gorm:"foreignKey:SubredditID"`
    Posts       []Post     `json:"posts" gorm:"foreignKey:FlairID"`
}

type Post struct {
    gorm.Model
    Title       string     `json:"title"`
    Author      string     `json:"author"`
    CreatedUTC  float64    `json:"created_utc"`
    Permalink   string     `json:"permalink"`
    RedditID    string     `json:"reddit_id"`
    SubredditID int        `json:"subreddit_id"`
    Subreddit   *Subreddit `json:"subreddit" gorm:"foreignKey:SubredditID"`
    FlairID     *int       `json:"flair_id"` // Pointer to allow NULL
    Flair       *Flair     `json:"flair" gorm:"foreignKey:FlairID"`
}
