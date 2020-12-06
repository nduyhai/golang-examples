package main

import (
	"github.com/go-chi/chi"
	"github.com/go-oauth2/oauth2/v4/errors"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-oauth2/oauth2/v4/store"
	httpSwagger "github.com/swaggo/http-swagger"
	_ "golang-examples/cmd/auth-service/docs"
	"log"
	"net/http"
)

// @title Auth Service API
// @version 1.0
// @description This is a auth server
// @termsOfService localhost:9096

// @contact.name API Support
// @contact.url localhost:9096
// @contact.email admin@localhost

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:9096
// @BasePath /api/v1
func main() {
	manager := manage.NewDefaultManager()
	// token memory store
	manager.MustTokenStorage(store.NewMemoryTokenStore())

	// client memory store
	clientStore := store.NewClientStore()
	_ = clientStore.Set("000000", &models.Client{
		ID:     "000000",
		Secret: "999999",
		Domain: "http://localhost",
	})
	manager.MapClientStorage(clientStore)

	srv := server.NewDefaultServer(manager)
	srv.SetAllowGetAccessRequest(true)
	srv.SetClientInfoHandler(server.ClientFormHandler)

	srv.SetInternalErrorHandler(func(err error) (re *errors.Response) {
		log.Println("Internal Error:", err.Error())
		return
	})

	srv.SetResponseErrorHandler(func(re *errors.Response) {
		log.Println("Response Error:", re.Error.Error())
	})

	router := chi.NewRouter()
	router.HandleFunc("/swagger/*", httpSwagger.Handler(httpSwagger.URL("http://localhost:9096/swagger/doc.json")))

	router.HandleFunc("/api/v1/authorize", authorize(srv))

	router.HandleFunc("/api/v1/token", token(srv))

	log.Fatal(http.ListenAndServe(":9096", router))
}

// Authorize
// @Summary Authorize
// @Description Authorize with Auth2
// @Accept  json
// @Produce  json
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} errors.Response
// @Failure 500 {object} errors.Response
// @Failure 404 {string} Not Found
// @Router /authorize [get]
func authorize(srv *server.Server) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := srv.HandleAuthorizeRequest(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}
}

// Retrieve access token
// @Summary Retrieve access token
// @Description  Retrieve auth2 access token
// @Accept  json
// @Produce  json
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} errors.Response
// @Failure 500 {object} errors.Response
// @Failure 404 {string} Not Found
// @Router /token [get]
func token(srv *server.Server) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := srv.HandleTokenRequest(w, r)
		log.Println("error during handle token request", err)
	}
}
