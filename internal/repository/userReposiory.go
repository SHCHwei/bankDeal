package repository

import (
	"bankDeal/internal/model"
	"bankDeal/internal/database"
	"database/sql"
	"fmt"
	"strings"
	"context"
)


type userRepository struct {
	userList map[int]*model.User
	sqlDB *sql.DB
}


func NewUserRepository(sqlDB *sql.DB) model.UserRepository {
	return &userRepository{
		userList: make(map[int]*model.User),
		sqlDB: sqlDB,
	}
}

func (r *userRepository) FindUserByID(id int) (*model.User, error) {

	rows, err := r.sqlDB.Query("SELECT * FROM users WHERE id = ?", id)

	if err != nil {
		return nil, fmt.Errorf("user id can not null")
	}

	defer rows.Close()

	var userData model.User

	if !rows.Next() {
		return nil, fmt.Errorf("user %d not found", id)
	}

	var birthDate sql.NullString

	if err := rows.Scan(&userData.ID, &userData.FirstName, &userData.LastName, &userData.Email, &userData.Phone, &birthDate, &userData.CreatedAt, &userData.UpdatedAt); err != nil {
		return nil, err
	}

	if birthDate.Valid {
		userData.BirthDate = birthDate.String
	}

	return &userData, nil

}

func (r *userRepository) Search(user model.User) (*model.User, error) {
	var clauses []string
	var args []interface{}

	if user.FirstName != "" {
		clauses = append(clauses, "FirstName = ?")
		args = append(args, user.FirstName)
	}

	if user.LastName != "" {
		clauses = append(clauses, "LastName = ?")
		args = append(args, user.LastName)
	}

	if user.Email != "" {
		clauses = append(clauses, "Email = ?")
		args = append(args, user.Email)
	}

	if user.Phone != "" {
		clauses = append(clauses, "Phone = ?")
		args = append(args, user.Phone)
	}

	if len(clauses) == 0 {
		return nil, fmt.Errorf("at least one search field is required")
	}

	query := "SELECT id, FirstName, LastName, Email, Phone, BirthDate, CreatedAt, UpdatedAt FROM users WHERE " + strings.Join(clauses, " AND ") + " LIMIT 1"
	row := r.sqlDB.QueryRow(query, args...)

	var result model.User
	if err := row.Scan(&result.ID, &result.FirstName, &result.LastName, &result.Email, &result.Phone, &result.BirthDate, &result.CreatedAt, &result.UpdatedAt); err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	return &result, nil

}

func (r *userRepository) InsertUser(ctx context.Context, inputData model.User) (int64, error) {
	const insertUserSQL = "INSERT INTO users (FirstName, LastName, Email, Phone, BirthDate) VALUES (?, ?, ?, ?, ?)"

	tx := database.GetExecutor(ctx, r.sqlDB)

	result, err := tx.ExecContext(ctx, insertUserSQL, inputData.FirstName, inputData.LastName, inputData.Email, inputData.Phone, inputData.BirthDate)
	
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}
