package postorage

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"microservice_t/internal/config"
	"microservice_t/internal/storage"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	UrlTable    = "urls"
	UrlColumn   = "url"
	AliasColumn = "alias"

	queryErr    = "can't do query"
	txCommitErr = "can't commit transaction"
	txBeginErr  = "can't start transaction"
)

type Postorage struct {
	p   *pgxpool.Pool
	log *slog.Logger
}

func New(ctx context.Context, log *slog.Logger, cfg *config.Config) (*Postorage, error) {

	connString := "postgres://" + cfg.Storage.Username + ":" + cfg.Storage.Password + "@" + cfg.Storage.Path + "/" + cfg.Storage.Name + "?sslmode=disable&sslmode=disable&lc_messages=C"

	conn, err := pgxpool.New(ctx, connString)
	if err != nil {
		log.Error("can't connect to pgx pool", "err", err.Error())

		return nil, err
	}

	repo := &Postorage{
		p:   conn,
		log: log,
	}

	stat := conn.Stat()

	log.Debug("DB info", "conns", stat.TotalConns())

	log.Debug("Postgress db is connected", "host", cfg.Storage.Path, "connected time", conn.Stat().AcquireDuration().Microseconds())

	return repo, nil
}

// каждые пять секунд проверят доступность базы данных и отправил в канал ошибку для gracefull shutdown
func (ps *Postorage) Ping(ctx context.Context, errChan chan error) {
	var err error
	for {

		select {
		case <-ctx.Done():
			return
		case <-time.After(5 * time.Second):
			err = ps.p.Ping(ctx)
			if err != nil {

				errChan <- err
				return

			}

		}

	}
}

func (ps *Postorage) Close() {
	ps.p.Close()
}

func (ps *Postorage) GetURL(ctx context.Context, alias string) (url string, err error) {

	tx, err := ps.p.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		ps.log.Error(txBeginErr, "err", err.Error())

		return "", fmt.Errorf("%s:%w", txBeginErr, err)
	}

	tx.Begin(ctx)
	defer tx.Rollback(ctx)

	query := fmt.Sprintf(`SELECT %s FROM %s WHERE %s = $1`, UrlTable, UrlColumn, AliasColumn)

	row := ps.p.QueryRow(ctx, query, alias)
	err = row.Scan(&url)
	if err != nil {
		ps.log.Error(queryErr, "err", err.Error())
		ps.log.Debug(queryErr, "err", err.Error(), "query", query)

		return "", fmt.Errorf("%s:%w", queryErr, err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		ps.log.Error(txCommitErr, "err", err.Error())
		ps.log.Debug(txCommitErr, "err", err.Error(), "query", query)

		return "", fmt.Errorf("%s:%w", txCommitErr, err)
	}

	return

}

func (ps *Postorage) SaveURL(ctx context.Context, url string, alias string) (err error) {

	tx, err := ps.p.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		ps.log.Error(txBeginErr, "err", err.Error())

		return fmt.Errorf("%s:%w", txBeginErr, err)
	}

	tx.Begin(ctx)
	defer tx.Rollback(ctx)

	query := fmt.Sprintf(`INSERT INTO %s (%s, %s) VALUES ($1, $2)`, UrlTable, UrlColumn, AliasColumn)

	_, err = ps.p.Exec(ctx, query, url, alias)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" { // Проверяем код ошибки уникальности
			return storage.ErrAliasAlreadyExists // Возвращаем кастомную ошибку (например, "alias already exists")
		}
		ps.log.Error(queryErr, "err", err.Error())
		ps.log.Debug(queryErr, "err", err.Error(), "query", query)

		return fmt.Errorf("%s:%w", queryErr, err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		ps.log.Error(txCommitErr, "err", err.Error())
		ps.log.Debug(txCommitErr, "err", err.Error(), "query", query)

		return fmt.Errorf("%s:%w", txCommitErr, err)
	}

	return nil
}

func (ps *Postorage) DeleteURL(ctx context.Context, alias string) error {

	tx, err := ps.p.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		ps.log.Error(txBeginErr, "err", err.Error())

		return fmt.Errorf("%s:%w", txBeginErr, err)
	}

	tx.Begin(ctx)
	defer tx.Rollback(ctx)

	query := fmt.Sprintf(`DELETE FROM %s WHERE %s = $1`, UrlTable, AliasColumn)

	com, err := ps.p.Exec(ctx, query, alias)
	if err != nil {

		ps.log.Error(queryErr, "err", err.Error())
		ps.log.Debug(queryErr, "err", err.Error(), "query", query)

		return fmt.Errorf("%s:%w", queryErr, err)
	}

	if com.RowsAffected() == 0 {
		return storage.ErrAliasNotFound
	}

	err = tx.Commit(ctx)
	if err != nil {
		ps.log.Error(txCommitErr, "err", err.Error())
		ps.log.Debug(txCommitErr, "err", err.Error(), "query", query)

		return fmt.Errorf("%s:%w", txCommitErr, err)
	}

	return nil
}
