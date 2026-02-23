package models

import "net/http"

// ContextKey is a custom type for context value keys to avoid collisions.
type ContextKey string

const (
	CtxTenantID       ContextKey = "tenant_id"
	CtxUserID         ContextKey = "user_id"
	CtxUserRole       ContextKey = "user_role"
	CtxUserEmail      ContextKey = "user_email"
	CtxResponseWriter ContextKey = "response_writer"
	CtxHTTPRequest    ContextKey = "http_request"
)

// PageInfo holds pagination metadata for list queries.
type PageInfo struct {
	TotalCount int `json:"total_count"`
	Page       int `json:"page"`
	PerPage    int `json:"per_page"`
	TotalPages int `json:"total_pages"`
}

// ResponseWriterFromContext extracts the http.ResponseWriter stored in context.
func ResponseWriterFromContext(ctx interface{ Value(any) any }) http.ResponseWriter {
	if w, ok := ctx.Value(CtxResponseWriter).(http.ResponseWriter); ok {
		return w
	}
	return nil
}
