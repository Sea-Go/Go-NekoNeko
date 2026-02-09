package favorite

import (
	"net/http"

	"sea-try-go/api/internal/logic/favorite"
	"sea-try-go/api/internal/svc"
	"sea-try-go/api/internal/utils"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// DeleteFavoriteHandler 删除收藏的处理器
func DeleteFavoriteHandler(serverCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 步骤 1: 解析请求数据
		var req favorite.DeleteFavoriteReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		// 步骤 2: 获取用户 ID
		userID, err := utils.GetUserIDFromRequest(r, serverCtx.Config.UserAuth.AccessSecret)
		if err != nil {
			utils.WriteErrorResponse(w, r, http.StatusUnauthorized, "invalid or missing authorization token")
			return
		}

		// 步骤 3: 执行业务逻辑
		logic := favorite.NewDeleteFavoriteLogic(serverCtx.FavoriteItemService)
		err = logic.Execute(r.Context(), req, userID)
		if err != nil {
			utils.WriteBusinessErrorResponse(w, r, err)
			return
		}

		// 步骤 5: 返回响应
		httpx.OkJsonCtx(r.Context(), w, favorite.DeleteFavoriteResp{
			Success: true,
			Message: "删除成功",
		})
	}
}
