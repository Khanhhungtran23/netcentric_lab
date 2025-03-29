package storage

import (
	"sync"
	"os"
	"errors"
	"encoding/json"
	"encoding/gob"
	"path/filepath"

	"socket-tcp/internal/model"
	"socket-tcp/internal/auth"
)

type StorageType string

const (
	JSONStorage StorageType = "json"
	GOBStorage	StorageType = "gob"
)

type UserStorage struct {
	filePath		string
	storageType		StorageType
	mu				sync.RWMutex
}

// NewUserStorage create a new UserStorage
func NewUserStorage(filePath string, storageType StorageType) *UserStorage {
	// Register type of GOB encoding
	gob.Register(model.User{})
	gob.Register(model.Address{})

	return &UserStorage{
		filePath: filePath,
		storageType: storageType,
		mu: sync.RWMutex{},
	}
}

// LoadUsers loads users from the storage file
func (us *UserStorage) LoadUsers() ([]*model.User, error) {
	us.mu.RLock()
	defer us.mu.RUnlock()

	// Check if the file exists
	if _, err := os.Stat(us.filePath); os.IsNotExist(err) {
		return us.CreateDefaultUsers()
	}

	// Open file
	file, err := os.Open(us.filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var users []*model.User

	switch us.storageType {
	case JSONStorage:
		decoder := json.NewDecoder(file)
		err = decoder.Decode(&users)
	case GOBStorage:
		decoder := gob.NewDecoder(file)
		err = decoder.Decode(&users)
	default: 
		return nil, errors.New("Unsupported storage type")
	}
	
	if err != nil {
		return nil, err
	}

	return users, nil
}

// CreateDefaultUsers creates default users
func (us *UserStorage) CreateDefaultUsers() ([]*model.User, error) {
	users := []*model.User{
		{
			Username: "admin",
			Password: auth.EncryptPassword("123"),
			Fullname: "Admin",
			Emails: []string{"admin@gmail.com"},
			Addresses: []model.Address{
				{
					Type: "work",
					Details: "Admin Office",
				},
			},
		},
		{
			Username: "user1",
			Password: auth.EncryptPassword("user123"),
			Fullname: "Test User",
			Emails: []string{"user1@example.com", "user1.alt@example.com"},
			Addresses: []model.Address{
				{
					Type: "home",
					Details: "123 Main St",
				},
				{
					Type: "work",
					Details: "456 Work Ave",
				},
			},
		},
	}

	if err := us.SaveUsers(users); err != nil {
		return nil, errors.New("Failed to save default users to file/dir!")
	}

	return users, nil
}


func (us *UserStorage) SaveUsers(users []*model.User) error {
	us.mu.RLock()
	defer us.mu.RUnlock()

	// create a new directory if it does not exists
	dir := filepath.Dir(us.filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return errors.New("Failed to make dir to save user")
	}

	// create or truncate file
	file, err := os.Create(us.filePath)
	if err != nil {
		return errors.New("Failed to create file to save user")
	}

	defer file.Close()

	// Encode based on the storage type
	switch us.storageType {
	case JSONStorage:
		encoder := json.NewEncoder(file)
		encoder.SetIndent("", " ")
		err = encoder.Encode(users)
	case GOBStorage:
		encoder := gob.NewEncoder(file)
		err = encoder.Encode(users)
	default:
		return errors.New("Unsupported storage type")
	}

	return err
}

