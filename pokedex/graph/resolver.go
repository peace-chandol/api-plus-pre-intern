package graph

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

import (
	"github.com/peace/pokedex/database"
)

// type Database struct {
// 	PokemonsTable map[string]model.Pokemon
// }

// func (db *Database) FindPokemonById(id string) (*model.Pokemon, error) {
// 	if rek, ok := db.PokemonsTable[id]; ok {
// 		return &rek, nil
// 	} else {
// 		return nil, fmt.Errorf("pokemon id: %s was not found", id)
// 	}
// }

// func (db *Database) FindPokemonByName(name string) (*model.Pokemon, error) {
// 	for _, p := range db.PokemonsTable {
// 		if p.Name == name {
// 			return &p, nil
// 		}
// 	}

// 	return nil, fmt.Errorf("pokemon name: %s was not found", name)
// }

// func (db *Database) FindAllPokemons() []*model.Pokemon {
// 	allPokemons := []*model.Pokemon{}
// 	for _, p := range db.PokemonsTable {
// 		newPokemon := p
// 		allPokemons = append(allPokemons, &newPokemon)
// 	}

// 	return allPokemons
// }

// func (db *Database) AddPokemon(input *model.Pokemon) error {
// 	newID := uuid.New().String()
// 	input.ID = newID
// 	db.PokemonsTable[newID] = *input

// 	return nil
// }

// func (db *Database) UpdatePokemon(input model.Pokemon) error {
// 	if _, ok := db.PokemonsTable[input.ID]; !ok {
// 		return fmt.Errorf("pokemon id: %s was not found", input.ID)
// 	}

// 	db.PokemonsTable[input.ID] = input
// 	return nil
// }

// func (db *Database) DeletePokemon(id string) error {
// 	if _, ok := db.PokemonsTable[id]; !ok {
// 		return fmt.Errorf("pokemon id: %s was not found", id)
// 	}

// 	delete(db.PokemonsTable, id)
// 	return nil
// }

type Resolver struct {
	DB *database.Database
}
