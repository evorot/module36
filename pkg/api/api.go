package api

import (
	"GoNews/pkg/storage"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type API struct {
	DB   storage.Interface
	Rout *mux.Router
}

// Создание объекта api
func New(db storage.Interface) *API {
	api := API{
		DB: db,
	}
	api.Rout = mux.NewRouter()
	api.endpoints()
	return &api

}

// Регистрация обработчиков API.
func (api *API) endpoints() {
	api.Rout.HandleFunc("/news/{quantity}", api.newsHandler).Methods(http.MethodGet, http.MethodOptions)
	api.Rout.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./webapp"))))
}

// Получение маршрутизатора запросов.
// Требуется для передачи маршрутизатора веб-серверу.
func (api *API) Router() *mux.Router {
	return api.Rout
}

// Handler,  который выводит заданное кол-во новостей.
// Требуемое количество публикаций указывается в пути запроса
func (api *API) newsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	v := mux.Vars(r)["quantity"]
	quantity, err := strconv.Atoi(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	posts, err := api.DB.Posts(quantity)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = json.NewEncoder(w).Encode(posts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
