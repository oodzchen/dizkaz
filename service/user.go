package service

import (
	"errors"

	"github.com/microcosm-cc/bluemonday"
	"github.com/oodzchen/dizkaz/model"
	"github.com/oodzchen/dizkaz/store"
)

type UserListType string

const (
	UserListAll        UserListType = "all"
	UserListSaved                   = "saved"
	UserListArticle                 = "article"
	UserListReply                   = "reply"
	UserListActivity                = "activity"
	UserListSubscribed              = "subscribed"
	UserListVoteUp                  = "vote_up"
)

var AuthRequiedUserTabMap = map[UserListType]bool{
	UserListSaved:      true,
	UserListSubscribed: true,
	UserListActivity:   true,
	UserListVoteUp:     true,
}

func CheckUserTabAuthRequired(tab UserListType) bool {
	return AuthRequiedUserTabMap[tab]
}

type User struct {
	Store         *store.Store
	SantizePolicy *bluemonday.Policy
}

func (u *User) Register(email string, password string, name string) (int, error) {
	if len(password) == 0 {
		return 0, errors.New("lack of password")
	}

	user := &model.User{
		Email: email,
		Name:  name,
	}
	err := user.Valid(false)
	if err != nil {
		return 0, err
	}

	return u.Store.User.Create(email, password, name, string(model.DefaultUserRoleCommon))
}

func (u *User) GetPosts(username string, listType UserListType) ([]*model.Article, error) {
	// fmt.Println("user tab:", listType)
	switch listType {
	case UserListSaved:
		return u.Store.User.GetSavedPosts(username)
	case UserListSubscribed:
		return u.Store.User.GetSubscribedPosts(username)
	case UserListVoteUp:
		return u.Store.User.GetVotedPosts(username, model.VoteTypeUp)
	default:
		return u.Store.User.GetPosts(username, string(listType))
	}
}
