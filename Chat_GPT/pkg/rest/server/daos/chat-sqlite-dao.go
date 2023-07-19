package daos

import (
	"database/sql"
	"errors"
	log "github.com/sirupsen/logrus"
	"github.com/unnagirirao/Imagevision/chat_gpt/pkg/rest/server/daos/clients/sqls"
	"github.com/unnagirirao/Imagevision/chat_gpt/pkg/rest/server/models"
)

type ChatDao struct {
	sqlClient *sqls.SQLiteClient
}

func migrateChats(r *sqls.SQLiteClient) error {
	query := `
	CREATE TABLE IF NOT EXISTS chats(
		Id INTEGER PRIMARY KEY AUTOINCREMENT,
        
		Output TEXT NOT NULL,
		Result TEXT NOT NULL,
		Input TEXT NOT NULL,
        CONSTRAINT id_unique_key UNIQUE (Id)
	)
	`
	_, err1 := r.DB.Exec(query)
	return err1
}

func NewChatDao() (*ChatDao, error) {
	sqlClient, err := sqls.InitSqliteDB()
	if err != nil {
		return nil, err
	}
	err = migrateChats(sqlClient)
	if err != nil {
		return nil, err
	}
	return &ChatDao{
		sqlClient,
	}, nil
}

func (chatDao *ChatDao) CreateChat(m *models.Chat) (*models.Chat, error) {
	insertQuery := "INSERT INTO chats(Output, Result, Input)values(?, ?, ?)"
	res, err := chatDao.sqlClient.DB.Exec(insertQuery, m.Output, m.Result, m.Input)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	m.Id = id

	log.Debugf("chat created")
	return m, nil
}

func (chatDao *ChatDao) UpdateChat(id int64, m *models.Chat) (*models.Chat, error) {
	if id == 0 {
		return nil, errors.New("invalid updated ID")
	}
	if id != m.Id {
		return nil, errors.New("id and payload don't match")
	}

	chat, err := chatDao.GetChat(id)
	if err != nil {
		return nil, err
	}
	if chat == nil {
		return nil, sql.ErrNoRows
	}

	updateQuery := "UPDATE chats SET Output = ?, Result = ?, Input = ? WHERE Id = ?"
	res, err := chatDao.sqlClient.DB.Exec(updateQuery, m.Output, m.Result, m.Input, id)
	if err != nil {
		return nil, err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rowsAffected == 0 {
		return nil, sqls.ErrUpdateFailed
	}

	log.Debugf("chat updated")
	return m, nil
}

func (chatDao *ChatDao) DeleteChat(id int64) error {
	deleteQuery := "DELETE FROM chats WHERE Id = ?"
	res, err := chatDao.sqlClient.DB.Exec(deleteQuery, id)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sqls.ErrDeleteFailed
	}

	log.Debugf("chat deleted")
	return nil
}

func (chatDao *ChatDao) ListChats() ([]*models.Chat, error) {
	selectQuery := "SELECT * FROM chats"
	rows, err := chatDao.sqlClient.DB.Query(selectQuery)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)
	var chats []*models.Chat
	for rows.Next() {
		m := models.Chat{}
		if err = rows.Scan(&m.Id, &m.Output, &m.Result, &m.Input); err != nil {
			return nil, err
		}
		chats = append(chats, &m)
	}
	if chats == nil {
		chats = []*models.Chat{}
	}

	log.Debugf("chat listed")
	return chats, nil
}

func (chatDao *ChatDao) GetChat(id int64) (*models.Chat, error) {
	selectQuery := "SELECT * FROM chats WHERE Id = ?"
	row := chatDao.sqlClient.DB.QueryRow(selectQuery, id)
	m := models.Chat{}
	if err := row.Scan(&m.Id, &m.Output, &m.Result, &m.Input); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sqls.ErrNotExists
		}
		return nil, err
	}

	log.Debugf("chat retrieved")
	return &m, nil
}
