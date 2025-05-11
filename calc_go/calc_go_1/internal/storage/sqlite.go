package storage

import (
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	ID           int
	Login        string
	PasswordHash string
}

type Calculation struct {
	ID         int
	UserID     int
	Expression string
	Result     string
	CreatedAt  time.Time
}

func InitDB(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	schema := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		login TEXT UNIQUE,
		password_hash TEXT
	);

	CREATE TABLE IF NOT EXISTS calculations (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER,
		expression TEXT,
		result TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY(user_id) REFERENCES users(id)
	);
	`
	_, err = db.Exec(schema)
	return db, err
}

func CreateUser(db *sql.DB, login, passwordHash string) error {
	_, err := db.Exec("INSERT INTO users (login, password_hash) VALUES (?, ?)", login, passwordHash)
	return err
}

func GetUserByLogin(db *sql.DB, login string) (*User, error) {
	row := db.QueryRow("SELECT id, login, password_hash FROM users WHERE login = ?", login)
	var u User
	if err := row.Scan(&u.ID, &u.Login, &u.PasswordHash); err != nil {
		return nil, err
	}
	return &u, nil
}

func SaveCalculation(db *sql.DB, userID int, expression, result string) error {
	_, err := db.Exec("INSERT INTO calculations (user_id, expression, result) VALUES (?, ?, ?)", userID, expression, result)
	return err
}

func GetCalculationsByUser(db *sql.DB, userID int) ([]Calculation, error) {
	rows, err := db.Query("SELECT id, user_id, expression, result, created_at FROM calculations WHERE user_id = ? ORDER BY created_at DESC", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []Calculation
	for rows.Next() {
		var c Calculation
		var created string
		if err := rows.Scan(&c.ID, &c.UserID, &c.Expression, &c.Result, &created); err != nil {
			return nil, err
		}
		c.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", created)
		list = append(list, c)
	}
	return list, nil
}

func GetCalculationByID(db *sql.DB, id int) (*Calculation, error) {
	row := db.QueryRow("SELECT id, user_id, expression, result, created_at FROM calculations WHERE id = ?", id)
	var c Calculation
	var created string
	if err := row.Scan(&c.ID, &c.UserID, &c.Expression, &c.Result, &created); err != nil {
		return nil, err
	}
	c.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", created)
	return &c, nil
}
