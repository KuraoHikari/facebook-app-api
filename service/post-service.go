package service

import (
	"fmt"
	"log"

	"github.com/KuraoHikari/facebook-app-res-api/dto"
	"github.com/KuraoHikari/facebook-app-res-api/entity"
	"github.com/KuraoHikari/facebook-app-res-api/repository"
	"github.com/mashingan/smapping"
)

type PostService interface {
	Insert(b dto.PostCreateDTO) entity.Post
	Update(b dto.PostUpdateDTO) entity.Post
	Delete(b entity.Post)
	All() []entity.Post
	FindByID(postID uint64) entity.Post
	IsAllowedToEdit(userID string, postID uint64) bool
}

type postService struct {
	postRepository repository.PostRepository
}

func NewBookService(postRepo repository.PostRepository) PostService {
	return &postService{
		postRepository: postRepo,
	}
}

func (service *postService) Insert(b dto.PostCreateDTO) entity.Post{
	post := entity.Post{}
	err := smapping.FillStruct(&post, smapping.MapFields(&b))
	if err != nil {
		log.Fatalf("Failed map %v: ", err)
	}
	res := service.postRepository.InsertPost(post)
	return res
}
func (service *postService)	Update(b dto.PostUpdateDTO) entity.Post{
	post := entity.Post{}
	err := smapping.FillStruct(&post, smapping.MapFields(&b))
	if err != nil {
		log.Fatalf("Failed map %v: ", err)
	}
	res := service.postRepository.UpdatePost(post)
	return res
}
func (service *postService)	Delete(b entity.Post){
	service.postRepository.DeletePost(b)
}
func (service *postService)	All() []entity.Post {
	return service.postRepository.AllPost()
}
func (service *postService)	FindByID(postID uint64) entity.Post{
	return service.postRepository.FindPostByID(postID)
}
func (service *postService)	IsAllowedToEdit(userID string, postID uint64) bool {
	post := service.postRepository.FindPostByID(postID)
	id := fmt.Sprintf("%v", post.UserID)
	return userID == id
}