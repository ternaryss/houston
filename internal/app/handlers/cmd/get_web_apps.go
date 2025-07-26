package cmd

import (
	"log/slog"
	"net/url"
	"strconv"

	"github.com/ternaryss/houston/internal/app/types"
)

type getWebAppsCmd struct {
	webAppsStore types.WebAppsStore
}

func NewGetWebAppsCmd(was types.WebAppsStore) *getWebAppsCmd {
	return &getWebAppsCmd{
		webAppsStore: was,
	}
}

func (c *getWebAppsCmd) Execute(qry url.Values, usr string) (types.Page[*types.WebApp], error) {
	slog.Info("Fetching web applications", "query", qry, "user", usr)
	sortable := []string{"name", "userEmail", "createdAt"}
	filter := types.NewFilter("", "name asc", "", sortable)
	filter.Params["userEmail"] = usr
	quantity, err := c.webAppsStore.CountByFilter(filter, nil)
	slog.Info("Web applications counted", "quantity", quantity, "filter", filter)

	if err != nil {
		return types.EmptyPage[*types.WebApp](), err
	}

	page, err := strconv.Atoi(qry.Get("page"))

	if err != nil {
		page = types.DefaultPage
	}

	size, err := strconv.Atoi(qry.Get("size"))

	if err != nil {
		size = types.DefaultPageSize
	}

	pagination := types.NewPagination(page, size, quantity)

	if !pagination.IsValid() {
		slog.Warn("Pagination not valid", "pagination", pagination)
		return types.EmptyPage[*types.WebApp](), nil
	}

	apps, err := c.webAppsStore.GetByFilter(filter, pagination, nil)

	if err != nil {
		return types.EmptyPage[*types.WebApp](), err
	}

	slog.Info("Web applications fetched", "filter", filter, "pagination", pagination, "contentSize", len(apps))

	return types.NewPage(pagination, apps), nil
}
