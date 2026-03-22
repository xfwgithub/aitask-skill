package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

var version = "0.2.8"

// Task 任务结构
type Task struct {
	ID            int64   `json:"id"`
	UUID          string  `json:"uuid"`
	Title         string  `json:"title"`
	Description   string  `json:"description,omitempty"`
	Status        string  `json:"status"`
	Priority      int     `json:"priority"`
	Project       *string `json:"project,omitempty"`
	Tags          *string `json:"tags,omitempty"`
	Assignee      *string `json:"assignee_name,omitempty"`
	AgentType     *string `json:"agent_type,omitempty"`
	AgentModel    *string `json:"agent_model,omitempty"`
	ReviewComment *string `json:"review_comment,omitempty"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     string  `json:"updated_at"`
}

// Skill 技能结构
type Skill struct {
	tasks []Task
	db    *Database
}

// NewSkill 创建新技能实例（内存模式）
func NewSkill() *Skill {
	return &Skill{
		tasks: make([]Task, 0),
		db:    nil,
	}
}

// NewSkillWithDB 创建新技能实例（数据库模式）
func NewSkillWithDB(db *Database) *Skill {
	return &Skill{
		tasks: make([]Task, 0),
		db:    db,
	}
}

// CreateTaskInput 创建任务输入
type CreateTaskInput struct {
	Title       string   `json:"title"`
	Description string   `json:"description,omitempty"`
	Priority    int      `json:"priority,omitempty"`
	Project     string   `json:"project,omitempty"`
	Tags        []string `json:"tags,omitempty"`
	Assignee    string   `json:"assignee_name,omitempty"`
	AgentType   string   `json:"agent_type,omitempty"`
	AgentModel  string   `json:"agent_model,omitempty"`
}

// CreateTask 创建任务
func (s *Skill) CreateTask(input CreateTaskInput) map[string]interface{} {
	uuid := generateUUID()
	tagsStr := ""
	if len(input.Tags) > 0 {
		tagsStr = strings.Join(input.Tags, ",")
	}

	var project *string
	if input.Project != "" {
		project = &input.Project
	}

	task := Task{
		UUID:        uuid,
		Title:       input.Title,
		Description: input.Description,
		Status:      "pending",
		Priority:    input.Priority,
		Project:     project,
		Tags:        &tagsStr,
		Assignee:    &input.Assignee,
		AgentType:   &input.AgentType,
		AgentModel:  &input.AgentModel,
		CreatedAt:   getCurrentTime(),
		UpdatedAt:   getCurrentTime(),
	}

	// 使用数据库存储
	if s.db != nil {
		err := s.db.CreateTask(task)
		if err != nil {
			return map[string]interface{}{
				"error": fmt.Sprintf("创建任务失败：%v", err),
			}
		}
	} else {
		// 内存模式（向后兼容）
		task.ID = int64(len(s.tasks) + 1)
		s.tasks = append(s.tasks, task)
	}

	return map[string]interface{}{
		"id":      task.ID,
		"uuid":    task.UUID,
		"title":   task.Title,
		"status":  task.Status,
		"message": fmt.Sprintf("任务已创建：%s", task.Title),
	}
}

// QueryTasksInput 查询任务输入
type QueryTasksInput struct {
	Status   string `json:"status,omitempty"`
	Priority int    `json:"priority,omitempty"`
	Project  string `json:"project,omitempty"`
	Assignee string `json:"assignee_name,omitempty"`
	Keyword  string `json:"keyword,omitempty"`
	Limit    int    `json:"limit,omitempty"`
}

// QueryTasks 查询任务
func (s *Skill) QueryTasks(input QueryTasksInput) map[string]interface{} {
	var tasks []Task

	// 使用数据库查询
	if s.db != nil {
		var err error
		tasks, err = s.db.QueryTasks(input.Status, input.Priority, input.Project, input.Assignee, input.Keyword, input.Limit)
		if err != nil {
			return map[string]interface{}{
				"error": fmt.Sprintf("查询任务失败：%v", err),
			}
		}
	} else {
		// 内存模式（向后兼容）
		var filtered []Task
		for _, task := range s.tasks {
			if input.Status != "" && task.Status != input.Status {
				continue
			}
			if input.Priority != 0 && task.Priority != input.Priority {
				continue
			}
			if input.Assignee != "" && (task.Assignee == nil || *task.Assignee != input.Assignee) {
				continue
			}
			if input.Keyword != "" && !strings.Contains(task.Title, input.Keyword) {
				continue
			}
			filtered = append(filtered, task)
		}
		limit := input.Limit
		if limit == 0 || limit > len(filtered) {
			limit = len(filtered)
		}
		tasks = filtered[:limit]
	}

	return map[string]interface{}{
		"total": len(tasks),
		"tasks": tasks,
	}
}

// UpdateTaskStatusInput 更新任务状态输入
type UpdateTaskStatusInput struct {
	TaskUUID      string `json:"task_uuid"`
	NewStatus     string `json:"new_status"`
	ReviewComment string `json:"review_comment,omitempty"`
}

// UpdateTaskInput 更新任务输入
type UpdateTaskInput struct {
	TaskUUID string `json:"task_uuid"`
	Title    string `json:"title,omitempty"`
	Priority int    `json:"priority,omitempty"`
	Project  string `json:"project,omitempty"`
}

// UpdateTask 更新任务
func (s *Skill) UpdateTask(input UpdateTaskInput) map[string]interface{} {
	// 使用数据库更新
	if s.db != nil {
		err := s.db.UpdateTask(input.TaskUUID, input.Title, input.Priority, input.Project)
		if err != nil {
			return map[string]interface{}{
				"error": fmt.Sprintf("更新失败：%v", err),
			}
		}
		return map[string]interface{}{
			"uuid":    input.TaskUUID,
			"message": "任务已更新",
		}
	}

	// 内存模式（向后兼容）
	for i, task := range s.tasks {
		if task.UUID == input.TaskUUID {
			if input.Title != "" {
				s.tasks[i].Title = input.Title
			}
			if input.Priority > 0 {
				s.tasks[i].Priority = input.Priority
			}
			if input.Project != "" {
				s.tasks[i].Project = &input.Project
			}
			s.tasks[i].UpdatedAt = getCurrentTime()

			return map[string]interface{}{
				"uuid":    task.UUID,
				"message": "任务已更新",
			}
		}
	}

	return map[string]interface{}{
		"error": "任务不存在",
	}
}

// UpdateTaskStatus 更新任务状态
func (s *Skill) UpdateTaskStatus(input UpdateTaskStatusInput) map[string]interface{} {
	// 使用数据库更新
	if s.db != nil {
		err := s.db.UpdateTaskStatusWithComment(input.TaskUUID, input.NewStatus, input.ReviewComment)
		if err != nil {
			return map[string]interface{}{
				"error": fmt.Sprintf("更新失败：%v", err),
			}
		}
		return map[string]interface{}{
			"uuid":       input.TaskUUID,
			"old_status": "",
			"new_status": input.NewStatus,
			"message":    fmt.Sprintf("任务状态已更新为 %s", input.NewStatus),
		}
	}

	// 内存模式（向后兼容）
	for i, task := range s.tasks {
		if task.UUID == input.TaskUUID {
			oldStatus := task.Status
			s.tasks[i].Status = input.NewStatus
			s.tasks[i].UpdatedAt = getCurrentTime()

			return map[string]interface{}{
				"uuid":       task.UUID,
				"old_status": oldStatus,
				"new_status": input.NewStatus,
				"message":    fmt.Sprintf("任务状态已更新为 %s", input.NewStatus),
			}
		}
	}

	return map[string]interface{}{
		"error": "任务不存在",
	}
}

// ClaimTaskInput 领取任务输入
type ClaimTaskInput struct {
	TaskUUID string `json:"task_uuid"`
}

// ClaimTask 领取任务（pending → agent_working）
func (s *Skill) ClaimTask(input ClaimTaskInput) map[string]interface{} {
	return s.UpdateTaskStatus(UpdateTaskStatusInput{
		TaskUUID:  input.TaskUUID,
		NewStatus: "agent_working",
	})
}

// CompleteTaskInput 提交初审输入
type CompleteTaskInput struct {
	TaskUUID string `json:"task_uuid"`
}

// CompleteTask 提交初审（agent_working → agent_review）
func (s *Skill) CompleteTask(input CompleteTaskInput) map[string]interface{} {
	return s.UpdateTaskStatus(UpdateTaskStatusInput{
		TaskUUID:  input.TaskUUID,
		NewStatus: "agent_review",
	})
}

// ReviewTaskInput 审查任务输入
type ReviewTaskInput struct {
	TaskUUID string `json:"task_uuid"`
}

// ReviewTask 审查任务（agent_review → human_review）
func (s *Skill) ReviewTask(input ReviewTaskInput) map[string]interface{} {
	return s.UpdateTaskStatus(UpdateTaskStatusInput{
		TaskUUID:  input.TaskUUID,
		NewStatus: "human_review",
	})
}

// ApproveTaskInput 人工审核输入
type ApproveTaskInput struct {
	TaskUUID string `json:"task_uuid"`
}

// ApproveTask 人工审核通过（human_review → done）
func (s *Skill) ApproveTask(input ApproveTaskInput) map[string]interface{} {
	return s.UpdateTaskStatus(UpdateTaskStatusInput{
		TaskUUID:  input.TaskUUID,
		NewStatus: "done",
	})
}

// CancelTaskInput 取消任务输入
type CancelTaskInput struct {
	TaskUUID string `json:"task_uuid"`
}

// CancelTask 取消任务（任意状态 → cancelled）
func (s *Skill) CancelTask(input CancelTaskInput) map[string]interface{} {
	return s.UpdateTaskStatus(UpdateTaskStatusInput{
		TaskUUID:  input.TaskUUID,
		NewStatus: "cancelled",
	})
}

// GetTaskDetailInput 获取任务详情输入
type GetTaskDetailInput struct {
	TaskUUID string `json:"task_uuid"`
}

// GetTaskDetail 获取任务详情
func (s *Skill) GetTaskDetail(input GetTaskDetailInput) map[string]interface{} {
	// 使用数据库查询
	if s.db != nil {
		task, err := s.db.GetTaskByUUID(input.TaskUUID)
		if err != nil {
			return map[string]interface{}{
				"error": "任务不存在",
			}
		}
		return map[string]interface{}{
			"id":          task.ID,
			"uuid":        task.UUID,
			"title":       task.Title,
			"description": task.Description,
			"status":      task.Status,
			"priority":    task.Priority,
			"assignee":    task.Assignee,
			"agent_type":  task.AgentType,
			"agent_model": task.AgentModel,
			"tags":        task.Tags,
			"created_at":  task.CreatedAt,
			"updated_at":  task.UpdatedAt,
		}
	}

	// 内存模式（向后兼容）
	for _, task := range s.tasks {
		if task.UUID == input.TaskUUID {
			return map[string]interface{}{
				"id":          task.ID,
				"uuid":        task.UUID,
				"title":       task.Title,
				"description": task.Description,
				"status":      task.Status,
				"priority":    task.Priority,
				"assignee":    task.Assignee,
				"agent_type":  task.AgentType,
				"agent_model": task.AgentModel,
				"tags":        task.Tags,
				"created_at":  task.CreatedAt,
				"updated_at":  task.UpdatedAt,
			}
		}
	}

	return map[string]interface{}{
		"error": "任务不存在",
	}
}

// GetVersion 获取版本号
func (s *Skill) GetVersion() map[string]interface{} {
	return map[string]interface{}{
		"version": version,
	}
}

// GetDashboardStats 获取仪表盘统计
func (s *Skill) GetDashboardStats() map[string]interface{} {
	var stats map[string]int

	// 使用数据库统计
	if s.db != nil {
		var err error
		stats, err = s.db.GetDashboardStats()
		if err != nil {
			return map[string]interface{}{
				"error": fmt.Sprintf("获取统计失败：%v", err),
			}
		}
	} else {
		// 内存模式（向后兼容）
		stats = map[string]int{
			"total":         len(s.tasks),
			"pending":       0,
			"agent_working": 0,
			"agent_review":  0,
			"human_review":  0,
			"done":          0,
			"cancelled":     0,
		}
		for _, task := range s.tasks {
			switch task.Status {
			case "pending":
				stats["pending"]++
			case "agent_working":
				stats["agent_working"]++
			case "agent_review":
				stats["agent_review"]++
			case "human_review":
				stats["human_review"]++
			case "done":
				stats["done"]++
			case "cancelled":
				stats["cancelled"]++
			}
		}
	}

	return map[string]interface{}{
		"total":         stats["total"],
		"pending":       stats["pending"],
		"agent_working": stats["agent_working"],
		"agent_review":  stats["agent_review"],
		"human_review":  stats["human_review"],
		"done":          stats["done"],
		"cancelled":     stats["cancelled"],
	}
}

// RunCLI 运行 CLI 模式
func (s *Skill) RunCLI() {
	// 读取输入
	var input map[string]interface{}
	decoder := json.NewDecoder(os.Stdin)
	if err := decoder.Decode(&input); err != nil {
		fmt.Printf(`{"error": "无效的输入: %v"}`, err)
		return
	}

	// 获取函数名
	funcName, ok := input["function"].(string)
	if !ok {
		fmt.Printf(`{"error": "缺少 function 字段"}`)
		return
	}

	// 获取参数
	params := input["parameters"]

	// 调用对应的函数
	var result map[string]interface{}

	switch funcName {
	case "create_task":
		var p CreateTaskInput
		jsonBytes, _ := json.Marshal(params)
		json.Unmarshal(jsonBytes, &p)
		if p.Priority == 0 {
			p.Priority = 3 // 默认中等优先级
		}
		result = s.CreateTask(p)

	case "query_tasks":
		var p QueryTasksInput
		jsonBytes, _ := json.Marshal(params)
		json.Unmarshal(jsonBytes, &p)
		result = s.QueryTasks(p)

	case "update_task_status":
		var p UpdateTaskStatusInput
		jsonBytes, _ := json.Marshal(params)
		json.Unmarshal(jsonBytes, &p)
		result = s.UpdateTaskStatus(p)

	case "claim_task":
		var p ClaimTaskInput
		jsonBytes, _ := json.Marshal(params)
		json.Unmarshal(jsonBytes, &p)
		result = s.ClaimTask(p)

	case "submit_initial_review":
		var p CompleteTaskInput
		jsonBytes, _ := json.Marshal(params)
		json.Unmarshal(jsonBytes, &p)
		result = s.CompleteTask(p)

	case "review_task":
		var p ReviewTaskInput
		jsonBytes, _ := json.Marshal(params)
		json.Unmarshal(jsonBytes, &p)
		result = s.ReviewTask(p)

	case "approve_task":
		var p ApproveTaskInput
		jsonBytes, _ := json.Marshal(params)
		json.Unmarshal(jsonBytes, &p)
		result = s.ApproveTask(p)

	case "cancel_task":
		var p CancelTaskInput
		jsonBytes, _ := json.Marshal(params)
		json.Unmarshal(jsonBytes, &p)
		result = s.CancelTask(p)

	case "get_task_detail":
		var p GetTaskDetailInput
		jsonBytes, _ := json.Marshal(params)
		json.Unmarshal(jsonBytes, &p)
		result = s.GetTaskDetail(p)

	case "get_version":
		result = s.GetVersion()

	case "get_dashboard_stats":
		result = s.GetDashboardStats()

	default:
		result = map[string]interface{}{
			"error": fmt.Sprintf("未知的函数: %s", funcName),
		}
	}

	// 输出结果
	output, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println(string(output))
}

func main() {
	// 检查命令行参数
	if len(os.Args) > 1 && os.Args[1] == "--server" {
		// 启动 Web 服务器
		startServer()
	} else {
		// CLI 模式 - 使用数据库存储
		db, err := NewDatabase("tasks.db")
		if err != nil {
			fmt.Printf(`{"error": "数据库连接失败：%v"}`, err)
			return
		}
		skill := NewSkillWithDB(db)
		skill.RunCLI()
	}
}

// startServer 启动 Web 服务器
func startServer() {
	initServer()
}
