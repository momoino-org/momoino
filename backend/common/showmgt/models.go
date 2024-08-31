package showmgt

import (
	"wano-island/common/core"

	"github.com/google/uuid"
)

type ShowModel struct {
	core.Model
	core.HasCreatedAtColumn
	core.HasUpdatedAtColumn

	Kind             string                 `gorm:"type:string;size:7;not null"`
	OriginalLanguage string                 `gorm:"type:string;size:256;not null"`
	OriginalTitle    string                 `gorm:"type:string;size:256;not null"`
	OriginalOverview *string                `gorm:"type:string;size:256"`
	Keywords         *[]string              `gorm:"type:text[];serializer:json"`
	IsReleased       bool                   `gorm:"type:boolean"`
	Seasons          []SeasonModel          `gorm:"foreignKey:ShowID;constraint:OnDelete:CASCADE"`
	Translations     []ShowTranslationModel `gorm:"foreignKey:ShowID;constraint:OnDelete:CASCADE"`
}

type ShowTranslationModel struct {
	core.Model
	core.HasCreatedAtColumn
	core.HasUpdatedAtColumn

	ShowID   uuid.UUID
	Locale   string `gorm:"type:string;size:256;not null"`
	Title    string `gorm:"type:string;size:256;not null"`
	Overview string `gorm:"type:string;size:256;not null"`
}

type SeasonModel struct {
	core.Model
	core.HasCreatedAtColumn
	core.HasUpdatedAtColumn

	ShowID       uuid.UUID
	Order        int
	Translations []SeasonTranslationModel `gorm:"foreignKey:SeasonID;constraint:OnDelete:CASCADE"`
}

type SeasonTranslationModel struct {
	core.Model
	core.HasCreatedAtColumn
	core.HasUpdatedAtColumn

	SeasonID uuid.UUID
	Locale   string `gorm:"type:string;size:256;not null"`
	Title    string `gorm:"type:string;size:256;not null"`
	Overview string `gorm:"type:string;size:256;not null"`
}

type EpisodeModel struct {
	core.Model
	core.HasCreatedAtColumn
	core.HasUpdatedAtColumn

	ShowID       uuid.UUID
	Order        int
	Title        string                    `gorm:"type:string;size:256;not null"`
	Overview     string                    `gorm:"type:string;size:256"`
	Translations []EpisodeTranslationModel `gorm:"foreignKey:EpisodeID;constraint:OnDelete:CASCADE"`
}

type EpisodeTranslationModel struct {
	core.Model
	core.HasCreatedAtColumn
	core.HasUpdatedAtColumn

	EpisodeID uuid.UUID
	Locale    string `gorm:"type:string;size:256;not null"`
	Title     string `gorm:"type:string;size:256;not null"`
	Overview  string `gorm:"type:string;size:256;not null"`
}

func (ShowModel) TableName() string {
	return "public.shows"
}

func (ShowTranslationModel) TableName() string {
	return "public.show_translations"
}

func (SeasonModel) TableName() string {
	return "public.seasons"
}

func (SeasonTranslationModel) TableName() string {
	return "public.season_translations"
}

func (EpisodeModel) TableName() string {
	return "public.episodes"
}

func (EpisodeTranslationModel) TableName() string {
	return "public.episode_translations"
}
