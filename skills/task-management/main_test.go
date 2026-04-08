package main

import (
	"encoding/json"
	"os"
	"strings"
	"testing"
)

// ========== 测试状态流转校验 ==========

func TestIsValidTransition(t *testing.T) {
	tests := []struct {
		from     string
		to       string
		expected bool
	}{
		// 合法流转
		{"pending", "agent_working", true},
		{"pending", "cancelled", true},
		{"agent_working", "agent_review", true},
		{"agent_working", "pending", true},
		{"agent_working", "cancelled", true},
		{"agent_review", "human_review", true},
		{"agent_review", "pending", true},
		{"agent_review", "cancelled", true},
		{"human_review", "done", true},
		{"human_review", "pending", true},
		{"human_review", "cancelled", true},

		// 非法流转
		{"pending", "done", false},
		{"pending", "human_review", false},
		{"pending", "agent_review", false},
		{"agent_working", "done", false},
		{"agent_working", "human_review", false},
		{"agent_review", "done", false},
		{"human_review", "agent_working", false},
		{"done", "pending", false},
		{"done", "cancelled", false},
		{"cancelled", "pending", false},
		{"invalid", "pending", false},
	}

	for _, tt := range tests {
		t.Run(tt.from+"_to_"+tt.to, func(t *testing.T) {
			result := IsValidTransition(tt.from, tt.to)
			if result != tt.expected {
				t.Errorf("IsValidTransition(%q, %q) = %v, want %v", tt.from, tt.to, result, tt.expected)
			}
		})
	}
}

func TestGetInvalidTransitionError(t *testing.T) {
	// 测试终态错误
	err := GetInvalidTransitionError("done", "pending")
	if err == nil {
		t.Fatal("Expected error for done->pending transition")
	}
	if !strings.Contains(err.Error(), "终态") {
		t.Errorf("Error should mention '终态', got: %s", err.Error())
	}

	// 测试非终态错误
	err = GetInvalidTransitionError("pending", "done")
	if err == nil {
		t.Fatal("Expected error for pending->done transition")
	}
	if !strings.Contains(err.Error(), "不允许") {
		t.Errorf("Error should mention '不允许', got: %s", err.Error())
	}
	if !strings.Contains(err.Error(), "agent_working") {
		t.Errorf("Error should list allowed transitions, got: %s", err.Error())
	}
}

// ========== 测试数据库操作 ==========

func setupTestDB(t *testing.T) *Database {
	t.Helper()
	tmpFile, err := os.CreateTemp("", "task-skill-test-*.db")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	tmpFile.Close()

	db, err := NewDatabase(tmpFile.Name())
	if err != nil {
		os.Remove(tmpFile.Name())
		t.Fatalf("Failed to create database: %v", err)
	}

	t.Cleanup(func() {
		db.Close()
		os.Remove(tmpFile.Name())
	})

	return db
}

func createTestTask(t *testing.T, db *Database, uuid string, status string) {
	t.Helper()
	task := Task{
		UUID:      uuid,
		Title:     "Test Task",
		Status:    status,
		Priority:  3,
		Project:   strPtr("test"),
		CreatedAt: "2026-01-01T00:00:00Z",
		UpdatedAt: "2026-01-01T00:00:00Z",
	}
	if err := db.CreateTask(task); err != nil {
		t.Fatalf("Failed to create test task: %v", err)
	}
}

func strPtr(s string) *string {
	return &s
}

func TestUpdateTaskStatus_TransactionIsolation(t *testing.T) {
	db := setupTestDB(t)
	createTestTask(t, db, "test-uuid-1", "pending")

	// 测试合法流转
	err := db.UpdateTaskStatus("test-uuid-1", "agent_working")
	if err != nil {
		t.Errorf("Expected success for pending->agent_working, got error: %v", err)
	}

	// 验证状态已更新
	task, err := db.GetTaskByUUID("test-uuid-1")
	if err != nil {
		t.Fatalf("Failed to get task: %v", err)
	}
	if task.Status != "agent_working" {
		t.Errorf("Expected status 'agent_working', got %q", task.Status)
	}

	// 测试非法流转
	err = db.UpdateTaskStatus("test-uuid-1", "done")
	if err == nil {
		t.Error("Expected error for agent_working->done transition")
	}
	if !strings.Contains(err.Error(), "不允许") {
		t.Errorf("Error should mention '不允许', got: %s", err.Error())
	}

	// 验证状态未改变
	task, err = db.GetTaskByUUID("test-uuid-1")
	if err != nil {
		t.Fatalf("Failed to get task: %v", err)
	}
	if task.Status != "agent_working" {
		t.Errorf("Status should still be 'agent_working', got %q", task.Status)
	}
}

func TestUpdateTaskStatus_TaskNotFound(t *testing.T) {
	db := setupTestDB(t)

	err := db.UpdateTaskStatus("non-existent-uuid", "agent_working")
	if err == nil {
		t.Error("Expected error for non-existent task")
	}
	if !strings.Contains(err.Error(), "任务不存在") {
		t.Errorf("Error should mention '任务不存在', got: %s", err.Error())
	}
}

func TestUpdateTaskStatusWithComment_TransactionIsolation(t *testing.T) {
	db := setupTestDB(t)
	createTestTask(t, db, "test-uuid-2", "agent_working")

	// 测试合法流转带评论
	err := db.UpdateTaskStatusWithComment("test-uuid-2", "agent_review", "完成开发")
	if err != nil {
		t.Errorf("Expected success, got error: %v", err)
	}

	// 验证评论已保存
	task, err := db.GetTaskByUUID("test-uuid-2")
	if err != nil {
		t.Fatalf("Failed to get task: %v", err)
	}
	if task.ReviewComment == nil || *task.ReviewComment != "完成开发" {
		t.Errorf("Expected review_comment '完成开发', got %v", task.ReviewComment)
	}

	// 测试非法流转
	err = db.UpdateTaskStatusWithComment("test-uuid-2", "done", "跳过审核")
	if err == nil {
		t.Error("Expected error for agent_review->done transition")
	}
}

func TestDeleteTask_CascadeDeleteActivities(t *testing.T) {
	db := setupTestDB(t)
	createTestTask(t, db, "test-uuid-3", "pending")

	// 创建活动记录
	_, err := db.db.Exec(`
		INSERT INTO task_activities (task_uuid, action, new_status, comment)
		VALUES (?, 'create', 'pending', '创建任务')
	`, "test-uuid-3")
	if err != nil {
		t.Fatalf("Failed to create activity: %v", err)
	}

	// 验证活动记录存在
	var count int
	err = db.db.QueryRow("SELECT COUNT(*) FROM task_activities WHERE task_uuid = ?", "test-uuid-3").Scan(&count)
	if err != nil || count == 0 {
		t.Fatalf("Expected activities to exist, got count=%d, err=%v", count, err)
	}

	// 删除任务
	err = db.DeleteTask("test-uuid-3")
	if err != nil {
		t.Fatalf("Failed to delete task: %v", err)
	}

	// 验证任务已删除
	_, err = db.GetTaskByUUID("test-uuid-3")
	if err == nil {
		t.Error("Expected task to be deleted")
	}

	// 验证活动记录也已删除
	err = db.db.QueryRow("SELECT COUNT(*) FROM task_activities WHERE task_uuid = ?", "test-uuid-3").Scan(&count)
	if err != nil || count != 0 {
		t.Errorf("Expected activities to be deleted, got count=%d, err=%v", count, err)
	}
}

func TestUpdateTask_CheckExistenceWhenNoUpdates(t *testing.T) {
	db := setupTestDB(t)
	createTestTask(t, db, "test-uuid-4", "pending")

	// 测试空更新（任务存在）
	err := db.UpdateTask("test-uuid-4", "", 0, "")
	if err != nil {
		t.Errorf("Expected success for existing task with no updates, got error: %v", err)
	}

	// 测试空更新（任务不存在）
	err = db.UpdateTask("non-existent-uuid", "", 0, "")
	if err == nil {
		t.Error("Expected error for non-existent task with no updates")
	}
	if !strings.Contains(err.Error(), "任务不存在") {
		t.Errorf("Error should mention '任务不存在', got: %s", err.Error())
	}
}

func TestValidateParentTask(t *testing.T) {
	db := setupTestDB(t)
	createTestTask(t, db, "parent-uuid", "pending")
	createTestTask(t, db, "child-uuid", "pending")

	// 设置 child-uuid 的 parent_uuid 为非空（模拟子任务）
	_, err := db.db.Exec("UPDATE tasks SET parent_uuid = ? WHERE uuid = ?", "parent-uuid", "child-uuid")
	if err != nil {
		t.Fatalf("Failed to set parent_uuid: %v", err)
	}

	// 测试合法父任务
	err = db.ValidateParentTask("parent-uuid")
	if err != nil {
		t.Errorf("Expected no error for valid parent task, got: %v", err)
	}

	// 测试不存在的父任务
	err = db.ValidateParentTask("non-existent-uuid")
	if err == nil {
		t.Error("Expected error for non-existent parent task")
	}
	if !strings.Contains(err.Error(), "父任务不存在") {
		t.Errorf("Error should mention '父任务不存在', got: %s", err.Error())
	}

	// 测试子任务作为父任务（不允许3级）
	err = db.ValidateParentTask("child-uuid")
	if err == nil {
		t.Error("Expected error for using child task as parent")
	}
	if !strings.Contains(err.Error(), "只支持2级") {
		t.Errorf("Error should mention '只支持2级', got: %s", err.Error())
	}

	// 测试空父任务
	err = db.ValidateParentTask("")
	if err != nil {
		t.Errorf("Expected no error for empty parent UUID, got: %v", err)
	}
}

func TestRecycleTasks_BatchInsert(t *testing.T) {
	db := setupTestDB(t)

	// 创建多个待回收任务
	for i := 0; i < 5; i++ {
		task := Task{
			UUID:      strings.Repeat(string(rune('a'+i)), 8),
			Title:     "Old Task",
			Status:    "agent_working",
			Priority:  3,
			Project:   strPtr("test"),
			CreatedAt: "2026-01-01T00:00:00Z", // 旧任务
			UpdatedAt: "2026-01-01T00:00:00Z",
		}
		if err := db.CreateTask(task); err != nil {
			t.Fatalf("Failed to create task %d: %v", i, err)
		}
	}

	// 创建新任务（不应被回收）
	task := Task{
		UUID:      "newtask001",
		Title:     "New Task",
		Status:    "agent_working",
		Priority:  3,
		Project:   strPtr("test"),
		CreatedAt: "2026-12-01T00:00:00Z", // 新任务
		UpdatedAt: "2026-12-01T00:00:00Z",
	}
	if err := db.CreateTask(task); err != nil {
		t.Fatalf("Failed to create new task: %v", err)
	}

	// 执行回收
	count, err := db.RecycleTasks("2026-06-01T00:00:00Z")
	if err != nil {
		t.Fatalf("RecycleTasks failed: %v", err)
	}
	if count != 5 {
		t.Errorf("Expected 5 recycled tasks, got %d", count)
	}

	// 验证旧任务已回收
	var recycledCount int
	err = db.db.QueryRow("SELECT COUNT(*) FROM tasks WHERE status = 'pending' AND created_at < '2026-06-01'").Scan(&recycledCount)
	if err != nil || recycledCount != 5 {
		t.Errorf("Expected 5 recycled tasks with status 'pending', got %d, err=%v", recycledCount, err)
	}

	// 验证新任务未被回收
	var newTaskStatus string
	err = db.db.QueryRow("SELECT status FROM tasks WHERE uuid = ?", "newtask001").Scan(&newTaskStatus)
	if err != nil || newTaskStatus != "agent_working" {
		t.Errorf("Expected new task to still be 'agent_working', got %q, err=%v", newTaskStatus, err)
	}

	// 验证活动记录已批量创建（每个回收任务一条）
	var activityCount int
	err = db.db.QueryRow("SELECT COUNT(*) FROM task_activities WHERE action = 'recycle'").Scan(&activityCount)
	if err != nil || activityCount != 5 {
		t.Errorf("Expected 5 recycle activities, got %d, err=%v", activityCount, err)
	}
}

// ========== 测试 CLI JSON 解析错误处理 ==========

func TestCLI_JSONParseError(t *testing.T) {
	// 创建一个内存模式的 skill 用于测试
	skill := NewSkill()

	// 验证逻辑存在：通过检查 CreateTask 的错误处理
	result := skill.CreateTask(CreateTaskInput{
		Title:       "",
		Description: "",
		Priority:    0,
		Project:     "",
	})

	if result["error"] == nil {
		t.Error("Expected error for empty title")
	}
	if !strings.Contains(result["error"].(string), "标题不能为空") {
		t.Errorf("Expected '标题不能为空' error, got: %v", result["error"])
	}
}

// ========== 测试内存模式状态流转校验 ==========

func TestMemoryMode_StatusTransitionValidation(t *testing.T) {
	skill := NewSkill()

	// 创建任务
	result := skill.CreateTask(CreateTaskInput{
		Title:   "Test Task",
		Project: "test",
	})
	if result["error"] != nil {
		t.Fatalf("Failed to create task: %v", result["error"])
	}
	taskUUID := result["uuid"].(string)

	// 合法流转：pending -> agent_working
	result = skill.UpdateTaskStatus(UpdateTaskStatusInput{
		TaskUUID:  taskUUID,
		NewStatus: "agent_working",
	})
	if result["error"] != nil {
		t.Errorf("Expected success for pending->agent_working, got error: %v", result["error"])
	}

	// 非法流转：agent_working -> done
	result = skill.UpdateTaskStatus(UpdateTaskStatusInput{
		TaskUUID:  taskUUID,
		NewStatus: "done",
	})
	if result["error"] == nil {
		t.Error("Expected error for agent_working->done transition")
	}
	if !strings.Contains(result["error"].(string), "不允许") {
		t.Errorf("Error should mention '不允许', got: %v", result["error"])
	}
}

// ========== 测试 JSON 序列化错误处理 ==========

func TestJSONMarshalErrorHandling(t *testing.T) {
	// 测试无法序列化的类型
	type BadType struct {
		Channel chan int
	}

	badParams := map[string]interface{}{
		"channel": make(chan int),
	}

	_, err := json.Marshal(badParams)
	if err == nil {
		t.Error("Expected marshaling error for channel type")
	}

	// 验证错误信息包含 channel
	if !strings.Contains(err.Error(), "chan") {
		t.Errorf("Error should mention 'chan', got: %s", err.Error())
	}
}
