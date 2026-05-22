package rag

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strings"
)

type HTTPHandler struct {
	service *Service
	mux     *http.ServeMux
}

func NewHTTPHandler(service *Service) http.Handler {
	handler := &HTTPHandler{
		service: service,
		mux:     http.NewServeMux(),
	}
	handler.routes()
	return handler
}

func (h *HTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	log.Printf("request method=%s path=%s", r.Method, r.URL.Path)
	h.mux.ServeHTTP(w, r)
}

func (h *HTTPHandler) routes() {
	h.mux.HandleFunc("GET /health", h.handleHealth)
	h.mux.HandleFunc("GET /capabilities", h.handleCapabilities)
	h.mux.HandleFunc("POST /documents", h.handleCreateDocument)
	h.mux.HandleFunc("GET /documents", h.handleListDocuments)
	h.mux.HandleFunc("POST /documents/", h.handleDocumentAction)
	h.mux.HandleFunc("POST /retrieval/search", h.handleSearch)
	h.mux.HandleFunc("POST /chat/sessions", h.handleCreateChatSession)
	h.mux.HandleFunc("POST /chat/sessions/", h.handleChatSessionAction)
	h.mux.HandleFunc("GET /chat/sessions/", h.handleChatSessionAction)
	h.mux.HandleFunc("GET /memory/users/", h.handleMemoryAction)
	h.mux.HandleFunc("POST /memory/users/", h.handleMemoryAction)
}

func (h *HTTPHandler) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, h.service.Health())
}

func (h *HTTPHandler) handleCapabilities(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, h.service.Capabilities())
}

func (h *HTTPHandler) handleCreateDocument(w http.ResponseWriter, r *http.Request) {
	var req CreateDocumentRequest
	if !decodeJSON(w, r, &req) {
		return
	}

	doc, err := h.service.CreateDocument(req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, doc)
}

func (h *HTTPHandler) handleListDocuments(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"documents": h.service.ListDocuments(),
	})
}

func (h *HTTPHandler) handleDocumentAction(w http.ResponseWriter, r *http.Request) {
	documentID, action, ok := splitTwoPartAction(r.URL.Path, "/documents/")
	if !ok || action != "index" {
		writeError(w, http.StatusNotFound, "not_found", "route was not found", nil)
		return
	}

	result, err := h.service.IndexDocument(documentID)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func (h *HTTPHandler) handleSearch(w http.ResponseWriter, r *http.Request) {
	var req SearchRequest
	if !decodeJSON(w, r, &req) {
		return
	}

	result, err := h.service.Search(req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func (h *HTTPHandler) handleCreateChatSession(w http.ResponseWriter, r *http.Request) {
	var req CreateChatSessionRequest
	if !decodeOptionalJSON(w, r, &req) {
		return
	}
	writeJSON(w, http.StatusCreated, h.service.CreateChatSession(req))
}

func (h *HTTPHandler) handleChatSessionAction(w http.ResponseWriter, r *http.Request) {
	sessionID, action, ok := splitOptionalAction(r.URL.Path, "/chat/sessions/")
	if !ok {
		writeError(w, http.StatusNotFound, "not_found", "route was not found", nil)
		return
	}

	if r.Method == http.MethodGet && action == "" {
		session, err := h.service.GetChatSession(sessionID)
		if err != nil {
			writeServiceError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, session)
		return
	}

	if r.Method == http.MethodPost && action == "messages" {
		var req ChatMessageRequest
		if !decodeJSON(w, r, &req) {
			return
		}

		response, err := h.service.AddChatMessage(sessionID, req)
		if err != nil {
			writeServiceError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, response)
		return
	}

	writeError(w, http.StatusNotFound, "not_found", "route was not found", nil)
}

func (h *HTTPHandler) handleMemoryAction(w http.ResponseWriter, r *http.Request) {
	userID, action, ok := splitOptionalAction(r.URL.Path, "/memory/users/")
	if !ok {
		writeError(w, http.StatusNotFound, "not_found", "route was not found", nil)
		return
	}

	if r.Method == http.MethodGet && action == "" {
		writeJSON(w, http.StatusOK, h.service.GetMemory(userID))
		return
	}

	if r.Method == http.MethodPost && action == "facts" {
		var req AddFactRequest
		if !decodeJSON(w, r, &req) {
			return
		}

		memory, err := h.service.AddFact(userID, req)
		if err != nil {
			writeServiceError(w, err)
			return
		}
		writeJSON(w, http.StatusCreated, memory)
		return
	}

	writeError(w, http.StatusNotFound, "not_found", "route was not found", nil)
}

func decodeJSON(w http.ResponseWriter, r *http.Request, target any) bool {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json", "request body must be valid JSON", map[string]any{
			"decode_error": err.Error(),
		})
		return false
	}
	return true
}

func decodeOptionalJSON(w http.ResponseWriter, r *http.Request, target any) bool {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid_body", "request body could not be read", map[string]any{
			"read_error": err.Error(),
		})
		return false
	}
	if len(strings.TrimSpace(string(body))) == 0 {
		return true
	}

	decoder := json.NewDecoder(strings.NewReader(string(body)))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json", "request body must be valid JSON", map[string]any{
			"decode_error": err.Error(),
		})
		return false
	}
	return true
}

func writeServiceError(w http.ResponseWriter, err error) {
	status := http.StatusInternalServerError
	var validation validationError
	if errors.Is(err, ErrDocumentNotFound) || errors.Is(err, ErrSessionNotFound) {
		status = http.StatusNotFound
	} else if errors.As(err, &validation) {
		status = http.StatusBadRequest
	}

	writeError(w, status, errorCode(err), err.Error(), nil)
}

func writeError(w http.ResponseWriter, status int, code string, message string, details map[string]any) {
	if details == nil {
		details = map[string]any{}
	}
	writeJSON(w, status, ErrorResponse{
		Error: ErrorBody{
			Code:    code,
			Message: message,
			Details: details,
		},
	})
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("response encode failed: %v", err)
	}
}

func splitTwoPartAction(path string, prefix string) (string, string, bool) {
	remaining := strings.Trim(strings.TrimPrefix(path, prefix), "/")
	parts := strings.Split(remaining, "/")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", false
	}
	return parts[0], parts[1], true
}

func splitOptionalAction(path string, prefix string) (string, string, bool) {
	remaining := strings.Trim(strings.TrimPrefix(path, prefix), "/")
	if remaining == "" {
		return "", "", false
	}

	parts := strings.Split(remaining, "/")
	if len(parts) == 1 && parts[0] != "" {
		return parts[0], "", true
	}
	if len(parts) == 2 && parts[0] != "" && parts[1] != "" {
		return parts[0], parts[1], true
	}
	return "", "", false
}
