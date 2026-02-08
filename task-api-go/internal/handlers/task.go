package handlers

import (
    "encoding/json"
    "net/http"
    "strconv"
    "strings"

    "taskapi/internal/models"
    "taskapi/internal/storage"
)

type TaskHandler struct {
    store *storage.MemoryStore
}

func NewTaskHandler(store *storage.MemoryStore) *TaskHandler {
    return &TaskHandler{store: store}
}

// Tasks handles:
// GET /tasks           -> list all tasks
// GET /tasks?id=1      -> get task by id
// POST /tasks          -> create task
// PATCH /tasks?id=1    -> update done
func (h *TaskHandler) Tasks(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case http.MethodGet:
        h.handleGet(w, r)
    case http.MethodPost:
        h.handlePost(w, r)
    case http.MethodPatch:
        h.handlePatch(w, r)
    default:
        writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
    }
}

func (h *TaskHandler) handleGet(w http.ResponseWriter, r *http.Request) {
    idStr := r.URL.Query().Get("id")
    if idStr == "" {
        tasks := h.store.List()
        writeJSON(w, http.StatusOK, tasks)
        return
    }

    id, err := strconv.Atoi(idStr)
    if err != nil || id <= 0 {
        writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
        return
    }

    t, err := h.store.Get(id)
    if err == storage.ErrNotFound {
        writeJSON(w, http.StatusNotFound, map[string]string{"error": "task not found"})
        return
    }
    if err != nil {
        writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
        return
    }

    writeJSON(w, http.StatusOK, t)
}

type createReq struct {
    Title string `json:"title"`
}

func (h *TaskHandler) handlePost(w http.ResponseWriter, r *http.Request) {
    var req createReq
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid title"})
        return
    }

    title := strings.TrimSpace(req.Title)
    if title == "" {
        writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid title"})
        return
    }

    t := h.store.Create(title)
    writeJSON(w, http.StatusCreated, t)
}

type patchReq struct {
    Done *bool `json:"done"`
}

func (h *TaskHandler) handlePatch(w http.ResponseWriter, r *http.Request) {
    idStr := r.URL.Query().Get("id")
    if idStr == "" {
        writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
        return
    }

    id, err := strconv.Atoi(idStr)
    if err != nil || id <= 0 {
        writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
        return
    }

    var req patchReq
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Done == nil {
        writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid done"})
        return
    }

    if err := h.store.UpdateDone(id, *req.Done); err == storage.ErrNotFound {
        writeJSON(w, http.StatusNotFound, map[string]string{"error": "task not found"})
        return
    } else if err != nil {
        writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
        return
    }

    writeJSON(w, http.StatusOK, map[string]bool{"updated": true})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    _ = json.NewEncoder(w).Encode(v)
}

// Ensure unused imports are not present
var _ models.Task
