// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package user

import (
	"net/http"
	"sea-try-go/service/common/response"
	"sea-try-go/service/user/common/errmsg"
	"sea-try-go/service/user/user/api/internal/logic/user"
	"sea-try-go/service/user/user/api/internal/svc"
	"sea-try-go/service/user/user/api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func UpdateuserHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UpdateUserReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := user.NewUpdateuserLogic(r.Context(), svcCtx)
		resp, code := l.Updateuser(&req)
		httpx.OkJson(w, &response.Response{
			Code: code,
			Msg:  errmsg.GetErrMsg(code),
			Data: resp,
		})
	}
}
