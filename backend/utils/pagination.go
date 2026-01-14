package utils

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// PaginationParams represents pagination parameters
type PaginationParams struct {
	Page     int
	PageSize int
	Search   string
}

// ParsePaginationParams parses pagination parameters from gin context
func ParsePaginationParams(c *gin.Context) PaginationParams {
	page := 1
	pageSize := 10

	if pageStr := c.Query("page"); pageStr != "" {
		if parsed, err := strconv.Atoi(pageStr); err == nil && parsed > 0 {
			page = parsed
		}
	}

	if pageSizeStr := c.Query("page_size"); pageSizeStr != "" {
		if parsed, err := strconv.Atoi(pageSizeStr); err == nil && parsed > 0 && parsed <= 100 {
			pageSize = parsed
		}
	}

	return PaginationParams{
		Page:     page,
		PageSize: pageSize,
		Search:   c.Query("search"),
	}
}

// BuildPaginatedResponse builds a paginated response
func BuildPaginatedResponse(data interface{}, total, page, pageSize int, search string) gin.H {
	totalPages := (total + pageSize - 1) / pageSize
	if totalPages == 0 {
		totalPages = 1
	}

	return gin.H{
		"data":        data,
		"total":       total,
		"page":        page,
		"page_size":   pageSize,
		"total_pages": totalPages,
		"search":      search,
	}
}
