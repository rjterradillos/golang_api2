package service

import (
	"fmt"
	"log"

	"github.com/mashingan/smapping"
	"github.com/rjterradillos/golang_api2/dto"
	"github.com/rjterradillos/golang_api2/entity"
	"github.com/rjterradillos/golang_api2/repository"
)

// CourseService is a ....
type CourseService interface {
	Insert(b dto.CourseCreateDTO) entity.Course
	Update(b dto.CourseUpdateDTO) entity.Course
	Delete(b entity.Course)
	All() []entity.Course
	FindByID(courseID uint64) entity.Course
	IsAllowedToEdit(userID string, courseID uint64) bool
}

type courseService struct {
	courseRepository repository.CourseRepository
}

// NewCourseService .....
func NewCourseService(courseRepo repository.CourseRepository) CourseService {
	return &courseService{
		courseRepository: courseRepo,
	}
}

func (service *courseService) Insert(b dto.CourseCreateDTO) entity.Course {
	course := entity.Course{}
	err := smapping.FillStruct(&course, smapping.MapFields(&b))
	if err != nil {
		log.Fatalf("Failed map %v: ", err)
	}
	res := service.courseRepository.InsertCourse(course)
	return res
}

func (service *courseService) Update(b dto.CourseUpdateDTO) entity.Course {
	course := entity.Course{}
	err := smapping.FillStruct(&course, smapping.MapFields(&b))
	if err != nil {
		log.Fatalf("Failed map %v: ", err)
	}
	res := service.courseRepository.UpdateCourse(course)
	return res
}

func (service *courseService) Delete(b entity.Course) {
	service.courseRepository.DeleteCourse(b)
}

func (service *courseService) All() []entity.Course {
	return service.courseRepository.AllCourse()
}

func (service *courseService) FindByID(courseID uint64) entity.Course {
	return service.courseRepository.FindCourseByID(courseID)
}

func (service *courseService) IsAllowedToEdit(userID string, courseID uint64) bool {
	b := service.courseRepository.FindCourseByID(courseID)
	id := fmt.Sprintf("%v", b.UserID)
	return userID == id
}
