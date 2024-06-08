package sqlite

import (
	"blog/config"
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Models struct {
	db     *sql.DB
	config config.DBSetting
}

func New(db *sql.DB, config config.DBSetting) *Models {
	return &Models{
		db:     db,
		config: config,
	}
}

func (m *Models) Prepare(ctx context.Context) error {
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(m.config.Timeout)*time.Second)
	defer cancel()

	// enable sqlite foreign key
	_, err := m.db.ExecContext(ctxTimeout, "PRAGMA foreign_keys = ON;")
	if err != nil {
		return fmt.Errorf("PrepareSqlite: enable foreign key failed: %w", err)
	}

	return nil
}

// ex: SELECT * FROM xxx WHERE bbb IN '(1,3,4,5)'
//
// (1,3,4,5) <-- this is what we will generate
func genInCondition(inpt []int) (string, error) {
	var condition strings.Builder
	if _, err := condition.WriteString("("); err != nil {
		return "", fmt.Errorf("genInCondition: write string '(' failed: %w", err)
	}

	for i, id := range inpt {
		if _, err := condition.WriteString(strconv.Itoa(id)); err != nil {
			return "", fmt.Errorf("genInCondition: write string 'id' failed: %w", err)
		}
		if i != len(inpt)-1 {
			if _, err := condition.WriteString(","); err != nil {
				return "", fmt.Errorf("genInCondition: write string ',' failed: %w", err)
			}
		}
	}
	if _, err := condition.WriteString(")"); err != nil {
		return "", fmt.Errorf("genInCondition: write string ')' failed: %w", err)
	}

	return condition.String(), nil
}

// ex: SELECT * FROM xxx WHERE bbb = xx AND bbb = yy;
//
// bbb = xx AND bbb = yy <-- this is what we will generate
func genEqualCondition(name string, inpt []int) (string, error) {
	var condition strings.Builder

	for i, id := range inpt {
		str := fmt.Sprintf("%s = %d", name, id)
		if _, err := condition.WriteString(str); err != nil {
			return "", fmt.Errorf("genInCondition: write string 'id' failed: %w", err)
		}
		if i != len(inpt)-1 {
			if _, err := condition.WriteString(" AND "); err != nil {
				return "", fmt.Errorf("genInCondition: write string ',' failed: %w", err)
			}
		}
	}

	return condition.String(), nil
}
