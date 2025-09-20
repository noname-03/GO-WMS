package service

import (
	"errors"
	"log"
	"myapp/internal/model"
	"myapp/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService() *UserService {
	return &UserService{
		userRepo: repository.NewUserRepository(),
	}
}

// Business logic methods
func (s *UserService) GetAllUsers() ([]model.User, error) {
	return s.userRepo.GetAllUsers()
}

func (s *UserService) GetUsersMinimal() ([]repository.UserMinimal, error) {
	return s.userRepo.GetUsersMinimal()
}

func (s *UserService) GetUserByID(id uint) (*model.User, error) {
	user, err := s.userRepo.GetUserByID(id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (s *UserService) CreateUser(name, email, password string) (*model.User, error) {
	// Check if email already exists using the new method
	exists, err := s.userRepo.CheckEmailExists(email, 0)
	if err != nil {
		log.Printf("Error checking email existence: %v", err)
		return nil, err
	}
	if exists {
		log.Printf("Email already exists: %s", email)
		return nil, errors.New("email already exists")
	}

	log.Printf("Email is available, creating new user for: %s", email)

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Name:     name,
		Email:    email,
		Password: string(hashedPassword),
	}

	err = s.userRepo.CreateUser(user)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		return nil, err
	}

	log.Printf("User created successfully with ID: %d", user.ID)
	return user, nil
}

func (s *UserService) AuthenticateUser(email, password string) (*model.User, error) {
	user, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Check password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}

func (s *UserService) SearchUsers(keyword string, limit, offset int) ([]model.User, error) {
	return s.userRepo.SearchUsersRaw(keyword, limit, offset)
}

func (s *UserService) GetUsersStats() (*repository.UserStats, error) {
	return s.userRepo.GetUsersStats()
}

func (s *UserService) GetUsersWithRawSQL() ([]model.User, error) {
	return s.userRepo.GetUsersWithRawSQL()
}

func (s *UserService) GetUsersWithStats() ([]repository.UserResult, error) {
	return s.userRepo.GetUsersWithStats()
}
