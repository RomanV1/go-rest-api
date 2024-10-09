package responses

import "github.com/gin-gonic/gin"

var (
	UserNotFound          = gin.H{"message": "User not found"}
	ErrorRetrievingUsers  = gin.H{"message": "Unable to retrieve users"}
	InvalidLimitParam     = gin.H{"message": "Invalid limit parameter; must be a number"}
	InvalidOffsetParam    = gin.H{"message": "Invalid offset parameter; must be a number"}
	InvalidIDParam        = gin.H{"message": "Invalid UUID parameter; must be a UUID"}
	ErrorParsingJSON      = gin.H{"message": "Unable to parse request body as JSON"}
	UserCreationError     = gin.H{"message": "Unable to create user"}
	UserUpdateError       = gin.H{"message": "Unable to update user"}
	UserDeleteError       = gin.H{"message": "Unable to delete user"}
	UserDeletionSuccess   = gin.H{"message": "User successfully deleted"}
	InternalServerError   = gin.H{"message": "An internal server error has occurred"}
	EmailAlreadyExists    = gin.H{"message": "Email already exists"}
	UsernameAlreadyExists = gin.H{"message": "Username already exists"}
)
