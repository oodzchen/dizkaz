package store

// func connectDB(appCfg *config.AppConfig) (*pgstore.PGStore, error) {
// 	pg := pgstore.New(&pgstore.DBConfig{
// 		DSN: appCfg.DB.GetDSN(),
// 	})

// 	err := pg.ConnectDB()
// 	if err != nil {
// 		return nil, err
// 	}

// 	return pg, nil
// }

// func registerNewUser(store *Store, appCfg *config.AppConfig) (int, error) {
// 	user := mt.GenUser()
// 	pwd, _ := bcrypt.GenerateFromPassword([]byte(appCfg.DB.UserDefaultPassword), 10)
// 	return store.User.Create(user.Email, string(pwd), user.Name, "common_user")
// }

// func createNewArticle(store *Store, userId int) (int, error) {
// 	article := mt.GenArticle()
// 	return store.Article.Create(article.Title, "", article.Content, userId, 0, "general", time.Now(), false)
// }

// func TestArticleVote(t *testing.T) {
// 	appCfg, err := config.NewTest()
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// fmt.Println("testing db config:", appCfg.DB)

// 	pg, err := connectDB(appCfg)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer pg.CloseDB()

// 	err = pg.InitModules()
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	store := New(pg.Article, pg.User, pg.Role, pg.Permission, pg.Activity, pg.Message, pg.Category)

// 	uId, err := registerNewUser(store, appCfg)
// 	mt.LogFailed(err)

// 	aId, err := createNewArticle(store, uId)
// 	mt.LogFailed(err)

// 	uBId, err := registerNewUser(store, appCfg)
// 	mt.LogFailed(err)

// 	// fmt.Println("uId: ", uId)
// 	// fmt.Println("aId: ", aId)
// 	// fmt.Println("uBId: ", uBId)

// 	t.Run("Vote up", func(t *testing.T) {
// 		_, err = store.Article.ToggleVote(aId, uBId, "up")
// 		if err != nil {
// 			t.Errorf("should vote up success but got %v", err)
// 		}
// 	})

// 	t.Run("Change vote to down", func(t *testing.T) {
// 		_, err = store.Article.ToggleVote(aId, uBId, "down")
// 		if err != nil {
// 			t.Errorf("should vote down success but got %v", err)
// 		}
// 	})

// 	t.Run("Revoke vote", func(t *testing.T) {
// 		_, err = store.Article.ToggleVote(aId, uBId, "down")
// 		if err != nil {
// 			t.Errorf("should revoke vote success but got %v", err)
// 		}
// 	})
// }

// func TestArticleCheckVote(t *testing.T) {
// 	appCfg, err := config.NewTest()
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	pg, err := connectDB(appCfg)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer pg.CloseDB()

// 	err = pg.InitModules()
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	store := &Store{
// 		Activity:   pg.Activity,
// 		Article:    pg.Article,
// 		Message:    pg.Message,
// 		Permission: pg.Permission,
// 		Role:       pg.Role,
// 		User:       pg.User,
// 	}

// 	uId, err := registerNewUser(store, appCfg)
// 	mt.LogFailed(err)

// 	aId, err := createNewArticle(store, uId)
// 	mt.LogFailed(err)

// 	t.Run("Check unvote article", func(t *testing.T) {
// 		err, voteType := store.Article.VoteCheck(aId, uId)
// 		if err != nil {
// 			t.Errorf("vote check error %v", err)
// 		}
// 		// fmt.Println("vote check result: ", err)
// 		if voteType != "up" {
// 			t.Errorf("vote check type should be %s, but got %v", "up", voteType)
// 		}
// 	})

// 	t.Run("Check voted article", func(t *testing.T) {
// 		_, err = store.Article.ToggleVote(aId, uId, "down")
// 		if err != nil {
// 			t.Errorf("vote down article failed: %v", err)
// 		}

// 		err, vt := store.Article.VoteCheck(aId, uId)
// 		if err != nil {
// 			t.Errorf("should check with no error, but got: %v", err)
// 		}

// 		if vt != "down" {
// 			t.Errorf("should get vote type as 'down' but got %s", vt)
// 		}
// 	})
// }
