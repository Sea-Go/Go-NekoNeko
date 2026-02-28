// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package folder

import (
	"net/http"

	"favorite-system/api/internal/logic/folder"
	"favorite-system/api/internal/svc"
	"favorite-system/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func CreateFolderHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CreateFolderReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := folder.NewCreateFolderLogic(r.Context(), svcCtx)
		resp, err := l.CreateFolder(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
