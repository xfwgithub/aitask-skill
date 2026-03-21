package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// Task 任务结构
type Task struct {
	ID          int64   `json:"id"`
	UUID        string  `json:"uuid"`
	Title       string  `json:"title"`
	Description string  `json:"description,omitempty"`
	Status      string  `json:"status"`
	Priority    int     `json:"priority"`
	Tags        *string `json:"tags,omitempty"`
	Assignee    *string `json:"assignee_name,omitempty"`
	AgentType   *string `json:"agent_type,omitempty"`
	AgentModel  *string `json:"agent_model,omitempty"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

// Skill 技能结构
type Skill struct {
	tasks []Task
}

// NewSkill 创建新技能实例
func NewSkill() *Skill {
	return &Skill{
		tasks: make([]Task, 0),
	}
}

// CreateTaskInput 创建任务输入
type CreateTaskInput struct {
	Title       string   `json:"title"`
	Description string   `json:"description,omitempty"`
	Priority    int      `json:"priority,omitempty"`
	Tags        []string `json:"tags,omitempty"`
	Assignee    string   `json:"assignee_name,omitempty"`
	AgentType   string   `json:"agent_type,omitempty"`
	AgentModel  string   `json:"agent_model,omitempty"`
}

// CreateTask 创建任务
func (s *Skill) CreateTask(input CreateTaskInput) map[string]interface{} {
	uuid := generateUUID()
	tagsStr := strings.Join(input.Tags, ",")
	
	task := Task{
		ID:          int64(len(s.tasks) + 1),
		UUID:        uuid,
		Title:       input.Title,
		Description: input.Description,
		Status:      "pending",
		Priority:    input.Priority,
		Tags:        &tagsStr,
		Assignee:    &input.Assignee,
		AgentType:   &input.AgentType,
		AgentModel:  &input.AgentModel,
		CreatedAt:   getCurrentTime(),
		UpdatedAt:   getCurrentTime(),
	}
	
	s.tasks = append(s.tasks, task)
	
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
	Assignee string `json:"assignee_name,omitempty"`
	Keyword  string `json:"keyword,omitempty"`
	Limit    int    `json:"limit,omitempty"`
}

// QueryTasks 查询任务
func (s *Skill) QueryTasks(input QueryTasksInput) map[string]interface{} {
	var filtered []Task
	
	for _, task := range s.tasks {
		// 状态筛选
		if input.Status != "" && task.Status != input.Status {
			continue
		}
		// 优先级筛选
		if input.Priority != 0 && task.Priority != input.Priority {
			continue
		}
		// 负责人筛选
		if input.Assignee != "" && (task.Assignee == nil || *task.Assignee != input.Assignee) {
			continue
		}
		// 关键词筛选
		if input.Keyword != "" && !strings.Contains(task.Title, input.Keyword) {
			continue
		}
		
		filtered = append(filtered, task)
	}
	
	// 限制数量
	limit := input.Limit
	if limit == 0 || limit > len(filtered) {
		limit = len(filtered)
	}
	
	resultTasks := filtered[:limit]
	
	return map[string]interface{}{
		"total": len(resultTasks),
		"tasks": resultTasks,
	}
}

// UpdateTaskStatusInput 更新任务状态输入
type UpdateTaskStatusInput struct {
	TaskUUID  string `json:"task_uuid"`
	NewStatus string `json:"new_status"`
}

// UpdateTaskStatus 更新任务状态
func (s *Skill) UpdateTaskStatus(input UpdateTaskStatusInput) map[string]interface{} {
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

// GetTaskDetailInput 获取任务详情输入
type GetTaskDetailInput struct {
	TaskUUID string `json:"task_uuid"`
}

// GetTaskDetail 获取任务详情
func (s *Skill) GetTaskDetail(input GetTaskDetailInput) map[string]interface{} {
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

// GetDashboardStats 获取仪表盘统计
func (s *Skill) GetDashboardStats() map[string]interface{} {
	stats := map[string]int{
		"total":         len(s.tasks),
		"pending":       0,
		"agent_working": 0,
		"done":          0,
		"cancelled":     0,
	}
	
	for _, task := range s.tasks {
		switch task.Status {
		case "pending":
			stats["pending"]++
		case "agent_working":
			stats["agent_working"]++
		case "done":
			stats["done"]++
		case "cancelled":
			stats["cancelled"]++
		}
	}
	
	return map[string]interface{}{
		"total":         stats["total"],
		"pending":       stats["pending"],
		"agent_working": stats["agent_working"],
		"done":          stats["done"],
		"cancelled":     stats["cancelled"],
	}
}

// Run 运行技能（读取 stdin，输出到 stdout）
func (s *Skill) Run() {
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
		
	case "get_task_detail":
		var p GetTaskDetailInput
		jsonBytes, _ := json.Marshal(params)
		json.Unmarshal(jsonBytes, &p)
		result = s.GetTaskDetail(p)
		
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
	skill := NewSkill()
	skill.Run()
}
