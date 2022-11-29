package main

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type DBHandler struct {
	dbPath         string
	contextTimeout time.Duration
}

func (d *DBHandler) Init(ctx context.Context) error {
	db, err := sql.Open("sqlite3", d.dbPath)
	if err != nil {
		return err
	}
	defer db.Close()
	ctx, cancel := context.WithTimeout(ctx, d.contextTimeout)
	defer cancel()
	_, err = db.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS "processes" (
		"id" INTEGER NOT NULL,
		"name" TEXT NOT NULL,
		"create_time" INTEGER NOT NULL,
		"command" TEXT NOT NULL,
		"status" INTEGER NOT NULL,
		"order_id" INTEGER NOT NULL,

		PRIMARY KEY("id" AUTOINCREMENT)
	)`)
	return err
}

func (d *DBHandler) InsertProcess(ctx context.Context, p *Process) (int64, error) {
	db, err := sql.Open("sqlite3", d.dbPath)
	if err != nil {
		return 0, err
	}
	defer db.Close()
	ctx, cancel := context.WithTimeout(ctx, d.contextTimeout)
	defer cancel()
	var id int64
	db.QueryRowContext(
		ctx,
		"insert into processes(name,create_time,command,status,order_id) values(?,?,?,?,(SELECT ifnull(max(order_id), 0) + 1 from processes)) returning id",
		p.Name,
		p.CreateTime,
		p.Command,
		p.Status,
		p.OrderId,
	).Scan(&id)
	return id, err
}

func (d *DBHandler) UpdateProcess(ctx context.Context, p *Process) error {
	db, err := sql.Open("sqlite3", d.dbPath)
	if err != nil {
		return err
	}
	defer db.Close()
	ctx, cancel := context.WithTimeout(ctx, d.contextTimeout)
	defer cancel()
	_, err = db.ExecContext(
		ctx,
		"update processes set name = ?, command = ?, status = ? where id = ?",
		p.Name,
		p.Command,
		p.Status,
		p.Id,
	)
	return err
}

func (d *DBHandler) UpdateProcessesOrderId(ctx context.Context, ids []int64) error {
	db, err := sql.Open("sqlite3", d.dbPath)
	if err != nil {
		return err
	}
	defer db.Close()
	ctx, cancel := context.WithTimeout(ctx, d.contextTimeout)
	defer cancel()
	tx, err := db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback()
	for i, v := range ids {
		_, err = tx.ExecContext(ctx, "update processes set order_id = ? where id = ?", i+1, v)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (d *DBHandler) DeleteProcess(ctx context.Context, id int64) error {
	db, err := sql.Open("sqlite3", d.dbPath)
	if err != nil {
		return err
	}
	defer db.Close()
	ctx, cancel := context.WithTimeout(ctx, d.contextTimeout)
	defer cancel()
	_, err = db.ExecContext(
		ctx,
		"delete from processes where id = ?",
		id,
	)
	return err
}

func (d *DBHandler) GetProcesses(ctx context.Context, onlyStatusOne bool) ([]Process, error) {
	db, err := sql.Open("sqlite3", d.dbPath)
	if err != nil {
		return nil, err
	}
	defer db.Close()
	ctx, cancel := context.WithTimeout(ctx, d.contextTimeout)
	defer cancel()
	query := "select id,name,create_time,command,status,order_id from processes"
	if onlyStatusOne {
		query = query + " where status = 1"
	}
	query += " order by order_id asc"
	rows, err := db.QueryContext(
		ctx,
		query,
	)
	if err != nil {
		return nil, err
	}
	ret := make([]Process, 0)
	for rows.Next() {
		p := Process{}
		err = rows.Scan(&p.Id, &p.Name, &p.CreateTime, &p.Command, &p.Status, &p.OrderId)
		if err != nil {
			return nil, err
		}
		ret = append(ret, p)
	}
	return ret, nil
}

func (d *DBHandler) GetProcess(ctx context.Context, id int64) (*Process, error) {
	db, err := sql.Open("sqlite3", d.dbPath)
	if err != nil {
		return nil, err
	}
	defer db.Close()
	ctx, cancel := context.WithTimeout(ctx, d.contextTimeout)
	defer cancel()
	row := db.QueryRowContext(
		ctx,
		"select id,name,create_time,command,status,order_id from processes where id = ?",
		id,
	)
	p := Process{}
	err = row.Scan(&p.Id, &p.Name, &p.CreateTime, &p.Command, &p.Status, &p.OrderId)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func NewDBHandler(dbName string, contextTimeout time.Duration) *DBHandler {
	return &DBHandler{
		dbPath:         dbName,
		contextTimeout: contextTimeout,
	}
}
