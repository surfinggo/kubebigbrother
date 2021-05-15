package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"k8s.io/klog/v2"
)

const (
	ReasonServerError    = "ServerError"
	ReasonInvalidRequest = "InvalidRequest"
)

type Detail struct {
	Message string `json:"message"`
}

type E struct {
	StatusCode  int      `json:"status_code"`
	Reason      string   `json:"reason"`
	Description string   `json:"description"`
	Details     []Detail `json:"details,omitempty"`
}

func (e *E) Error() string {
	return fmt.Sprintf("[%s] %s", e.Reason, e.Description)
}

func e(statusCode int, reason, description string) *E {
	return &E{
		StatusCode:  statusCode,
		Reason:      reason,
		Description: description,
	}
}

func (app *App) HandleError(c *gin.Context, err error) {
	switch t := err.(type) {
	case *E:
		if t.StatusCode != 0 {
			c.JSON(t.StatusCode, t)
		} else {
			c.JSON(400, t)
		}
	default:
		klog.Errorf("server error occurred: %s", err)
		c.JSON(500, e(500, ReasonServerError, err.Error()))
	}
}

func ShouldBindQuery(c *gin.Context, i interface{}) error {
	if err := c.ShouldBindQuery(i); err != nil {
		return e(400, ReasonInvalidRequest, errors.Wrap(err, "query decode error").Error())
	}
	return nil
}

// MustBindQuery used to bind query data to structure,
//  returns false when failed, and response is written automatically
func (app *App) MustBindQuery(c *gin.Context, i interface{}) bool {
	if err := ShouldBindQuery(c, i); err != nil {
		app.HandleError(c, err)
		return false
	}
	return true
}

func ShouldBindJSON(c *gin.Context, i interface{}) error {
	if err := c.ShouldBindJSON(i); err != nil {
		return e(400, ReasonInvalidRequest, errors.Wrap(err, "JSON decode error").Error())
	}
	return nil
}

// MustBindJSON used to bind query data to structure,
// returns false when failed, and response is written automatically
func (app *App) MustBindJSON(c *gin.Context, i interface{}) bool {
	if err := ShouldBindJSON(c, i); err != nil {
		app.HandleError(c, err)
		return false
	}
	return true
}
