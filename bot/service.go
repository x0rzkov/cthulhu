package bot

import (
	"context"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"

	"cthulhu/telegram"
)

type Service interface {
	Update(ctx context.Context, req *telegram.Update) error
	GetToken() Token
}

type service struct {
	Logger log.Logger
	Token  Token
	Config Config
}

func NewService(logger log.Logger, config Config, token Token) *service {
	return &service{
		Logger: logger,
		Token:  token,
		Config: config,
	}
}

func (s *service) GetToken() Token {
	return s.Token
}

func (s *service) Update(ctx context.Context, updateReq *telegram.Update) error {
	if updateReq.Message == nil {
		return nil
	}

	if !s.checkOrigin(ctx, updateReq) {
		level.Error(s.Logger).Log("msg", "group is not part of the network")
		return nil
	}

	if updateReq.Message.NewChatMembers != nil {
		s.handleNewUsers(ctx, updateReq)
	}

	if command := updateReq.Message.Command(); command != "" {
		level.Info(s.Logger).Log("msg", "received new command", "command", command)
		switch command {
		case banCommand:
			return s.handleBan(ctx, updateReq)
		case unbanCommand:
			return s.handleUnban(ctx, updateReq)
		}
	}

	return s.handleCrossposts(ctx, updateReq)
}

func (s *service) checkOrigin(ctx context.Context, updateReq *telegram.Update) bool {
	var originID int64 = updateReq.Message.Chat.ID

	for _, g := range s.Config.Bot.AccessControl.Groups {
		if g.Group.ID == originID {
			return true
		}
	}
	return false
}

func (s *service) handleNewUsers(ctx context.Context, updateReq *telegram.Update) {
	var originID int64 = updateReq.Message.Chat.ID

	for _, g := range s.Config.Bot.AccessControl.Groups {
		if g.Group.ID == originID {
			if g.Group.WelcomeMessage != "" {
				for range *updateReq.Message.NewChatMembers {
					telegram.Reply(ctx, string(s.Token), g.Group.ID, g.Group.WelcomeMessage, updateReq.Message.MessageID)
				}
			}
		}
	}
}
