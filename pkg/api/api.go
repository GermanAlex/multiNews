package api

import (
	"encoding/json"
	"multiNews/pkg/storage"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type API struct {
	db *storage.DB
	r  *mux.Router
}

// new api метод как точка старта
func New(db *storage.DB) *API {
	api := API{db: db, r: mux.NewRouter()}
	api.endpoints() // определяем хендлеры
	return &api
}

func (api *API) Router() *mux.Router {
	return api.r
}

// определяем эндпоинты, у нас их 2 - получение списка из rss-потока и приложение, созданное экспертом курса
func (api *API) endpoints() {
	api.r.HandleFunc("/news/{n}", api.listnews).Methods(http.MethodGet)
	api.r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./cmd/app/webapp"))))
}

// handle для добавления новостей

func (api *API) listnews(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	param := mux.Vars(r)["n"]
	n, _ := strconv.Atoi(param)

	news, err := api.db.ListNews(n)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(news)
}
