package example

// Status представляет собой структуру для хранения данных о статусе твита
type Status struct {
	Coordinates      interface{} `json:"coordinates"`
	Favorited        bool        `json:"favorited"`
	Truncated        bool        `json:"truncated"`
	CreatedAt        string      `json:"created_at"`
	IDStr            string      `json:"id_str"`
	Entities         Entities    `json:"entities"`
	InReplyToUserID  interface{} `json:"in_reply_to_user_id"`
	Contributors     interface{} `json:"contributors"`
	Text             string      `json:"text"`
	Metadata         Metadata    `json:"metadata"`
	RetweetCount     int         `json:"retweet_count"`
	InReplyToStatus  interface{} `json:"in_reply_to_status_id"`
	ID               int64       `json:"id"`
	Geo              interface{} `json:"geo"`
	Retweeted        bool        `json:"retweeted"`
	InReplyToUser    interface{} `json:"in_reply_to_user_id_str"`
	Place            interface{} `json:"place"`
	User             User        `json:"user"`
	InReplyToScreen  interface{} `json:"in_reply_to_screen_name"`
	Source           string      `json:"source"`
	InReplyToStatusI interface{} `json:"in_reply_to_status_id_str"`
}

// Entities представляет собой структуру для хранения данных о сущностях, упомянутых в твите
type Entities struct {
	URLs         []interface{} `json:"urls"`
	Hashtags     []Hashtag     `json:"hashtags"`
	UserMentions []interface{} `json:"user_mentions"`
}

// Hashtag представляет собой структуру для хранения данных о хештеге
type Hashtag struct {
	Text    string `json:"text"`
	Indices []int  `json:"indices"`
}

// Metadata представляет собой структуру для хранения дополнительной информации о твите
type Metadata struct {
	ISOLanguageCode string `json:"iso_language_code"`
	ResultType      string `json:"result_type"`
}

// User представляет собой структуру для хранения данных о пользователе, который опубликовал твит
type User struct {
	ProfileSidebarFillColor   string      `json:"profile_sidebar_fill_color"`
	ProfileSidebarBorderColor string      `json:"profile_sidebar_border_color"`
	ProfileBackgroundTile     bool        `json:"profile_background_tile"`
	Name                      string      `json:"name"`
	ProfileImageURL           string      `json:"profile_image_url"`
	CreatedAt                 string      `json:"created_at"`
	Location                  string      `json:"location"`
	FollowRequestSent         interface{} `json:"follow_request_sent"`
	ProfileLinkColor          string      `json:"profile_link_color"`
	IsTranslator              bool        `json:"is_translator"`
	IDStr                     string      `json:"id_str"`
	Entities                  struct {
		URL struct {
			URLs []interface{} `json:"urls"`
		} `json:"url"`
		Description struct {
			URLs []interface{} `json:"urls"`
		} `json:"description"`
	} `json:"entities"`
	DefaultProfile                 bool        `json:"default_profile"`
	ContributorsEnabled            bool        `json:"contributors_enabled"`
	FavouritesCount                int         `json:"favourites_count"`
	URL                            interface{} `json:"url"`
	ProfileImageURLHTTPS           string      `json:"profile_image_url_https"`
	UTCOffset                      int         `json:"utc_offset"`
	ID                             int64       `json:"id"`
	ProfileUseBackgroundImage      bool        `json:"profile_use_background_image"`
	ListedCount                    int         `json:"listed_count"`
	ProfileTextColor               string      `json:"profile_text_color"`
	Lang                           string      `json:"lang"`
	FollowersCount                 int         `json:"followers_count"`
	Protected                      bool        `json:"protected"`
	Notifications                  interface{} `json:"notifications"`
	ProfileBackgroundImageURLHTTPS string      `json:"profile_background_image_url_https"`
	ProfileBackgroundColor         string      `json:"profile_background_color"`
	Verified                       bool        `json:"verified"`
	GeoEnabled                     bool        `json:"geo_enabled"`
	TimeZone                       string      `json:"time_zone"`
	Description                    string      `json:"description"`
	DefaultProfileImage            bool        `json:"default_profile_image"`
	ProfileBackgroundImageURL      string      `json:"profile_background_image_url"`
	StatusesCount                  int         `json:"statuses_count"`
	FriendsCount                   int         `json:"friends_count"`
	Following                      interface{} `json:"following"`
	ShowAllInlineMedia             bool        `json:"show_all_inline_media"`
	ScreenName                     string      `json:"screen_name"`
}

// SearchMetadata представляет собой структуру для хранения данных о поиске
type SearchMetadata struct {
	MaxID       int64   `json:"max_id"`
	SinceID     int64   `json:"since_id"`
	RefreshURL  string  `json:"refresh_url"`
	NextResults string  `json:"next_results"`
	Count       int     `json:"count"`
	CompletedIn float64 `json:"completed_in"`
	SinceIDStr  string  `json:"since_id_str"`
	Query       string  `json:"query"`
	MaxIDStr    string  `json:"max_id_str"`
}

// Response представляет собой структуру для хранения ответа от сервера с результатами поиска
type Response struct {
	Statuses       []Status       `json:"statuses"`
	SearchMetadata SearchMetadata `json:"search_metadata"`
}
