package database

import (
	"database/sql"
)

func CreateTables(db *sql.DB) error {
	createTableSQL := `
		CREATE TABLE IF NOT EXISTS pokemons (
    		id TEXT PRIMARY KEY,
    		name TEXT NOT NULL UNIQUE,
    		description TEXT NOT NULL,
    		category TEXT NOT NULL
		);

		CREATE TABLE IF NOT EXISTS pokemon_types (
			pokemon_id TEXT NOT NULL,
			type TEXT NOT NULL,
			FOREIGN KEY (pokemon_id) REFERENCES pokemons (id) ON DELETE CASCADE
		);

		CREATE TABLE IF NOT EXISTS pokemon_abilities (
			pokemon_id TEXT NOT NULL,
			ability TEXT NOT NULL,
			FOREIGN KEY (pokemon_id) REFERENCES pokemons (id) ON DELETE CASCADE
		);
	`

	_, err := db.Exec(createTableSQL)
	return err
}
