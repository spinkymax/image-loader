package server

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
	"image-loader/internal/model"
	"io"
	"net/http"
	"strconv"
)

type User struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type controller interface {
	AddUser(user model.User) error
	GetUser(id int) (model.User, error)
}

type Server struct {
	listenURI string
	logger    *logrus.Logger
	r         chi.Router
	c         controller
}

func NewServer(listenURI string, logger *logrus.Logger, c controller) *Server {
	return &Server{
		listenURI: listenURI,
		logger:    logger,
		r:         chi.NewRouter(),
		c:         c,
	}
}

func (s *Server) RegisterRoutes() {
	s.r.Get("/user/{userID}", s.HandleGetUser)
	s.r.Post("/user/add", s.HandleAddUser)
}

func (s *Server) StartRouter() {
	srv := http.Server{
		Addr:    s.listenURI,
		Handler: s.r,
	}
	s.logger.Info("server is running!")
	err := srv.ListenAndServe()
	if err != nil {
		s.logger.Fatal(err)
	}
}

func (s *Server) HandleAddUser(w http.ResponseWriter, r *http.Request) {
	var user User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		s.handleError(err, http.StatusBadRequest, w)
		return
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			s.logger.Error(err)
		}
	}(r.Body)

	err = s.c.AddUser(model.User(user))

	if err != nil {
		s.handleError(err, http.StatusInternalServerError, w)
		return
	}

	w.WriteHeader(http.StatusOK)

}

func (s *Server) HandleGetUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "userID")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		s.handleError(err, http.StatusBadRequest, w)
		return
	}

	user, err := s.c.GetUser(id)
	if err != nil {
		s.handleError(err, http.StatusInternalServerError, w)
		return
	}

	b, err := json.Marshal(&user)
	if err != nil {
		s.handleError(err, http.StatusInternalServerError, w)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(b)
	if err != nil {
		s.logger.Error(err)
	}
}

func (s *Server) handleError(err error, status int, w http.ResponseWriter) {
	s.logger.Error(err)
	w.WriteHeader(status)
	_, err = w.Write([]byte(err.Error()))
	if err != nil {
		s.logger.Error(err)
	}

}
