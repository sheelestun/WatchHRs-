package pg

import (
	"context"

	"github.com/google/uuid"
	"github.com/sheelestun/WatchHRs-/internal/domain"
	log "github.com/sirupsen/logrus"
)

func (p *PostgresRepository) AddEmployee(ctx context.Context, employee domain.Employee) (uuid.UUID, error) {
	query := `
		INSERT INTO employees (id, name, email, manager_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id;
	`
	var id uuid.UUID
	err := p.db.GetContext(ctx, &id, query, employee.ID, employee.Name, employee.Email, employee.ManagerID)
	if err != nil {
		return uuid.Nil, err
	}
	log.Debugf("added employee in database: %+v", employee)
	return id, nil
}

func (p *PostgresRepository) GetAllEmployeesByManagerID(ctx context.Context, managerID uuid.UUID) ([]domain.Employee, error) {
	var employees []domain.Employee

	query := `
		SELECT id, name, email, manager_id
		FROM employees
		WHERE manager_id = $1;
	`

	err := p.db.SelectContext(ctx, &employees, query, managerID)
	if err != nil {
		return nil, err
	}

	log.Debugf("Get all employees in database: %+v", employees)
	return employees, nil
}

func (p *PostgresRepository) RemoveEmployee(ctx context.Context, employeeID uuid.UUID) error {
	query := `
		DELETE FROM employees WHERE id = $1;
	`
	_, err := p.db.ExecContext(ctx, query, employeeID)
	log.Debugf("Removed employee %s, with err %s", employeeID, err.Error())
	return err

}
