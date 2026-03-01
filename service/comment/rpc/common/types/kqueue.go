package kqtypes

type CommentKafkaMsg struct {
	CommentId  int64  `json:"comment_id"`
	TargetType string `json:"target_type"`
	TargetId   string `json:"target_id"`
	UserId     int64  `json:"user_id"`
	RootId     int64  `json:"root_id"`
	ParentId   int64  `json:"parent_id"`
	Content    string `json:"content"`
	Meta       string `json:"meta"`
	OwnerId    int64  `json:"owner_id"`
	Attribute  int64  `json:"attribute"`
	CreateTime int64  `json:"create_time"`
}
