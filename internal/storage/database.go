package storage

import (
	"database/sql"
	"embed"
	"fmt"
	"log"
	"path/filepath"
	"strings"

	mapset "github.com/deckarep/golang-set/v2"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/viper"
)

//go:embed words/*.txt
var wordsFS embed.FS

func getWords(fs *embed.FS) (map[string]string, error) {
	result := make(map[string]string)
	files, err := fs.ReadDir("words")
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		fileName := file.Name()
		val, err := fs.ReadFile("words/" + fileName)
		if err != nil {
			return nil, err
		}
		content := string(val)
		nameWithoutSuffix := strings.TrimSuffix(fileName, filepath.Ext(fileName))
		result[nameWithoutSuffix] = content
	}
	return result, nil
}

type Database struct {
	db *sql.DB
}

func New() *Database {

	const file string = "words.db"
	handler, err := sql.Open("sqlite3", file)
	if err != nil {
		log.Print(err)
		panic("can't connect to database")
	}

	err = handler.Ping()
	if err != nil {
		log.Print(err)
		panic("can't ping to database")
	}
	database := &Database{
		db: handler,
	}
	database.configing()
	return database
}

func (d *Database) configing() {
	reset := viper.Get("reset").(bool)
	if reset {
		fmt.Println("reset all data...")
		err := d.resetDb()
		if err != nil {
			panic(err)
		}
	} else {
		fmt.Println("check all data...")
		err := d.checkDb()
		if err != nil {
			panic(err)
		}
	}
}

func (d *Database) GetWordTypes() ([]string, error) {
	typeSet, err := d.requestAllTablesName()
	if err != nil {
		return nil, err
	}
	return typeSet.ToSlice(), nil
}

func (d *Database) resetDb() error {
	words, err := getWords(&wordsFS)
	if err != nil {
		return err
	}
	err = d.removeAllTables()
	if err != nil {
		return err
	}
	for wordsType, content := range words {
		fmt.Printf("add data to table %s ... \n", wordsType)
		err = d.createTableWithContent(wordsType, strings.Fields(content))
		if err != nil {
			return err
		}
	}
	fmt.Println(len(words))
	return nil
}

func (d *Database) checkDb() error {
	words, err := getWords(&wordsFS)
	if err != nil {
		return err
	}
	existingTable, err := d.requestAllTablesName()
	if err != nil {
		return err
	}
	for wordsType, content := range words {
		if existingTable.Contains(wordsType) {
			fmt.Printf("table %s already exist exits \n", wordsType)
		} else {
			fmt.Printf("table %s not exits \n", wordsType)
			err := d.createTableWithContent(wordsType, strings.Fields(content))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (d *Database) requestAllTablesName() (mapset.Set[string], error) {
	rows, err := d.db.Query("SELECT name FROM sqlite_master WHERE type='table'")

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := mapset.NewSet[string]()
	for rows.Next() {
		var tableName string
		err = rows.Scan(&tableName)
		if err != nil {
			return nil, err
		}
		result.Add(tableName)
	}
	fmt.Println(result)
	return result, nil
}

func (d *Database) createTableWithContent(name string, content []string) error {
	queryString := fmt.Sprintf("CREATE TABLE %s ( word TEXT PRIMARY KEY);", name)
	fmt.Println(queryString)
	statement, err := d.db.Prepare(queryString)
	if err != nil {
		return err
	}
	_, err = statement.Exec()
	if err != nil {
		return err
	}

	for _, word := range content {
		queryString := fmt.Sprintf("INSERT INTO %s (word) VALUES (\"%s\")", name, word)

		fmt.Println(queryString)
		statement, err := d.db.Prepare(queryString)
		if err != nil {
			return err
		}
		_, err = statement.Exec()
		if err != nil {
			return err
		}

		// fmt.Println(queryString)
		// intertResult, err := d.db.Query(queryString)
		// if intertResult != nil {
		// 	intertResult.Close()
		// }
		// if err != nil {
		// 	fmt.Println(err)
		// 	panic(err)
		// }
	}
	return nil
}

func (d *Database) removeAllTables() error {
	dropResult, err := d.db.Query("DROP SCHEMA public CASCADE")
	if err != nil {
		return err
	}
	defer dropResult.Close()
	createResult, err := d.db.Query("CREATE SCHEMA public")
	if err != nil {
		return err
	}
	defer createResult.Close()
	return nil
}
