package handler

import (
	"net/http"
	"strconv"

	"favorite-system/internal/app"
	"favorite-system/internal/pkg/httpx"
	foldersvc "favorite-system/internal/service/folder"

	"github.com/gin-gonic/gin"
)

type FolderHandler struct {
	svc *foldersvc.Service
}

func NewFolderHandler(a *app.App) *FolderHandler {
	return &FolderHandler{
		svc: foldersvc.New(a.FolderRepo),
	}
}

type createFolderReq struct {
	UserID   int64  `json:"user_id" binding:"required"`
	Name     string `json:"name" binding:"required"`
	IsPublic bool   `json:"is_public"`
}

// POST /api/v1/folders
func (h *FolderHandler) Create(c *gin.Context) {
	var req createFolderReq

	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Fail(c, http.StatusBadRequest, 40001, "invalid request body")
		return
	}

	f, err := h.svc.Create(c.Request.Context(), foldersvc.CreateInput{
		UserID:   req.UserID,
		Name:     req.Name,
		IsPublic: req.IsPublic,
	})
	if err != nil {
		httpx.Fail(c, http.StatusBadRequest, 40002, err.Error())
		return
	}

	httpx.OK(c, f)
}

// GET /api/v1/folders?user_id=1
func (h *FolderHandler) ListByUser(c *gin.Context) {
	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		httpx.Fail(c, http.StatusBadRequest, 40003, "missing query param: user_id")
		return
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil || userID <= 0 {
		httpx.Fail(c, http.StatusBadRequest, 40004, "invalid user_id")
		return
	}

	list, err := h.svc.ListByUser(c.Request.Context(), userID)
	if err != nil {
		httpx.Fail(c, http.StatusInternalServerError, 50001, err.Error())
		return
	}

	httpx.OK(c, list)
}

// GET /api/v1/folders/:id
func (h *FolderHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		httpx.Fail(c, http.StatusBadRequest, 40005, "invalid id")
		return
	}

	f, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		httpx.Fail(c, http.StatusNotFound, 40401, err.Error())
		return
	}

	httpx.OK(c, f)
}

// DELETE /api/v1/folders/:id
func (h *FolderHandler) SoftDelete(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		httpx.Fail(c, http.StatusBadRequest, 40006, "invalid id")
		return
	}

	if err := h.svc.SoftDelete(c.Request.Context(), id); err != nil {
		httpx.Fail(c, http.StatusBadRequest, 40007, err.Error())
		return
	}

	httpx.OK(c, gin.H{"deleted": true})
}
