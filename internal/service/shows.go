package service

import (
	"context"
	"database/sql"

	db "github.com/lepidoptera/lepidoptera/internal/db/generated"
	"github.com/lepidoptera/lepidoptera/internal/viewmodels"
)

type ShowService struct {
	queries *db.Queries
}

func NewShowService(queries *db.Queries) *ShowService {
	return &ShowService{queries: queries}
}

func (s *ShowService) ListPublishedShows(ctx context.Context, genre string) ([]viewmodels.ShowListViewModel, error) {
	rows, err := s.queries.ListPublishedShows(ctx, sql.NullString{
		String: genre,
		Valid:  genre != "",
	})
	if err != nil {
		return nil, ErrInternalError
	}

	return viewmodels.NewShowListViewModels(rows), nil
}

func (s *ShowService) GetShowDetail(ctx context.Context, slug string) (viewmodels.ShowDetailViewModel, error) {
	row, err := s.queries.GetShowBySlug(ctx, slug)
	if err != nil {
		return viewmodels.ShowDetailViewModel{}, ErrShowNotFound
	}

	bands, err := s.queries.GetBandsByShow(ctx, row.ID)
	if err != nil {
		return viewmodels.ShowDetailViewModel{}, ErrInternalError
	}

	return viewmodels.NewShowDetailViewModel(row, bands), nil
}

func (s *ShowService) ListGenres(ctx context.Context) ([]string, error) {
	rows, err := s.queries.ListGenres(ctx)
	if err != nil {
		return nil, ErrInternalError
	}

	genres := make([]string, 0, len(rows))
	for _, g := range rows {
		if g.Valid {
			genres = append(genres, g.String)
		}
	}

	return genres, nil
}
