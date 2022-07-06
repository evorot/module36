package main

import (
	"GoNews/pkg/api"
	"GoNews/pkg/rss"
	"GoNews/pkg/storage"
	"GoNews/pkg/storage/postgres"
	"log"
	"net/http"
)

const (
	configURL = "./config.json"
	dbURL     = "postgres://postgres:postgrespw@localhost:55000"
)

func main() {

	// инициализация db
	psgr, err := postgres.New(dbURL)
	if err != nil {
		log.Fatal(err)
	}

	// инициализация api
	a := api.New(psgr)

	chP := make(chan []storage.Post)
	chE := make(chan error)

	// Чтение RSS-лент из конфига с заданным интервалом
	go func() {
		err := rss.GoNews(configURL, chP, chE)
		if err != nil {
			log.Fatal(err)
		}
	}()

	// запись публикаций в db
	go func() {
		for posts := range chP {
			if err := a.DB.PostsMany(posts); err != nil {
				chE <- err
			}
		}
	}()

	// вывод ошибок
	go func() {
		for err := range chE {
			log.Println(err)
		}
	}()

	// запуск сервера с api
	err = http.ListenAndServe(":80", a.Router())
	if err != nil {
		log.Fatal(err)
	}
}
