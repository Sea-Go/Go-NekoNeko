package model

import "time"

type Subject struct {
	Id         int64     `gorm:"column:id;primaryKey;autoIncrement;comment:主键ID"`
	TargetType string    `gorm:"column:target_type;type:varchar(32);not null;uniqueIndex:idx_target;comment:内容类型"`
	TargetId   string    `gorm:"column:target_id;type:varchar(64);not null;uniqueIndex:idx_target;comment:内容ID"`
	TotalCount int64     `gorm:"column:total_count;not null;default:0;comment:总评论数"`
	RootCount  int64     `gorm:"column:root_count;not null;default:0;comment:根评论数"`
	State      int32     `gorm:"column:state;not null;default:0;comment:状态:0正常,1关闭,2仅粉丝"`
	Attribute  int64     `gorm:"column:attribute;not null;default:0;comment:属性位图"`
	CreatedAt  time.Time `gorm:"column:created_at;autoCreateTime;comment:创建时间"`
	UpdatedAt  time.Time `gorm:"column:updated_at;autoUpdateTime;comment:更新时间"`
}

func (Subject) TableName() string {
	return "subject"
}

type CommentIndex struct {
	Id           int64     `gorm:"column:id;primaryKey;autoIncrement:false;comment:评论ID"`
	TargetType   string    `gorm:"column:target_type;type:varchar(32);not null;index:idx_target_root_time,priority:1;comment:内容类型"`
	TargetId     string    `gorm:"column:target_id;type:varchar(64);not null;index:idx_target_root_time,priority:2;comment:内容Id"`
	UserId       int64     `gorm:"column:user_id;not null;comment:发布者Uid"`
	RootId       int64     `gorm:"column:root_id;not null;default:0;index:idx_target_root_time,priority:3;comment:根评论ID"`
	ParentId     int64     `gorm:"column:parent_id;not null;default:0;comment:父评论Id"`
	LikeCount    int64     `gorm:"column:like_count;not null;default:0;comment:点赞数"`
	DislikeCount int64     `gorm:"column:dislike_count;not null;default:0;comment:点踩数"`
	ReplyCount   int64     `gorm:"column:reply_count;not null;default:0;comment:子评论数"`
	State        int32     `gorm:"column:state;not null;default:0;comment:状态:0正常,1审核,2删除,3封禁"`
	Attribute    int64     `gorm:"column:attribute;not null;default:0;comment:属性位图"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime;index:idx_target_root_time,priority:4;comment:创建时间"`
	UpdatedAt    time.Time `gorm:"column:updated_at;autoUpdateTime;comment:更新时间"`
}

func (CommentIndex) TableName() string {
	return "comment_index"
}

type CommentContent struct {
	//CommentId就是CommentIndex中的Id
	CommentId int64     `gorm:"column:comment_id;primaryKey;autoIncrement:false;comment:评论Id"`
	Content   string    `gorm:"column:content;type:text;not null;comment:评论正文"`
	Meta      string    `gorm:"column:meta;type:jsonb;comment:扩展内容"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime;comment:创建时间"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime;comment:更新时间"`
}

func (CommentContent) TableName() string {
	return "comment_content"
}
