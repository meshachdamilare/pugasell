package errors

// import (
// 	"net/http"

// 	"github.com/gin-gonic/gin"
// )

// type ErrorResponse struct {
// 	code    int
// 	messsge string
// 	ctx     *gin.Context
// }

// func (e *ErrorResponse) BadRequestError(msg string, err error) {
// 	e.code = http.StatusBadRequest
// 	if msg != "" {
// 		e.messsge = msg
// 		e.ctx.JSON(e.code, gin.H{"error": e.messsge})
// 	} else {
// 		e.ctx.JSON(e.code, gin.H{"error": err.Error()})
// 	}
// }
