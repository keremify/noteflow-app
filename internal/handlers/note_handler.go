package handlers

import (
	"net/http"
	"strconv"

	"saasproject/internal/models"
	"saasproject/internal/services"

	"github.com/gin-gonic/gin"
)

type NoteHandler struct {
	Service *services.NoteService
}

func NewNoteHandler(s *services.NoteService) *NoteHandler {
	return &NoteHandler{Service: s}
}

// CREATE
func (h *NoteHandler) Create(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req struct {
		Title   string `json:"title" binding:"required"`
		Content string `json:"content" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	note := models.Note{
		UserID:  userID,
		Title:   req.Title,
		Content: req.Content,
	}

	if err := h.Service.CreateNote(&note); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "note oluÅŸturulamadÄ±"})
		return
	}

	c.JSON(http.StatusCreated, note)
}

// LIST
func (h *NoteHandler) List(c *gin.Context) {
	userID := c.GetUint("user_id")
	role := c.GetString("role")

	notes, err := h.Service.ListNotes(userID, role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "notlar alÄ±namadÄ±"})
		return
	}

	c.JSON(http.StatusOK, notes)
}

// UPDATE
func (h *NoteHandler) Update(c *gin.Context) {
	userID := c.GetUint("user_id")
	role := c.GetString("role")
	id, _ := strconv.Atoi(c.Param("id"))

	var req struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}

	c.ShouldBindJSON(&req)

	err := h.Service.UpdateNote(uint(id), userID, role, req.Title, req.Content)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "gÃ¼ncellendi"})
}

// DELETE
func (h *NoteHandler) Delete(c *gin.Context) {
	userID := c.GetUint("user_id")
	role := c.GetString("role")
	id, _ := strconv.Atoi(c.Param("id"))

	err := h.Service.DeleteNote(uint(id), userID, role)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "silindi"})
}
