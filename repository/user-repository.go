package repository

import (
	"log"

	"github.com/KuraoHikari/facebook-app-res-api/entity"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

//UserRepo
//Auth Service
//Auth Controller
//api login regis
//user service
//user module
//POstRepo
///Post service
//Post Controller
//post hastMany relation
//book module

type UserRepository interface {
	InsertUser(user entity.User) entity.User
	UpdateUser(user entity.User) entity.User
	VerifyCredential(email string, password string) interface{}
	IsDuplicateEmail(email string) (tx *gorm.DB)
	FindByEmail(email string) entity.User
	ProfileUser(userID string) entity.User
}

type userConnection struct {
	connection *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userConnection{
		connection: db,
	}
}


func(db *userConnection) InsertUser(user entity.User) entity.User{
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.MinCost)
	if err != nil {
		log.Println(err)
		panic("Failed to hash a password")
	}
	user.Password = string(hash)
	db.connection.Save(&user)
	return user
}
func(db *userConnection) UpdateUser(user entity.User) entity.User{
	if user.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.MinCost)
		if err != nil {
			log.Println(err)
			panic("Failed to hash a password")
		}
		user.Password = string(hash)
	}else {
		var tempUser entity.User
		db.connection.Find(&tempUser, user.ID)
		user.Password = tempUser.Password
	}
	db.connection.Save(&user)
	return user
}
func(db *userConnection) VerifyCredential(email string, password string) interface{} {
	var user entity.User
	res := db.connection.Where("email = ?", email).Take(&user)
	if res.Error == nil {
		return user
	}
	return nil
}
func(db *userConnection) IsDuplicateEmail(email string) (tx *gorm.DB){
	var user entity.User
	return db.connection.Where("email = ?", email).Take(&user)
}
func(db *userConnection) FindByEmail(email string) entity.User{
	var user entity.User
	 db.connection.Where("email = ?", email).Take(&user)
	 return user
}
func(db *userConnection) ProfileUser(userID string) entity.User{
	var user entity.User
	db.connection.Preload("Posts").Preload("Posts.User").Find(&user, userID)
	return user
}