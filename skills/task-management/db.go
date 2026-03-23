package main

import (
	"database/sql"
	"strings"
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
	err = database.migrate()
	if err != nil {
		return nil, err
	}

	return database, nil
}

// migrate 数据库迁移
func (d *Database) migrate() error {
	// 先初始化表（如果不存在）
	err := d.initTables()
	if err != nil {
		return err
	}

	// 检查 project 列是否存在
	var count int
	err = d.db.QueryRow(`SELECT COUNT(*) FROM pragma_table_info('tasks') WHERE name='project'`).Scan(&count)
	if err != nil {
		return err
	}

	// 如果不存在，添加 project 列
	if count == 0 {
		_, err = d.db.Exec(`ALTER TABLE tasks ADD COLUMN project TEXT`)
		if err != nil {
			return err
		}
	}

	// 检查 review_comment 列是否存在
	err = d.db.QueryRow(`SELECT COUNT(*) FROM pragma_table_info('tasks') WHERE name='review_comment'`).Scan(&count)
	if err != nil {
		return err
	}

	// 如果不存在，添加 review_comment 列
	if count == 0 {
		_, err = d.db.Exec(`ALTER TABLE tasks ADD COLUMN review_comment TEXT`)
		if err != nil {
			return err
		}
	}
	// 检查 parent_uuid 列是否存在
	err = d.db.QueryRow(`SELECT COUNT(*) FROM pragma_table_info('tasks') WHERE name='parent_uuid'`).Scan(&count)
	if err != nil {
		return err
	}
	if count == 0 {
		_, err = d.db.Exec(`ALTER TABLE tasks ADD COLUMN parent_uuid TEXT`)
		if err != nil {
			return err
		}
	}

	// 创建索引
	_, err = d.db.Exec(`CREATE INDEX IF NOT EXISTS idx_status ON tasks(status)`)
	if err != nil {
		return err
	}
	_, err = d.db.Exec(`CREATE INDEX IF NOT EXISTS idx_priority ON tasks(priority)`)
	if err != nil {
		return err
	}
	_, err = d.db.Exec(`CREATE INDEX IF NOT EXISTS idx_assignee ON tasks(assignee_name)`)
	if err != nil {
		return err
	}
	_, err = d.db.Exec(`CREATE INDEX IF NOT EXISTS idx_project ON tasks(project)`)
	if err != nil {
		return err
	}
	_, err = d.db.Exec(`CREATE INDEX IF NOT EXISTS idx_parent_uuid ON tasks(parent_uuid)`)
	if err != nil {
		return err
	}

	return nil
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
		project TEXT,
		parent_uuid TEXT,
		tags TEXT,
		assignee_name TEXT,
		agent_type TEXT,
		agent_model TEXT,
		review_comment TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	
	CREATE TABLE IF NOT EXISTS task_activities (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		task_uuid TEXT NOT NULL,
		action TEXT NOT NULL,
		old_status TEXT,
		new_status TEXT,
		comment TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	
	CREATE INDEX IF NOT EXISTS idx_status ON tasks(status);
	CREATE INDEX IF NOT EXISTS idx_priority ON tasks(priority);
	CREATE INDEX IF NOT EXISTS idx_assignee ON tasks(assignee_name);
	CREATE INDEX IF NOT EXISTS idx_project ON tasks(project);
	CREATE INDEX IF NOT EXISTS idx_activity_task ON task_activities(task_uuid);
	`

	_, err := d.db.Exec(query)
	return err
}

// CreateTask 创建任务
func (d *Database) CreateTask(task Task) error {
	tx, err := d.db.Begin()
	if err != nil {
		return err
	}

	query := `
	INSERT INTO tasks (uuid, title, description, status, priority, project, parent_uuid, tags, assignee_name, agent_type, agent_model, created_at, updated_at)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err = tx.Exec(query,
		task.UUID,
		task.Title,
		task.Description,
		task.Status,
		task.Priority,
		task.Project,
		task.ParentUUID,
		task.Tags,
		task.Assignee,
		task.AgentType,
		task.AgentModel,
		task.CreatedAt,
		task.UpdatedAt,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	activityQuery := `
	INSERT INTO task_activities (task_uuid, action, new_status, comment)
	VALUES (?, 'create', ?, ?)
	`
	_, err = tx.Exec(activityQuery, task.UUID, task.Status, "创建任务")
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

// GetTaskByUUID 根据 UUID 获取任务
func (d *Database) GetTaskByUUID(uuid string) (*Task, error) {
	query := `
	SELECT id, uuid, title, description, status, priority, project, parent_uuid, tags, assignee_name, agent_type, agent_model, review_comment, created_at, updated_at
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
		&task.Project,
		&task.ParentUUID,
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
func (d *Database) QueryTasks(status string, priority int, project string, assignee string, keyword string, limit int) ([]Task, error) {
	query := `
	SELECT id, uuid, title, description, status, priority, project, parent_uuid, tags, assignee_name, agent_type, agent_model, review_comment, created_at, updated_at
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

	if project != "" {
		query += " AND project = ?"
		args = append(args, project)
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
			&task.Project,
			&task.ParentUUID,
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
	// 获取旧状态
	var oldStatus string
	err := d.db.QueryRow("SELECT status FROM tasks WHERE uuid = ?", uuid).Scan(&oldStatus)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	tx, err := d.db.Begin()
	if err != nil {
		return err
	}

	query := `
	UPDATE tasks
	SET status = ?, updated_at = ?
	WHERE uuid = ?
	`
	_, err = tx.Exec(query, newStatus, time.Now().Format(time.RFC3339), uuid)
	if err != nil {
		tx.Rollback()
		return err
	}

	if oldStatus != newStatus {
		activityQuery := `
		INSERT INTO task_activities (task_uuid, action, old_status, new_status)
		VALUES (?, 'status_change', ?, ?)
		`
		_, err = tx.Exec(activityQuery, uuid, oldStatus, newStatus)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

// UpdateTaskStatusWithComment 更新任务状态并添加审核意见
func (d *Database) UpdateTaskStatusWithComment(uuid string, newStatus string, comment string) error {
	// 获取旧状态
	var oldStatus string
	err := d.db.QueryRow("SELECT status FROM tasks WHERE uuid = ?", uuid).Scan(&oldStatus)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	tx, err := d.db.Begin()
	if err != nil {
		return err
	}

	query := `
	UPDATE tasks
	SET status = ?, review_comment = ?, updated_at = ?
	WHERE uuid = ?
	`
	_, err = tx.Exec(query, newStatus, comment, time.Now().Format(time.RFC3339), uuid)
	if err != nil {
		tx.Rollback()
		return err
	}

	activityQuery := `
	INSERT INTO task_activities (task_uuid, action, old_status, new_status, comment)
	VALUES (?, 'status_change_with_comment', ?, ?, ?)
	`
	_, err = tx.Exec(activityQuery, uuid, oldStatus, newStatus, comment)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

// UpdateTask 更新任务信息
func (d *Database) UpdateTask(uuid string, title string, priority int, project string) error {
	updates := []string{}
	args := []interface{}{}

	if title != "" {
		updates = append(updates, "title = ?")
		args = append(args, title)
	}
	if priority > 0 {
		updates = append(updates, "priority = ?")
		args = append(args, priority)
	}
	if project != "" {
		updates = append(updates, "project = ?")
		args = append(args, project)
	}

	if len(updates) == 0 {
		return nil
	}

	updates = append(updates, "updated_at = ?")
	args = append(args, time.Now().Format(time.RFC3339))
	args = append(args, uuid)

	query := `
	UPDATE tasks
	SET ` + strings.Join(updates, ", ") + `
	WHERE uuid = ?
	`

	_, err := d.db.Exec(query, args...)
	return err
}

type TaskActivity struct {
	ID        int    `json:"id"`
	TaskUUID  string `json:"task_uuid"`
	Action    string `json:"action"`
	OldStatus string `json:"old_status"`
	NewStatus string `json:"new_status"`
	Comment   string `json:"comment"`
	CreatedAt string `json:"created_at"`
}

// GetTaskActivities 获取任务的活动记录
func (d *Database) GetTaskActivities(uuid string) ([]TaskActivity, error) {
	query := `
	SELECT id, task_uuid, action, IFNULL(old_status, ''), IFNULL(new_status, ''), IFNULL(comment, ''), created_at
	FROM task_activities
	WHERE task_uuid = ?
	ORDER BY created_at DESC
	`
	rows, err := d.db.Query(query, uuid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var activities []TaskActivity
	for rows.Next() {
		var act TaskActivity
		err := rows.Scan(
			&act.ID,
			&act.TaskUUID,
			&act.Action,
			&act.OldStatus,
			&act.NewStatus,
			&act.Comment,
			&act.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		activities = append(activities, act)
	}
	return activities, nil
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

// RecycleTasks 回收到期未完成的任务
// 将 due_date 之前创建的、状态为 agent_working 的任务重置为 pending
func (d *Database) RecycleTasks(dueDate string) (int64, error) {
	// 找出所有需要回收的任务 UUID
	querySelect := `
	SELECT uuid FROM tasks
	WHERE created_at < ? AND status = 'agent_working'
	`
	rows, err := d.db.Query(querySelect, dueDate)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	var uuids []string
	for rows.Next() {
		var uuid string
		if err := rows.Scan(&uuid); err != nil {
			return 0, err
		}
		uuids = append(uuids, uuid)
	}

	if len(uuids) == 0 {
		return 0, nil
	}

	tx, err := d.db.Begin()
	if err != nil {
		return 0, err
	}

	queryUpdate := `
	UPDATE tasks
	SET status = 'pending',
		updated_at = ?
	WHERE created_at < ?
	  AND status = 'agent_working'
	`
	result, err := tx.Exec(queryUpdate, getCurrentTime(), dueDate)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	activityQuery := `
	INSERT INTO task_activities (task_uuid, action, old_status, new_status, comment)
	VALUES (?, 'recycle', 'agent_working', 'pending', '任务超时被系统回收')
	`
	for _, uuid := range uuids {
		_, err = tx.Exec(activityQuery, uuid)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

// DeleteTask 物理删除任务
func (d *Database) DeleteTask(uuid string) error {
	query := `
	DELETE FROM tasks
	WHERE uuid = ?
	`

	_, err := d.db.Exec(query, uuid)
	return err
}

// Close 关闭数据库连接
func (d *Database) Close() error {
	return d.db.Close()
}
