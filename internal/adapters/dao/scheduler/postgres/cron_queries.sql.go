// Code generated by pggen. DO NOT EDIT.

package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	"time"
)

// Querier is a typesafe Go interface backed by SQL queries.
//
// Methods ending with Batch enqueue a query to run later in a pgx.Batch. After
// calling SendBatch on pgx.Conn, pgxpool.Pool, or pgx.Tx, use the Scan methods
// to parse the results.
type Querier interface {
	CreateJob(ctx context.Context, params CreateJobParams) (int64, error)
	// CreateJobBatch enqueues a CreateJob query into batch to be executed
	// later by the batch.
	CreateJobBatch(batch genericBatch, params CreateJobParams)
	// CreateJobScan scans the result of an executed CreateJobBatch query.
	CreateJobScan(results pgx.BatchResults) (int64, error)

	FindByID(ctx context.Context, id int64) (FindByIDRow, error)
	// FindByIDBatch enqueues a FindByID query into batch to be executed
	// later by the batch.
	FindByIDBatch(batch genericBatch, id int64)
	// FindByIDScan scans the result of an executed FindByIDBatch query.
	FindByIDScan(results pgx.BatchResults) (FindByIDRow, error)

	FindInPeriod(ctx context.Context, at time.Time, period time.Duration) ([]FindInPeriodRow, error)
	// FindInPeriodBatch enqueues a FindInPeriod query into batch to be executed
	// later by the batch.
	FindInPeriodBatch(batch genericBatch, at time.Time, period time.Duration)
	// FindInPeriodScan scans the result of an executed FindInPeriodBatch query.
	FindInPeriodScan(results pgx.BatchResults) ([]FindInPeriodRow, error)

	CountInChat(ctx context.Context, chatID int64) (int64, error)
	// CountInChatBatch enqueues a CountInChat query into batch to be executed
	// later by the batch.
	CountInChatBatch(batch genericBatch, chatID int64)
	// CountInChatScan scans the result of an executed CountInChatBatch query.
	CountInChatScan(results pgx.BatchResults) (int64, error)

	FindByChat(ctx context.Context, chatID int64) ([]FindByChatRow, error)
	// FindByChatBatch enqueues a FindByChat query into batch to be executed
	// later by the batch.
	FindByChatBatch(batch genericBatch, chatID int64)
	// FindByChatScan scans the result of an executed FindByChatBatch query.
	FindByChatScan(results pgx.BatchResults) ([]FindByChatRow, error)

	FindAtTime(ctx context.Context, at time.Time) ([]FindAtTimeRow, error)
	// FindAtTimeBatch enqueues a FindAtTime query into batch to be executed
	// later by the batch.
	FindAtTimeBatch(batch genericBatch, at time.Time)
	// FindAtTimeScan scans the result of an executed FindAtTimeBatch query.
	FindAtTimeScan(results pgx.BatchResults) ([]FindAtTimeRow, error)

	Delete(ctx context.Context, id int64) (pgconn.CommandTag, error)
	// DeleteBatch enqueues a Delete query into batch to be executed
	// later by the batch.
	DeleteBatch(batch genericBatch, id int64)
	// DeleteScan scans the result of an executed DeleteBatch query.
	DeleteScan(results pgx.BatchResults) (pgconn.CommandTag, error)
}

type DBQuerier struct {
	conn  genericConn   // underlying Postgres transport to use
	types *typeResolver // resolve types by name
}

var _ Querier = &DBQuerier{}

// genericConn is a connection to a Postgres database. This is usually backed by
// *pgx.Conn, pgx.Tx, or *pgxpool.Pool.
type genericConn interface {
	// Query executes sql with args. If there is an error the returned Rows will
	// be returned in an error state. So it is allowed to ignore the error
	// returned from Query and handle it in Rows.
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)

	// QueryRow is a convenience wrapper over Query. Any error that occurs while
	// querying is deferred until calling Scan on the returned Row. That Row will
	// error with pgx.ErrNoRows if no rows are returned.
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row

	// Exec executes sql. sql can be either a prepared statement name or an SQL
	// string. arguments should be referenced positionally from the sql string
	// as $1, $2, etc.
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
}

// genericBatch batches queries to send in a single network request to a
// Postgres server. This is usually backed by *pgx.Batch.
type genericBatch interface {
	// Queue queues a query to batch b. query can be an SQL query or the name of a
	// prepared statement. See Queue on *pgx.Batch.
	Queue(query string, arguments ...interface{})
}

// NewQuerier creates a DBQuerier that implements Querier. conn is typically
// *pgx.Conn, pgx.Tx, or *pgxpool.Pool.
func NewQuerier(conn genericConn) *DBQuerier {
	return NewQuerierConfig(conn, QuerierConfig{})
}

type QuerierConfig struct {
	// DataTypes contains pgtype.Value to use for encoding and decoding instead
	// of pggen-generated pgtype.ValueTranscoder.
	//
	// If OIDs are available for an input parameter type and all of its
	// transitive dependencies, pggen will use the binary encoding format for
	// the input parameter.
	DataTypes []pgtype.DataType
}

// NewQuerierConfig creates a DBQuerier that implements Querier with the given
// config. conn is typically *pgx.Conn, pgx.Tx, or *pgxpool.Pool.
func NewQuerierConfig(conn genericConn, cfg QuerierConfig) *DBQuerier {
	return &DBQuerier{conn: conn, types: newTypeResolver(cfg.DataTypes)}
}

// WithTx creates a new DBQuerier that uses the transaction to run all queries.
func (q *DBQuerier) WithTx(tx pgx.Tx) (*DBQuerier, error) {
	return &DBQuerier{conn: tx}, nil
}

// preparer is any Postgres connection transport that provides a way to prepare
// a statement, most commonly *pgx.Conn.
type preparer interface {
	Prepare(ctx context.Context, name, sql string) (sd *pgconn.StatementDescription, err error)
}

// PrepareAllQueries executes a PREPARE statement for all pggen generated SQL
// queries in querier files. Typical usage is as the AfterConnect callback
// for pgxpool.Config
//
// pgx will use the prepared statement if available. Calling PrepareAllQueries
// is an optional optimization to avoid a network round-trip the first time pgx
// runs a query if pgx statement caching is enabled.
func PrepareAllQueries(ctx context.Context, p preparer) error {
	if _, err := p.Prepare(ctx, createJobSQL, createJobSQL); err != nil {
		return fmt.Errorf("prepare query 'CreateJob': %w", err)
	}
	if _, err := p.Prepare(ctx, findByIDSQL, findByIDSQL); err != nil {
		return fmt.Errorf("prepare query 'FindByID': %w", err)
	}
	if _, err := p.Prepare(ctx, findInPeriodSQL, findInPeriodSQL); err != nil {
		return fmt.Errorf("prepare query 'FindInPeriod': %w", err)
	}
	if _, err := p.Prepare(ctx, countInChatSQL, countInChatSQL); err != nil {
		return fmt.Errorf("prepare query 'CountInChat': %w", err)
	}
	if _, err := p.Prepare(ctx, findByChatSQL, findByChatSQL); err != nil {
		return fmt.Errorf("prepare query 'FindByChat': %w", err)
	}
	if _, err := p.Prepare(ctx, findAtTimeSQL, findAtTimeSQL); err != nil {
		return fmt.Errorf("prepare query 'FindAtTime': %w", err)
	}
	if _, err := p.Prepare(ctx, deleteSQL, deleteSQL); err != nil {
		return fmt.Errorf("prepare query 'Delete': %w", err)
	}
	return nil
}

// typeResolver looks up the pgtype.ValueTranscoder by Postgres type name.
type typeResolver struct {
	connInfo *pgtype.ConnInfo // types by Postgres type name
}

func newTypeResolver(types []pgtype.DataType) *typeResolver {
	ci := pgtype.NewConnInfo()
	for _, typ := range types {
		if txt, ok := typ.Value.(textPreferrer); ok && typ.OID != unknownOID {
			typ.Value = txt.ValueTranscoder
		}
		ci.RegisterDataType(typ)
	}
	return &typeResolver{connInfo: ci}
}

// findValue find the OID, and pgtype.ValueTranscoder for a Postgres type name.
func (tr *typeResolver) findValue(name string) (uint32, pgtype.ValueTranscoder, bool) {
	typ, ok := tr.connInfo.DataTypeForName(name)
	if !ok {
		return 0, nil, false
	}
	v := pgtype.NewValue(typ.Value)
	return typ.OID, v.(pgtype.ValueTranscoder), true
}

// setValue sets the value of a ValueTranscoder to a value that should always
// work and panics if it fails.
func (tr *typeResolver) setValue(vt pgtype.ValueTranscoder, val interface{}) pgtype.ValueTranscoder {
	if err := vt.Set(val); err != nil {
		panic(fmt.Sprintf("set ValueTranscoder %T to %+v: %s", vt, val, err))
	}
	return vt
}

const createJobSQL = `INSERT INTO crons(chat_id, send_at, flags, title)
VALUES ($1,
        $2,
        $3,
        $4)
RETURNING id;`

type CreateJobParams struct {
	ChatID int64
	SendAt time.Time
	Flags  pgtype.Int2
	Title  pgtype.Varchar
}

// CreateJob implements Querier.CreateJob.
func (q *DBQuerier) CreateJob(ctx context.Context, params CreateJobParams) (int64, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "CreateJob")
	row := q.conn.QueryRow(ctx, createJobSQL, params.ChatID, params.SendAt, params.Flags, params.Title)
	var item int64
	if err := row.Scan(&item); err != nil {
		return item, fmt.Errorf("query CreateJob: %w", err)
	}
	return item, nil
}

// CreateJobBatch implements Querier.CreateJobBatch.
func (q *DBQuerier) CreateJobBatch(batch genericBatch, params CreateJobParams) {
	batch.Queue(createJobSQL, params.ChatID, params.SendAt, params.Flags, params.Title)
}

// CreateJobScan implements Querier.CreateJobScan.
func (q *DBQuerier) CreateJobScan(results pgx.BatchResults) (int64, error) {
	row := results.QueryRow()
	var item int64
	if err := row.Scan(&item); err != nil {
		return item, fmt.Errorf("scan CreateJobBatch row: %w", err)
	}
	return item, nil
}

const findByIDSQL = `SELECT id, chat_id, title, send_at, flags
FROM crons
WHERE id = $1;`

type FindByIDRow struct {
	ID     int64          `json:"id"`
	ChatID int64          `json:"chat_id"`
	Title  pgtype.Varchar `json:"title"`
	SendAt time.Time      `json:"send_at"`
	Flags  pgtype.Int2    `json:"flags"`
}

// FindByID implements Querier.FindByID.
func (q *DBQuerier) FindByID(ctx context.Context, id int64) (FindByIDRow, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "FindByID")
	row := q.conn.QueryRow(ctx, findByIDSQL, id)
	var item FindByIDRow
	if err := row.Scan(&item.ID, &item.ChatID, &item.Title, &item.SendAt, &item.Flags); err != nil {
		return item, fmt.Errorf("query FindByID: %w", err)
	}
	return item, nil
}

// FindByIDBatch implements Querier.FindByIDBatch.
func (q *DBQuerier) FindByIDBatch(batch genericBatch, id int64) {
	batch.Queue(findByIDSQL, id)
}

// FindByIDScan implements Querier.FindByIDScan.
func (q *DBQuerier) FindByIDScan(results pgx.BatchResults) (FindByIDRow, error) {
	row := results.QueryRow()
	var item FindByIDRow
	if err := row.Scan(&item.ID, &item.ChatID, &item.Title, &item.SendAt, &item.Flags); err != nil {
		return item, fmt.Errorf("scan FindByIDBatch row: %w", err)
	}
	return item, nil
}

const findInPeriodSQL = `SELECT id, chat_id, title, send_at, flags
FROM crons
WHERE send_at > $1::time - $2::interval
  AND send_at <= $1::time
ORDER BY id, send_at;`

type FindInPeriodRow struct {
	ID     int64          `json:"id"`
	ChatID int64          `json:"chat_id"`
	Title  pgtype.Varchar `json:"title"`
	SendAt time.Time      `json:"send_at"`
	Flags  pgtype.Int2    `json:"flags"`
}

// FindInPeriod implements Querier.FindInPeriod.
func (q *DBQuerier) FindInPeriod(ctx context.Context, at time.Time, period time.Duration) ([]FindInPeriodRow, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "FindInPeriod")
	rows, err := q.conn.Query(ctx, findInPeriodSQL, at, period)
	if err != nil {
		return nil, fmt.Errorf("query FindInPeriod: %w", err)
	}
	defer rows.Close()
	items := []FindInPeriodRow{}
	for rows.Next() {
		var item FindInPeriodRow
		if err := rows.Scan(&item.ID, &item.ChatID, &item.Title, &item.SendAt, &item.Flags); err != nil {
			return nil, fmt.Errorf("scan FindInPeriod row: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("close FindInPeriod rows: %w", err)
	}
	return items, err
}

// FindInPeriodBatch implements Querier.FindInPeriodBatch.
func (q *DBQuerier) FindInPeriodBatch(batch genericBatch, at time.Time, period time.Duration) {
	batch.Queue(findInPeriodSQL, at, period)
}

// FindInPeriodScan implements Querier.FindInPeriodScan.
func (q *DBQuerier) FindInPeriodScan(results pgx.BatchResults) ([]FindInPeriodRow, error) {
	rows, err := results.Query()
	if err != nil {
		return nil, fmt.Errorf("query FindInPeriodBatch: %w", err)
	}
	defer rows.Close()
	items := []FindInPeriodRow{}
	for rows.Next() {
		var item FindInPeriodRow
		if err := rows.Scan(&item.ID, &item.ChatID, &item.Title, &item.SendAt, &item.Flags); err != nil {
			return nil, fmt.Errorf("scan FindInPeriodBatch row: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("close FindInPeriodBatch rows: %w", err)
	}
	return items, err
}

const countInChatSQL = `SELECT count(*)
FROM crons
WHERE chat_id = $1::bigint;`

// CountInChat implements Querier.CountInChat.
func (q *DBQuerier) CountInChat(ctx context.Context, chatID int64) (int64, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "CountInChat")
	row := q.conn.QueryRow(ctx, countInChatSQL, chatID)
	var item int64
	if err := row.Scan(&item); err != nil {
		return item, fmt.Errorf("query CountInChat: %w", err)
	}
	return item, nil
}

// CountInChatBatch implements Querier.CountInChatBatch.
func (q *DBQuerier) CountInChatBatch(batch genericBatch, chatID int64) {
	batch.Queue(countInChatSQL, chatID)
}

// CountInChatScan implements Querier.CountInChatScan.
func (q *DBQuerier) CountInChatScan(results pgx.BatchResults) (int64, error) {
	row := results.QueryRow()
	var item int64
	if err := row.Scan(&item); err != nil {
		return item, fmt.Errorf("scan CountInChatBatch row: %w", err)
	}
	return item, nil
}

const findByChatSQL = `SELECT id, chat_id, title, send_at, flags
FROM crons
WHERE chat_id = $1::bigint
ORDER BY id;`

type FindByChatRow struct {
	ID     int64          `json:"id"`
	ChatID int64          `json:"chat_id"`
	Title  pgtype.Varchar `json:"title"`
	SendAt time.Time      `json:"send_at"`
	Flags  pgtype.Int2    `json:"flags"`
}

// FindByChat implements Querier.FindByChat.
func (q *DBQuerier) FindByChat(ctx context.Context, chatID int64) ([]FindByChatRow, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "FindByChat")
	rows, err := q.conn.Query(ctx, findByChatSQL, chatID)
	if err != nil {
		return nil, fmt.Errorf("query FindByChat: %w", err)
	}
	defer rows.Close()
	items := []FindByChatRow{}
	for rows.Next() {
		var item FindByChatRow
		if err := rows.Scan(&item.ID, &item.ChatID, &item.Title, &item.SendAt, &item.Flags); err != nil {
			return nil, fmt.Errorf("scan FindByChat row: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("close FindByChat rows: %w", err)
	}
	return items, err
}

// FindByChatBatch implements Querier.FindByChatBatch.
func (q *DBQuerier) FindByChatBatch(batch genericBatch, chatID int64) {
	batch.Queue(findByChatSQL, chatID)
}

// FindByChatScan implements Querier.FindByChatScan.
func (q *DBQuerier) FindByChatScan(results pgx.BatchResults) ([]FindByChatRow, error) {
	rows, err := results.Query()
	if err != nil {
		return nil, fmt.Errorf("query FindByChatBatch: %w", err)
	}
	defer rows.Close()
	items := []FindByChatRow{}
	for rows.Next() {
		var item FindByChatRow
		if err := rows.Scan(&item.ID, &item.ChatID, &item.Title, &item.SendAt, &item.Flags); err != nil {
			return nil, fmt.Errorf("scan FindByChatBatch row: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("close FindByChatBatch rows: %w", err)
	}
	return items, err
}

const findAtTimeSQL = `SELECT id, chat_id, title, send_at, flags
FROM crons
WHERE send_at = $1::time
ORDER BY id, send_at;`

type FindAtTimeRow struct {
	ID     int64          `json:"id"`
	ChatID int64          `json:"chat_id"`
	Title  pgtype.Varchar `json:"title"`
	SendAt time.Time      `json:"send_at"`
	Flags  pgtype.Int2    `json:"flags"`
}

// FindAtTime implements Querier.FindAtTime.
func (q *DBQuerier) FindAtTime(ctx context.Context, at time.Time) ([]FindAtTimeRow, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "FindAtTime")
	rows, err := q.conn.Query(ctx, findAtTimeSQL, at)
	if err != nil {
		return nil, fmt.Errorf("query FindAtTime: %w", err)
	}
	defer rows.Close()
	items := []FindAtTimeRow{}
	for rows.Next() {
		var item FindAtTimeRow
		if err := rows.Scan(&item.ID, &item.ChatID, &item.Title, &item.SendAt, &item.Flags); err != nil {
			return nil, fmt.Errorf("scan FindAtTime row: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("close FindAtTime rows: %w", err)
	}
	return items, err
}

// FindAtTimeBatch implements Querier.FindAtTimeBatch.
func (q *DBQuerier) FindAtTimeBatch(batch genericBatch, at time.Time) {
	batch.Queue(findAtTimeSQL, at)
}

// FindAtTimeScan implements Querier.FindAtTimeScan.
func (q *DBQuerier) FindAtTimeScan(results pgx.BatchResults) ([]FindAtTimeRow, error) {
	rows, err := results.Query()
	if err != nil {
		return nil, fmt.Errorf("query FindAtTimeBatch: %w", err)
	}
	defer rows.Close()
	items := []FindAtTimeRow{}
	for rows.Next() {
		var item FindAtTimeRow
		if err := rows.Scan(&item.ID, &item.ChatID, &item.Title, &item.SendAt, &item.Flags); err != nil {
			return nil, fmt.Errorf("scan FindAtTimeBatch row: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("close FindAtTimeBatch rows: %w", err)
	}
	return items, err
}

const deleteSQL = `DELETE
FROM crons
WHERE id = $1;`

// Delete implements Querier.Delete.
func (q *DBQuerier) Delete(ctx context.Context, id int64) (pgconn.CommandTag, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "Delete")
	cmdTag, err := q.conn.Exec(ctx, deleteSQL, id)
	if err != nil {
		return cmdTag, fmt.Errorf("exec query Delete: %w", err)
	}
	return cmdTag, err
}

// DeleteBatch implements Querier.DeleteBatch.
func (q *DBQuerier) DeleteBatch(batch genericBatch, id int64) {
	batch.Queue(deleteSQL, id)
}

// DeleteScan implements Querier.DeleteScan.
func (q *DBQuerier) DeleteScan(results pgx.BatchResults) (pgconn.CommandTag, error) {
	cmdTag, err := results.Exec()
	if err != nil {
		return cmdTag, fmt.Errorf("exec DeleteBatch: %w", err)
	}
	return cmdTag, err
}

// textPreferrer wraps a pgtype.ValueTranscoder and sets the preferred encoding
// format to text instead binary (the default). pggen uses the text format
// when the OID is unknownOID because the binary format requires the OID.
// Typically occurs if the results from QueryAllDataTypes aren't passed to
// NewQuerierConfig.
type textPreferrer struct {
	pgtype.ValueTranscoder
	typeName string
}

// PreferredParamFormat implements pgtype.ParamFormatPreferrer.
func (t textPreferrer) PreferredParamFormat() int16 { return pgtype.TextFormatCode }

func (t textPreferrer) NewTypeValue() pgtype.Value {
	return textPreferrer{ValueTranscoder: pgtype.NewValue(t.ValueTranscoder).(pgtype.ValueTranscoder), typeName: t.typeName}
}

func (t textPreferrer) TypeName() string {
	return t.typeName
}

// unknownOID means we don't know the OID for a type. This is okay for decoding
// because pgx call DecodeText or DecodeBinary without requiring the OID. For
// encoding parameters, pggen uses textPreferrer if the OID is unknown.
const unknownOID = 0
