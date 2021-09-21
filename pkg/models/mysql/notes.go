package mysql

import (
	"YaroslavSivash/web-site/pkg/models"
	"database/sql"
	"errors"
)

// SnippetModel - Определяем тип который обертывает пул подключения sql.DB
type NoteModel struct {
DB *sql.DB
}

// Insert - Метод для создания новой заметки в базе дынных.
func (m *NoteModel) Insert(title, content, expires string) (int, error) {

	result, err := m.DB.Exec("INSERT INTO notes (title, content, created, expires) " +
		"VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))", title, content,expires)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	id_note := int(id)

return id_note, err
}

// Get - Метод для возвращения данных заметки по её идентификатору ID.
func (m *NoteModel) Get(id int) (*models.Note, error) {
	note := &models.Note{}

	row := m.DB.QueryRow("SELECT id, title, content, created, expires " +
		"FROM notes WHERE expires > UTC_TIMESTAMP() AND id = ?", id)

	err := row.Scan(&note.ID, &note.Title, &note.Content, &note.Created, &note.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows ) {
			return nil, models.ErrNoRecord
		} else{
			return nil, err
		}
	}

return note, nil
}

// Latest - Метод возвращает 10 наиболее часто используемые заметки.
func (m *NoteModel) Latest() ([]*models.Note, error) {

	rows, err := m.DB.Query("SELECT id, title, content, created, expires FROM notes " +
		"WHERE expires > UTC_TIMESTAMP() ORDER BY created ASC LIMIT 10")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var notes []*models.Note

	for rows.Next() {
		n := &models.Note{}
		err = rows.Scan(&n.ID, &n.Title, &n.Content, &n.Created, &n.Expires)
		if err != nil {
			return nil, err
		}
		notes = append(notes, n)
	}
	if err = rows.Err(); err !=nil {
		return nil, err
	}
return notes, nil
}