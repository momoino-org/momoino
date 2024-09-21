package showmgt

import (
	"time"

	"github.com/google/uuid"
)

type ShowDTO struct {
	ID               uuid.UUID `json:"id"`
	Kind             string    `json:"kind"`
	OriginalLanguage string    `json:"originalLanguage"`
	OriginalTitle    string    `json:"originalTitle"`
	OriginalOverview *string   `json:"originalOverview"`
	Keywords         []string  `json:"keywords"`
	IsReleased       bool      `json:"isReleased"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
}

// ToShowDTO converts a ShowModel to a ShowDTO.
func ToShowDTO(showModel *ShowModel) *ShowDTO {
	if showModel == nil {
		return nil
	}

	return &ShowDTO{
		ID:               showModel.ID,
		Kind:             showModel.Kind,
		OriginalLanguage: showModel.OriginalLanguage,
		OriginalTitle:    showModel.OriginalTitle,
		OriginalOverview: showModel.OriginalOverview,
		Keywords:         []string(showModel.Keywords),
		IsReleased:       showModel.IsReleased,
		CreatedAt:        showModel.CreatedAt,
		UpdatedAt:        showModel.UpdatedAt,
	}
}
