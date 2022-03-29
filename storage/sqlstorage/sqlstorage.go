package sqlstorage

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"program/model"
	"program/storage"
	"program/tools"
	"program/users"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type SqlStorage struct {
	mydb *sql.DB
}

func NewSqlStorage(url string) (*SqlStorage, error) {
	// time.Sleep(time.Second * 15)
	db, err := sql.Open("mysql", url)
	if err != nil {
		return nil, fmt.Errorf(" error while connecting to MYSQL database: %v", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf(" error while pinging  MYSQL database: %v", err)
	}
	conn := &SqlStorage{
		mydb: db,
	}
	err = conn.CreateUsersTable()
	if err != nil {
		return nil, fmt.Errorf(" error creating table Users: %v", err)
	}
	err = conn.CreateJokesTable()
	if err != nil {
		return nil, fmt.Errorf(" error creating table Jokes: %v", err)
	}

	// _, err = db.Exec("CREATE INDEX  Jokes_score_IDX  ON Jokes (score DESC);")
	// if err != nil {
	// 	return nil, fmt.Errorf(" error creating index Score: %v", err)
	// }

	// _, err = db.Exec("ALTER TABLE Jokes ADD FULLTEXT(title,body)")
	// if err != nil {
	// 	return nil, fmt.Errorf(" error creating fulltext index: %v", err)
	// }

	return conn, nil
}
func (db *SqlStorage) CreateJokesTable() error {
	_, err := db.mydb.Exec(`
	CREATE TABLE IF NOT EXISTS Jokes (
		title TEXT NOT NULL, 
		body TEXT NOT NULL, 
		score INTEGER NOT NULL DEFAULT 0, 
		id VARCHAR(6) ,
		created_at DATETIME, 
		user_id INT ,
		FOREIGN KEY (user_id) REFERENCES Users(user_id) 
		) ;`)
	if err != nil {
		return fmt.Errorf(" error creating table Jokes: %v", err)
	}
	sqlJokes := `
	SELECT  COUNT(*) 
	FROM      Jokes
	;`
	rows := db.mydb.QueryRow(sqlJokes)

	var l int
	err = rows.Scan(&l)
	if err != nil {
		return fmt.Errorf(" error scan from table Jokes: %v", err)
	}
	if l == 0 {
		file := "reddit_jokes.json"
		byteValues, err := ioutil.ReadFile(file)
		if err != nil {

			fmt.Println("ioutil.ReadFile ERROR:", err)
		} else {
			var docs []model.Joke
			err = json.Unmarshal(byteValues, &docs)
			if err != nil {
				fmt.Println("Unmarshal eerror :", err)

			}
			stmt, err := db.mydb.Prepare(`
			INSERT INTO Jokes 
			(title, body,score,id,created_at,user_id) 
			VALUES (?, ?, ?, ?, ?, ?);`)
			if err != nil {
				log.Fatal(err)
			}
			for _, row := range docs[1:] {
				randomTime, randUserID := tools.RandTimeAndUserID()

				_, err := stmt.Exec(row.Title, row.Body, row.Score, row.ID, randomTime, randUserID)
				if err != nil {
					log.Fatal(err)
				}

			}
			fmt.Println("Table Jokes created")
		}

	}
	return nil
}
func (db *SqlStorage) CreateUsersTable() error {
	_, err := db.mydb.Exec(`
	CREATE TABLE IF NOT EXISTS Users (
		user_id INT AUTO_INCREMENT PRIMARY KEY,
		username TEXT, 
		password TEXT , 
		token TEXT , 
		refresh_token TEXT, 
		created_at DATETIME, 
		updated_at DATETIME 
		);`)
	if err != nil {
		return fmt.Errorf(" error creating table Users: %v", err)
	}
	sqlUsers := `
	SELECT  COUNT(*) 
	FROM      Users
	;`
	rows := db.mydb.QueryRow(sqlUsers)

	var l int
	err = rows.Scan(&l)
	if err != nil {
		return fmt.Errorf(" error scan from table Users: %v", err)
	}
	if l == 0 {
		userServer := users.NewUserServer(db)
		var u model.User
		u.Username = "Denys"
		u.Password = "111111"

		for i := 0; i < 10; i++ {
			u.Username += "1"
			u.Password += "1"
			userServer.SignUpUser(u)
		}
		fmt.Println("Table Users created")
	}

	return nil
}

func (db *SqlStorage) MonthAndCount(year, count int) (int, int, error) {
	sqlStatement := `
	SELECT  MONTH(created_at),  COUNT(*) 
	FROM      Jokes
	WHERE     YEAR(created_at) = ? 
	GROUP BY  MONTH(created_at)
	HAVING COUNT(*) > ?
	;`
	rows, err := db.mydb.Query(sqlStatement, year, count)
	if err != nil {
		return -1, -1, fmt.Errorf(" error : %v", err)
	}
	var r int
	var t int

	for rows.Next() {

		err = rows.Scan(&r, &t)
		if err != nil {
			return -1, -1, fmt.Errorf(" error scan month and count: %v", err)
		}

	}
	return r, t, nil
}
func (db *SqlStorage) JokesByMonth(monthNumber int) (int, error) {
	sqlStatement := `
	SELECT  COUNT(*)
	FROM Jokes
	WHERE MONTH(created_at)= ?
	;`
	row := db.mydb.QueryRow(sqlStatement, monthNumber)
	var r int
	err := row.Scan(&r)
	if err != nil {
		return -1, fmt.Errorf(" error scan month: %v", err)
	}

	return r, nil
}
func (db *SqlStorage) UsersWithoutJokes() ([]string, error) {
	var users []string
	sq := `select username
	from Users LEFT JOIN Jokes 
	ON Users.user_id = Jokes.user_id
	WHERE Jokes.user_id is NULL
	;
	`
	rows, err := db.mydb.Query(sq)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var user model.User
		err = rows.Scan(&user.Username)
		if err != nil {
			return nil, fmt.Errorf(" error scan user: %v", err)
		}

		users = append(users, user.Username)
	}

	return users, nil
}
func (db *SqlStorage) FindID(id string) (model.Joke, error) {

	var j model.Joke
	sqlStatement := `
	SELECT 
	title,body,score,id 
	FROM Jokes 
	WHERE id =?;`

	row := db.mydb.QueryRow(sqlStatement, id)
	err := row.Scan(&j.Title, &j.Body, &j.Score, &j.ID)

	switch err {
	case sql.ErrNoRows:
		return model.Joke{}, storage.ErrNoJokes
	case nil:
		return j, nil
	default:
		return model.Joke{}, err

	}

}
func (db *SqlStorage) Funniest(limit int) ([]model.Joke, error) {

	var jokes []model.Joke

	rows, err := db.mydb.Query(`
	SELECT 
	title,body,score,id 
	FROM Jokes 
	ORDER BY Score DESC LIMIT ?`, limit)

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var joke model.Joke
		err := rows.Scan(&joke.Title, &joke.Body, &joke.Score, &joke.ID)
		if err != nil {
			return nil, err
		}
		jokes = append(jokes, joke)
	}

	return jokes, nil

}
func (db *SqlStorage) Random(limit int) ([]model.Joke, error) {

	var jokes []model.Joke

	rows, err := db.mydb.Query(`
	SELECT 
	title,body,score,id 
	FROM Jokes 
	ORDER BY RAND() LIMIT ?;`, limit)

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var joke model.Joke
		err := rows.Scan(&joke.Title, &joke.Body, &joke.Score, &joke.ID)
		if err != nil {
			return nil, err
		}
		jokes = append(jokes, joke)
	}

	return jokes, nil

}
func (db *SqlStorage) AddJoke(j model.Joke) error {

	sqlStatement := `
	INSERT INTO Jokes 
	(title, body,score,id,created_at) 
	VALUES (?, ?, ?, ?, ?);`

	_, err := db.mydb.Exec(sqlStatement, j.Title, j.Body, j.Score, j.ID, j.Created_at)
	if err != nil {
		return err
	}
	return nil

}

func (db *SqlStorage) UpdateByID(text string, id string) error {
	sqlStatement := `
	UPDATE Jokes
	SET body = ?
	WHERE id = ?;`
	_, err := db.mydb.Exec(sqlStatement, text, id)
	if err != nil {
		return err
	}

	return nil
}
func (db *SqlStorage) TextSearch(text string) ([]model.Joke, error) {
	var jokes []model.Joke
	sqlStatement := `
	SELECT 
	title,body,score,id 
	FROM Jokes 
	WHERE MATCH(title,body) AGAINST(?);`
	rows, err := db.mydb.Query(sqlStatement, text)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var joke model.Joke
		err := rows.Scan(&joke.Title, &joke.Body, &joke.Score, &joke.ID)
		if err != nil {
			return nil, err
		}
		jokes = append(jokes, joke)
	}
	if len(jokes) == 0 {
		return nil, storage.ErrNoMatches
	}
	return jokes, nil

}
func (conn *SqlStorage) CloseClientDB() error {

	err := conn.mydb.Close()
	if err != nil {
		return err
	}
	return nil
}
func (db *SqlStorage) IsExists(user model.User) error {

	sqlStatement := `
	SELECT 
	username 
	FROM Users 
	WHERE username =?;`

	row := db.mydb.QueryRow(sqlStatement, user.Username)
	err := row.Scan(&user.Username)
	switch err {
	case sql.ErrNoRows:
		return nil
	case nil:
		return errors.New("THIS username already exists")
	default:
		return err

	}

}
func (db *SqlStorage) CreateUser(user model.User) error {

	sqlStatement := `
	INSERT INTO Users 
	(username, password,token,refresh_token,created_at,updated_at) 
	VALUES (?, ?, ?, ?, ?, ?);`

	_, err := db.mydb.Exec(sqlStatement, user.Username, user.Password, user.Token, user.Refresh_token, user.Created_at, user.Updated_at)
	if err != nil {
		return err
	}
	return nil
}
func (db *SqlStorage) LoginUser(user model.User) (model.User, error) {

	var foundUser model.User
	sqlStatement := `
	SELECT 
	password,token,refresh_token 
	FROM Users 
	WHERE username =?;`

	row := db.mydb.QueryRow(sqlStatement, user.Username)
	err := row.Scan(&foundUser.Password, &foundUser.Token,
		&foundUser.Refresh_token)
	fmt.Println(err)
	if err != nil {
		if err == sql.ErrNoRows {
			return model.User{}, storage.ErrUserValidate
		}
		return model.User{}, err
	}
	return foundUser, nil

}
func (db *SqlStorage) UpdateTokens(signedToken string, signedRefreshToken string, username string) error {
	updatedTime, err := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	if err != nil {
		return err
	}
	sqlStatement := `
	UPDATE Users
	SET token = ? ,
	refresh_token = ? ,
	updated_at = ?
	WHERE username = ?;`
	_, err = db.mydb.Exec(sqlStatement, signedToken, signedRefreshToken, updatedTime, username)
	if err != nil {
		return err
	}

	return nil

}
