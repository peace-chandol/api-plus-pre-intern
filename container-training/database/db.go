package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/peace/pokedex/graph/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Database struct {
	DB *gorm.DB
}

const (
	host     = "postgres"
	port     = 5432
	user     = "pokedex"
	password = "password"
	dbname   = "pokedex_db"
)

func ConnectDB() (*Database, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			LogLevel: logger.Info,
			Colorful: true,
		},
	)

	var db *gorm.DB
	var err error

	for i := 0; i < 5; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: newLogger})
		if err == nil {
			break
		}
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&Pokemon{}, &PokemonType{}, &Ability{}); err != nil {
		return nil, err
	}

	return &Database{DB: db}, nil
}

func InitDB() (*Database, error) {
	db, err := ConnectDB()
	if err != nil {
		return nil, err
	}

	return db, nil
}

// Interacts with database

func (db *Database) AddPokemon(pokemon *model.Pokemon) error {
	p := Pokemon{
		Name:        pokemon.Name,
		Description: pokemon.Description,
		Category:    pokemon.Category,
	}
	p.ID = uuid.New()

	for _, t := range pokemon.Type {
		pt := PokemonType{Type: string(t)}
		pt.ID = uuid.New()
		if result := db.DB.Create(&pt); result.Error != nil {
			return result.Error
		}
		p.Types = append(p.Types, pt)
	}

	for _, a := range pokemon.Abilities {
		pa := Ability{Ability: a}
		pa.ID = uuid.New()
		if result := db.DB.Create(&pa); result.Error != nil {
			return result.Error
		}
		p.Abilities = append(p.Abilities, pa)
	}

	if result := db.DB.Create(&p); result.Error != nil {
		return result.Error
	}

	pokemon.ID = p.ID.String()

	return nil
}

func (db *Database) UpdatePokemon(pokemon *model.Pokemon) error {
	p := Pokemon{}
	if result := db.DB.Where("id = ?", pokemon.ID).First(&p); result.Error != nil {
		return result.Error
	}
	p.Name = pokemon.Name
	p.Description = pokemon.Description
	p.Category = pokemon.Category

	if result := db.DB.Save(&p); result.Error != nil {
		return result.Error
	}

	if result := db.DB.Model(&p).Association("Types").Clear(); result != nil {
		return result
	}

	if result := db.DB.Model(&p).Association("Abilities").Clear(); result != nil {
		return result
	}

	for _, t := range pokemon.Type {
		pt := PokemonType{Type: string(t)}
		pt.ID = uuid.New()
		if result := db.DB.Create(&pt); result.Error != nil {
			return result.Error
		}
		p.Types = append(p.Types, pt)
	}

	if result := db.DB.Model(&p).Association("Types").Replace(p.Types); result != nil {
		return result
	}

	for _, a := range pokemon.Abilities {
		pa := Ability{Ability: a}
		pa.ID = uuid.New()
		if result := db.DB.Create(&pa); result.Error != nil {
			return result.Error
		}
		p.Abilities = append(p.Abilities, pa)
	}

	if result := db.DB.Model(&p).Association("Abilities").Replace(p.Abilities); result != nil {
		return result
	}

	return nil
}

func (db *Database) DeletePokemon(id string) error {
	p := Pokemon{}
	if result := db.DB.Preload("Types").Preload("Abilities").Where("id = ?", id).First(&p); result.Error != nil {
		return result.Error
	}

	if err := db.DB.Model(&p).Association("Types").Clear(); err != nil {
		return err
	}

	if err := db.DB.Model(&p).Association("Abilities").Clear(); err != nil {
		return err
	}

	if result := db.DB.Delete(&p); result.Error != nil {
		return result.Error
	}

	return nil
}

func convertPokemonTypes(types []PokemonType) []model.PokemonType {
	pokemonTypes := []model.PokemonType{}
	for _, t := range types {
		pokemonTypes = append(pokemonTypes, model.PokemonType(t.Type))
	}

	return pokemonTypes
}

func convertPokemonAbilities(abilities []Ability) []string {
	var pokemonAbilites []string
	for _, a := range abilities {
		pokemonAbilites = append(pokemonAbilites, a.Ability)
	}

	return pokemonAbilites
}

func (db *Database) FindAllPokemons() ([]*model.Pokemon, error) {
	pokemons := []Pokemon{}

	if result := db.DB.Preload("Types").Preload("Abilities").Find(&pokemons); result.Error != nil {
		return nil, result.Error
	}

	pokemonsModel := []*model.Pokemon{}
	for _, p := range pokemons {
		pokemonsModel = append(pokemonsModel, &model.Pokemon{
			ID:          p.ID.String(),
			Name:        p.Name,
			Description: p.Description,
			Category:    p.Category,
			Type:        convertPokemonTypes(p.Types),
			Abilities:   convertPokemonAbilities(p.Abilities),
		})
	}

	return pokemonsModel, nil
}

func (db *Database) FindPokemonById(id string) (*model.Pokemon, error) {
	pokemon := Pokemon{}

	if result := db.DB.Preload("Types").Preload("Abilities").Where("ID = ?", id).First(&pokemon); result.Error != nil {
		return nil, result.Error
	}

	modelPokemon := model.Pokemon{
		ID:          pokemon.ID.String(),
		Name:        pokemon.Name,
		Description: pokemon.Description,
		Category:    pokemon.Category,
		Type:        convertPokemonTypes(pokemon.Types),
		Abilities:   convertPokemonAbilities(pokemon.Abilities),
	}

	return &modelPokemon, nil
}

func (db *Database) FindPokemonByName(name string) (*model.Pokemon, error) {
	var pokemon Pokemon

	if result := db.DB.Preload("Types").Preload("Abilities").Where("name = ?", name).First(&pokemon); result.Error != nil {
		return nil, result.Error
	}

	pokemonModel := model.Pokemon{
		ID:          pokemon.ID.String(),
		Name:        pokemon.Name,
		Description: pokemon.Description,
		Category:    pokemon.Category,
		Type:        convertPokemonTypes(pokemon.Types),
		Abilities:   convertPokemonAbilities(pokemon.Abilities),
	}

	return &pokemonModel, nil
}
