package viewmodels

import (
	"fmt"

	db "github.com/bradonstudie/lepidoptera/internal/db/generated"
	"github.com/bradonstudie/lepidoptera/internal/timeutil"
)

type ShowListViewModel struct {
	ID        string
	Title     string
	Slug      string
	Date      string
	Venue     string
	City      string
	TicketURL string
	HasTicket bool
}

type ShowDetailViewModel struct {
	ID          string
	Title       string
	Slug        string
	Date        string
	Venue       string
	City        string
	Address     string
	Description string
	FlyerURL    string
	HasFlyer    bool
	TicketURL   string
	HasTicket   bool
	Bands       []BandViewModel
}

type AdminShowViewModel struct {
	ID          string
	Title       string
	Slug        string
	Date        string
	Venue       string
	IsPublished bool
	PublishedAt string
}

func NewShowListViewModel(row db.ListPublishedShowsRow) ShowListViewModel {
	return ShowListViewModel{
		ID:        row.ID.String(),
		Title:     row.Title,
		Slug:      row.Slug,
		Date:      row.Date.Format(timeutil.DisplayDate),
		Venue:     row.VenueName.String,
		City:      row.VenueCity.String,
		TicketURL: row.TicketUrl.String,
		HasTicket: row.TicketUrl.Valid,
	}
}

func NewShowListViewModels(rows []db.ListPublishedShowsRow) []ShowListViewModel {
	vms := make([]ShowListViewModel, len(rows))
	for i, row := range rows {
		vms[i] = NewShowListViewModel(row)
	}
	return vms
}

func NewShowDetailViewModel(row db.GetShowBySlugRow, bands []db.GetBandsByShowRow) ShowDetailViewModel {
	return ShowDetailViewModel{
		ID:          row.ID.String(),
		Title:       row.Title,
		Slug:        row.Slug,
		Date:        row.Date.Format(timeutil.DisplayFull),
		Venue:       row.VenueName.String,
		City:        row.VenueCity.String,
		Address:     row.VenueAddress.String,
		Description: row.Description.String,
		FlyerURL:    row.FlyerUrl.String,
		HasFlyer:    row.FlyerUrl.Valid,
		TicketURL:   row.TicketUrl.String,
		HasTicket:   row.TicketUrl.Valid,
		Bands:       NewBandViewModels(bands),
	}
}

func NewAdminShowViewModel(row db.ListAllShowsRow) AdminShowViewModel {
	publishedAt := ""
	if row.PublishedAt.Valid {
		publishedAt = fmt.Sprintf("published %s", row.PublishedAt.Time.Format(timeutil.DisplayDate))
	}

	return AdminShowViewModel{
		ID:          row.ID.String(),
		Title:       row.Title,
		Slug:        row.Slug,
		Date:        row.Date.Format(timeutil.DisplayDate),
		Venue:       row.VenueName.String,
		IsPublished: row.IsPublished,
		PublishedAt: publishedAt,
	}
}

func NewAdminShowViewModels(rows []db.ListAllShowsRow) []AdminShowViewModel {
	vms := make([]AdminShowViewModel, len(rows))
	for i, row := range rows {
		vms[i] = NewAdminShowViewModel(row)
	}
	return vms
}
