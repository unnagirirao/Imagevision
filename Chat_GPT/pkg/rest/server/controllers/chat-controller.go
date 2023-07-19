package controllers

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/unnagirirao/Imagevision/chat_gpt/pkg/rest/server/models"
	"github.com/unnagirirao/Imagevision/chat_gpt/pkg/rest/server/services"
	"net/http"
	"strconv"
)

type ChatController struct {
	chatService *services.ChatService
}

func NewChatController() (*ChatController, error) {
	chatService, err := services.NewChatService()
	if err != nil {
		return nil, err
	}
	return &ChatController{
		chatService: chatService,
	}, nil
}

func (chatController *ChatController) CreateChat(context *gin.Context) {
	// validate input
	var input models.Chat
	if err := context.ShouldBindJSON(&input); err != nil {
		log.Error(err)
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// trigger chat creation
	if _, err := chatController.chatService.CreateChat(&input); err != nil {
		log.Error(err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"message": "Chat created successfully"})
}

func (chatController *ChatController) UpdateChat(context *gin.Context) {
	// validate input
	var input models.Chat
	if err := context.ShouldBindJSON(&input); err != nil {
		log.Error(err)
		context.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	id, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		log.Error(err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// trigger chat update
	if _, err := chatController.chatService.UpdateChat(id, &input); err != nil {
		log.Error(err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Chat updated successfully"})
}

func (chatController *ChatController) FetchChat(context *gin.Context) {
	id, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		log.Error(err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// trigger chat fetching
	chat, err := chatController.chatService.GetChat(id)
	if err != nil {
		log.Error(err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, chat)
}

func (chatController *ChatController) DeleteChat(context *gin.Context) {
	id, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		log.Error(err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// trigger chat deletion
	if err := chatController.chatService.DeleteChat(id); err != nil {
		log.Error(err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"message": "Chat deleted successfully",
	})
}

func (chatController *ChatController) ListChats(context *gin.Context) {
	// trigger all chats fetching
	chats, err := chatController.chatService.ListChats()
	if err != nil {
		log.Error(err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, chats)
}

func (*ChatController) PatchChat(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{
		"message": "PATCH",
	})
}

func (*ChatController) OptionsChat(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{
		"message": "OPTIONS",
	})
}

func (*ChatController) HeadChat(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{
		"message": "HEAD",
	})
}
