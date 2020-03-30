package gproxy

import (
	"context"
	"github.com/dgraph-io/ristretto"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"net"
	"time"
)

type Store struct {
	log                *zap.SugaredLogger
	db                 *pgxpool.Pool
	cache              *ristretto.Cache
	nextProxyQueue     chan ProxyItem
	nextProxyLastQuery time.Time
}

const (
	nextProxyQueryTimeout      = 30 * time.Second
	selectProxyCountryIdByCode = `select country_id from country where country_code=$1`
	selectProxyIdByUI          = `select proxy_id from proxy where proxy_ip=$1 and proxy_port=$2`

	//getNextProxyItem = `select proxy_id, proxy_ip, proxy_port from proxy_stat_view limit 1`
	//getNextProxyItem      = `select proxy_id, proxy_ip, proxy_port from get_next_proxy_for_check limit 1`
	getNextProxyItemBatch = `select proxy_id, proxy_ip, proxy_port from get_next_proxy_for_check(60, $1);`

	insertProxyCountry = `INSERT INTO country (country_name, country_code) VALUES ($1, $2) returning country_id`

	insertProxyStat = `INSERT INTO stat (proxy_id, conn_time, conn_status) VALUES ($1, $2, $3) returning created_at`
	insertProxyItem = `INSERT INTO proxy_service.public.proxy (proxy_ip, proxy_port, country_id) VALUES ($1, $2, $3) returning proxy_id`
)

var ErrNoRows = errors.New("no rows in result set")
var ErrQueryTooOften = errors.New("query too often, need to wait")

func (s *Store) GetOrCreateProxyItem(ctx context.Context, p *ProxyItem) error {
	// Ищем прокси в базе, если есть сразу возвращаем
	err := s.SelectProxyIdByUI(p)
	if err == nil {
		return nil
	}
	s.log.Debug("Proxy not found")
	err = s.GetOrCreateProxyCountry(p.ProxyCountry)
	if err != nil {
		return errors.Wrap(err, "GetOrCreateProxyCountry")
	}
	err = s.InsertProxyItem(p)
	if err != nil {
		return errors.Wrap(err, "InsertProxyItem")
	}
	return nil
}

func (s *Store) InsertProxyItem(p *ProxyItem) error {
	return s.db.QueryRow(context.Background(), insertProxyItem,
		p.ProxyIp, p.ProxyPort,
		p.ProxyCountry.CountryId).Scan(&p.ProxyId)
}

func (s *Store) GetOrCreateProxyCountry(c *ProxyCountry) error {
	err := s.SelectProxyCountryIdByCode(c)
	if err == pgx.ErrNoRows {
		return s.InsertProxyCountry(c)
	}
	return err
}

func (s *Store) SelectProxyIdByUI(p *ProxyItem) error {
	return s.db.QueryRow(context.Background(), selectProxyIdByUI, p.ProxyIp, p.ProxyPort).Scan(&p.ProxyId)
}

func (s *Store) SelectProxyCountryIdByCode(c *ProxyCountry) error {
	return s.db.QueryRow(context.Background(), selectProxyCountryIdByCode, c.CountryCode).Scan(&c.CountryId)
}
func (s *Store) InsertProxyCountry(c *ProxyCountry) error {
	return s.db.QueryRow(context.Background(), insertProxyCountry, c.CountryName, c.CountryCode).Scan(&c.CountryId)
}

func (s *Store) GetNextProxyItemBatch(ctx context.Context, size int) error {
	lastQueryTimeout := time.Now().UTC().Sub(s.nextProxyLastQuery) > nextProxyQueryTimeout
	if !lastQueryTimeout {
		s.log.Info("QueryTooOften")
		return ErrQueryTooOften
	}
	s.nextProxyLastQuery = time.Now().UTC()
	rows, err := s.db.Query(ctx, getNextProxyItemBatch, size)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var ip net.IP
		p := ProxyItem{}
		err := rows.Scan(&p.ProxyId, &ip, &p.ProxyPort)
		if err != nil {
			s.log.Error(errors.Wrap(err, "error wile scan getNextProxyItemBatch"))
			continue
		}
		p.ProxyIp = ip.String()
		s.nextProxyQueue <- p
	}
	return nil
}

// NextProxyItemProducer вынимает ProxyItem из nextProxyQueue
// Если элемента нет, вызывает функцию пополнения, и повторяет попытку взять элемент
func (s *Store) NextProxyItemProducer(ctx context.Context) (ProxyItem, error) {
	for {
		select {
		case p := <-s.nextProxyQueue:
			return p, nil
		default:
			err := s.GetNextProxyItemBatch(ctx, 100)
			switch err {
			case nil:
			case ErrQueryTooOften:
				time.Sleep(100 * time.Millisecond)
			default:
				return ProxyItem{}, err
			}
		}
	}
}

func (s *Store) GetNextProxyItem(ctx context.Context) (*ProxyItem, error) {
	proxyItem, err := s.NextProxyItemProducer(ctx)
	switch err {
	case nil:
		return &proxyItem, nil
	case pgx.ErrNoRows:
		return nil, ErrNoRows
	default:
		return nil, errors.Wrap(err, "Store: GetNextProxyItem error")
	}
}

//func (s *Store) GetNextProxyItem(p *ProxyItem) error {
//	var ip net.IP
//	err := s.db.QueryRow(context.Background(), getNextProxyItem).Scan(&p.ProxyId, &ip, &p.ProxyPort)
//	switch err {
//	case nil:
//		p.ProxyIp = ip.String()
//		return nil
//	case pgx.ErrNoRows:
//		return ErrNoRows
//	default:
//		return errors.Wrap(err, "Store: GetNextProxyItem error")
//	}
//}

func (s *Store) CreateProxyStat(ctx context.Context, stat *ProxyStat) error {
	return s.db.QueryRow(ctx, insertProxyStat, stat.ProxyId, stat.ConnTime, stat.ConnStatus).Scan(&stat.CreatedAt)
}

func New(log *zap.SugaredLogger, db *pgxpool.Pool) *Store {
	return &Store{log: log, db: db, nextProxyQueue: make(chan ProxyItem, 200)}
}
