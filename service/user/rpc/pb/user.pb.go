package pb

type UserInfo struct {
	Uid       int64
	Score     uint32
	Username  string
	Email     string
	ExtraInfo map[string]string
}

type GetUserReq struct {
	Uid int64
}

type GetUserResp struct {
	User  *UserInfo
	Found bool
}

type CreateUserReq struct {
	Username  string
	Password  string
	Email     string
	ExtraInfo map[string]string
}

type CreateUserResp struct {
	Uid int64
}

type LoginReq struct {
	Username string
	Password string
}

type LoginResp struct {
	Uid    int64
	Status uint64
}

type LogoutReq struct {
	Token string
}

type LogoutResp struct {
	Success bool
}

type UpdateUserReq struct {
	Uid       int64
	Username  string
	Password  string
	Email     string
	ExtraInfo map[string]string
}

type UpdateUserResp struct {
	User *UserInfo
}

type DeleteUserReq struct {
	Uid int64
}

type DeleteUserResp struct {
	Success bool
}
