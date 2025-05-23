package web

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgconn"
	mdw "github.com/oodzchen/dizkaz/middleware"
	"github.com/oodzchen/dizkaz/model"
	"github.com/pkg/errors"
)

type ManageResource struct {
	*Renderer
	ur *UserResource
}

func NewManageResource(renderer *Renderer, ur *UserResource) *ManageResource {
	return &ManageResource{
		renderer,
		ur,
	}
}

func (mr *ManageResource) Routes() http.Handler {
	rt := chi.NewRouter()

	rt.With(mdw.AuthCheck(mr.sessStore), mdw.PermitCheck(mr.srv.Permission, []string{
		"manage.access",
	}, mr)).Route("/", func(r chi.Router) {
		r.With(mdw.PermitCheck(mr.srv.Permission, []string{
			"permission.access",
		}, mr)).Group(func(r chi.Router) {
			r.Get("/", mr.PermissionListPage)

			r.Route("/permissions", func(r chi.Router) {
				r.Get("/", mr.PermissionListPage)
				// r.Post("/", mr.PermissionSubmit)
				// r.Get("/new", mr.PermissionCreatePage)
			})
		})

		// r.With(middlewares ...func(http.Handler) http.Handler)
		r.With(mdw.PermitCheck(mr.srv.Permission, []string{
			"role.access",
		}, mr)).Route("/roles", func(r chi.Router) {
			r.Get("/", mr.RoleListPage)

			r.With(mdw.PermitCheck(mr.srv.Permission, []string{
				"role.add",
			}, mr)).Group(func(r chi.Router) {
				r.With(mdw.UserLogger(
					mr.uLogger, model.AcTypeManage, model.AcActionAddRole, model.AcModelEmpty, mdw.ULogEmpty),
				).Post("/", mr.RoleSubmit)
				r.Get("/new", mr.RoleCreatePage)
			})

			r.With(mdw.PermitCheck(mr.srv.Permission, []string{
				"role.edit",
			}, mr)).Group(func(r chi.Router) {
				r.Get("/{roleId}/edit", mr.RoleEditPage)
				r.With(mdw.UserLogger(
					mr.uLogger, model.AcTypeManage, model.AcActionEditRole, model.AcModelRole, mdw.ULogRoleId),
				).Post("/{roleId}/edit", mr.RoleEditSubmit)
			})
		})

		// r.Get("/roles", mr.RoleListPage)
		// r.Get("/users", mr.ur.List)
		r.With(mdw.PermitCheck(mr.srv.Permission, []string{
			"manage.access",
			"user.list_access",
		}, mr)).Get("/users", mr.ur.List)

		r.With(mdw.PermitCheck(mr.srv.Permission, []string{
			"activity.access",
		}, mr)).Get("/activities", mr.ActivityList)

		r.Get("/trash", mr.TrashPage)

		rootDir, _ := os.Getwd()
		manageStaticPath := filepath.Join(rootDir, "/manage_static")
		fmt.Println("manage static path:", manageStaticPath)
		err := FileServer(r, "/static", http.Dir(manageStaticPath))
		if err != nil {
			fmt.Println("manage static route init error:", err)
		}
	})

	return rt
}

func FileServer(r chi.Router, path string, root http.FileSystem) error {
	if strings.ContainsAny(path, "{}*") {
		return errors.New("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})

	return nil
}

func (mr *ManageResource) PermissionListPage(w http.ResponseWriter, r *http.Request) {
	mr.handlePermissionList(w, r, PermissionPageList)
}

func (mr *ManageResource) PermissionCreatePage(w http.ResponseWriter, r *http.Request) {
	mr.handlePermissionList(w, r, PermissionPageCreate)
}

type PermissionPageType string

const (
	PermissionPageList   PermissionPageType = "list"
	PermissionPageCreate                    = "create"
)

func (mr *ManageResource) handlePermissionList(w http.ResponseWriter, r *http.Request, pageType PermissionPageType) {
	r.ParseForm()

	paramPage := r.Form.Get("page")

	tab := r.Form.Get("tab")
	// fmt.Println("paramPage:", paramPage)
	if !mr.srv.Permission.PermissionData.Valid(tab) {
		tab = "all"
	}

	page, err := strconv.Atoi(paramPage)
	if err != nil {
		// fmt.Printf("page err %v\n", err)
		page = 1
	}

	pageSize, err := strconv.Atoi(r.Form.Get("page_size"))
	if err != nil {
		pageSize = 999
	}

	list, err := mr.store.Permission.List(page, pageSize, tab)
	if err != nil {
		mr.Error("", err, w, r, http.StatusInternalServerError)
	}

	total := len(list)

	type PermissionListPage struct {
		List          []*model.Permission
		Total         int
		CurrPage      int
		TotalPage     int
		PageSize      int
		PageType      string
		ModuleOptions []string
		CurrTab       string
	}

	title := mr.Local("List", "Name", mr.Local("Permission", "Count", 1))
	breadCrumbs := []*model.BreadCrumb{
		{
			Path: "/manage/permissions",
			Name: title,
		},
	}

	if pageType == PermissionPageCreate {
		title = mr.Local("AddItem", "Name", mr.Local("Permission", "Count", 1))
		breadCrumbs = append(breadCrumbs, &model.BreadCrumb{
			Path: "",
			Name: title,
		})
	}

	if pageType == PermissionPageList {
		mr.SavePrevPage(w, r)
	}

	mr.Render(w, r, "permission_list", &model.PageData{
		Title: title,
		Data: &PermissionListPage{
			list,
			total,
			page,
			CeilInt(total, pageSize),
			pageSize,
			string(pageType),
			mr.srv.Permission.PermissionData.GetModuleList(),
			tab,
		},
		BreadCrumbs: breadCrumbs,
	})
}

func (mr *ManageResource) PermissionSubmit(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	module := r.Form.Get("module")
	frontId := r.Form.Get("front_id")
	name := r.Form.Get("name")

	permission := &model.Permission{
		Module:  module,
		FrontId: frontId,
		Name:    name,
	}

	permission.TrimSpace()

	moduleValid := mr.srv.Permission.PermissionData.Valid(module)
	if !moduleValid {
		mr.Error("module dose not exist", errors.New("module dose not exist"), w, r, http.StatusBadRequest)
		return
	}

	err := permission.Valid()
	if err != nil {
		mr.Error(err.Error(), err, w, r, http.StatusBadRequest)
		return
	}

	_, err = mr.store.Permission.Create(string(permission.Module), permission.FrontId, permission.Name)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == PGErrUniqueViolation {
			mr.Error("the permission already existed", err, w, r, http.StatusBadRequest)
		} else {
			mr.Error("", errors.WithStack(err), w, r, http.StatusInternalServerError)
		}

		return
	}

	mr.Session("one", w, r).Flash("Add permission successfully")
	// http.Redirect(w, r, "/manage/permissions", http.StatusFound)
	mr.ToPrevPage(w, r)
}

func (mr *ManageResource) RoleListPage(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	paramPage := r.Form.Get("page")

	page, err := strconv.Atoi(paramPage)
	if err != nil {
		// fmt.Printf("page err %v\n", err)
		page = 1
	}

	pageSize, err := strconv.Atoi(r.Form.Get("page_size"))
	if err != nil {
		pageSize = 999
	}

	list, err := mr.store.Role.List(page, pageSize)
	if err != nil {
		mr.Error("", err, w, r, http.StatusInternalServerError)
	}

	for _, item := range list {
		item.FormattedPermissions = formatPermissionList(item.Permissions, mr.srv.Permission.PermissionData.GetModuleList())
	}

	total := len(list)

	type RoleListPageData struct {
		List      []*model.Role
		Total     int
		CurrPage  int
		TotalPage int
		PageSize  int
	}

	title := mr.Local("List", "Name", mr.Local("Role", "Count", 1))
	breadCrumbs := []*model.BreadCrumb{
		{
			Path: "/manage/roles",
			Name: title,
		},
	}

	mr.SavePrevPage(w, r)

	mr.Render(w, r, "role_list", &model.PageData{
		Title: title,
		Data: &RoleListPageData{
			list,
			total,
			page,
			CeilInt(total, pageSize),
			pageSize,
		},
		BreadCrumbs: breadCrumbs,
	})
}

type RoleFormPageType string

const (
	RoleFormPageAdd  RoleFormPageType = "add"
	RoleFormPageEdit                  = "edit"
)

type RoleFormPageData struct {
	Role                 *model.Role
	RolePermissionIdList []int
	PermissionList       []*model.PermissionListItem
	PageType             RoleFormPageType
}

func (mr *ManageResource) RoleCreatePage(w http.ResponseWriter, r *http.Request) {
	// permissionList, err := mr.store.Permission.List(1, 999, "all")
	// if err != nil {
	// 	mr.Error("", err, w, r, http.StatusInternalServerError)
	// 	return
	// }

	// type RoleCreatePageData struct {
	// 	PermissionList []*model.PermissionListItem
	// }

	filteredPermissionList := mr.getFilteredPermissionList(w, r)

	formattedPermissionList := formatPermissionList(filteredPermissionList, mr.srv.Permission.PermissionData.GetModuleList())

	upLevelTitle := mr.Local("List", "Name", mr.Local("Role", "Count", 1))
	title := mr.Local("AddItem", "Name", mr.Local("Role", "Count", 1))
	breadCrumbs := []*model.BreadCrumb{
		{
			Path: "/manage/roles",
			Name: upLevelTitle,
		},
		{
			Path: "",
			Name: title,
		},
	}

	mr.Render(w, r, "role_form", &model.PageData{
		Title: title,
		Data: &RoleFormPageData{
			PermissionList: formattedPermissionList,
			PageType:       RoleFormPageAdd,
		},
		BreadCrumbs: breadCrumbs,
	})
}

func formatPermissionList(rawList []*model.Permission, moduleOptions []string) []*model.PermissionListItem {
	var list []*model.PermissionListItem
	listMap := make(map[string][]*model.Permission)

	for _, item := range rawList {
		if mList, ok := listMap[item.Module]; !ok {
			listMap[item.Module] = []*model.Permission{item}
		} else {
			listMap[item.Module] = append(mList, item)
		}
	}

	for _, module := range moduleOptions {
		if pList, ok := listMap[module]; ok {
			list = append(list, &model.PermissionListItem{
				Module: module,
				List:   pList,
			})
		}
	}

	for _, item := range list {
		sort.Slice(item.List, func(i, j int) bool {
			return rune(item.List[i].Name[0]) < rune(item.List[j].Name[0])
		})
	}

	return list
}

func (mr *ManageResource) RoleSubmit(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	frontId := r.Form.Get("front_id")
	name := r.Form.Get("name")
	permissions := r.Form["permissions"]

	// fmt.Println("permissions: ", permissions)

	role := &model.Role{
		FrontId: frontId,
		Name:    name,
	}

	role.TrimSpace()

	err := role.Valid(false)
	if err != nil {
		mr.Error(err.Error(), err, w, r, http.StatusBadRequest)
		return
	}

	filteredPermissionList := mr.getFilteredPermissionList(w, r)
	filteredPermissionMap := make(map[string]bool)
	for _, pItem := range filteredPermissionList {
		filteredPermissionMap[pItem.FrontId] = true
	}

	var permissionFrontIds []string
	for _, pId := range permissions {
		if _, ok := filteredPermissionMap[pId]; ok {
			permissionFrontIds = append(permissionFrontIds, pId)
		}
	}

	// var permissionIds []int

	// for _, idStr := range permissions {
	// 	id, err := strconv.Atoi(idStr)
	// 	if err == nil {
	// 		permissionIds = append(permissionIds, id)
	// 	}
	// }

	// _, err = mr.store.Role.Create(role.FrontId, role.Name, permissionIds)

	// fmt.Println("permissionFrontIds: ", permissionFrontIds)
	_, err = mr.store.Role.CreateWithFrontId(role.FrontId, role.Name, permissionFrontIds)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == PGErrUniqueViolation {
			mr.Error("the role is existing", err, w, r, http.StatusBadRequest)
		} else {
			mr.Error("", errors.WithStack(err), w, r, http.StatusInternalServerError)
		}

		return
	}

	mr.Session("one", w, r).Flash("Add role successfully")
	http.Redirect(w, r, "/manage/roles", http.StatusFound)
	// mr.ToPrevPage(w, r)
}

// Fitler permmmsions that not belong to current user
func (mr *ManageResource) getFilteredPermissionList(w http.ResponseWriter, r *http.Request) []*model.Permission {
	permissionList, err := mr.store.Permission.List(1, 999, "all")
	if err != nil {
		mr.Error("", err, w, r, http.StatusInternalServerError)
		return nil
	}

	// userPermittedIdList := mr.Session("one", w, r).GetValue("user_permitted_id_list")
	userPermittedIdList := mr.srv.Permission.GetEnabledIdList(mr.GetLoginedUserData(r))

	userPermittedIdMap := make(map[string]bool)
	for _, frontId := range userPermittedIdList {
		userPermittedIdMap[frontId] = true
	}

	// if uPList, ok := userPermittedIdList.([]string); ok {
	// 	for _, frontId := range uPList {
	// 		userPermittedIdMap[frontId] = true
	// 	}
	// }

	var filteredPermissionList []*model.Permission
	for _, item := range permissionList {
		if _, ok := userPermittedIdMap[item.FrontId]; ok {
			filteredPermissionList = append(filteredPermissionList, item)
		}
	}
	return filteredPermissionList
}

func (mr *ManageResource) RoleEditPage(w http.ResponseWriter, r *http.Request) {
	roleIdStr := chi.URLParam(r, "roleId")
	// fmt.Println("roleId: ", roleIdStr)

	roleId, err := strconv.Atoi(roleIdStr)
	if err != nil {
		mr.Error("", err, w, r, http.StatusBadRequest)
		return
	}

	role, err := mr.store.Role.Item(roleId)
	if err != nil {
		mr.Error("", err, w, r, http.StatusInternalServerError)
		return
	}

	filteredPermissionList := mr.getFilteredPermissionList(w, r)
	formattedPermissionList := formatPermissionList(filteredPermissionList, mr.srv.Permission.PermissionData.GetModuleList())

	var rolePermissionIdList []int
	if role.Permissions != nil {
		for _, item := range role.Permissions {
			rolePermissionIdList = append(rolePermissionIdList, item.Id)
		}
	}

	upLevelTitle := mr.Local("List", "Name", mr.Local("Role", "Count", 1))
	title := mr.Local("EditItem", "Name", mr.Local("Role", "Count", 1))
	breadCrumbs := []*model.BreadCrumb{
		{
			Path: "/manage/roles",
			Name: upLevelTitle,
		},
		{
			Path: "",
			Name: title,
		},
	}

	mr.Render(w, r, "role_form", &model.PageData{
		Title: title,
		Data: &RoleFormPageData{
			Role:                 role,
			RolePermissionIdList: rolePermissionIdList,
			PermissionList:       formattedPermissionList,
			PageType:             RoleFormPageEdit,
		},
		BreadCrumbs: breadCrumbs,
	})
}

func (mr *ManageResource) RoleEditSubmit(w http.ResponseWriter, r *http.Request) {
	roleIdStr := chi.URLParam(r, "roleId")
	// fmt.Println("roleId: ", roleIdStr)

	roleId, err := strconv.Atoi(roleIdStr)
	if err != nil {
		mr.Error("", err, w, r, http.StatusBadRequest)
		return
	}

	// isDefault := r.Form.Get("is_default")
	role, err := mr.store.Role.Item(roleId)
	if err != nil {
		mr.Error("", err, w, r, http.StatusInternalServerError)
		return
	}

	currUser := mr.GetLoginedUserData(r)

	if role.IsDefault && !currUser.Super {
		mr.Error("", nil, w, r, http.StatusForbidden)
		return
	}

	name := r.Form.Get("name")
	permissions := r.Form["permissions"]

	// fmt.Println("permissions: ", permissions)

	role = &model.Role{
		Id:   roleId,
		Name: name,
	}

	role.TrimSpace()

	err = role.Valid(true)
	if err != nil {
		mr.Error(err.Error(), err, w, r, http.StatusBadRequest)
		return
	}

	// var permissionIds []int

	// for _, idStr := range permissions {
	// 	id, err := strconv.Atoi(idStr)
	// 	if err == nil {
	// 		permissionIds = append(permissionIds, id)
	// 	}
	// }

	// _, err = mr.store.Role.Update(role.Id, role.Name, permissionIds)

	filteredPermissionList := mr.getFilteredPermissionList(w, r)
	filteredPermissionMap := make(map[string]bool)
	for _, pItem := range filteredPermissionList {
		filteredPermissionMap[pItem.FrontId] = true
	}

	var permissionFrontIds []string
	for _, pId := range permissions {
		if _, ok := filteredPermissionMap[pId]; ok {
			permissionFrontIds = append(permissionFrontIds, pId)
		}
	}

	_, err = mr.store.Role.UpdateWithFrontId(role.Id, role.Name, permissionFrontIds)

	if err != nil {
		mr.Error("", errors.WithStack(err), w, r, http.StatusInternalServerError)

		return
	}

	mr.Session("one", w, r).Flash("Update role successfully")

	http.Redirect(w, r, "/manage/roles", http.StatusFound)

}

func (mr *ManageResource) ActivityList(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	userName := query.Get("username")
	actType := query.Get("type")
	action := query.Get("action")
	pageStr := query.Get("page")
	pageSizeStr := query.Get("page_size")

	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)

	if page < DefaultPage {
		page = DefaultPage
	}

	if pageSize < DefaultPageSize {
		pageSize = DefaultPageSize
	}

	userName = strings.TrimSpace(userName)
	actType = strings.TrimSpace(actType)
	action = strings.TrimSpace(action)

	list, total, err := mr.store.Activity.List(0, userName, actType, action, page, pageSize)
	if err != nil {
		mr.Error("", err, w, r, http.StatusInternalServerError)
		return
	}

	for _, item := range list {
		item.Format(mr.i18nCustom)
	}

	type QueryData struct {
		UserName, Type, Action string
		Total, Page, TotalPage int
	}

	type QctivityPageData struct {
		List            []*model.Activity
		AcTypeOptions   []*model.OptionItem
		AcActionOptions []*model.OptionItem
		Query           *QueryData
	}

	// acTypeVals := model.AcTypeValues()
	var acTypeStrEnums []model.StringEnum
	for _, item := range model.AcTypeValues() {
		acTypeStrEnums = append(acTypeStrEnums, item)
	}
	acTypeOptons := model.ConvertEnumToOPtions(acTypeStrEnums, true, "AcType", mr.i18nCustom)

	var acActionStrEnums []model.StringEnum
	for _, item := range model.AcActionValues() {
		acActionStrEnums = append(acActionStrEnums, item)
	}
	acActionOptons := model.ConvertEnumToOPtions(acActionStrEnums, true, "AcAction", mr.i18nCustom)

	title := mr.Local("List", "Name", mr.Local("Activity", "Count", 1))
	breadCrumbs := []*model.BreadCrumb{
		{
			Path: "/manage/activities",
			Name: title,
		},
	}

	mr.Render(w, r, "activity", &model.PageData{
		Title: title,
		Data: &QctivityPageData{
			List:            list,
			AcTypeOptions:   acTypeOptons,
			AcActionOptions: acActionOptons,
			Query: &QueryData{
				UserName:  userName,
				Type:      actType,
				Action:    action,
				Total:     total,
				Page:      page,
				TotalPage: CeilInt(total, pageSize),
			},
		},
		BreadCrumbs: breadCrumbs,
	})
}

func (mr *ManageResource) TrashPage(w http.ResponseWriter, r *http.Request) {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		// fmt.Printf("page err %v\n", err)
		page = DefaultPage
	}

	pageSize, err := strconv.Atoi(r.URL.Query().Get("page_size"))
	if err != nil {
		pageSize = DefaultPageSize
	}

	sort := r.URL.Query().Get("sort")
	keywords := r.URL.Query().Get("keywords")
	categoryFrontId := r.URL.Query().Get("category")

	var sortType model.ArticleSortType
	if model.ValidArticleSort(sort) {
		sortType = model.ArticleSortType(sort)
	} else {
		sortType = model.ListSortLatest
	}

	deletedList, total, err := mr.store.Article.List(page, pageSize, sortType, categoryFrontId, false, true, true, keywords)
	if err != nil {
		mr.ServerErrorp("", err, w, r)
		return
	}

	var categoryList []*model.Category
	var wg sync.WaitGroup
	categoryMap := make(map[string]*model.Category)
	ch := make(chan any, 1)

	wg.Add(1)
	go func() {
		defer wg.Done()
		list, err := mr.store.Category.List(model.CategoryStateAll)
		if err != nil {
			ch <- err
			return
		}

		ch <- list
	}()

	go func() {
		wg.Wait()
		close(ch)
	}()

	for v := range ch {
		switch val := v.(type) {
		case error:
			mr.ServerErrorp("", val, w, r)
			return
		case []*model.Category:
			categoryList = val
		}
	}

	for _, item := range categoryList {
		categoryMap[item.FrontId] = item
	}

	type pageData struct {
		CurrPage, PageSize, Total, TotalPage int
		SortType, Keywords, CategoryFrontId  string
		List                                 []*model.Article
		CategoryList                         []*model.Category
		CategoryMap                          map[string]*model.Category
	}

	mr.Render(w, r, "trash", &model.PageData{
		Title: mr.Local("Trash"),
		Data: &pageData{
			CurrPage:        page,
			PageSize:        pageSize,
			Total:           total,
			TotalPage:       CeilInt(total, pageSize),
			List:            deletedList,
			SortType:        string(sortType),
			Keywords:        keywords,
			CategoryFrontId: categoryFrontId,
			CategoryList:    categoryList,
			CategoryMap:     categoryMap,
		},
		BreadCrumbs: []*model.BreadCrumb{
			{
				Path: "/manage/trash",
				Name: mr.Local("Trash"),
			},
		},
	})
}
