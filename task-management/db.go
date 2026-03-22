package main

import (
	"database/sql"
	"time"

	_ "modernc.org/sqlite"
)

// Database 数据库结构
type Database struct {
	db *sql.DB
}

// NewDatabase 创建数据库连接
func NewDatabase(dbPath string) (*Database, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	database := &Database{db: db}
	err = database.initTables()
	if err != nil {
		return nil, err
	}

	return database, nil
}

// initTables 初始化数据库表
func (d *Database) initTables() error {
	query := `
	CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		uuid TEXT UNIQUE NOT NULL,
		title TEXT NOT NULL,
		description TEXT,
		status TEXT DEFAULT 'pending',
		priority INTEGER DEFAULT 3,
		tags TEXT,
		assignee_name TEXT,
		agent_type TEXT,
		agent_model TEXT,
		review_comment TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	
	CREATE INDEX IF NOT EXISTS idx_status ON tasks(status);
	CREATE INDEX IF NOT EXISTS idx_priority ON tasks(priority);
	CREATE INDEX IF NOT EXISTS idx_assignee ON tasks(assignee_name);
	`

	_, err := d.db.Exec(query)
	return err
}

// CreateTask 创建任务
func (d *Database) CreateTask(task Task) error {
	query := `
	INSERT INTO tasks (uuid, title, description, status, priority, tags, assignee_name, agent_type, agent_model, created_at, updated_at)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := d.db.Exec(query,
		task.UUID,
		task.Title,
		task.Description,
		task.Status,
		task.Priority,
		task.Tags,
		task.Assignee,
		task.AgentType,
		task.AgentModel,
		task.CreatedAt,
		task.UpdatedAt,
	)

	return err
}

// GetTaskByUUID 根据 UUID 获取任务
func (d *Database) GetTaskByUUID(uuid string) (*Task, error) {
	query := `
	SELECT id, uuid, title, description, status, priority, tags, assignee_name, agent_type, agent_model, review_comment, created_at, updated_at
	FROM tasks
	WHERE uuid = ?
	`

	task := &Task{}
	err := d.db.QueryRow(query, uuid).Scan(
		&task.ID,
		&task.UUID,
		&task.Title,
		&task.Description,
		&task.Status,
		&task.Priority,
		&task.Tags,
		&task.Assignee,
		&task.AgentType,
		&task.AgentModel,
		&task.ReviewComment,
		&task.CreatedAt,
		&task.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return task, nil
}

// QueryTasks 查询任务
func (d *Database) QueryTasks(status string, priority int, assignee string, keyword string, limit int) ([]Task, error) {
	query := `
	SELECT id, uuid, title, description, status, priority, tags, assignee_name, agent_type, agent_model, review_comment, created_at, updated_at
	FROM tasks
	WHERE 1=1
	`

	args := []interface{}{}

	if status != "" {
		query += " AND status = ?"
		args = append(args, status)
	}

	if priority != 0 {
		query += " AND priority = ?"
		args = append(args, priority)
	}

	if assignee != "" {
		query += " AND assignee_name = ?"
		args = append(args, assignee)
	}

	if keyword != "" {
		query += " AND title LIKE ?"
		args = append(args, "%"+keyword+"%")
	}

	query += " ORDER BY created_at DESC"

	if limit > 0 {
		query += " LIMIT ?"
		args = append(args, limit)
	}

	rows, err := d.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		err := rows.Scan(
			&task.ID,
			&task.UUID,
			&task.Title,
			&task.Description,
			&task.Status,
			&task.Priority,
			&task.Tags,
			&task.Assignee,
			&task.AgentType,
			&task.AgentModel,
			&task.ReviewComment,
			&task.CreatedAt,
			&task.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, rows.Err()
}

// UpdateTaskStatus 更新任务状态
func (d *Database) UpdateTaskStatus(uuid string, newStatus string) error {
	query := `
	UPDATE tasks
	SET status = ?, updated_at = ?
	WHERE uuid = ?
	`

	_, err := d.db.Exec(query, newStatus, time.Now().Format(time.RFC3339), uuid)
	return err
}

// UpdateTaskStatusWithComment 更新任务状态并添加审核意见
func (d *Database) UpdateTaskStatusWithComment(uuid string, newStatus string, comment string) error {
	query := `
	UPDATE tasks
	SET status = ?, review_comment = ?, updated_at = ?
	WHERE uuid = ?
	`

	_, err := d.db.Exec(query, newStatus, comment, time.Now().Format(time.RFC3339), uuid)
	return err
}

// GetDashboardStats 获取统计信息
func (d *Database) GetDashboardStats() (map[string]int, error) {
	stats := map[string]int{
		"total":         0,
		"pending":       0,
		"agent_working": 0,
		"agent_review":  0,
		"human_review":  0,
		"done":          0,
		"cancelled":     0,
	}

	// 总数
	var total int
	err := d.db.QueryRow("SELECT COUNT(*) FROM tasks").Scan(&total)
	if err != nil {
		return nil, err
	}
	stats["total"] = total

	// 按状态分组
	query := `
	SELECT status, COUNT(*) as count
	FROM tasks
	GROUP BY status
	`

	rows, err := d.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var status string
		var count int
		err := rows.Scan(&status, &count)
		if err != nil {
			return nil, err
		}

		switch status {
		case "pending":
			stats["pending"] = count
		case "agent_working":
			stats["agent_working"] = count
		case "agent_review":
			stats["agent_review"] = count
		case "human_review":
			stats["human_review"] = count
		case "done":
			stats["done"] = count
		case "cancelled":
			stats["cancelled"] = count
		}
	}

	return stats, rows.Err()
}

// Close 关闭数据库连接
func (d *Database) Close() error {
	return d.db.Close()
}
