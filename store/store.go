package store

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/figment-networks/avalanche-indexer/util"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var ErrNotFound = gorm.ErrRecordNotFound

type DB struct {
	db *gorm.DB

	Addresses  AddressesStore
	Validators ValidatorsStore
	Delegators DelegatorsStore
	Networks   NetworksStore
}

func NewRaw(connStr string) (*gorm.DB, error) {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second * 5,
			LogLevel:      logger.Warn,
			Colorful:      true,
		},
	)

	return gorm.Open(postgres.Open(connStr), &gorm.Config{
		Logger: newLogger,
	})
}

func New(connStr string) (*DB, error) {
	conn, err := NewRaw(connStr)
	if err != nil {
		return nil, err
	}

	return &DB{
		db: conn,

		Addresses:  AddressesStore{conn},
		Validators: ValidatorsStore{conn},
		Delegators: DelegatorsStore{conn},
		Networks:   NetworksStore{conn},
	}, nil
}

func (s DB) Test() error {
	return s.db.Exec("SELECT 1").Error
}

func (s DB) ResetTableSeqCounters() error {
	seqmap := map[string]string{
		"network_stats_id_seq":   "network_stats",
		"validator_stats_id_seq": "validator_stats",
		"addresses_id_seq":       "addresses",
		"validators_id_seq":      "validators",
		"delegations_id_seq":     "delegations",
	}

	for k, v := range seqmap {
		q := fmt.Sprintf("SELECT SETVAL('%s', (SELECT COALESCE(MAX(id), 1) FROM %s));", k, v)
		if err := s.db.Exec(q).Error; err != nil {
			return err
		}
	}

	return nil
}

func checkErr(err error) error {
	if err == gorm.ErrRecordNotFound {
		return ErrNotFound
	}
	return err
}

func getTimeRange(t time.Time, bucket string) (time.Time, time.Time) {
	switch bucket {
	case "h":
		return util.HourInterval(t)
	case "d":
		return util.DayInterval(t)
	default:
		panic("invalid time bucket")
	}
}

func prepareBucket(q, bucket string) string {
	return strings.ReplaceAll(q, "@bucket", bucket)
}
