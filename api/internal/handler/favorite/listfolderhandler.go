// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package favorite

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"sea-try-go/api/internal/logic/favorite"
	"sea-try-go/api/internal/svc"
	"sea-try-go/api/internal/types"
)

// 获取收藏夹列表
func ListFolderHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ListFolderReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := favorite.NewListFolderLogic(r.Context(), svcCtx)
		resp, err := l.ListFolder(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
