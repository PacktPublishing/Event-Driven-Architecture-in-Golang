package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgtype"
	"github.com/stackus/errors"

	"eda-in-golang/internal/am"
	"eda-in-golang/internal/tm"
)

type OutboxStore struct {
	tableName string
	db        DB
}

type outboxMessage struct {
	id      string
	name    string
	subject string
	data    []byte
}

var _ tm.OutboxStore = (*OutboxStore)(nil)
var _ am.RawMessage = (*outboxMessage)(nil)

func NewOutboxStore(tableName string, db DB) OutboxStore {
	return OutboxStore{
		tableName: tableName,
		db:        db,
	}
}

func (s OutboxStore) Save(ctx context.Context, msg am.RawMessage) error {
	const query = "INSERT INTO %s (id, name, subject, data) VALUES ($1, $2, $3, $4)"

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

func (s OutboxStore) FindUnpublished(ctx context.Context, limit int) ([]am.RawMessage, error) {
	const query = "SELECT id, name, subject, data FROM %s WHERE published_at IS NULL LIMIT %d"

	rows, err := s.db.QueryContext(ctx, s.table(query, limit))
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			err = errors.Wrap(err, "closing event rows")
		}
	}(rows)

	var msgs []am.RawMessage

	for rows.Next() {
		msg := outboxMessage{}
		err = rows.Scan(&msg.id, &msg.name, &msg.subject, &msg.data)
		if err != nil {
			return msgs, err
		}

		msgs = append(msgs, msg)
	}

	return msgs, rows.Err()
}

func (s OutboxStore) MarkPublished(ctx context.Context, ids ...string) error {
	const query = "UPDATE %s SET published_at = CURRENT_TIMESTAMP WHERE id = ANY ($1)"

	msgIDs := &pgtype.TextArray{}
	err := msgIDs.Set(ids)
	if err != nil {
		return err
	}

	_, err = s.db.ExecContext(ctx, s.table(query), msgIDs)

	return err
}

func (s OutboxStore) table(query string, args ...any) string {
	params := []any{s.tableName}
	params = append(params, args...)
	return fmt.Sprintf(query, params...)
}

func (m outboxMessage) ID() string {
	return m.id
}

func (m outboxMessage) Subject() string {
	return m.subject
}

func (m outboxMessage) MessageName() string {
	return m.name
}

func (m outboxMessage) Data() []byte {
	return m.data
}
