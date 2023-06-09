package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/MrTomSawyer/url-shortener/internal/app/apperrors"
	"github.com/MrTomSawyer/url-shortener/internal/app/config"
	"github.com/MrTomSawyer/url-shortener/internal/app/logger"
	"github.com/MrTomSawyer/url-shortener/internal/app/models"
	"github.com/jmoiron/sqlx"
)

type PostgresURLrepo struct {
	Table            string
	Postgres         *sqlx.DB
	ctx              context.Context
	cfg              config.AppConfig
	urlsToDeleteChan chan models.UserURLs
}

func NewPostgresURLrepo(ctx context.Context, db *sqlx.DB, cfg config.AppConfig, urlsToDeleteChan chan models.UserURLs) *PostgresURLrepo {
	return &PostgresURLrepo{
		Table:            "urls",
		Postgres:         db,
		ctx:              ctx,
		cfg:              cfg,
		urlsToDeleteChan: urlsToDeleteChan,
	}
}

func (u PostgresURLrepo) Create(shortURL, originalURL, userID string) error {
	cxt, cancel := context.WithTimeout(u.ctx, 2000*time.Millisecond)
	defer cancel()

	query := fmt.Sprintf("INSERT INTO %s (shortURL, OriginalURL, userID) VALUES ($1, $2, $3) ON CONFLICT (OriginalURL) DO NOTHING RETURNING id", u.Table)
	row := u.Postgres.QueryRowContext(cxt, query, shortURL, originalURL, userID)
	var res string
	err := row.Scan(&res)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			query := fmt.Sprintf("SELECT shortURL FROM %s WHERE OriginalURL=$1", u.Table)
			row := u.Postgres.QueryRowContext(cxt, query, originalURL)
			err := row.Scan(&res)
			if err != nil {
				return err
			}
			return apperrors.NewURLConflict(res, err)
		}
		return err
	}
	return nil
}

func (u PostgresURLrepo) OriginalURL(shortURL string) (string, error) {
	cxt, cancel := context.WithTimeout(u.ctx, 2000000*time.Millisecond)
	defer cancel()

	query := fmt.Sprintf("SELECT originalurl, isdeleted from %s WHERE shorturl=$1", u.Table)
	row := u.Postgres.QueryRowContext(cxt, query, shortURL)
	var originalURL string
	var isDeleted bool
	err := row.Scan(&originalURL, &isDeleted)
	if err != nil {
		return "", nil
	}
	if isDeleted {
		return "", apperrors.ErrURLDeleted
	}
	return originalURL, nil
}

func (u PostgresURLrepo) BatchCreate(data []models.TempURLBatchRequest, userID string) ([]models.BatchURLResponce, error) {
	tx, err := u.Postgres.Begin()
	if err != nil {
		logger.Log.Infof("Failed to start transaction")
		return []models.BatchURLResponce{}, err
	}

	var response []models.BatchURLResponce

	for _, req := range data {
		query := "INSERT INTO urls (correlationid, shorturl, originalurl, userid) VALUES ($1, $2, $3, $4)"
		_, err = tx.Exec(query, req.CorrelationID, req.ShortURL, req.OriginalURL, userID)

		if err != nil {
			logger.Log.Infof("Failed to insert a shortened URL", err)
			tx.Rollback()
			continue
		}

		response = append(response, models.BatchURLResponce{
			CorrelationID: req.CorrelationID,
			ShortURL:      fmt.Sprintf("%s/%s", u.cfg.Server.DefaultAddr, req.ShortURL),
		})
	}
	err = tx.Commit()
	if err != nil {
		logger.Log.Infof("Failed to commit a transaction")
		return []models.BatchURLResponce{}, err
	}
	return response, nil
}

func (u PostgresURLrepo) GetAll(userid string) ([]models.URLJsonResponse, error) {
	ctx, cancel := context.WithTimeout(u.ctx, 2000*time.Millisecond)
	defer cancel()

	query := fmt.Sprintf("SELECT shorturl, originalurl FROM %s WHERE userid=$1", u.Table)
	logger.Log.Infof("SELECT shorturl, originalurl FROM %s WHERE userid=%s", u.Table, userid)

	rows, err := u.Postgres.QueryContext(ctx, query, userid)
	if err != nil {
		return []models.URLJsonResponse{}, err
	}
	defer rows.Close()

	var responce []models.URLJsonResponse
	for rows.Next() {
		var urlResp models.URLJsonResponse
		err := rows.Scan(&urlResp.ShortURL, &urlResp.OriginalURL)
		if err != nil {
			return []models.URLJsonResponse{}, err
		}
		urlResp.ShortURL = fmt.Sprintf("%s/%s", u.cfg.Server.DefaultAddr, urlResp.ShortURL)
		responce = append(responce, urlResp)
	}

	if err = rows.Err(); err != nil {
		return []models.URLJsonResponse{}, err
	}

	return responce, nil
}

func WorkerDeleteURLs(deleteChan chan models.UserURLs, db *sqlx.DB) {
	for userURL := range deleteChan {
		logger.Log.Infof("Deleting url of %s, userId id %s", userURL.URLs, userURL.UserID)
		// TODO Оптимизировать, убрать вложенный цикл
		for _, url := range userURL.URLs {
			_, err := db.Exec("UPDATE urls SET isdeleted=true WHERE (shorturl = $1 AND userid=$2)", url, userURL.UserID)
			if err != nil {
				logger.Log.Infof("error while deleting: %e", err)
			}
		}
	}
}

func (u PostgresURLrepo) DeleteAll(shortURLs []string, userid string) error {
	u.urlsToDeleteChan <- models.UserURLs{UserID: userid, URLs: shortURLs}
	return nil
}
