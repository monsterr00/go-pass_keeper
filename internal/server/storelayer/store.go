package storelayer

import (
	"context"
	"database/sql"
	"errors"
	"log"

	config "github.com/monsterr00/go-pass-keeper/configs/server"
	"github.com/monsterr00/go-pass-keeper/helpers"
	"github.com/monsterr00/go-pass-keeper/internal/models"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

type store struct {
	conn *sql.DB
}

type Store interface {
	UserCreate(ctx context.Context, user models.User) error
	UserGet(ctx context.Context, login string) (models.User, error)
	FileCreate(ctx context.Context, file models.File) error
	FileUpdate(ctx context.Context, file models.File) error
	FileGet(ctx context.Context, login string, title string) (models.File, error)
	CheckDuplicateFile(ctx context.Context, login string, title string) error
	CardCreate(ctx context.Context, card models.Card) error
	CardUpdate(ctx context.Context, card models.Card) error
	CardGet(ctx context.Context, login string, title string) (models.Card, error)
	CheckDuplicateCard(ctx context.Context, login string, title string) error
	TextCreate(ctx context.Context, text models.Text) error
	TextUpdate(ctx context.Context, text models.Text) error
	TextGet(ctx context.Context, login string, title string) (models.Text, error)
	CheckDuplicateText(ctx context.Context, login string, title string) error
	PasswordCreate(ctx context.Context, password models.Password) error
	PasswordUpdate(ctx context.Context, password models.Password) error
	PasswordGet(ctx context.Context, login string, title string) (models.Password, error)
	CheckDuplicatePassword(ctx context.Context, login string, title string) error
	Close() error
}

const (
	migrationsPath = "db/migrations"
)

// New инициализирует соединение с БД и соотвествующие настройки.
func New() *store {
	db, err := sql.Open("postgres", config.ServerOptions.DBaddress)
	if err != nil {
		log.Fatal(err)
	}

	filePath := helpers.AbsolutePath("file:///", migrationsPath)

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		filePath,
		"postgres", driver)
	if err != nil {
		log.Fatal(err)
	}

	if err = m.Up(); err != nil {
		if err != migrate.ErrNoChange {
			log.Fatal(err)
		}
	}

	storl := &store{
		conn: db,
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	return storl
}

// UserCreate создает новую запись в таблице БД users.
func (storl *store) UserCreate(ctx context.Context, user models.User) error {
	errPing := storl.conn.Ping()
	if errPing != nil {
		return errPing
	}

	tx, err := storl.conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `
	INSERT INTO users 
	(Login, Password)
	VALUES
	($1, $2);
    `, user.Login, user.Password)
	if err != nil {
		return errors.Join(err, tx.Rollback())
	}
	return tx.Commit()
}

// UserGet возвращает запись из таблицы БД users.
func (storl *store) UserGet(ctx context.Context, login string) (models.User, error) {
	row := storl.conn.QueryRowContext(ctx, `	
	SELECT 
		Login,
		Password
	FROM users
	WHERE login = $1
    `,
		login,
	)

	var user models.User
	err := row.Scan(&user.Login, &user.Password)
	if err != nil {
		return user, err
	}

	return user, nil
}

// FileCreate создает новую запись в таблице БД files.
func (storl *store) FileCreate(ctx context.Context, file models.File) error {
	errPing := storl.conn.Ping()
	if errPing != nil {
		return errPing
	}

	tx, err := storl.conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `
	INSERT INTO files 
	(UserLogin, Title, FileName, File, DataType)
	VALUES
	($1, $2, $3, $4, $5);
    `, file.UserLogin, file.Title, file.FileName, file.File, file.DataType)
	if err != nil {
		return errors.Join(err, tx.Rollback())
	}
	return tx.Commit()
}

// FileUpdate обновляет запись в таблице БД files.
func (storl *store) FileUpdate(ctx context.Context, file models.File) error {
	tx, err := storl.conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `
	UPDATE files
	SET FileName = $3,
		File = $4,
		DataType = $5
	WHERE UserLogin = $1 and Title = $2;		
    `, file.UserLogin, file.Title, file.FileName, file.File, file.DataType)
	if err != nil {
		return errors.Join(err, tx.Rollback())
	}
	return tx.Commit()
}

// FileGet возвращает запись из таблицы БД files.
func (storl *store) FileGet(ctx context.Context, login string, title string) (models.File, error) {
	row := storl.conn.QueryRowContext(ctx, `	
	SELECT
		Title, 
		FileName,
		File,
		DataType
	FROM files
	WHERE UserLogin = $1 and title = $2
    `,
		login, title,
	)

	var file models.File
	err := row.Scan(&file.Title, &file.FileName, &file.File, &file.DataType)
	if err != nil {
		return file, err
	}

	return file, nil
}

// CheckDuplicateFile проверяет наличие записи в таблице БД files.
func (storl *store) CheckDuplicateFile(ctx context.Context, login string, title string) error {
	row := storl.conn.QueryRowContext(ctx, `	
	SELECT
		UserLogin 		
	FROM files
	WHERE UserLogin = $1 and title = $2
    `,
		login, title,
	)

	var user string
	err := row.Scan(&user)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	if err == sql.ErrNoRows {
		return nil
	}

	return errors.New("record already exists")
}

// CardCreate создает новую запись в таблице БД cards.
func (storl *store) CardCreate(ctx context.Context, card models.Card) error {
	errPing := storl.conn.Ping()
	if errPing != nil {
		return errPing
	}

	tx, err := storl.conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `
	INSERT INTO cards 
	(UserLogin, Title, Bank, CardNumber, CVV, DateExpire, CardHolder)
	VALUES
	($1, $2, $3, $4, $5, $6, $7);
    `, card.UserLogin, card.Title, card.Bank, card.CardNumber, card.CVV, card.CardHolder)
	if err != nil {
		return errors.Join(err, tx.Rollback())
	}
	return tx.Commit()
}

// CardUpdate обновляет запись в таблице БД cards.
func (storl *store) CardUpdate(ctx context.Context, card models.Card) error {
	tx, err := storl.conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `
	UPDATE cards
	SET Bank = $3,
		CardNumber = $4,
		CVV = $5,
		CardHolder = $6
	WHERE UserLogin = $1 and Title = $2;		
    `, card.UserLogin, card.Title, card.Bank, card.CardNumber, card.CVV, card.CardHolder)
	if err != nil {
		return errors.Join(err, tx.Rollback())
	}
	return tx.Commit()
}

// CardGet возвращает запись из таблицы БД cards.
func (storl *store) CardGet(ctx context.Context, login string, title string) (models.Card, error) {
	row := storl.conn.QueryRowContext(ctx, `	
	SELECT
		Title, 
		Bank,
		CardNumber,
		CVV,
		CardHolder
	FROM cards
	WHERE UserLogin = $1 and title = $2
    `,
		login, title,
	)

	var card models.Card
	err := row.Scan(&card.Title, &card.Bank, &card.CardNumber, &card.CVV, &card.CardHolder)
	if err != nil {
		return card, err
	}

	return card, nil
}

// CheckDuplicateCard проверяет наличие записи в таблице БД cards.
func (storl *store) CheckDuplicateCard(ctx context.Context, login string, title string) error {
	row := storl.conn.QueryRowContext(ctx, `	
	SELECT
		UserLogin 		
	FROM cards
	WHERE UserLogin = $1 and title = $2
    `,
		login, title,
	)

	var user string
	err := row.Scan(&user)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	if err == sql.ErrNoRows {
		return nil
	}

	return errors.New("record already exists")
}

// TextCreate создает новую запись в таблице БД texts.
func (storl *store) TextCreate(ctx context.Context, text models.Text) error {
	errPing := storl.conn.Ping()
	if errPing != nil {
		return errPing
	}

	tx, err := storl.conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `
	INSERT INTO texts 
	(UserLogin, Title, Text)
	VALUES
	($1, $2, $3);
    `, text.UserLogin, text.Title, text.Text)
	if err != nil {
		return errors.Join(err, tx.Rollback())
	}
	return tx.Commit()
}

// TextUpdate обновляет запись в таблице БД texts.
func (storl *store) TextUpdate(ctx context.Context, text models.Text) error {
	tx, err := storl.conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `
	UPDATE texts
	SET Text = $3
	WHERE UserLogin = $1 and Title = $2;		
    `, text.UserLogin, text.Title, text.Text)
	if err != nil {
		return errors.Join(err, tx.Rollback())
	}
	return tx.Commit()
}

// TextGet возвращает запись из таблицы БД texts.
func (storl *store) TextGet(ctx context.Context, login string, title string) (models.Text, error) {
	row := storl.conn.QueryRowContext(ctx, `	
	SELECT
		UserLogin,
		Title, 
		Text
	FROM texts
	WHERE UserLogin = $1 and title = $2
    `,
		login, title,
	)

	var text models.Text
	err := row.Scan(&text.UserLogin, &text.Title, &text.Text)
	if err != nil {
		return text, err
	}

	return text, nil
}

// CheckDuplicateText проверяет наличие записи в таблице БД texts.
func (storl *store) CheckDuplicateText(ctx context.Context, login string, title string) error {
	row := storl.conn.QueryRowContext(ctx, `	
	SELECT
		UserLogin 		
	FROM texts
	WHERE UserLogin = $1 and title = $2
    `,
		login, title,
	)

	var user string
	err := row.Scan(&user)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	if err == sql.ErrNoRows {
		return nil
	}

	return errors.New("record already exists")
}

// PasswordCreate создает новую запись в таблице БД passwords.
func (storl *store) PasswordCreate(ctx context.Context, password models.Password) error {
	errPing := storl.conn.Ping()
	if errPing != nil {
		return errPing
	}

	tx, err := storl.conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `
	INSERT INTO passwords 
	(UserLogin, Title, Login, Password)
	VALUES
	($1, $2, $3, $4);
    `, password.UserLogin, password.Title, password.Login, password.Password)
	if err != nil {
		return errors.Join(err, tx.Rollback())
	}
	return tx.Commit()
}

// PasswordUpdate обновляет запись в таблице БД passwords.
func (storl *store) PasswordUpdate(ctx context.Context, password models.Password) error {
	tx, err := storl.conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `
	UPDATE passwords
	SET Login = $3,
		Password = $4
	WHERE UserLogin = $1 and Title = $2;		
    `, password.UserLogin, password.Title, password.Login, password.Password)
	if err != nil {
		return errors.Join(err, tx.Rollback())
	}
	return tx.Commit()
}

// PasswordGet возвращает запись из таблицы БД passwords.
func (storl *store) PasswordGet(ctx context.Context, login string, title string) (models.Password, error) {
	row := storl.conn.QueryRowContext(ctx, `	
	SELECT
		Title, 
		Login,
		Password
	FROM passwords
	WHERE UserLogin = $1 and title = $2
    `,
		login, title,
	)

	var password models.Password
	err := row.Scan(&password.Title, &password.Login, &password.Password)
	if err != nil {
		return password, err
	}

	return password, nil
}

// CheckDuplicatePassword проверяет наличие записи в таблице БД passwords.
func (storl *store) CheckDuplicatePassword(ctx context.Context, login string, title string) error {
	row := storl.conn.QueryRowContext(ctx, `	
	SELECT
		UserLogin 		
	FROM passwords
	WHERE UserLogin = $1 and title = $2
    `,
		login, title,
	)

	var user string
	err := row.Scan(&user)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	if err == sql.ErrNoRows {
		return nil
	}

	return errors.New("record already exists")
}

// Close закрывает соединение с БД.
func (storl *store) Close() error {
	return storl.conn.Close()
}
