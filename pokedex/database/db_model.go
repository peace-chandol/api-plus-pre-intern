package database

import "gorm.io/gorm"

type Pokemon struct {
	gorm.Model
	Name        string        `gorm:"unique;not null"`
	Description string        `gorm:"not null"`
	Category    string        `gorm:"not null"`
	Types       []PokemonType `gorm:"many2many:pokemon_pokemon_types;"`
	Abilities   []Ability     `gorm:"many2many:pokemon_abilities;"`
}

type PokemonType struct {
	gorm.Model
	Type     string     `gorm:"not null"`
	Pokemons []*Pokemon `gorm:"many2many:pokemon_pokemon_types;"`
}

type Ability struct {
	gorm.Model
	Ability  string     `gorm:"not null"`
	Pokemons []*Pokemon `gorm:"many2many:pokemon_abilities;"`
}
