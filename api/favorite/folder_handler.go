package favorite

import (
	"encoding/json"
	"net/http"
	"strconv"

	folder "sea-try-go/service/favorite/folder"
)

type FolderHandler struct {
	Svc *folder.Service
}

func NewFolderHandler(svc *folder.Service) *FolderHandler {
	return &FolderHandler{Svc: svc}
}

type createFolderReq struct {
	Name     string `json:"name"`
	IsPublic bool   `json:"is_public"`
}

// 先用 Header 传用户 id：X-User-Id: 1
func getUserID(r *http.Request) (int64, error) {
	return strconv.ParseInt(r.Header.Get("X-User-Id"), 10, 64)
}

// POST /api/favorite/folders
func (h *FolderHandler) CreateFolder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, err := getUserID(r)
	if err != nil || userID <= 0 {
		http.Error(w, "invalid X-User-Id", http.StatusBadRequest)
		return
	}

	var req createFolderReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	f, err := h.Svc.CreateFolder(r.Context(), userID, req.Name, req.IsPublic)
	if err != nil {
		switch err {
		case folder.ErrFolderNameEmpty, folder.ErrFolderNameTooLong:
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		case folder.ErrFolderNameExists:
			http.Error(w, err.Error(), http.StatusConflict)
			return
		default:
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(f)
}

// GET /api/favorite/folders
func (h *FolderHandler) ListFolders(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, err := getUserID(r)
	if err != nil || userID <= 0 {
		http.Error(w, "invalid X-User-Id", http.StatusBadRequest)
		return
	}

	list, err := h.Svc.ListByUser(r.Context(), userID)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(list)
}
