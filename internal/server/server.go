package server

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
	"github.com/spinkymax/image-loader/internal/config"
	"github.com/spinkymax/image-loader/internal/middleware"
	"github.com/spinkymax/image-loader/internal/model"
	"github.com/spinkymax/image-loader/internal/response"
	"io"
	"net/http"
	"strconv"
)

type User struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Login       string `json:"login"`
	Password    string `json:"password"`
	Description string `json:"description"`
}

type controller interface {
	AddUser(ctx context.Context, user model.User) error
	GetUser(ctx context.Context, id int64) (model.User, error)
	UpdateUser(ctx context.Context, modelUser model.User) error
	DeleteUser(ctx context.Context, id int64) error
	GetAllUsers(ctx context.Context) ([]model.User, error)
	Authorize(ctx context.Context, login, password string) (string, error)
	AddFile(ctx context.Context, filename string, file io.Reader) error
}

type Server struct {
	listenURI string
	logger    *logrus.Logger
	r         chi.Router
	c         controller
	cfg       *config.Config
}

func NewServer(listenURI string, logger *logrus.Logger, c controller, cfg *config.Config) *Server {
	return &Server{
		listenURI: listenURI,
		logger:    logger,
		r:         chi.NewRouter(),
		c:         c,
		cfg:       cfg,
	}
}

func (s *Server) RegisterRoutes() {
	s.r.Use(middleware.Logger(s.logger))

	s.r.Get("/user/auth", s.HandleAuthorize)
	s.r.Post("/user/add", s.HandleAddUser)
	s.r.Post("/image/add", s.HandleAddFile)

	s.r.Group(func(r chi.Router) {
		r.Use(middleware.Auth(s.cfg.JWTKeyword, s.logger))

		r.Get("/user/{userID}", s.HandleGetUser)
		r.Put("/user/update", s.HandleUpdateUser)
		r.Delete("/user/{userID}", s.HandleDeleteUser)
		r.Get("/user", s.HandleGetAllUsers)

	})

}

func (s *Server) StartServer() {
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

	err = s.c.AddUser(r.Context(), user.toModel())
	if err != nil {
		s.handleError(err, http.StatusInternalServerError, w)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (s *Server) HandleGetUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "userID")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		s.handleError(err, http.StatusBadRequest, w)
		return
	}

	user, err := s.c.GetUser(context.Background(), id)
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

func (s *Server) HandleUpdateUser(w http.ResponseWriter, r *http.Request) {
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

	err = s.c.UpdateUser(r.Context(), user.toModel())
	if err != nil {
		s.handleError(err, http.StatusInternalServerError, w)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) HandleDeleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "userID")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		s.handleError(err, http.StatusBadRequest, w)
		return
	}

	err = s.c.DeleteUser(r.Context(), id)
	if err != nil {
		s.handleError(err, http.StatusInternalServerError, w)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) HandleGetAllUsers(w http.ResponseWriter, r *http.Request) {
	user, err := s.c.GetAllUsers(context.Background())
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
		s.handleError(err, http.StatusInternalServerError, w)
		return
	}

}

func (s *Server) HandleAuthorize(w http.ResponseWriter, r *http.Request) {
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

	token, err := s.c.Authorize(r.Context(), user.Login, user.Password)
	if err != nil {
		s.handleError(err, http.StatusInternalServerError, w)
		return
	}

	b, err := response.ParseResponse(token, false)

	w.WriteHeader(http.StatusOK)

	_, err = w.Write(b)
	if err != nil {
		s.handleError(err, http.StatusInternalServerError, w)
		return
	}
}

func (s *Server) HandleAddFile(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("fileKey")
	if err != nil {
		s.handleError(err, http.StatusBadRequest, w)
	}

	err = s.c.AddFile(r.Context(), header.Filename, file)
	if err != nil {
		s.handleError(err, http.StatusInternalServerError, w)
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) handleError(err error, status int, w http.ResponseWriter) {
	s.logger.Error(err)
	w.WriteHeader(status)

	b, err := response.ParseResponse(err.Error(), true)
	if err != nil {
		s.logger.Error(err)
	}

	_, err = w.Write(b)
	if err != nil {
		s.logger.Error(err)
	}

}

func (u User) toModel() model.User {
	return model.User{
		ID:          u.ID,
		Name:        u.Name,
		Description: u.Description,
		Login:       u.Login,
		Password:    u.Password,
	}
}
