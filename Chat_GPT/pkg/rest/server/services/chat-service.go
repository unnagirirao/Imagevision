package services

import (
	"github.com/unnagirirao/Imagevision/chat_gpt/pkg/rest/server/daos"
	"github.com/unnagirirao/Imagevision/chat_gpt/pkg/rest/server/models"
)

type ChatService struct {
	chatDao *daos.ChatDao
}

func NewChatService() (*ChatService, error) {
	chatDao, err := daos.NewChatDao()
	if err != nil {
		return nil, err
	}
	return &ChatService{
		chatDao: chatDao,
	}, nil
}

func (chatService *ChatService) CreateChat(chat *models.Chat) (*models.Chat, error) {
	return chatService.chatDao.CreateChat(chat)
}

func (chatService *ChatService) UpdateChat(id int64, chat *models.Chat) (*models.Chat, error) {
	return chatService.chatDao.UpdateChat(id, chat)
}

func (chatService *ChatService) DeleteChat(id int64) error {
	return chatService.chatDao.DeleteChat(id)
}

func (chatService *ChatService) ListChats() ([]*models.Chat, error) {
	return chatService.chatDao.ListChats()
}

func (chatService *ChatService) GetChat(id int64) (*models.Chat, error) {
	return chatService.chatDao.GetChat(id)
}
