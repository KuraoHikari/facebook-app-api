package service

import (
	"log"

	"github.com/KuraoHikari/facebook-app-res-api/entity"
	"github.com/KuraoHikari/facebook-app-res-api/repository"
	"github.com/mashingan/smapping"
	"golang.org/x/crypto/bcrypt"

	"github.com/KuraoHikari/facebook-app-res-api/dto"
)

type UserService interface {
	VerifyCredential(email string, password string) interface{}
	CreateUser(user dto.RegisterDTO) entity.User
	FindByEmail(email string) entity.User
	IsDuplicateEmail(email string) bool
	Update(user dto.UserUpdateDTO) entity.User
	Profile(userID string) entity.User
}
type userService struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepository: userRepo,
	}
}
func (service *userService) VerifyCredential(email string, password string) interface{}{
	res := service.userRepository.VerifyCredential(email, password)
	if v, ok := res.(entity.User); ok {
		err := bcrypt.CompareHashAndPassword([]byte(v.Password), []byte(password) )
		if err != nil {
			log.Println(err)
			return false
		}else if v.Email != email {
			return false
		}
		return res
	}else {
		return false
	}
}
func (service *userService) CreateUser(user dto.RegisterDTO) entity.User{
	userToCreate := entity.User{}
	err := smapping.FillStruct(&userToCreate, smapping.MapFields(&user))
	if err != nil {
		log.Fatalf("Failed map %v", err)
	}
	res := service.userRepository.InsertUser(userToCreate)
	return res
}
func (service *userService) FindByEmail(email string) entity.User{
	return service.userRepository.FindByEmail(email)
}
func (service *userService) IsDuplicateEmail(email string) bool{
	res := service.userRepository.IsDuplicateEmail(email)
	return !(res.Error == nil)
}
func (service *userService) Update(user dto.UserUpdateDTO) entity.User {
	userToUpdate := entity.User{}
	err := smapping.FillStruct(&userToUpdate, smapping.MapFields(&user))
	if err != nil {
		log.Fatalf("Failed map %v:", err)
	}
	updatedUser := service.userRepository.UpdateUser(userToUpdate)
	return updatedUser
}

func (service *userService) Profile(userID string) entity.User {
	return service.userRepository.ProfileUser(userID)
}
