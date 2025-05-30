package web

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
	"github.com/gorilla/sessions"
	"github.com/microcosm-cc/bluemonday"
	"github.com/oodzchen/dizkaz/config"
	i18nc "github.com/oodzchen/dizkaz/i18n"
	"github.com/oodzchen/dizkaz/model"
	"github.com/oodzchen/dizkaz/service"
	"github.com/oodzchen/dizkaz/store"
	"github.com/oodzchen/dizkaz/utils"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"github.com/sergi/go-diff/diffmatchpatch"
)

type Renderer struct {
	tmpl           *template.Template
	sessStore      *sessions.CookieStore
	router         *chi.Mux
	store          *store.Store
	uLogger        *service.UserLogger
	sanitizePolicy *bluemonday.Policy
	i18nCustom     *i18nc.I18nCustom
	srv            *service.Service
	rdb            *redis.Client
	dmp            *diffmatchpatch.DiffMatchPatch
}

func NewRenderer(
	tmpl *template.Template,
	sessStore *sessions.CookieStore,
	router *chi.Mux,
	store *store.Store,
	sp *bluemonday.Policy,
	ic *i18nc.I18nCustom,
	srv *service.Service,
	rdb *redis.Client,
	dmp *diffmatchpatch.DiffMatchPatch,
) *Renderer {
	return &Renderer{
		tmpl,
		sessStore,
		router,
		store,
		srv.UserLogger,
		sp,
		ic,
		srv,
		rdb,
		dmp,
	}
}

func (rd *Renderer) Render(w http.ResponseWriter, r *http.Request, name string, data *model.PageData) {
	rd.doRender(w, r, name, data, http.StatusOK)
}

func (rd *Renderer) ServerErrorp(msg string, err error, w http.ResponseWriter, r *http.Request) {
	rd.Error(msg, err, w, r, http.StatusInternalServerError)
}

func (rd *Renderer) ToLogin(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/login", http.StatusFound)
}

func (rd *Renderer) NotFound(w http.ResponseWriter, r *http.Request) {
	rd.Error("", nil, w, r, http.StatusNotFound)
}

func (rd *Renderer) ServerError(w http.ResponseWriter, r *http.Request) {
	rd.Error("", nil, w, r, http.StatusInternalServerError)
}

func (rd *Renderer) Forbidden(err error, w http.ResponseWriter, r *http.Request) {
	rd.Error("", err, w, r, http.StatusForbidden)
}

func (rd *Renderer) GetLoginedUserId(w http.ResponseWriter, r *http.Request) int {
	userId := rd.Session("one", w, r).GetValue("user_id")
	if userId, ok := userId.(int); ok {
		return userId
	}
	return 0
}

func (rd *Renderer) GetLoginedUserData(r *http.Request) *model.User {
	if userData, ok := r.Context().Value("user_data").(*model.User); ok {
		return userData
	}
	return nil
}

func (rd *Renderer) CheckPermit(r *http.Request, module, action string) bool {
	return rd.srv.Permission.Permit(rd.GetLoginedUserData(r), module, action)
}

func (rd *Renderer) Error(msg string, err error, w http.ResponseWriter, r *http.Request, code int) {
	fmt.Printf("render err: %+v\n", err)
	fmt.Println("msg: ", msg)

	referer := r.Referer()
	refererUrl, _ := url.Parse(referer)
	prevUrl := ""

	if IsRegisterdPage(refererUrl, rd.router) {
		prevUrl = referer
	}

	errText := http.StatusText(code)

	type errPageData struct {
		HttpStatusCode int
		AppErrCode     model.AppErrCode
		ErrCode        int
		ErrText        string
		PrevUrl        string
	}
	var pageData errPageData
	pageData = errPageData{0, 0, 0, "", prevUrl}

	data := &model.PageData{
		Title: errText,
		Data:  &pageData,
	}

	if len(msg) > 0 {
		text := []rune(msg)
		errText += " - " + strings.ToUpper(string(text[:1])) + string(text[1:])
	}
	pageData.ErrText = errText

	if err, ok := err.(model.AppError); ok {
		pageData.ErrCode = int(err.ErrCode)
	}

	pageData.HttpStatusCode = code
	rd.doRender(w, r, "error", data, code)
}

func (rd *Renderer) doRender(w http.ResponseWriter, r *http.Request, name string, data *model.PageData, code int) {
	sess := rd.Session("one", w, r).Raw

	if flashes := sess.Flashes(); len(flashes) > 0 {
		for _, item := range flashes {
			if value, ok := item.(string); ok {
				data.TipMsg = append(data.TipMsg, value)
			}
		}
	}

	data.LoginedUser = rd.GetLoginedUserData(r)

	// fmt.Printf("renderer ui settings: %v\n", data.UISettings)

	if uiSettings, ok := r.Context().Value("ui_settings").(*model.UISettings); ok {
		// fmt.Println("uiSettings in renderer: ", uiSettings)
		data.UISettings = uiSettings
	} else {
		data.UISettings = model.DefaultUiSettings
	}

	// fmt.Printf("renderer ui settings: %v\n", data.UISettings)

	err := sess.Save(r, w)
	if err != nil {
		fmt.Printf("session save error: %+v", err)
	}
	// fmt.Println("currLang: ", rd.i18nCustom.CurrLang)

	data.CSRFField = string(csrf.TemplateField(r))
	data.RoutePath = r.URL.Path
	data.Debug = config.Config.Debug
	data.BrandDomainName = config.Config.BrandDomainName
	data.Slogan = config.Config.Slogan
	data.PermissionEnabledList = rd.srv.Permission.GetEnabledIdList(data.LoginedUser)
	data.RouteRawQuery = r.URL.RawQuery
	data.RouteQuery = r.URL.Query()
	data.Host = config.Config.GetServerURL()
	data.CFSiteKey = config.Config.CloudflareSiteKey

	loginedUseId := rd.GetLoginedUserId(w, r)
	if loginedUseId > 0 {
		messageCount, err := rd.store.Message.UnreadCount(loginedUseId)
		if err != nil {
			fmt.Printf("get message count error: %v", err)
		}
		data.MessageCount = messageCount
	}

	var startTime time.Time
	if v, ok := r.Context().Value("req_duration_start").(time.Time); ok {
		startTime = v
	}

	// fmt.Println("start time before render: ", startTime)

	rd.tmpl = rd.tmpl.Funcs(template.FuncMap{
		"permit": func(module, action string) bool {
			return rd.srv.Permission.Permit(data.LoginedUser, module, action)
		},
		"local":   rd.i18nCustom.LocalTpl,
		"timeAgo": rd.i18nCustom.TimeAgo.Format,
	})

	data.BrandName = rd.Local("BrandName")
	if data.Title != "" {
		data.Title += fmt.Sprintf(" - %s", data.BrandName)
	} else {
		data.Title = fmt.Sprintf("%s", data.BrandName)
	}

	if data.Debug {
		users, _, err := rd.store.User.List(1, 50, true, "", "", "")
		if err != nil {
			fmt.Println("get debug user data error: ", err)
		}
		data.DebugUsers = users

		jsonData, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			HttpError("", errors.WithStack(err), w, http.StatusInternalServerError)
			return
		}

		data.JSONStr = string(jsonData)
	}

	header := w.Header()

	contentSecurity := []string{
		"default-src 'self'",
		"img-src 'self' https://*",
		"style-src 'self' 'unsafe-inline'",
		"child-src  https://challenges.cloudflare.com/",
	}

	if data.Debug {
		contentSecurity = append(contentSecurity, "script-src 'self' 'unsafe-inline' https://challenges.cloudflare.com/")
	} else {
		contentSecurity = append(contentSecurity, "script-src 'self' https://challenges.cloudflare.com/")
	}
	// header.Set("Content-Type", "text/html")
	header.Add("Content-Security-Policy", strings.Join(contentSecurity, ";"))

	// fmt.Println("header", header.Values("Content-Security-Policy"))
	// if code == http.StatusInternalServerError {
	// 	ClearSession(rd.sessStore, w, r)
	// }
	w.WriteHeader(code)

	data.RespStart = startTime
	data.RenderStart = time.Now()

	err = rd.tmpl.ExecuteTemplate(w, name, data)
	if err != nil {
		HttpError("", errors.WithStack(err), w, http.StatusInternalServerError)
	}
}

func (rd *Renderer) isHuman(w http.ResponseWriter, r *http.Request) (bool, error) {
	var cfTurnstileResponse string

	cfTurnstileCookie, _ := r.Cookie("cf_ts_resp")
	if cfTurnstileCookie != nil {
		cfTurnstileResponse, _ = url.QueryUnescape(cfTurnstileCookie.Value)

		// fmt.Println("turnstile response:", cfTurnstileResponse)

		http.SetCookie(w, &http.Cookie{
			Name:    "cf_ts_resp",
			Value:   "",
			Expires: time.Now().Add(-1 * time.Hour),
			Path:    "/",
		})
	}

	// cfTurnstileResponse := r.URL.Query().Get("cf_ts_resp")
	// fmt.Println("turnstile response:", cfTurnstileResponse)

	var isHuman bool
	if !config.Config.Debug && !config.Config.Testing && cfTurnstileResponse != "" {
		client := &http.Client{
			Timeout: 5 * time.Second,
		}

		payload := []byte(`{
"secret":"` + config.Config.CloudflareSecret + `",
"response":"` + cfTurnstileResponse + `",
"remoteip":"` + utils.GetRealIP(r) + `"
    	        }`)

		// fmt.Println("payload: ", string(payload))

		req, err := http.NewRequest("POST", "https://challenges.cloudflare.com/turnstile/v0/siteverify", bytes.NewBuffer(payload))
		if err != nil {
			rd.ServerErrorp("", err, w, r)
			return false, err
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			rd.ServerErrorp("", err, w, r)
			return false, err
		}
		defer resp.Body.Close()

		// fmt.Println("cloudflare turnstile response status: ", resp.Status)

		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)

		// fmt.Println("cloudflare turnstile response body: ", buf.String())
		var trutileVerifyResult TurnstileVerifyResult
		err = json.Unmarshal(buf.Bytes(), &trutileVerifyResult)
		if err != nil {
			rd.ServerErrorp("", err, w, r)
			return false, err
		}

		if !trutileVerifyResult.Success {
			err = errors.New("cloudflare challenge result: " + buf.String())
			rd.Error("", err, w, r, http.StatusBadRequest)
			return false, err
		}

		isHuman = true
	}

	if config.Config.Debug || config.Config.Testing {
		isHuman = true
	}

	return isHuman, nil
}

// func (rd *Renderer) getUserPermittedFrontIds(r *http.Request) []string {
// 	var frontIdList []string
// 	if userData, ok := r.Context().Value("user_data").(*model.User); ok {
// 		if userData.Super {
// 			return rd.srv.Permission.PermissionData.EnabledFrondIdList
// 		}

// 		if userData.Permissions != nil && len(userData.Permissions) > 0 {

// 			for _, item := range userData.Permissions {
// 				frontIdList = append(frontIdList, item.FrontId)
// 			}

// 			return frontIdList
// 		}
// 	}
// 	return nil
// }

func (rd *Renderer) Session(name string, w http.ResponseWriter, r *http.Request) *Session {
	sess, err := rd.sessStore.Get(r, name)
	if err != nil {
		logSessError(name, errors.WithStack(err))
		ClearSession(rd.sessStore, w, r)
	}
	return &Session{rd, sess, w, r}
}

func (rd *Renderer) ToRefererUrl(w http.ResponseWriter, r *http.Request) {
	targetUrl := "/"
	refererUrl, err := url.Parse(r.Referer())
	if err != nil {
		http.Redirect(w, r, targetUrl, http.StatusFound)
		return
	}

	if IsRegisterdPage(refererUrl, rd.router) {
		// fmt.Println("Matched!")
		targetUrl = r.Referer()
	}

	http.Redirect(w, r, targetUrl, http.StatusFound)
}

func (rd *Renderer) ToTargetUrl(w http.ResponseWriter, r *http.Request) {
	target := rd.Session("one", w, r).GetValue("target_url")

	rd.Session("one", w, r).SetValue("target_url", "")

	if targetUrl, ok := target.(string); ok && len(targetUrl) > 0 {
		http.Redirect(w, r, targetUrl, http.StatusFound)
	} else {
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func (rd *Renderer) SavePrevPage(w http.ResponseWriter, r *http.Request) {
	referer := r.Referer()
	refererUrl, _ := url.Parse(referer)

	if refererUrl != nil && IsRegisterdPage(refererUrl, rd.router) {
		rd.Session("one", w, r).SetValue("prev_url", referer)
	}
}

func (rd *Renderer) ToPrevPage(w http.ResponseWriter, r *http.Request) {
	prevPgaeUrl := rd.Session("one", w, r).GetStringValue("prev_url")
	if prevPgaeUrl != "" {
		http.Redirect(w, r, prevPgaeUrl, http.StatusFound)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

func (rd *Renderer) Local(id string, data ...any) string {
	return rd.i18nCustom.LocalTpl(id, data...)
}

func (rd *Renderer) GetPaginationData(r *http.Request) (int, int) {
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("page_size")
	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)

	if page < DefaultPage {
		page = DefaultPage
	}

	if pageSize < DefaultPageSize {
		pageSize = DefaultPageSize
	}

	return page, pageSize
}

// func (rd *Renderer) SaveUserInfo(u *model.User, w http.ResponseWriter, r *http.Request) {
// 	ss := rd.Session("one", w, r)
// 	ss.SetValue("user_id", u.Id)
// 	ss.SetValue("user_name", u.Name)

// 	// gob.Register(model.User{})
// 	ss.SetValue("user_info", *u)
// }

// func (rd *Renderer) GetUserInfo(w http.ResponseWriter, r *http.Request) (*model.User, error) {
// 	ss := rd.Session("one", w, r)
// 	u := ss.GetValue("user_info")
// 	if v, ok := u.(model.User); ok {
// 		return &v, nil
// 	}
// 	return nil, errors.New("no user info stored in cookie")
// }

type Session struct {
	rd  *Renderer
	Raw *sessions.Session
	w   http.ResponseWriter
	r   *http.Request
}

// Get value from *sessions.Session.Values
func (ss *Session) GetValue(key string) any {
	return ss.Raw.Values[key]
}

func (ss *Session) GetStringValue(key string) string {
	val := ss.GetValue(key)
	if v, ok := val.(string); ok {
		return v
	}
	return ""
}

// Set data to *sessons.Session.Values and auto save, handle save error
func (ss *Session) SetValue(key string, val any) {
	ss.Raw.Values[key] = val

	ss.Raw.Options.HttpOnly = true
	ss.Raw.Options.Secure = !utils.IsDebug()
	ss.Raw.Options.SameSite = http.SameSiteLaxMode
	ss.Raw.Options.Path = "/"

	err := ss.Raw.Save(ss.r, ss.w)
	if err != nil {
		fmt.Println("ss.SetValue save session error: ", err)
		// if ss.w != nil {
		// 	ss.rd.Error("", errors.WithStack(err), ss.w, ss.r, http.StatusInternalServerError)
		// }
	}
}

func (ss *Session) Flash(data any, vars ...string) {
	ss.Raw.AddFlash(data, vars...)
	err := ss.Raw.Save(ss.r, ss.w)
	if err != nil {
		fmt.Println("add flash save session error: ", err)
	}
}
