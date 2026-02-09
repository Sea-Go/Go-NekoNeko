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

// 创建收藏夹
func CreateFolderHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CreateFolderReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := favorite.NewCreateFolderLogic(r.Context(), svcCtx)
		resp, err := l.CreateFolder(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
