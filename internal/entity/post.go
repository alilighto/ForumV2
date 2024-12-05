package entity

type Post struct {
	PostID   uint      `json:"post_id"`
	UserID   uint      `json:"user_id"`
	UserName string    `json:"username"`
	Tags     []string  `json:"tags"`
	Title    string    `json:"title"`
	Data     string    `json:"data"`
	Likes    uint      `json:"likes"`
	Dislikes uint      `json:"dislikes"`
	Comments []Comment `json:"comments"`
}

type Tag struct {
	TagID uint
	Name  string
}

type TagAndPost struct {
	TagID  uint
	PostID uint
}

type PostVote struct {
	UserID uint `json:"user_id"`
	PostID uint `json:"post_id"`
	Vote   int  `json:"vote"`
}
