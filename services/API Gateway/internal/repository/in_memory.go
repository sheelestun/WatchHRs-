package repository

import (
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/sheelestun/WatchHRs-/internal/entity"
)

type InMemoryStorage struct {
	managers     map[uuid.UUID]*entity.Manager
	employees    map[uuid.UUID]*entity.Employee
	photos       map[uuid.UUID]*entity.Photo
	screenshots  map[uuid.UUID][]*entity.ScreenshotStatistic
	workSessions map[uuid.UUID][]*entity.WorkSession
	validate     *validator.Validate
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{validate: validator.New()}
}

func (i *InMemoryStorage) AddManager(manager entity.Manager) (uuid.UUID, error) {
	manager.Id = uuid.New()
	if err := i.validate.Struct(manager); err != nil {
		return uuid.Nil, err
	}
	i.managers[manager.Id] = &manager
	return manager.Id, nil
}

func (i *InMemoryStorage) RemoveManager(managerId uuid.UUID) error {
	if _, exists := i.managers[managerId]; !exists {
		return errors.New("manager does not exist")
	}
	delete(i.managers, managerId)
	delete(i.photos, managerId)
	return nil
}

func (i *InMemoryStorage) AddEmployee(employee entity.Employee) (uuid.UUID, error) {
	employee.Id = uuid.New()
	if err := i.validate.Struct(employee); err != nil {
		return uuid.Nil, err
	}
	i.employees[employee.Id] = &employee
	return employee.Id, nil
}

func (i *InMemoryStorage) RemoveEmployee(employeeId uuid.UUID) error {
	if _, exists := i.employees[employeeId]; !exists {
		return errors.New("employee does not exist")
	}
	delete(i.employees, employeeId)
	delete(i.photos, employeeId)
	return nil
}

func (i *InMemoryStorage) AddPhoto(photo entity.Photo) (uuid.UUID, error) {
	photo.Id = photo.UserId
	if err := i.validate.Struct(photo); err != nil {
		return uuid.Nil, err
	}
	i.photos[photo.Id] = &photo
	return photo.Id, nil
}

func (i *InMemoryStorage) AddScreenshotStatistic(screenshot entity.ScreenshotStatistic) (uuid.UUID, error) {
	screenshot.Id = uuid.New()
	screenshot.Timestamp = time.Now()
	if err := i.validate.Struct(screenshot); err != nil {
		return uuid.Nil, err
	}
	screenshots := i.screenshots[screenshot.Id]
	screenshots = append(screenshots, &screenshot)
	i.screenshots[screenshot.EmployeeId] = screenshots
	return screenshot.Id, nil
}

func (i *InMemoryStorage) GetScreenshotsStatistic(employeeId uuid.UUID, date time.Time) ([]entity.ScreenshotStatistic, error) {
	resScreenshots := make([]entity.ScreenshotStatistic, 0)

	if screenshots, exists := i.screenshots[employeeId]; exists {
		for _, screenshot := range screenshots {
			if date == screenshot.Timestamp {
				resScreenshots = append(resScreenshots, *screenshot)
			}
		}
		return resScreenshots, nil
	}
	return resScreenshots, errors.New("screenshot does not exist by this employee and date")
}

func (i *InMemoryStorage) StartWorkSession(employeeId uuid.UUID) (uuid.UUID, error) {
	newSession := entity.WorkSession{Id: uuid.New(), EmployeeId: employeeId, StartTime: time.Now()}
	if err := i.validate.Struct(newSession); err != nil {
		return uuid.Nil, err
	}
	workSessions := i.workSessions[employeeId]
	if lastWorkSession := workSessions[len(workSessions)-1]; lastWorkSession.EndTime.IsZero() {
		return uuid.Nil, errors.New("last work session did not stop yet")
	}
	workSessions = append(workSessions, &newSession)
	i.workSessions[newSession.EmployeeId] = workSessions
	return newSession.Id, nil
}

func (i *InMemoryStorage) StopWorkSession(employeeId uuid.UUID) (uuid.UUID, error) {
	workSessions := i.workSessions[employeeId]
	lastWorkSession := workSessions[len(workSessions)-1]
	if !lastWorkSession.EndTime.IsZero() {
		return uuid.Nil, errors.New("last work session did not start yet")
	}
	lastWorkSession.EndTime = time.Now()
	lastWorkSession.TotalTime = lastWorkSession.EndTime.Sub(lastWorkSession.StartTime)
	workSessions[len(workSessions)-1] = lastWorkSession
	i.workSessions[employeeId] = workSessions
	return lastWorkSession.Id, nil
}

func (i *InMemoryStorage) GetWorkSessions(employeeId uuid.UUID, date time.Time) ([]entity.WorkSession, error) {
	resWorkSessions := make([]entity.WorkSession, 0)
	if workSessions, exists := i.workSessions[employeeId]; exists {
		for _, workSession := range workSessions {
			if date.Year() == workSession.StartTime.Year() && date.Month() == workSession.StartTime.Month() && date.Day() == workSession.StartTime.Day() {
				resWorkSessions = append(resWorkSessions, *workSession)
			}
		}
		return resWorkSessions, nil
	}
	return resWorkSessions, errors.New("work session does not exist by this employee and date")
}
