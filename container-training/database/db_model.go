package database

import (
	"github.com/google/uuid"
)

type Pokemon struct {
	ID          uuid.UUID     `gorm:"type:uuid;primaryKey"`
	Name        string        `gorm:"unique;not null"`
	Description string        `gorm:"not null"`
	Category    string        `gorm:"not null"`
	Types       []PokemonType `gorm:"many2many:pokemon_pokemon_types;"`
	Abilities   []Ability     `gorm:"many2many:pokemon_abilities;"`
}

type PokemonType struct {
	ID       uuid.UUID  `gorm:"type:uuid;primaryKey;"`
	Type     string     `gorm:"not null"`
	Pokemons []*Pokemon `gorm:"many2many:pokemon_pokemon_types;"`
}

type Ability struct {
	ID       uuid.UUID  `gorm:"type:uuid;primaryKey;"`
	Ability  string     `gorm:"not null"`
	Pokemons []*Pokemon `gorm:"many2many:pokemon_abilities;"`
}
