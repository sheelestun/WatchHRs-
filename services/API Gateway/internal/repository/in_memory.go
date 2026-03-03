package repository

import (
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/sheelestun/WatchHRs-/internal/entity"
	log "github.com/sirupsen/logrus"
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
	return &InMemoryStorage{
		managers:     make(map[uuid.UUID]*entity.Manager),
		employees:    make(map[uuid.UUID]*entity.Employee),
		photos:       make(map[uuid.UUID]*entity.Photo),
		screenshots:  make(map[uuid.UUID][]*entity.ScreenshotStatistic),
		workSessions: make(map[uuid.UUID][]*entity.WorkSession),
		validate:     validator.New()}
}

func (i *InMemoryStorage) AddManager(manager entity.Manager) (uuid.UUID, error) {
	manager.ID = uuid.New()
	if err := i.validate.Struct(manager); err != nil {
		return uuid.Nil, err
	}
	i.managers[manager.ID] = &manager
	log.Debugf("Added manager: %+v", manager)
	return manager.ID, nil
}

func (i *InMemoryStorage) RemoveManager(managerID uuid.UUID) error {
	if _, exists := i.managers[managerID]; !exists {
		return errors.New("manager does not exist")
	}
	delete(i.managers, managerID)
	delete(i.photos, managerID)
	log.Debugf("Removed manager: %+v", managerID)
	return nil
}

func (i *InMemoryStorage) AddEmployee(employee entity.Employee) (uuid.UUID, error) {
	employee.ID = uuid.New()
	if err := i.validate.Struct(employee); err != nil {
		return uuid.Nil, err
	}
	i.employees[employee.ID] = &employee
	log.Debugf("Added employee: %+v", employee)
	return employee.ID, nil
}

func (i *InMemoryStorage) GetAllEmployeesByManagerID(managerID uuid.UUID) ([]entity.Employee, error) {
	employees := make([]entity.Employee, 0)
	for _, employee := range i.employees {
		if employee.ManagerID == managerID {
			employees = append(employees, *employee)
		}
	}
	if len(employees) == 0 {
		return employees, errors.New("no employees found")
	}
	log.Debugf("Found employees: %+v", employees)
	return employees, nil
}

func (i *InMemoryStorage) RemoveEmployee(employeeID uuid.UUID) error {
	if _, exists := i.employees[employeeID]; !exists {
		return errors.New("employee does not exist")
	}
	delete(i.employees, employeeID)
	delete(i.photos, employeeID)
	log.Debugf("Removed employee: %+v", employeeID)
	return nil
}

func (i *InMemoryStorage) AddPhoto(photo entity.Photo) (uuid.UUID, error) {
	photo.ID = photo.UserID
	if err := i.validate.Struct(photo); err != nil {
		return uuid.Nil, err
	}
	i.photos[photo.ID] = &photo
	log.Debugf("Added photo: %+v", photo)
	return photo.ID, nil
}

func (i *InMemoryStorage) AddScreenshotStatistic(screenshot entity.ScreenshotStatistic) (uuid.UUID, error) {
	screenshot.ID = uuid.New()
	screenshot.Timestamp = time.Now()
	if err := i.validate.Struct(screenshot); err != nil {
		return uuid.Nil, err
	}
	screenshots := i.screenshots[screenshot.ID]
	screenshots = append(screenshots, &screenshot)
	i.screenshots[screenshot.EmployeeID] = screenshots
	log.Debugf("Added screenshot: %+v", screenshot)
	return screenshot.ID, nil
}

func (i *InMemoryStorage) GetScreenshotsStatistic(employeeID uuid.UUID, date time.Time) ([]entity.ScreenshotStatistic, error) {
	resScreenshots := make([]entity.ScreenshotStatistic, 0)

	if screenshots, exists := i.screenshots[employeeID]; exists {
		for _, screenshot := range screenshots {
			if date == screenshot.Timestamp {
				resScreenshots = append(resScreenshots, *screenshot)
			}
		}
		log.Debugf("Found screenshots: %+v", resScreenshots)
		return resScreenshots, nil
	}
	return resScreenshots, errors.New("screenshot does not exist by this employee and date")
}

func (i *InMemoryStorage) StartWorkSession(employeeID uuid.UUID) (uuid.UUID, error) {
	newSession := entity.WorkSession{ID: uuid.New(), EmployeeID: employeeID, StartTime: time.Now()}
	if err := i.validate.Struct(newSession); err != nil {
		return uuid.Nil, err
	}
	workSessions := i.workSessions[employeeID]
	if len(workSessions) != 0 {
		if lastWorkSession := workSessions[len(workSessions)-1]; lastWorkSession.EndTime.IsZero() {
			return uuid.Nil, errors.New("last work session did not stop yet")
		}
	}
	workSessions = append(workSessions, &newSession)
	i.workSessions[newSession.EmployeeID] = workSessions
	log.Debugf("Started work session: %+v", newSession)
	return newSession.ID, nil
}

func (i *InMemoryStorage) StopWorkSession(employeeID uuid.UUID) (uuid.UUID, error) {
	workSessions := i.workSessions[employeeID]
	if len(workSessions) == 0 {
		return uuid.Nil, errors.New("work sessions do not exist")
	}

	lastWorkSession := workSessions[len(workSessions)-1]
	if !lastWorkSession.EndTime.IsZero() {
		return uuid.Nil, errors.New("last work session did not start yet")
	}
	lastWorkSession.EndTime = time.Now()
	lastWorkSession.TotalTime = lastWorkSession.EndTime.Sub(lastWorkSession.StartTime)
	workSessions[len(workSessions)-1] = lastWorkSession
	i.workSessions[employeeID] = workSessions
	log.Debugf("Stopped work session: %+v", lastWorkSession)
	return lastWorkSession.ID, nil
}

func (i *InMemoryStorage) GetWorkSessions(employeeID uuid.UUID, date time.Time) ([]entity.WorkSession, error) {
	resWorkSessions := make([]entity.WorkSession, 0)
	if workSessions, exists := i.workSessions[employeeID]; exists {
		for _, workSession := range workSessions {
			if date.Year() == workSession.StartTime.Year() && date.Month() == workSession.StartTime.Month() && date.Day() == workSession.StartTime.Day() {
				resWorkSessions = append(resWorkSessions, *workSession)
			}
		}
		log.Debugf("Found work sessions: %+v", resWorkSessions)
		return resWorkSessions, nil
	}
	return resWorkSessions, errors.New("work session does not exist by this employee and date")
}
