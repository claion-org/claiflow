package server

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/claion-org/claiflow/pkg/server/config"
	"github.com/claion-org/claiflow/pkg/server/database"
	"github.com/claion-org/claiflow/pkg/server/migrations"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/golang-migrate/migrate/v4/source"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/pkg/errors"
)

func Migrate(cfg config.Config) error {
	latest, ok := migrations.Latests[cfg.Migrate.Source]
	if !ok {
		return errors.New("cannot found latest version info in map")
	}

	latestVersion, err := strconv.Atoi(latest.Version())
	if err != nil {
		return errors.Wrapf(err, "strconv.Atoi %s=%v", "s", latest.Version())
	}

	sourceInstance, err := iofs.New(migrations.Migrations, cfg.Migrate.Source)
	if err != nil {
		return errors.Wrapf(err, "iofs.New %s=%v", "source", cfg.Migrate.Source)
	}

	dest := fmt.Sprintf("%v://%v", cfg.Database.Type, database.FormatDSN(cfg.Database))
	mgrt, err := NewWithSourceInstance("iofs", sourceInstance, dest)
	if err != nil {
		return errors.Wrapf(err, "failed to migrate src=%q dest=%q", cfg.Migrate.Source, dest)
	}

	defer mgrt.Close()

	// get migreate version (current)
	cur_ver, cur_dirty, err := mgrt.Version()
	if err != nil && err != migrate.ErrNilVersion {
		err = errors.Wrapf(err, "failed to get current version")
		return err
	}

	// check dirty state
	if cur_dirty {
		return &migrate.ErrDirty{Version: int(cur_ver)}
	}

	if cur_ver < uint(latestVersion) {
		// do migrate goto V
		err = mgrt.Migrate(uint(latestVersion))
		if err != nil && err != migrate.ErrNoChange {
			err = errors.Wrapf(err, "failed to migrate goto version=\"%v\"", latestVersion)
			return err
		}
	}

	// get migreate version (latest)
	new_ver, new_dirty, err := mgrt.Version()
	if err != nil && err != migrate.ErrNilVersion {
		err = errors.Wrapf(err, "failed to get current version")
		return err
	}

	cols := []string{
		"",
		"driver",
		"database",
		"version",
		"status",
		"dirty",
	}

	vals := []string{
		"-",
		cfg.Database.Type,
		cfg.Database.DBName,
		fmt.Sprintf("v%v", new_ver),
		func() string {
			if cur_ver == new_ver {
				return "no change"
			} else {
				return fmt.Sprintf("v%v->v%v", cur_ver, new_ver)
			}
		}(),
		strconv.FormatBool(new_dirty),
	}

	// print migrate info
	w := os.Stdout
	defer fmt.Fprintln(w, strings.Repeat("_", 40))

	tw := tabwriter.NewWriter(w, 0, 0, 3, ' ', 0)
	defer tw.Flush()
	fmt.Fprintln(w, "database migration:")
	tw.Write([]byte(strings.Join(cols, "\t") + "\n"))
	tw.Write([]byte(strings.Join(vals, "\t") + "\n"))

	return nil
}

func NewWithSourceInstance(sourceName string, sourceInstance source.Driver, dest string) (m *migrate.Migrate, err error) {
	m, err = migrate.NewWithSourceInstance(sourceName, sourceInstance, dest)
	err = errors.Wrapf(err, "failed to new migrate")
	return
}
