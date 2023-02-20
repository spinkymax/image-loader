package repository

import (
	"encoding/json"
	"fmt"
	"image-loader/internal/model"
	"io"
	"os"
	"sync"
)

type User struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type UserRepo struct {
	mut      *sync.RWMutex
	filename string
}

func NewUserRepo(filename string) *UserRepo {
	return &UserRepo{
		filename: filename,
		mut:      &sync.RWMutex{},
	}
}

func (u *UserRepo) AddUser(user model.User) error {
	u.mut.Lock()
	defer u.mut.Unlock()
	file, err := os.OpenFile(u.filename, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return fmt.Errorf("couldn't open  file: %w", err)
	}

	defer file.Close()

	b, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("couldn't read file: %w", err)
	}

	users := make([]User, 0)
	if len(b) != 0 {
		err = json.Unmarshal(b, &users)
		if err != nil {
			return fmt.Errorf("failed to unmarshal users: %w", err)
		}
	}

	users = append(users, User(user))

	_, err = file.Seek(0, 0)
	if err != nil {
		return fmt.Errorf("failed to return beggining of the file: %w", err)
	}

	b, err = json.MarshalIndent(&users, "\t", "")
	if err != nil {
		return fmt.Errorf("failed to marshal users: %w", err)
	}

	_, err = file.Write(b)
	if err != nil {
		return fmt.Errorf("failed to write users to file: %w", err)
	}

	return nil
}

func (u *UserRepo) GetUser(id int) (model.User, error) {
	u.mut.RLock()
	defer u.mut.RUnlock()
	file, err := os.OpenFile(u.filename, os.O_RDONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return model.User{}, fmt.Errorf("couldn't open file: %w", err)
	}

	defer file.Close()

	decoder := json.NewDecoder(file)

	_, err = decoder.Token()
	if err != nil {
		return model.User{}, fmt.Errorf("failed to get first json token: %w", err)
	}

	for decoder.More() {
		var user User
		err = decoder.Decode(&user)
		if err != nil {
			return model.User{}, fmt.Errorf("failed to decode user: %w", err)
		}
		if user.Id == id {
			return model.User(user), nil
		}
	}
	return model.User{}, fmt.Errorf("couldn't find the user")
}
