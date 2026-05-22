package service

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	db "github.com/lepidoptera/lepidoptera/internal/db/generated"
	"github.com/lepidoptera/lepidoptera/internal/mailer"
	"github.com/lepidoptera/lepidoptera/internal/timeutil"
	"github.com/lepidoptera/lepidoptera/internal/viewmodels"
)

type AdminService struct {
	queries *db.Queries
	mailer  mailer.Mailer
}

func NewAdminService(queries *db.Queries, m mailer.Mailer) *AdminService {
	return &AdminService{queries: queries, mailer: m}
}

func (s *AdminService) ListAllShows(ctx context.Context) ([]viewmodels.AdminShowViewModel, error) {
	rows, err := s.queries.ListAllShows(ctx)
	if err != nil {
		return nil, ErrInternalError
	}

	return viewmodels.NewAdminShowViewModels(rows), nil
}

func (s *AdminService) ListVenues(ctx context.Context) ([]db.Venue, error) {
	venues, err := s.queries.ListVenues(ctx)
	if err != nil {
		return nil, ErrInternalError
	}

	return venues, nil
}

func (s *AdminService) ListBands(ctx context.Context) ([]db.Band, error) {
	bands, err := s.queries.ListBands(ctx)
	if err != nil {
		return nil, ErrInternalError
	}

	return bands, nil
}

func (s *AdminService) CreateVenue(ctx context.Context, params db.CreateVenueParams) error {
	_, err := s.queries.CreateVenue(ctx, params)
	if err != nil {
		return ErrInternalError
	}

	return nil
}

func (s *AdminService) CreateBand(ctx context.Context, params db.CreateBandParams) error {
	_, err := s.queries.CreateBand(ctx, params)
	if err != nil {
		return ErrInternalError
	}

	return nil
}

func (s *AdminService) CreateShow(ctx context.Context, params db.CreateShowParams, bandIDs []string, headlinerID string) error {
	show, err := s.queries.CreateShow(ctx, params)
	if err != nil {
		return ErrInternalError
	}

	for i, bandIDStr := range bandIDs {
		bandID, err := uuid.Parse(bandIDStr)
		if err != nil {
			continue
		}

		s.queries.AddBandToShow(ctx, db.AddBandToShowParams{
			ShowID:       show.ID,
			BandID:       bandID,
			BillingOrder: int32(i),
			IsHeadliner:  bandIDStr == headlinerID,
		})
	}

	return nil
}

func (s *AdminService) PublishShow(ctx context.Context, id uuid.UUID) error {
	show, err := s.queries.PublishShow(ctx, id)
	if err != nil {
		return ErrShowNotFound
	}

	err = s.queries.CreateNotificationsForShow(ctx, db.CreateNotificationsForShowParams{
		ShowID: uuid.NullUUID{UUID: show.ID, Valid: true},
		Type:   "published",
	})
	if err != nil {
		return ErrInternalError
	}

	return nil
}

func (s *AdminService) GenerateSlug(title string, date time.Time) string {
	return slugify(title + " " + date.Format(timeutil.SlugDate))
}

func (s *AdminService) Mailer() mailer.Mailer {
	return s.mailer
}

func slugify(s string) string {
	s = strings.ToLower(s)

	var result strings.Builder
	for _, r := range s {
		switch {
		case r >= 'a' && r <= 'z':
			result.WriteRune(r)
		case r >= '0' && r <= '9':
			result.WriteRune(r)
		case r == ' ' || r == '-' || r == '_':
			result.WriteRune('-')
		}
	}

	slug := result.String()
	for strings.Contains(slug, "--") {
		slug = strings.ReplaceAll(slug, "--", "-")
	}

	return strings.Trim(slug, "-")
}
