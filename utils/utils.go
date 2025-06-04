package utils

import (
	"context"
	// "fmt"
	"net/http"

	"github.com/a-h/templ"

	// "github.com/tinkerbaj/moje/utils"

		"github.com/fatih/color"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
)


//Almost all functions here I add from my old project but I left it for reference later 

// TemplRender implements the render.Render interface.
type TemplRender struct {
	Code int
	Data templ.Component
}

// Render implements the render.Render interface.
func (t TemplRender) Render(w http.ResponseWriter) error {
	t.WriteContentType(w)
	w.WriteHeader(t.Code)
	if t.Data != nil {
		return t.Data.Render(context.Background(), w)
	}
	return nil
}

// WriteContentType implements the render.Render interface.
func (t TemplRender) WriteContentType(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
}

// Instance implements the render.Render interface.
func (t *TemplRender) Instance(name string, data interface{}) render.Render {
	if templData, ok := data.(templ.Component); ok {
		return &TemplRender{
			Code: http.StatusOK,
			Data: templData,
		}
	}
	return nil
}

func Render(comp templ.Component, c *gin.Context) {
	comp.Render(c.Request.Context(), c.Writer)
}

func ColorError(message string, err ...error) {
var e = color.New(color.FgHiRed, color.Bold)
e.Printf(message, err)
}
