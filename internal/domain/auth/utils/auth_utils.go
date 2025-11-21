package auth

import (
	"github.com/gin-gonic/gin"
)

// authorizeUserAccess checks if the authenticated user can access the requested user's data
func AuthorizeUserAccess(c *gin.Context, requestedUserID int) bool {
	// Get authenticated user ID from context (set by auth middleware)
	authenticatedUserID, exists := c.Get("user_id")
	if !exists {
		return false
	}

	// Convert to int for comparison
	authUserID, ok := authenticatedUserID.(int)
	if !ok {
		return false
	}

	// Users can only access their own data
	return authUserID == requestedUserID
}
