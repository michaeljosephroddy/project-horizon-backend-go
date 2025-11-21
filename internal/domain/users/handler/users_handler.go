// internal/domain/users/handler/users_handler.go
package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	authutils "github.com/michaeljosephroddy/project-horizon-backend-go/internal/domain/auth/utils"
	"github.com/michaeljosephroddy/project-horizon-backend-go/internal/domain/users/service"
)

type UserHandler struct {
	usersService *service.UserService
}

func NewUsersHandler(usersService *service.UserService) *UserHandler {
	return &UserHandler{
		usersService: usersService,
	}
}

// GetUserMedications handles GET /users/:userID/user-medication
func (uh *UserHandler) GetUserMedications(c *gin.Context) {
	userIDStr := c.Param("userID")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil || userID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid userID"})
		return
	}

	// Authorization check - ensure user can only access their own medications
	if !authutils.AuthorizeUserAccess(c, userID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	medications, err := uh.usersService.GetUserMedications(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve user medications"})
		return
	}

	c.JSON(http.StatusOK, medications)
}
