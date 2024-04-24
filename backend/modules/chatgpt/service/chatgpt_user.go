package service

import (
	"backend/modules/chatgpt/model"
	"github.com/cool-team-official/cool-admin-go/cool"
	"github.com/gogf/gf/v2/frame/g"
)

type ChatgptUserService struct {
	*cool.Service
}

func NewChatgptUserService() *ChatgptUserService {
	return &ChatgptUserService{
		&cool.Service{
			Model: model.NewChatgptUser(),
			PageQueryOp: &cool.QueryOp{
				FieldEQ:      []string{"usertoken"},
				KeyWordField: []string{"usertoken"},
			},
		},
	}
}

func QueryUserByToken(ctx g.Ctx, userToken string) *model.ChatgptUser {
	userRecord, err := cool.DBM(model.NewChatgptUser()).Where("userToken=?", userToken).WhereNull("deleted_at").One()
	if err != nil || userRecord == nil {
		return nil
	}
	result := model.NewChatgptUser()
	userRecord.Struct(result)
	return result
}
