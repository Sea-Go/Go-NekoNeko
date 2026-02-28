package Model

// 任务定义表
/*type TaskDef struct {
	TaskID           int64     `gorm:"column:task_id;primaryKey"`
	Name             string    `gorm:"column:name"`
	Desc             string    `gorm:"column:desc"`
	RequiredProgress int64     `gorm:"column:required_progress"`
	UpdatedAt        time.Time `gorm:"column:updated_at"`
}

func (TaskDef) TableName() string { return "task_def" }*/
type TaskDef struct {
	TaskID           int64
	Name             string
	Desc             string
	RequiredProgress int64
}

var AllTaskDefs = []TaskDef{
	{TaskID: 190, Name: "点赞任务", Desc: "给文章点赞5次", RequiredProgress: 3},
	{TaskID: 1902, Name: "文章点赞任务", Desc: "有一篇文章点赞大于3", RequiredProgress: 5},
}

// 用户任务进度表
type TaskProgress struct {
	UserID   int64  `gorm:"primary_key;column:user_id"`
	TaskID   int64  `gorm:"primary_key;column:task_id"`
	Status   string `gorm:"column:status"`
	Progress int64  `gorm:"column:progress"`
	Target   int64  `gorm:"column:target"`
}

func (TaskProgress) TableName() string {
	return "task_progress"
}
