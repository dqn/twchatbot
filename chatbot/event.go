package chatbot

type userID string

type Event struct {
	ForUserID           string               `json:"for_user_id"`
	DirectMessageEvents []DirectMessageEvent `json:"direct_message_events"`
	Users               map[userID]User      `json:"users"`
}

type Target struct {
	RecipientID string `json:"recipient_id"`
}

type Entities struct {
	Hashtags     []string `json:"hashtags"`
	Symbols      []string `json:"symbols"`
	UserMentions []string `json:"user_mentions"`
	URLs         []string `json:"urls"`
}

type MessageData struct {
	Text     string   `json:"text"`
	Entities Entities `json:"entities"`
}

type MessageCreate struct {
	Target      Target      `json:"target"`
	SenderID    string      `json:"sender_id"`
	MessageData MessageData `json:"message_data"`
}

type DirectMessageEvent struct {
	Type             string        `json:"type"`
	ID               string        `json:"id"`
	CreatedTimestamp string        `json:"created_timestamp"`
	MessageCreate    MessageCreate `json:"message_create"`
}

type User struct {
	ID                   string `json:"id"`
	CreatedTimestamp     string `json:"created_timestamp"`
	Name                 string `json:"name"`
	ScreenName           string `json:"screen_name"`
	Description          string `json:"description"`
	Protected            bool   `json:"protected"`
	Verified             bool   `json:"verified"`
	FollowersCount       int    `json:"followers_count"`
	FriendsCount         int    `json:"friends_count"`
	StatusesCount        int    `json:"statuses_count"`
	ProfileImageURL      string `json:"profile_image_url"`
	ProfileImageURLHTTPS string `json:"profile_image_url_https"`
}
