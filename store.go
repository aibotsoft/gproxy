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
	log   *zap.SugaredLogger
	db    *pgxpool.Pool
	cache *ristretto.Cache
}

const (
	selectProxyCountryIdByCode = `select country_id from country where country_code=$1`
	selectProxyIdByUI          = `select proxy_id from proxy where proxy_ip=$1 and proxy_port=$2`

	getNextProxyItem = `select proxy_id, proxy_ip, proxy_port from proxy_stat_view limit 1`

	insertProxyCountry = `INSERT INTO country (country_name, country_code) VALUES ($1, $2) returning country_id`

	insertProxyStat = `INSERT INTO stat (proxy_id, conn_time, conn_status) VALUES ($1, $2, $3) returning created_at`
	insertProxyItem = `INSERT INTO proxy_service.public.proxy (proxy_ip, proxy_port, country_id) VALUES ($1, $2, $3) returning proxy_id`
)
var ErrNoRows = errors.New("no rows in result set")

func (s *Store) GetOrCreateProxyItem(p *ProxyItem) error {
	s.log.Debug("GetOrCreateProxyItem: ", p)
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

func (s *Store) GetNextProxyItem(p *ProxyItem) error {
	var ip net.IP
	err := s.db.QueryRow(context.Background(), getNextProxyItem).Scan(&p.ProxyId, &ip, &p.ProxyPort)
	switch err {
	case nil:
		p.ProxyIp = ip.String()
		return nil
	case pgx.ErrNoRows:
		return ErrNoRows
	default:
		return errors.Wrap(err, "Store: GetNextProxyItem error")
	}
}

func (s *Store) CreateProxyStat(stat *ProxyStat) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return s.db.QueryRow(ctx, insertProxyStat, stat.ProxyId, stat.ConnTime, stat.ConnStatus).Scan(&stat.CreatedAt)
}

func New(log *zap.SugaredLogger, db *pgxpool.Pool) *Store {
	return &Store{log: log, db: db}
}
