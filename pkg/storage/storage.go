package storage

import (
	"context"
	"errors"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
)

// Структура для создания БД
type DB struct {
	pool *pgxpool.Pool
}

// Структура для нашей модели новости
type NewsRecord struct {
	ID          int
	Title       string // заголовок
	Description string // описание
	PublicTime  int64  // время публикации
	Link        string // ссылка
}

// Создаем подключения к БД
func New() (*DB, error) {
	connectionStr := os.Getenv("dbnews")
	if connectionStr == "" {
		return nil, errors.New("не указана строка подключения к базе данных")
	}
	pool, err := pgxpool.Connect(context.Background(), connectionStr)
	if err != nil {
		return nil, err
	}
	db := DB{
		pool: pool,
	}
	return &db, nil
}

// вставка записи
func (db *DB) InsertNews(news []NewsRecord) error {
	var id int
	// при вставке будем работать через транзакции
	tx, err := db.pool.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	for _, newsRecord := range news {
		err = tx.QueryRow(context.Background(), `
			INSERT INTO news (title, descr, public_time, link)
			VALUES ($1, $2, $3, $4);
			`,
			newsRecord.Title,
			newsRecord.Description,
			newsRecord.PublicTime,
			newsRecord.Link,
		).Scan(&id)

		if err != nil {
			return err
		}
	}
	return nil
}

// возврат записей
func (db *DB) ListNews(n int) ([]NewsRecord, error) {
	if n == 0 {
		err := errors.New("не могу вернуть 0 записей")
		return nil, err
	}
	rows, err := db.pool.Query(context.Background(), `
		SELECT id, title, descr, public_time, link FROM news
		ORDER BY public_time DESC
		LIMIT $1
		`,
		n,
	)
	if err != nil {
		return nil, err
	}
	var news []NewsRecord
	for rows.Next() {
		var r NewsRecord
		err = rows.Scan(
			&r.ID,
			&r.Title,
			&r.Description,
			&r.PublicTime,
			&r.Link,
		)
		if err != nil {
			return nil, err
		}
		news = append(news, r)
	}
	return news, rows.Err()
}
