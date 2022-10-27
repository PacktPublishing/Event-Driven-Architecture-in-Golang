package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/stackus/errors"

	"eda-in-golang/internal/am"
	"eda-in-golang/internal/tm"
)

type InboxStore struct {
	tableName string
	db        DB
}

var _ tm.InboxStore = (*InboxStore)(nil)

func NewInboxStore(tableName string, db DB) InboxStore {
	return InboxStore{
		tableName: tableName,
		db:        db,
	}
}

func (s InboxStore) Save(ctx context.Context, msg am.RawMessage) error {
	const query = "INSERT INTO %s (id, name, subject, data, received_at) VALUES ($1, $2, $3, $4, CURRENT_TIMESTAMP)"

	_, err := s.db.ExecContext(ctx, s.table(query), msg.ID(), msg.MessageName(), msg.Subject(), msg.Data())
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				return tm.ErrDuplicateMessage(msg.ID())
			}
		}
	}

	return err
}

func (s InboxStore) table(query string) string {
	return fmt.Sprintf(query, s.tableName)
}
