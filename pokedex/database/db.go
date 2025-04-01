package database

import (
	"context"
	"database/sql"
	"log"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"github.com/peace/pokedex/graph/model"
)

type Database struct {
	Conn *sql.DB
}

func ConnectDB() (*Database, error) {
	db, err := sql.Open("sqlite3", "./pokedex.db")
	if err != nil {
		return nil, err
	}

	return &Database{Conn: db}, nil
}

func InitDB() (*Database, error) {
	dbStruct, err := ConnectDB()
	if err != nil {
		return nil, err
	}

	if err := CreateTables(dbStruct.Conn); err != nil {
		log.Fatal(err)
	}

	return dbStruct, nil
}

// Database Function

func (db *Database) AddPokemon(ctx context.Context, pokemon *model.Pokemon) error {
	newID := uuid.New().String()
	pokemon.ID = newID

	_, err := db.Conn.ExecContext(ctx, `INSERT INTO pokemons (id, name, description, category) VALUES (?,?,?,?)`, newID, pokemon.Name, pokemon.Description, pokemon.Category)
	if err != nil {
		return err
	}

	for _, typeValue := range pokemon.Type {
		_, err = db.Conn.ExecContext(ctx, `INSERT INTO pokemon_types (pokemon_id, type) VALUES (?,?)`, newID, typeValue.String())
		if err != nil {
			return err
		}
	}

	for _, ability := range pokemon.Abilities {
		_, err = db.Conn.ExecContext(ctx, `INSERT INTO pokemon_abilities (pokemon_id, ability) VALUES (?,?)`, newID, ability)
		if err != nil {
			return err
		}
	}

	return nil
}

func (db *Database) UpdatePokemon(ctx context.Context, pokemon *model.Pokemon) error {
	_, err := db.Conn.ExecContext(ctx, `UPDATE pokemons SET name = ?, description = ?, category = ? WHERE id = ?`, pokemon.Name, pokemon.Description, pokemon.Category, pokemon.ID)
	if err != nil {
		return err
	}

	_, err = db.Conn.ExecContext(ctx, `DELETE FROM pokemon_types WHERE pokemon_id = ?`, pokemon.ID)
	if err != nil {
		return err
	}

	_, err = db.Conn.ExecContext(ctx, `DELETE FROM pokemon_abilities WHERE pokemon_id = ?`, pokemon.ID)
	if err != nil {
		return err
	}

	for _, typeValue := range pokemon.Type {
		_, err = db.Conn.ExecContext(ctx, `INSERT INTO pokemon_types (pokemon_id, type) VALUES (?,?)`, pokemon.ID, typeValue.String())
		if err != nil {
			return err
		}
	}

	for _, ability := range pokemon.Abilities {
		_, err = db.Conn.ExecContext(ctx, `INSERT INTO pokemon_abilities (pokemon_id, ability) VALUES (?,?)`, pokemon.ID, ability)
		if err != nil {
			return err
		}
	}

	return nil
}

func (db *Database) DeletePokemon(ctx context.Context, id string) error {
	_, err := db.Conn.ExecContext(ctx, `DELETE FROM pokemons WHERE id = ?`, id)
	return err
}

func (db *Database) getPokemonTypes(ctx context.Context, id string) ([]model.PokemonType, error) {
	types := []model.PokemonType{}

	rows, err := db.Conn.QueryContext(ctx, `SELECT type FROM pokemon_types WHERE pokemon_id = ?`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var pokemonType string
		if err := rows.Scan(&pokemonType); err != nil {
			return nil, err
		}

		types = append(types, model.PokemonType(pokemonType))
	}

	return types, nil
}

func (db *Database) getPokemonAbilities(ctx context.Context, id string) ([]string, error) {
	abilities := []string{}

	rows, err := db.Conn.QueryContext(ctx, `SELECT ability FROM pokemon_abilities WHERE pokemon_id = ?`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var ability string
		if err := rows.Scan(&ability); err != nil {
			return nil, err
		}

		abilities = append(abilities, ability)
	}

	return abilities, nil
}

func (db *Database) FindAllPokemons(ctx context.Context) ([]*model.Pokemon, error) {
	pokemons := []*model.Pokemon{}

	rows, err := db.Conn.QueryContext(ctx, `SELECT id, name, description, category FROM pokemons`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		pokemon := model.Pokemon{}
		if err := rows.Scan(&pokemon.ID, &pokemon.Name, &pokemon.Description, &pokemon.Category); err != nil {
			return nil, err
		}
		pokemon.Type, err = db.getPokemonTypes(ctx, pokemon.ID)
		if err != nil {
			return nil, err
		}
		pokemon.Abilities, err = db.getPokemonAbilities(ctx, pokemon.ID)
		if err != nil {
			return nil, err
		}

		pokemons = append(pokemons, &pokemon)
	}

	return pokemons, nil
}

func (db *Database) FindPokemonById(ctx context.Context, id string) (*model.Pokemon, error) {
	row := db.Conn.QueryRowContext(ctx, `SELECT id, name, description, category FROM pokemons WHERE id = ?`, id)

	pokemon := model.Pokemon{}
	err := row.Scan(&pokemon.ID, &pokemon.Name, &pokemon.Description, &pokemon.Category)
	if err != nil {
		return nil, err
	}
	pokemon.Type, err = db.getPokemonTypes(ctx, pokemon.ID)
	if err != nil {
		return nil, err
	}
	pokemon.Abilities, err = db.getPokemonAbilities(ctx, pokemon.ID)
	if err != nil {
		return nil, err
	}

	return &pokemon, nil
}

func (db *Database) FindPokemonByName(ctx context.Context, name string) (*model.Pokemon, error) {
	row := db.Conn.QueryRowContext(ctx, `SELECT id, name, description, category FROM pokemons WHERE name = ?`, name)

	pokemon := model.Pokemon{}
	err := row.Scan(&pokemon.ID, &pokemon.Name, &pokemon.Description, &pokemon.Category)
	if err != nil {
		return nil, err
	}
	pokemon.Type, err = db.getPokemonTypes(ctx, pokemon.ID)
	if err != nil {
		return nil, err
	}
	pokemon.Abilities, err = db.getPokemonAbilities(ctx, pokemon.ID)
	if err != nil {
		return nil, err
	}

	return &pokemon, nil
}
