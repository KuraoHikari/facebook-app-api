package repository

import (
	"github.com/KuraoHikari/facebook-app-res-api/entity"
	"gorm.io/gorm"
)

type PostRepository interface {
	InsertPost(b entity.Post)entity.Post
	UpdatePost(b entity.Post)entity.Post
	FindPostByID(postID uint64) entity.Post
	DeletePost(b entity.Post)
	AllPost()[]entity.Post
}
type postConnection struct {
	connection *gorm.DB
}

func NewPostRepository(dbConn *gorm.DB) PostRepository {
	return &postConnection{
		connection: dbConn,
	}
}

func(db *postConnection)InsertPost(b entity.Post)entity.Post {
	db.connection.Save(&b)
	db.connection.Preload("User").Find(&b)
	return b
}
func(db *postConnection)UpdatePost(b entity.Post)entity.Post{
	db.connection.Save(&b)
	db.connection.Preload("User").Find(&b)
	return b
}
func(db *postConnection)FindPostByID(postID uint64)entity.Post{
	var post entity.Post
	db.connection.Preload("User").Find(&post, postID)
	return post
}
func(db *postConnection)DeletePost(b entity.Post){
	db.connection.Delete(&b)
}
func(db *postConnection)AllPost()[]entity.Post{
	var posts []entity.Post
	db.connection.Preload("User").Find(&posts)
	return posts
}
