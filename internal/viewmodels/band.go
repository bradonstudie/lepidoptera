package viewmodels

import db "github.com/bradonstudie/lepidoptera/internal/db/generated"

type BandViewModel struct {
	ID          string
	Name        string
	Genre       string
	Description string
	LogoURL     string
	HasLogo     bool
	WebsiteURL  string
	HasWebsite  bool
	IsHeadliner bool
	Order       int
}

func NewBandViewModel(row db.GetBandsByShowRow) BandViewModel {
	return BandViewModel{
		ID:          row.ID.String(),
		Name:        row.Name,
		Genre:       row.Genre.String,
		Description: row.Description.String,
		LogoURL:     row.LogoUrl.String,
		HasLogo:     row.LogoUrl.Valid,
		WebsiteURL:  row.WebsiteUrl.String,
		HasWebsite:  row.WebsiteUrl.Valid,
		IsHeadliner: row.IsHeadliner,
		Order:       int(row.BillingOrder),
	}
}

func NewBandViewModels(rows []db.GetBandsByShowRow) []BandViewModel {
	vms := make([]BandViewModel, len(rows))
	for i, row := range rows {
		vms[i] = NewBandViewModel(row)
	}
	return vms
}
