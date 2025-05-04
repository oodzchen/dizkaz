package main

import (
	"encoding/gob"
	"fmt"
	"net/http"
	"os"
	"path"
	"text/template"
	"time"

	"github.com/Masterminds/sprig/v3"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
	"github.com/gorilla/csrf"
	"github.com/gorilla/sessions"
	"github.com/microcosm-cc/bluemonday"
	"github.com/oodzchen/dizkaz/config"
	i18nc "github.com/oodzchen/dizkaz/i18n"
	mdw "github.com/oodzchen/dizkaz/middleware"
	"github.com/oodzchen/dizkaz/model"
	"github.com/oodzchen/dizkaz/service"
	"github.com/oodzchen/dizkaz/store"
	"github.com/oodzchen/dizkaz/utils"
	"github.com/oodzchen/dizkaz/web"
	"github.com/oschwald/geoip2-golang"
	"github.com/redis/go-redis/v9"
	"github.com/sergi/go-diff/diffmatchpatch"
)

func Server(
	sessSecret string,
	csrfSecret string,
	store *store.Store,
	permisisonSrv *service.Permission,
	sanitizePolicy *bluemonday.Policy,
	i18nCustom *i18nc.I18nCustom,
	rdb *redis.Client,
	mail *service.Mail,
	geoDB *geoip2.Reader,
) http.Handler {
	wd, _ := os.Getwd()
	// fmt.Println("work directory: ", wd)
	// fmt.Println("templates directory: ", path.Join(wd, "./views/*.tmpl"))
	tmplPath := path.Join(wd, "./views/*.tmpl")
	tmplFuncs := template.FuncMap{
		"permit": func(module, action string) bool {
			return permisisonSrv.Permit(nil, module, action)
		},
		"local":   i18nCustom.LocalTpl,
		"timeAgo": i18nCustom.TimeAgo.Format,
	}

	baseTmpl := template.New("base").Funcs(TmplFuncs).Funcs(tmplFuncs).Funcs(sprig.FuncMap())
	baseTmpl = template.Must(baseTmpl.ParseGlob(tmplPath))

	r := chi.NewRouter()
	r.Use(mdw.RequestDuration)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.AllowContentEncoding("default", "gzip"))
	r.Use(middleware.AllowContentType("application/x-www-form-urlencoded"))
	r.Use(middleware.Compress(5, "text/html", "text/css", "text/plain", "text/javascript"))
	r.Use(middleware.GetHead)
	// r.Use(middleware.RedirectSlashes)
	r.Use(mdw.CreateGeoDetect(geoDB))

	gob.Register(model.Lang(""))

	sessStore := sessions.NewCookieStore([]byte(sessSecret))
	sessStore.Options.HttpOnly = true
	sessStore.Options.Secure = !utils.IsDebug()
	sessStore.Options.SameSite = http.SameSiteLaxMode

	userLogger := &service.UserLogger{
		Store: store,
	}

	settingsManager := &service.SettingsManager{
		Rdb:      rdb,
		LifeTime: service.DefaultSettingsLifeTime,
	}
	srv := &service.Service{
		Article: &service.Article{
			Store:         store,
			SantizePolicy: sanitizePolicy,
		},
		User: &service.User{
			Store:         store,
			SantizePolicy: sanitizePolicy,
		},
		Permission: permisisonSrv,
		UserLogger: userLogger,
		Verifier: &service.Verifier{
			CodeLifeTime: service.DefaultCodeLifeTime,
			Rdb:          rdb,
		},
		Mail:            mail,
		SettingsManager: settingsManager,
	}

	dmp := diffmatchpatch.New()

	renderer := web.NewRenderer(
		baseTmpl,
		sessStore,
		r,
		store,
		sanitizePolicy,
		i18nCustom,
		srv,
		rdb,
		dmp,
	)

	r.Use(mdw.FetchUserData(store, sessStore, permisisonSrv, renderer))
	r.Use(mdw.CreateUISettingsMiddleware(sessStore, settingsManager, i18nCustom))

	articleResource := web.NewArticleResource(renderer)
	userResource := web.NewUserResource(renderer)
	mainResource := web.NewMainResource(renderer, articleResource)
	manageResource := web.NewManageResource(renderer, userResource)
	rssResource := web.NewRSSResource(renderer, articleResource)

	rateLimit := 100
	if utils.IsDebug() {
		rateLimit = 10000
	}

	r.Use(httprate.Limit(
		rateLimit,
		1*time.Minute,
		httprate.WithKeyByIP(),
		httprate.WithLimitHandler(func(w http.ResponseWriter, r *http.Request) {
			mainResource.Error("", nil, w, r, http.StatusTooManyRequests)
		}),
	))

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		mainResource.Error("", nil, w, r, http.StatusNotFound)
	})
	r.HandleFunc("/403", func(w http.ResponseWriter, r *http.Request) {
		mainResource.Error("", nil, w, r, http.StatusForbidden)
	})
	r.HandleFunc("/500", func(w http.ResponseWriter, r *http.Request) {
		mainResource.Error("", nil, w, r, http.StatusInternalServerError)
	})

	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		mainResource.Error("", nil, w, r, http.StatusMethodNotAllowed)
	})

	if config.Config.Debug {
		r.Mount("/debug", middleware.Profiler())
	}

	// FileServer(r, "/static", http.Dir("./static"))
	// r.Get("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
	// 	http.Redirect(w, r, "/static/favicon.ico", http.StatusFound)
	// })

	r.Mount("/", mainResource.Routes())
	r.Mount("/articles", articleResource.Routes())
	r.Mount("/u/{username}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username := chi.URLParam(r, "username")
		http.Redirect(w, r, fmt.Sprintf("/users/%s", username), http.StatusFound)
	}))
	r.Mount("/users", userResource.Routes())
	r.Mount("/manage", manageResource.Routes())
	r.Mount("/feed", rssResource.Routes())

	// chi.Walk(r, func(method, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
	// 	////
	// 	fmt.Println("walk:", method, route)
	// 	return nil
	// })

	CSRF := csrf.Protect([]byte(csrfSecret),
		csrf.FieldName("tk"),
		csrf.CookieName("sc"),
		csrf.HttpOnly(true),
		csrf.Secure(!utils.IsDebug()),
		csrf.Path("/"),
		// csrf.ErrorHandler(r),
	)
	return CSRF(r)
}
