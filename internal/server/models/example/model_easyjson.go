// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package example

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjsonC80ae7adDecodeBanchModel(in *jlexer.Lexer, out *User) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "profile_sidebar_fill_color":
			out.ProfileSidebarFillColor = string(in.String())
		case "profile_sidebar_border_color":
			out.ProfileSidebarBorderColor = string(in.String())
		case "profile_background_tile":
			out.ProfileBackgroundTile = bool(in.Bool())
		case "name":
			out.Name = string(in.String())
		case "profile_image_url":
			out.ProfileImageURL = string(in.String())
		case "created_at":
			out.CreatedAt = string(in.String())
		case "location":
			out.Location = string(in.String())
		case "follow_request_sent":
			if m, ok := out.FollowRequestSent.(easyjson.Unmarshaler); ok {
				m.UnmarshalEasyJSON(in)
			} else if m, ok := out.FollowRequestSent.(json.Unmarshaler); ok {
				_ = m.UnmarshalJSON(in.Raw())
			} else {
				out.FollowRequestSent = in.Interface()
			}
		case "profile_link_color":
			out.ProfileLinkColor = string(in.String())
		case "is_translator":
			out.IsTranslator = bool(in.Bool())
		case "id_str":
			out.IDStr = string(in.String())
		case "entities":
			easyjsonC80ae7adDecode(in, &out.Entities)
		case "default_profile":
			out.DefaultProfile = bool(in.Bool())
		case "contributors_enabled":
			out.ContributorsEnabled = bool(in.Bool())
		case "favourites_count":
			out.FavouritesCount = int(in.Int())
		case "url":
			if m, ok := out.URL.(easyjson.Unmarshaler); ok {
				m.UnmarshalEasyJSON(in)
			} else if m, ok := out.URL.(json.Unmarshaler); ok {
				_ = m.UnmarshalJSON(in.Raw())
			} else {
				out.URL = in.Interface()
			}
		case "profile_image_url_https":
			out.ProfileImageURLHTTPS = string(in.String())
		case "utc_offset":
			out.UTCOffset = int(in.Int())
		case "id":
			out.ID = int64(in.Int64())
		case "profile_use_background_image":
			out.ProfileUseBackgroundImage = bool(in.Bool())
		case "listed_count":
			out.ListedCount = int(in.Int())
		case "profile_text_color":
			out.ProfileTextColor = string(in.String())
		case "lang":
			out.Lang = string(in.String())
		case "followers_count":
			out.FollowersCount = int(in.Int())
		case "protected":
			out.Protected = bool(in.Bool())
		case "notifications":
			if m, ok := out.Notifications.(easyjson.Unmarshaler); ok {
				m.UnmarshalEasyJSON(in)
			} else if m, ok := out.Notifications.(json.Unmarshaler); ok {
				_ = m.UnmarshalJSON(in.Raw())
			} else {
				out.Notifications = in.Interface()
			}
		case "profile_background_image_url_https":
			out.ProfileBackgroundImageURLHTTPS = string(in.String())
		case "profile_background_color":
			out.ProfileBackgroundColor = string(in.String())
		case "verified":
			out.Verified = bool(in.Bool())
		case "geo_enabled":
			out.GeoEnabled = bool(in.Bool())
		case "time_zone":
			out.TimeZone = string(in.String())
		case "description":
			out.Description = string(in.String())
		case "default_profile_image":
			out.DefaultProfileImage = bool(in.Bool())
		case "profile_background_image_url":
			out.ProfileBackgroundImageURL = string(in.String())
		case "statuses_count":
			out.StatusesCount = int(in.Int())
		case "friends_count":
			out.FriendsCount = int(in.Int())
		case "following":
			if m, ok := out.Following.(easyjson.Unmarshaler); ok {
				m.UnmarshalEasyJSON(in)
			} else if m, ok := out.Following.(json.Unmarshaler); ok {
				_ = m.UnmarshalJSON(in.Raw())
			} else {
				out.Following = in.Interface()
			}
		case "show_all_inline_media":
			out.ShowAllInlineMedia = bool(in.Bool())
		case "screen_name":
			out.ScreenName = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonC80ae7adEncodeBanchModel(out *jwriter.Writer, in User) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"profile_sidebar_fill_color\":"
		out.RawString(prefix[1:])
		out.String(string(in.ProfileSidebarFillColor))
	}
	{
		const prefix string = ",\"profile_sidebar_border_color\":"
		out.RawString(prefix)
		out.String(string(in.ProfileSidebarBorderColor))
	}
	{
		const prefix string = ",\"profile_background_tile\":"
		out.RawString(prefix)
		out.Bool(bool(in.ProfileBackgroundTile))
	}
	{
		const prefix string = ",\"name\":"
		out.RawString(prefix)
		out.String(string(in.Name))
	}
	{
		const prefix string = ",\"profile_image_url\":"
		out.RawString(prefix)
		out.String(string(in.ProfileImageURL))
	}
	{
		const prefix string = ",\"created_at\":"
		out.RawString(prefix)
		out.String(string(in.CreatedAt))
	}
	{
		const prefix string = ",\"location\":"
		out.RawString(prefix)
		out.String(string(in.Location))
	}
	{
		const prefix string = ",\"follow_request_sent\":"
		out.RawString(prefix)
		if m, ok := in.FollowRequestSent.(easyjson.Marshaler); ok {
			m.MarshalEasyJSON(out)
		} else if m, ok := in.FollowRequestSent.(json.Marshaler); ok {
			out.Raw(m.MarshalJSON())
		} else {
			out.Raw(json.Marshal(in.FollowRequestSent))
		}
	}
	{
		const prefix string = ",\"profile_link_color\":"
		out.RawString(prefix)
		out.String(string(in.ProfileLinkColor))
	}
	{
		const prefix string = ",\"is_translator\":"
		out.RawString(prefix)
		out.Bool(bool(in.IsTranslator))
	}
	{
		const prefix string = ",\"id_str\":"
		out.RawString(prefix)
		out.String(string(in.IDStr))
	}
	{
		const prefix string = ",\"entities\":"
		out.RawString(prefix)
		easyjsonC80ae7adEncode(out, in.Entities)
	}
	{
		const prefix string = ",\"default_profile\":"
		out.RawString(prefix)
		out.Bool(bool(in.DefaultProfile))
	}
	{
		const prefix string = ",\"contributors_enabled\":"
		out.RawString(prefix)
		out.Bool(bool(in.ContributorsEnabled))
	}
	{
		const prefix string = ",\"favourites_count\":"
		out.RawString(prefix)
		out.Int(int(in.FavouritesCount))
	}
	{
		const prefix string = ",\"url\":"
		out.RawString(prefix)
		if m, ok := in.URL.(easyjson.Marshaler); ok {
			m.MarshalEasyJSON(out)
		} else if m, ok := in.URL.(json.Marshaler); ok {
			out.Raw(m.MarshalJSON())
		} else {
			out.Raw(json.Marshal(in.URL))
		}
	}
	{
		const prefix string = ",\"profile_image_url_https\":"
		out.RawString(prefix)
		out.String(string(in.ProfileImageURLHTTPS))
	}
	{
		const prefix string = ",\"utc_offset\":"
		out.RawString(prefix)
		out.Int(int(in.UTCOffset))
	}
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix)
		out.Int64(int64(in.ID))
	}
	{
		const prefix string = ",\"profile_use_background_image\":"
		out.RawString(prefix)
		out.Bool(bool(in.ProfileUseBackgroundImage))
	}
	{
		const prefix string = ",\"listed_count\":"
		out.RawString(prefix)
		out.Int(int(in.ListedCount))
	}
	{
		const prefix string = ",\"profile_text_color\":"
		out.RawString(prefix)
		out.String(string(in.ProfileTextColor))
	}
	{
		const prefix string = ",\"lang\":"
		out.RawString(prefix)
		out.String(string(in.Lang))
	}
	{
		const prefix string = ",\"followers_count\":"
		out.RawString(prefix)
		out.Int(int(in.FollowersCount))
	}
	{
		const prefix string = ",\"protected\":"
		out.RawString(prefix)
		out.Bool(bool(in.Protected))
	}
	{
		const prefix string = ",\"notifications\":"
		out.RawString(prefix)
		if m, ok := in.Notifications.(easyjson.Marshaler); ok {
			m.MarshalEasyJSON(out)
		} else if m, ok := in.Notifications.(json.Marshaler); ok {
			out.Raw(m.MarshalJSON())
		} else {
			out.Raw(json.Marshal(in.Notifications))
		}
	}
	{
		const prefix string = ",\"profile_background_image_url_https\":"
		out.RawString(prefix)
		out.String(string(in.ProfileBackgroundImageURLHTTPS))
	}
	{
		const prefix string = ",\"profile_background_color\":"
		out.RawString(prefix)
		out.String(string(in.ProfileBackgroundColor))
	}
	{
		const prefix string = ",\"verified\":"
		out.RawString(prefix)
		out.Bool(bool(in.Verified))
	}
	{
		const prefix string = ",\"geo_enabled\":"
		out.RawString(prefix)
		out.Bool(bool(in.GeoEnabled))
	}
	{
		const prefix string = ",\"time_zone\":"
		out.RawString(prefix)
		out.String(string(in.TimeZone))
	}
	{
		const prefix string = ",\"description\":"
		out.RawString(prefix)
		out.String(string(in.Description))
	}
	{
		const prefix string = ",\"default_profile_image\":"
		out.RawString(prefix)
		out.Bool(bool(in.DefaultProfileImage))
	}
	{
		const prefix string = ",\"profile_background_image_url\":"
		out.RawString(prefix)
		out.String(string(in.ProfileBackgroundImageURL))
	}
	{
		const prefix string = ",\"statuses_count\":"
		out.RawString(prefix)
		out.Int(int(in.StatusesCount))
	}
	{
		const prefix string = ",\"friends_count\":"
		out.RawString(prefix)
		out.Int(int(in.FriendsCount))
	}
	{
		const prefix string = ",\"following\":"
		out.RawString(prefix)
		if m, ok := in.Following.(easyjson.Marshaler); ok {
			m.MarshalEasyJSON(out)
		} else if m, ok := in.Following.(json.Marshaler); ok {
			out.Raw(m.MarshalJSON())
		} else {
			out.Raw(json.Marshal(in.Following))
		}
	}
	{
		const prefix string = ",\"show_all_inline_media\":"
		out.RawString(prefix)
		out.Bool(bool(in.ShowAllInlineMedia))
	}
	{
		const prefix string = ",\"screen_name\":"
		out.RawString(prefix)
		out.String(string(in.ScreenName))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v User) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonC80ae7adEncodeBanchModel(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v User) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonC80ae7adEncodeBanchModel(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *User) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonC80ae7adDecodeBanchModel(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *User) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonC80ae7adDecodeBanchModel(l, v)
}
func easyjsonC80ae7adDecode(in *jlexer.Lexer, out *struct {
	URL struct {
		URLs []interface{} `json:"urls"`
	} `json:"url"`
	Description struct {
		URLs []interface{} `json:"urls"`
	} `json:"description"`
}) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "url":
			easyjsonC80ae7adDecode1(in, &out.URL)
		case "description":
			easyjsonC80ae7adDecode1(in, &out.Description)
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonC80ae7adEncode(out *jwriter.Writer, in struct {
	URL struct {
		URLs []interface{} `json:"urls"`
	} `json:"url"`
	Description struct {
		URLs []interface{} `json:"urls"`
	} `json:"description"`
}) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"url\":"
		out.RawString(prefix[1:])
		easyjsonC80ae7adEncode1(out, in.URL)
	}
	{
		const prefix string = ",\"description\":"
		out.RawString(prefix)
		easyjsonC80ae7adEncode1(out, in.Description)
	}
	out.RawByte('}')
}
func easyjsonC80ae7adDecode1(in *jlexer.Lexer, out *struct {
	URLs []interface{} `json:"urls"`
}) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "urls":
			if in.IsNull() {
				in.Skip()
				out.URLs = nil
			} else {
				in.Delim('[')
				if out.URLs == nil {
					if !in.IsDelim(']') {
						out.URLs = make([]interface{}, 0, 4)
					} else {
						out.URLs = []interface{}{}
					}
				} else {
					out.URLs = (out.URLs)[:0]
				}
				for !in.IsDelim(']') {
					var v1 interface{}
					if m, ok := v1.(easyjson.Unmarshaler); ok {
						m.UnmarshalEasyJSON(in)
					} else if m, ok := v1.(json.Unmarshaler); ok {
						_ = m.UnmarshalJSON(in.Raw())
					} else {
						v1 = in.Interface()
					}
					out.URLs = append(out.URLs, v1)
					in.WantComma()
				}
				in.Delim(']')
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonC80ae7adEncode1(out *jwriter.Writer, in struct {
	URLs []interface{} `json:"urls"`
}) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"urls\":"
		out.RawString(prefix[1:])
		if in.URLs == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v2, v3 := range in.URLs {
				if v2 > 0 {
					out.RawByte(',')
				}
				if m, ok := v3.(easyjson.Marshaler); ok {
					m.MarshalEasyJSON(out)
				} else if m, ok := v3.(json.Marshaler); ok {
					out.Raw(m.MarshalJSON())
				} else {
					out.Raw(json.Marshal(v3))
				}
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}
func easyjsonC80ae7adDecodeBanchModel1(in *jlexer.Lexer, out *Status) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "coordinates":
			if m, ok := out.Coordinates.(easyjson.Unmarshaler); ok {
				m.UnmarshalEasyJSON(in)
			} else if m, ok := out.Coordinates.(json.Unmarshaler); ok {
				_ = m.UnmarshalJSON(in.Raw())
			} else {
				out.Coordinates = in.Interface()
			}
		case "favorited":
			out.Favorited = bool(in.Bool())
		case "truncated":
			out.Truncated = bool(in.Bool())
		case "created_at":
			out.CreatedAt = string(in.String())
		case "id_str":
			out.IDStr = string(in.String())
		case "entities":
			(out.Entities).UnmarshalEasyJSON(in)
		case "in_reply_to_user_id":
			if m, ok := out.InReplyToUserID.(easyjson.Unmarshaler); ok {
				m.UnmarshalEasyJSON(in)
			} else if m, ok := out.InReplyToUserID.(json.Unmarshaler); ok {
				_ = m.UnmarshalJSON(in.Raw())
			} else {
				out.InReplyToUserID = in.Interface()
			}
		case "contributors":
			if m, ok := out.Contributors.(easyjson.Unmarshaler); ok {
				m.UnmarshalEasyJSON(in)
			} else if m, ok := out.Contributors.(json.Unmarshaler); ok {
				_ = m.UnmarshalJSON(in.Raw())
			} else {
				out.Contributors = in.Interface()
			}
		case "text":
			out.Text = string(in.String())
		case "metadata":
			(out.Metadata).UnmarshalEasyJSON(in)
		case "retweet_count":
			out.RetweetCount = int(in.Int())
		case "in_reply_to_status_id":
			if m, ok := out.InReplyToStatus.(easyjson.Unmarshaler); ok {
				m.UnmarshalEasyJSON(in)
			} else if m, ok := out.InReplyToStatus.(json.Unmarshaler); ok {
				_ = m.UnmarshalJSON(in.Raw())
			} else {
				out.InReplyToStatus = in.Interface()
			}
		case "id":
			out.ID = int64(in.Int64())
		case "geo":
			if m, ok := out.Geo.(easyjson.Unmarshaler); ok {
				m.UnmarshalEasyJSON(in)
			} else if m, ok := out.Geo.(json.Unmarshaler); ok {
				_ = m.UnmarshalJSON(in.Raw())
			} else {
				out.Geo = in.Interface()
			}
		case "retweeted":
			out.Retweeted = bool(in.Bool())
		case "in_reply_to_user_id_str":
			if m, ok := out.InReplyToUser.(easyjson.Unmarshaler); ok {
				m.UnmarshalEasyJSON(in)
			} else if m, ok := out.InReplyToUser.(json.Unmarshaler); ok {
				_ = m.UnmarshalJSON(in.Raw())
			} else {
				out.InReplyToUser = in.Interface()
			}
		case "place":
			if m, ok := out.Place.(easyjson.Unmarshaler); ok {
				m.UnmarshalEasyJSON(in)
			} else if m, ok := out.Place.(json.Unmarshaler); ok {
				_ = m.UnmarshalJSON(in.Raw())
			} else {
				out.Place = in.Interface()
			}
		case "user":
			(out.User).UnmarshalEasyJSON(in)
		case "in_reply_to_screen_name":
			if m, ok := out.InReplyToScreen.(easyjson.Unmarshaler); ok {
				m.UnmarshalEasyJSON(in)
			} else if m, ok := out.InReplyToScreen.(json.Unmarshaler); ok {
				_ = m.UnmarshalJSON(in.Raw())
			} else {
				out.InReplyToScreen = in.Interface()
			}
		case "source":
			out.Source = string(in.String())
		case "in_reply_to_status_id_str":
			if m, ok := out.InReplyToStatusI.(easyjson.Unmarshaler); ok {
				m.UnmarshalEasyJSON(in)
			} else if m, ok := out.InReplyToStatusI.(json.Unmarshaler); ok {
				_ = m.UnmarshalJSON(in.Raw())
			} else {
				out.InReplyToStatusI = in.Interface()
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonC80ae7adEncodeBanchModel1(out *jwriter.Writer, in Status) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"coordinates\":"
		out.RawString(prefix[1:])
		if m, ok := in.Coordinates.(easyjson.Marshaler); ok {
			m.MarshalEasyJSON(out)
		} else if m, ok := in.Coordinates.(json.Marshaler); ok {
			out.Raw(m.MarshalJSON())
		} else {
			out.Raw(json.Marshal(in.Coordinates))
		}
	}
	{
		const prefix string = ",\"favorited\":"
		out.RawString(prefix)
		out.Bool(bool(in.Favorited))
	}
	{
		const prefix string = ",\"truncated\":"
		out.RawString(prefix)
		out.Bool(bool(in.Truncated))
	}
	{
		const prefix string = ",\"created_at\":"
		out.RawString(prefix)
		out.String(string(in.CreatedAt))
	}
	{
		const prefix string = ",\"id_str\":"
		out.RawString(prefix)
		out.String(string(in.IDStr))
	}
	{
		const prefix string = ",\"entities\":"
		out.RawString(prefix)
		(in.Entities).MarshalEasyJSON(out)
	}
	{
		const prefix string = ",\"in_reply_to_user_id\":"
		out.RawString(prefix)
		if m, ok := in.InReplyToUserID.(easyjson.Marshaler); ok {
			m.MarshalEasyJSON(out)
		} else if m, ok := in.InReplyToUserID.(json.Marshaler); ok {
			out.Raw(m.MarshalJSON())
		} else {
			out.Raw(json.Marshal(in.InReplyToUserID))
		}
	}
	{
		const prefix string = ",\"contributors\":"
		out.RawString(prefix)
		if m, ok := in.Contributors.(easyjson.Marshaler); ok {
			m.MarshalEasyJSON(out)
		} else if m, ok := in.Contributors.(json.Marshaler); ok {
			out.Raw(m.MarshalJSON())
		} else {
			out.Raw(json.Marshal(in.Contributors))
		}
	}
	{
		const prefix string = ",\"text\":"
		out.RawString(prefix)
		out.String(string(in.Text))
	}
	{
		const prefix string = ",\"metadata\":"
		out.RawString(prefix)
		(in.Metadata).MarshalEasyJSON(out)
	}
	{
		const prefix string = ",\"retweet_count\":"
		out.RawString(prefix)
		out.Int(int(in.RetweetCount))
	}
	{
		const prefix string = ",\"in_reply_to_status_id\":"
		out.RawString(prefix)
		if m, ok := in.InReplyToStatus.(easyjson.Marshaler); ok {
			m.MarshalEasyJSON(out)
		} else if m, ok := in.InReplyToStatus.(json.Marshaler); ok {
			out.Raw(m.MarshalJSON())
		} else {
			out.Raw(json.Marshal(in.InReplyToStatus))
		}
	}
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix)
		out.Int64(int64(in.ID))
	}
	{
		const prefix string = ",\"geo\":"
		out.RawString(prefix)
		if m, ok := in.Geo.(easyjson.Marshaler); ok {
			m.MarshalEasyJSON(out)
		} else if m, ok := in.Geo.(json.Marshaler); ok {
			out.Raw(m.MarshalJSON())
		} else {
			out.Raw(json.Marshal(in.Geo))
		}
	}
	{
		const prefix string = ",\"retweeted\":"
		out.RawString(prefix)
		out.Bool(bool(in.Retweeted))
	}
	{
		const prefix string = ",\"in_reply_to_user_id_str\":"
		out.RawString(prefix)
		if m, ok := in.InReplyToUser.(easyjson.Marshaler); ok {
			m.MarshalEasyJSON(out)
		} else if m, ok := in.InReplyToUser.(json.Marshaler); ok {
			out.Raw(m.MarshalJSON())
		} else {
			out.Raw(json.Marshal(in.InReplyToUser))
		}
	}
	{
		const prefix string = ",\"place\":"
		out.RawString(prefix)
		if m, ok := in.Place.(easyjson.Marshaler); ok {
			m.MarshalEasyJSON(out)
		} else if m, ok := in.Place.(json.Marshaler); ok {
			out.Raw(m.MarshalJSON())
		} else {
			out.Raw(json.Marshal(in.Place))
		}
	}
	{
		const prefix string = ",\"user\":"
		out.RawString(prefix)
		(in.User).MarshalEasyJSON(out)
	}
	{
		const prefix string = ",\"in_reply_to_screen_name\":"
		out.RawString(prefix)
		if m, ok := in.InReplyToScreen.(easyjson.Marshaler); ok {
			m.MarshalEasyJSON(out)
		} else if m, ok := in.InReplyToScreen.(json.Marshaler); ok {
			out.Raw(m.MarshalJSON())
		} else {
			out.Raw(json.Marshal(in.InReplyToScreen))
		}
	}
	{
		const prefix string = ",\"source\":"
		out.RawString(prefix)
		out.String(string(in.Source))
	}
	{
		const prefix string = ",\"in_reply_to_status_id_str\":"
		out.RawString(prefix)
		if m, ok := in.InReplyToStatusI.(easyjson.Marshaler); ok {
			m.MarshalEasyJSON(out)
		} else if m, ok := in.InReplyToStatusI.(json.Marshaler); ok {
			out.Raw(m.MarshalJSON())
		} else {
			out.Raw(json.Marshal(in.InReplyToStatusI))
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Status) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonC80ae7adEncodeBanchModel1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Status) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonC80ae7adEncodeBanchModel1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Status) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonC80ae7adDecodeBanchModel1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Status) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonC80ae7adDecodeBanchModel1(l, v)
}
func easyjsonC80ae7adDecodeBanchModel2(in *jlexer.Lexer, out *SearchMetadata) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "max_id":
			out.MaxID = int64(in.Int64())
		case "since_id":
			out.SinceID = int64(in.Int64())
		case "refresh_url":
			out.RefreshURL = string(in.String())
		case "next_results":
			out.NextResults = string(in.String())
		case "count":
			out.Count = int(in.Int())
		case "completed_in":
			out.CompletedIn = float64(in.Float64())
		case "since_id_str":
			out.SinceIDStr = string(in.String())
		case "query":
			out.Query = string(in.String())
		case "max_id_str":
			out.MaxIDStr = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonC80ae7adEncodeBanchModel2(out *jwriter.Writer, in SearchMetadata) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"max_id\":"
		out.RawString(prefix[1:])
		out.Int64(int64(in.MaxID))
	}
	{
		const prefix string = ",\"since_id\":"
		out.RawString(prefix)
		out.Int64(int64(in.SinceID))
	}
	{
		const prefix string = ",\"refresh_url\":"
		out.RawString(prefix)
		out.String(string(in.RefreshURL))
	}
	{
		const prefix string = ",\"next_results\":"
		out.RawString(prefix)
		out.String(string(in.NextResults))
	}
	{
		const prefix string = ",\"count\":"
		out.RawString(prefix)
		out.Int(int(in.Count))
	}
	{
		const prefix string = ",\"completed_in\":"
		out.RawString(prefix)
		out.Float64(float64(in.CompletedIn))
	}
	{
		const prefix string = ",\"since_id_str\":"
		out.RawString(prefix)
		out.String(string(in.SinceIDStr))
	}
	{
		const prefix string = ",\"query\":"
		out.RawString(prefix)
		out.String(string(in.Query))
	}
	{
		const prefix string = ",\"max_id_str\":"
		out.RawString(prefix)
		out.String(string(in.MaxIDStr))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v SearchMetadata) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonC80ae7adEncodeBanchModel2(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v SearchMetadata) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonC80ae7adEncodeBanchModel2(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *SearchMetadata) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonC80ae7adDecodeBanchModel2(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *SearchMetadata) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonC80ae7adDecodeBanchModel2(l, v)
}
func easyjsonC80ae7adDecodeBanchModel3(in *jlexer.Lexer, out *Response) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "statuses":
			if in.IsNull() {
				in.Skip()
				out.Statuses = nil
			} else {
				in.Delim('[')
				if out.Statuses == nil {
					if !in.IsDelim(']') {
						out.Statuses = make([]Status, 0, 0)
					} else {
						out.Statuses = []Status{}
					}
				} else {
					out.Statuses = (out.Statuses)[:0]
				}
				for !in.IsDelim(']') {
					var v4 Status
					(v4).UnmarshalEasyJSON(in)
					out.Statuses = append(out.Statuses, v4)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "search_metadata":
			(out.SearchMetadata).UnmarshalEasyJSON(in)
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonC80ae7adEncodeBanchModel3(out *jwriter.Writer, in Response) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"statuses\":"
		out.RawString(prefix[1:])
		if in.Statuses == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v5, v6 := range in.Statuses {
				if v5 > 0 {
					out.RawByte(',')
				}
				(v6).MarshalEasyJSON(out)
			}
			out.RawByte(']')
		}
	}
	{
		const prefix string = ",\"search_metadata\":"
		out.RawString(prefix)
		(in.SearchMetadata).MarshalEasyJSON(out)
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Response) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonC80ae7adEncodeBanchModel3(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Response) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonC80ae7adEncodeBanchModel3(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Response) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonC80ae7adDecodeBanchModel3(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Response) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonC80ae7adDecodeBanchModel3(l, v)
}
func easyjsonC80ae7adDecodeBanchModel4(in *jlexer.Lexer, out *Metadata) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "iso_language_code":
			out.ISOLanguageCode = string(in.String())
		case "result_type":
			out.ResultType = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonC80ae7adEncodeBanchModel4(out *jwriter.Writer, in Metadata) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"iso_language_code\":"
		out.RawString(prefix[1:])
		out.String(string(in.ISOLanguageCode))
	}
	{
		const prefix string = ",\"result_type\":"
		out.RawString(prefix)
		out.String(string(in.ResultType))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Metadata) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonC80ae7adEncodeBanchModel4(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Metadata) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonC80ae7adEncodeBanchModel4(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Metadata) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonC80ae7adDecodeBanchModel4(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Metadata) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonC80ae7adDecodeBanchModel4(l, v)
}
func easyjsonC80ae7adDecodeBanchModel5(in *jlexer.Lexer, out *Hashtag) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "text":
			out.Text = string(in.String())
		case "indices":
			if in.IsNull() {
				in.Skip()
				out.Indices = nil
			} else {
				in.Delim('[')
				if out.Indices == nil {
					if !in.IsDelim(']') {
						out.Indices = make([]int, 0, 8)
					} else {
						out.Indices = []int{}
					}
				} else {
					out.Indices = (out.Indices)[:0]
				}
				for !in.IsDelim(']') {
					var v7 int
					v7 = int(in.Int())
					out.Indices = append(out.Indices, v7)
					in.WantComma()
				}
				in.Delim(']')
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonC80ae7adEncodeBanchModel5(out *jwriter.Writer, in Hashtag) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"text\":"
		out.RawString(prefix[1:])
		out.String(string(in.Text))
	}
	{
		const prefix string = ",\"indices\":"
		out.RawString(prefix)
		if in.Indices == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v8, v9 := range in.Indices {
				if v8 > 0 {
					out.RawByte(',')
				}
				out.Int(int(v9))
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Hashtag) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonC80ae7adEncodeBanchModel5(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Hashtag) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonC80ae7adEncodeBanchModel5(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Hashtag) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonC80ae7adDecodeBanchModel5(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Hashtag) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonC80ae7adDecodeBanchModel5(l, v)
}
func easyjsonC80ae7adDecodeBanchModel6(in *jlexer.Lexer, out *Entities) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "urls":
			if in.IsNull() {
				in.Skip()
				out.URLs = nil
			} else {
				in.Delim('[')
				if out.URLs == nil {
					if !in.IsDelim(']') {
						out.URLs = make([]interface{}, 0, 4)
					} else {
						out.URLs = []interface{}{}
					}
				} else {
					out.URLs = (out.URLs)[:0]
				}
				for !in.IsDelim(']') {
					var v10 interface{}
					if m, ok := v10.(easyjson.Unmarshaler); ok {
						m.UnmarshalEasyJSON(in)
					} else if m, ok := v10.(json.Unmarshaler); ok {
						_ = m.UnmarshalJSON(in.Raw())
					} else {
						v10 = in.Interface()
					}
					out.URLs = append(out.URLs, v10)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "hashtags":
			if in.IsNull() {
				in.Skip()
				out.Hashtags = nil
			} else {
				in.Delim('[')
				if out.Hashtags == nil {
					if !in.IsDelim(']') {
						out.Hashtags = make([]Hashtag, 0, 1)
					} else {
						out.Hashtags = []Hashtag{}
					}
				} else {
					out.Hashtags = (out.Hashtags)[:0]
				}
				for !in.IsDelim(']') {
					var v11 Hashtag
					(v11).UnmarshalEasyJSON(in)
					out.Hashtags = append(out.Hashtags, v11)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "user_mentions":
			if in.IsNull() {
				in.Skip()
				out.UserMentions = nil
			} else {
				in.Delim('[')
				if out.UserMentions == nil {
					if !in.IsDelim(']') {
						out.UserMentions = make([]interface{}, 0, 4)
					} else {
						out.UserMentions = []interface{}{}
					}
				} else {
					out.UserMentions = (out.UserMentions)[:0]
				}
				for !in.IsDelim(']') {
					var v12 interface{}
					if m, ok := v12.(easyjson.Unmarshaler); ok {
						m.UnmarshalEasyJSON(in)
					} else if m, ok := v12.(json.Unmarshaler); ok {
						_ = m.UnmarshalJSON(in.Raw())
					} else {
						v12 = in.Interface()
					}
					out.UserMentions = append(out.UserMentions, v12)
					in.WantComma()
				}
				in.Delim(']')
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonC80ae7adEncodeBanchModel6(out *jwriter.Writer, in Entities) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"urls\":"
		out.RawString(prefix[1:])
		if in.URLs == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v13, v14 := range in.URLs {
				if v13 > 0 {
					out.RawByte(',')
				}
				if m, ok := v14.(easyjson.Marshaler); ok {
					m.MarshalEasyJSON(out)
				} else if m, ok := v14.(json.Marshaler); ok {
					out.Raw(m.MarshalJSON())
				} else {
					out.Raw(json.Marshal(v14))
				}
			}
			out.RawByte(']')
		}
	}
	{
		const prefix string = ",\"hashtags\":"
		out.RawString(prefix)
		if in.Hashtags == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v15, v16 := range in.Hashtags {
				if v15 > 0 {
					out.RawByte(',')
				}
				(v16).MarshalEasyJSON(out)
			}
			out.RawByte(']')
		}
	}
	{
		const prefix string = ",\"user_mentions\":"
		out.RawString(prefix)
		if in.UserMentions == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v17, v18 := range in.UserMentions {
				if v17 > 0 {
					out.RawByte(',')
				}
				if m, ok := v18.(easyjson.Marshaler); ok {
					m.MarshalEasyJSON(out)
				} else if m, ok := v18.(json.Marshaler); ok {
					out.Raw(m.MarshalJSON())
				} else {
					out.Raw(json.Marshal(v18))
				}
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Entities) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonC80ae7adEncodeBanchModel6(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Entities) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonC80ae7adEncodeBanchModel6(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Entities) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonC80ae7adDecodeBanchModel6(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Entities) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonC80ae7adDecodeBanchModel6(l, v)
}