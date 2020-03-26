package gproxy

import (
	"context"
	"github.com/dgraph-io/ristretto"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type Store struct {
	log   *zap.SugaredLogger
	db    *pgx.Conn
	cache *ristretto.Cache
}

const (
	selectProxyCountryIdByCode = `select country_id from country where country_code=$1`
	selectProxyIdByUI          = `select proxy_id from proxy where proxy_ip=$1 and proxy_port=$2`
	insertProxyCountry         = `INSERT INTO country (created_at, country_name, country_code) VALUES ($1, $2, $3) returning country_id`
	insertProxyItem            = `INSERT INTO proxy_service.public.proxy (created_at, updated_at, proxy_ip, proxy_port, country_id) 
					VALUES ($1, $2, $3, $4, $5) returning proxy_id`
)

func (s Store) GetOrCreateProxyItem(p *ProxyItem) error {
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

func (s Store) InsertProxyItem(p *ProxyItem) error {
	return s.db.QueryRow(context.Background(), insertProxyItem,
		MustTimestampOrNow(p.GetCreatedAt()),
		MustTimestampOrNow(p.GetUpdatedAt()),
		p.ProxyIp, p.ProxyPort,
		p.ProxyCountry.CountryId).Scan(&p.ProxyId)
}

func (s Store) GetOrCreateProxyCountry(c *ProxyCountry) error {
	err := s.SelectProxyCountryIdByCode(c)
	if err == pgx.ErrNoRows {
		return s.InsertProxyCountry(c)
	}
	return err
}

func (s Store) SelectProxyIdByUI(p *ProxyItem) error {
	return s.db.QueryRow(context.Background(), selectProxyIdByUI, p.ProxyIp, p.ProxyPort).Scan(&p.ProxyId)
}

func (s Store) SelectProxyCountryIdByCode(c *ProxyCountry) error {
	return s.db.QueryRow(context.Background(), selectProxyCountryIdByCode, c.CountryCode).Scan(&c.CountryId)
}
func (s Store) InsertProxyCountry(c *ProxyCountry) error {
	return s.db.QueryRow(context.Background(), insertProxyCountry,
		MustTimestampOrNow(c.GetCreatedAt()),
		c.CountryName, c.CountryCode).Scan(&c.CountryId)
}

func (s Store) GetNextProxyItem(item *ProxyItem) error {
	return errors.New("Not implemented")
}

func New(log *zap.SugaredLogger, db *pgx.Conn) *Store {
	return &Store{log: log, db: db}
}
