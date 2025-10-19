package points

import (
	"context"
	"database/sql"
	"fmt"
)

func AddPointsForAll(ctx context.Context, db *sql.DB, rate string) error {
	tx, err := db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil { return fmt.Errorf("begin tx: %w", err) }
	defer func() { _ = tx.Rollback() }()

	// Step 1: 先确保 points 表里有对应地址的行（初值 0），避免 UPDATE 关联不到
	if _, err := tx.ExecContext(ctx, `
INSERT INTO points (user_address, points, updated_at)
SELECT b.user_address, 0, NOW()
FROM balances b
ON CONFLICT (user_address) DO NOTHING
`); err != nil {
		return fmt.Errorf("ensure rows: %w", err)
	}

	// Step 2: 使用 UPDATE ... FROM 累加（避免 EXCLUDED + 列/表名歧义）
	if _, err := tx.ExecContext(ctx, `
WITH r AS (SELECT CAST($1 AS NUMERIC) AS rate)
UPDATE points p
SET points = p.points + acc.delta,
    updated_at = NOW()
FROM (
  SELECT b.user_address,
         (CAST(b.balance AS NUMERIC) * (SELECT rate FROM r)) AS delta
  FROM balances b
) acc
WHERE p.user_address = acc.user_address
`, rate); err != nil {
		return fmt.Errorf("update add points: %w", err)
	}

	// 记录流水（可选）
	if _, err := tx.ExecContext(ctx, `
INSERT INTO points_events (user_address, points_delta, reason)
SELECT b.user_address, (CAST(b.balance AS NUMERIC) * CAST($1 AS NUMERIC)), 'scheduler'
FROM balances b
`, rate); err != nil {
		return fmt.Errorf("insert points_events: %w", err)
	}

	if err := tx.Commit(); err != nil { return fmt.Errorf("commit: %w", err) }
	return nil
}
