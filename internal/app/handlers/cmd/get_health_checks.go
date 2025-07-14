package cmd

import (
	"log/slog"
	"net/url"
	"strconv"

	"github.com/ternaryss/houston/internal/app/types"
)

type getHealthChecksCmd struct {
	healthChecksStore types.HealthChecksStore
}

func NewGetHealthChecksCmd(hcs types.HealthChecksStore) *getHealthChecksCmd {
	return &getHealthChecksCmd{
		healthChecksStore: hcs,
	}
}

func (c *getHealthChecksCmd) Execute(wid string, qry url.Values) (types.Page, error) {
	slog.Info("Fetching web application health checks", "webAppId", wid, "query", qry)
	sortable := []string{"createdAt"}
	filter := types.NewFilter("", "createdAt desc", "", sortable)
	filter.Params["webAppId"] = wid
	quantity, err := c.healthChecksStore.CountByFilter(filter)

	if err != nil {
		return types.EmptyPage(), err
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
		return types.EmptyPage(), nil
	}

	checks, err := c.healthChecksStore.GetByFilter(filter, pagination)

	if err != nil {
		return types.EmptyPage(), err
	}

	content := []types.PageContent{}

	for _, check := range checks {
		content = append(content, check)
	}

	slog.Info("Web application health checks fetched", "filter", filter, "pagination", pagination, "contentSize", len(content))

	return types.NewPage(pagination, content), nil
}
