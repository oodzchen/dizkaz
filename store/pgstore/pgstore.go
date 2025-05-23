package pgstore

import (
	"context"
	"errors"
)

type PGStore struct {
	Activity   *Activity
	Article    *Article
	Message    *Message
	Permission *Permission
	Role       *Role
	User       *User
	Category   *Category
}

type DBConfig struct {
	DSN string
}

var pgDB *DB

const DefaultPage = 1
const DefaultPageSize = 50

func New(config *DBConfig) *PGStore {
	pgDB = &DB{config.DSN, nil}
	return &PGStore{}
}

// func (pg *PGStore) ConfigDB(dsn string) *PGStore {
// 	pgDB = &DB{dsn, nil}
// 	return &PGStore{}
// }

func (pg *PGStore) ConnectDB() error {
	err := CheckDB(true)
	if err != nil {
		return err
	}
	return pgDB.Connect()
}

func (pg *PGStore) CloseDB() {
	pgDB.Close()
}

func (pg *PGStore) Ping(ctx context.Context) error {
	return pgDB.Pool.Ping(ctx)
}

func CheckDB(beforeConnect bool) error {
	if pgDB == nil {
		return errors.New("Database config is required")
	}

	if beforeConnect {
		return nil
	}

	if pgDB.Pool == nil {
		return errors.New("No database connection")
	}

	return nil
}

func (pg *PGStore) InitModules() error {
	err := CheckDB(false)
	if err != nil {
		return err
	}

	pg.Activity = &Activity{pgDB.Pool}
	pg.Article = &Article{pgDB.Pool}
	pg.Message = &Message{pgDB.Pool}
	pg.Permission = &Permission{pgDB.Pool}
	pg.Role = &Role{pgDB.Pool}
	pg.User = &User{pgDB.Pool}
	pg.Category = &Category{pgDB.Pool}

	return nil
}

// func (pg *PGStore) NewArticleStore() (any, error) {
// 	err := CheckDB(false)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &Article{pgDB.Pool}, nil
// }

// func (pg *PGStore) NewUserStore() (any, error) {
// 	err := CheckDB(false)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &User{pgDB.Pool}, nil
// }

// func (pg *PGStore) NewRoleStore() (any, error) {
// 	err := CheckDB(false)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &Role{pgDB.Pool}, nil
// }

// func (pg *PGStore) NewPermissionStore() (any, error) {
// 	err := CheckDB(false)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &Permission{pgDB.Pool}, nil
// }

// func (pg *PGStore) NewActivity() (any, error) {
// 	err := CheckDB(false)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &Activity{pgDB.Pool}, nil
// }

// func (pg *PGStore) NewMessage() (any, error) {
// 	err := CheckDB(false)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &Message{pgDB.Pool}, nil
// }
