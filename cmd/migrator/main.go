package main

import (
	"errors"
	"flag"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	var storaePath, migrationsPath, migrationsTable string

	flag.StringVar(&storaePath, "storage-path", "", "путь до хранилища с базой данных sqlite")
	flag.StringVar(&migrationsPath, "migrations-path", "", "путь до директории с файлами миграции")
	flag.StringVar(&migrationsTable, "migrations-table", "migrations", "имя таблицы с миграциями")
	flag.Parse()

	if storaePath == "" {
		panic("не задан параметр storage-path: путь до хранилища с базой данных sqlite")
	}
	if migrationsPath == "" {
		panic("не задан параметр migrations-path: путь до директории с файлами миграции")
	}

	m, err := migrate.New(
		"file://"+migrationsPath,
		fmt.Sprintf("sqlite3://%s?x-migrations-table=%s", storaePath, migrationsTable),
	)
	if err != nil {
		panic(err)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("нет новых миграций для применения")
			return
		}
		panic(err)
	}

	fmt.Println("все миграции успешно применены")
}
