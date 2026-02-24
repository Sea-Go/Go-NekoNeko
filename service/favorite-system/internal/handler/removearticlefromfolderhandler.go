// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package handler

import (
	"net/http"

	"favorite-system/internal/logic"
	"favorite-system/internal/svc"
	"favorite-system/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func removeArticleFromFolderHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RemoveArticleReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewRemoveArticleFromFolderLogic(r.Context(), svcCtx)
		resp, err := l.RemoveArticleFromFolder(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
