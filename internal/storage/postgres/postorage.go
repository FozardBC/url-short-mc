package postorage

import (
	"context"
	"log/slog"
	"microservice_t/internal/config"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	UrlTable    = "urls"
	UrlColumn   = "url"
	AliasColumn = "alias"
)

type Postorage struct {
	p *pgxpool.Pool
}

func New(ctx context.Context, log *slog.Logger, cfg *config.Config) (*Postorage, error) {

	connString := "postgres://" + cfg.Storage.Username + ":" + cfg.Storage.Password + "@" + cfg.Storage.Path + "/" + cfg.Storage.Name + "?sslmode=disable&sslmode=disable&lc_messages=C"

	conn, err := pgxpool.New(ctx, connString)
	if err != nil {
		log.Error("can't connect to pgx pool", "err", err.Error())

		return nil, err
	}

	repo := &Postorage{
		p: conn,
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
