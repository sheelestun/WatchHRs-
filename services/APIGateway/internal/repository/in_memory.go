package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sheelestun/WatchHRs-/internal/entity"
	log "github.com/sirupsen/logrus"
)

type refreshToken struct {
	UserID    string
	ExpiresAt time.Time
}

type InMemoryStorage struct {
	managers     map[uuid.UUID]*entity.Manager
	employees    map[uuid.UUID]*entity.Employee
	photos       map[uuid.UUID]*entity.Photo
	screenshots  map[uuid.UUID][]*entity.ScreenshotStatistic
	workSessions map[uuid.UUID][]*entity.WorkSession

	refreshTokens map[string]refreshToken
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		managers:      make(map[uuid.UUID]*entity.Manager),
		employees:     make(map[uuid.UUID]*entity.Employee),
		photos:        make(map[uuid.UUID]*entity.Photo),
		screenshots:   make(map[uuid.UUID][]*entity.ScreenshotStatistic),
		workSessions:  make(map[uuid.UUID][]*entity.WorkSession),
		refreshTokens: make(map[string]refreshToken),
	}
}

func (i *InMemoryStorage) FindUser(ctx context.Context, userId uuid.UUID) (string, error) {
	_, ok := i.managers[userId]
	if ok {
		return "manager", nil
	}
	_, ok = i.employees[userId]
	if ok {
		return "employee", nil
	}

	return "", fmt.Errorf(`user "%s" not found`, userId)
}

func (i *InMemoryStorage) AddManager(ctx context.Context, manager entity.Manager) (uuid.UUID, error) {
	i.managers[manager.ID] = &manager
	log.Debugf("Added manager: %+v", manager)
	return manager.ID, nil
}

func (i *InMemoryStorage) RemoveManager(ctx context.Context, managerID uuid.UUID) error {
	if _, exists := i.managers[managerID]; !exists {
		return errors.New("manager does not exist")
	}
	delete(i.managers, managerID)
	delete(i.photos, managerID)
	log.Debugf("Removed manager: %+v", managerID)
	return nil
}

func (i *InMemoryStorage) AddEmployee(ctx context.Context, employee entity.Employee) (uuid.UUID, error) {
	i.employees[employee.ID] = &employee
	log.Debugf("Added employee: %+v", employee)
	return employee.ID, nil
}

func (i *InMemoryStorage) GetAllEmployeesByManagerID(ctx context.Context, managerID uuid.UUID) ([]entity.Employee, error) {
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

func (i *InMemoryStorage) RemoveEmployee(ctx context.Context, employeeID uuid.UUID) error {
	if _, exists := i.employees[employeeID]; !exists {
		return errors.New("employee does not exist")
	}
	delete(i.employees, employeeID)
	delete(i.photos, employeeID)
	log.Debugf("Removed employee: %+v", employeeID)
	return nil
}

func (i *InMemoryStorage) AddPhoto(ctx context.Context, photo entity.Photo) (uuid.UUID, error) {
	i.photos[photo.ID] = &photo
	log.Debugf("Added photo: %+v", photo)
	return photo.ID, nil
}

func (i *InMemoryStorage) AddScreenshotStatistic(ctx context.Context, screenshot entity.ScreenshotStatistic) (uuid.UUID, error) {
	screenshots := i.screenshots[screenshot.ID]
	screenshots = append(screenshots, &screenshot)
	i.screenshots[screenshot.EmployeeID] = screenshots
	log.Debugf("Added screenshot: %+v", screenshot)
	return screenshot.ID, nil
}

func (i *InMemoryStorage) GetScreenshotsStatistic(ctx context.Context, employeeID uuid.UUID, date time.Time) ([]entity.ScreenshotStatistic, error) {
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

func (i *InMemoryStorage) StartWorkSession(ctx context.Context, session entity.WorkSession) (uuid.UUID, error) {
	workSessions := i.workSessions[session.EmployeeID]
	if len(workSessions) != 0 {
		if lastWorkSession := workSessions[len(workSessions)-1]; lastWorkSession.EndTime.IsZero() {
			return uuid.Nil, errors.New("last work session did not stop yet")
		}
	}
	workSessions = append(workSessions, &session)
	i.workSessions[session.EmployeeID] = workSessions
	log.Debugf("Started work session: %+v", session)
	return session.ID, nil
}

func (i *InMemoryStorage) StopWorkSession(ctx context.Context, employeeID uuid.UUID) (uuid.UUID, error) {
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

func (i *InMemoryStorage) GetWorkSessions(ctx context.Context, employeeID uuid.UUID, date time.Time) ([]entity.WorkSession, error) {
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

func (i *InMemoryStorage) SaveTokenInCache(ctx context.Context, tokenID, userID string, expiresAt time.Time) error {
	i.refreshTokens[tokenID] = refreshToken{
		UserID:    userID,
		ExpiresAt: expiresAt,
	}

	return nil
}

func (i *InMemoryStorage) ExistsTokenInCache(ctx context.Context, tokenID string) (bool, error) {
	token, ok := i.refreshTokens[tokenID]
	if !ok {
		return false, nil
	}

	if time.Now().After(token.ExpiresAt) {
		return false, nil
	}

	return true, nil
}

func (i *InMemoryStorage) DeleteTokenInCache(ctx context.Context, tokenID string) error {
	delete(i.refreshTokens, tokenID)
	return nil
}
