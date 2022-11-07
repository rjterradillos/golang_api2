package repository

import (
	"github.com/rjterradillos/golang_api2/entity"
	"gorm.io/gorm"
)

// CourseRepository is a ....
type CourseRepository interface {
	InsertCourse(b entity.Course) entity.Course
	UpdateCourse(b entity.Course) entity.Course
	DeleteCourse(b entity.Course)
	AllCourse() []entity.Course
	FindCourseByID(courseID uint64) entity.Course
}

type courseConnection struct {
	connection *gorm.DB
}

// NewCourseRepository creates an instance CourseRepository
func NewCourseRepository(dbConn *gorm.DB) CourseRepository {
	return &courseConnection{
		connection: dbConn,
	}
}

func (db *courseConnection) InsertCourse(b entity.Course) entity.Course {
	db.connection.Save(&b)
	db.connection.Preload("User").Find(&b)
	return b
}

func (db *courseConnection) UpdateCourse(b entity.Course) entity.Course {
	db.connection.Save(&b)
	db.connection.Preload("User").Find(&b)
	return b
}

func (db *courseConnection) DeleteCourse(b entity.Course) {
	db.connection.Delete(&b)
}

func (db *courseConnection) FindCourseByID(courseID uint64) entity.Course {
	var course entity.Course
	db.connection.Preload("User").Find(&course, courseID)
	return course
}

func (db *courseConnection) AllCourse() []entity.Course {
	var courses []entity.Course
	db.connection.Preload("User").Find(&courses)
	return courses
}
