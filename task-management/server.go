package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var database *Database

// initServer 初始化并启动 Web 服务器
func initServer() {
	// 初始化数据库
	var err error
	database, err = NewDatabase("tasks.db")
	if err != nil {
		fmt.Printf("数据库初始化失败：%v\n", err)
		os.Exit(1)
	}
	defer database.Close()

	// 创建 Echo 实例
	e := echo.New()

	// 中间件
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// 静态文件
	e.Static("/static", "static")

	// 模板函数
	funcMap := template.FuncMap{
		"formatDate": func(dateStr string) string {
			t, err := time.Parse(time.RFC3339, dateStr)
			if err != nil {
				return dateStr
			}
			return t.Format("2006-01-02 15:04")
		},
	}

	// 模板
	t := &Template{
		templates: template.Must(template.New("").Funcs(funcMap).ParseGlob("templates/*.html")),
	}
	e.Renderer = t

	// 路由
	e.GET("/", indexHandler)
	e.GET("/tasks", tasksHandler)
	e.GET("/tasks/:uuid", taskDetailHandler)
	e.POST("/api/tasks", createTaskAPI)
	e.GET("/api/tasks", queryTasksAPI)
	e.PUT("/api/tasks/:uuid/status", updateTaskStatusAPI)
	e.GET("/api/stats", getStatsAPI)

	// 启动服务器
	fmt.Println("服务器启动在 http://localhost:8080")
	e.Logger.Fatal(e.Start(":8080"))
}

// 模板渲染器
type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

// 页面处理器
func indexHandler(c echo.Context) error {
	return c.Render(http.StatusOK, "base.html", nil)
}

func tasksHandler(c echo.Context) error {
	return c.Render(http.StatusOK, "base.html", nil)
}

func taskDetailHandler(c echo.Context) error {
	uuid := c.Param("uuid")
	task, err := database.GetTaskByUUID(uuid)
	if err != nil {
		return c.String(http.StatusNotFound, "任务不存在")
	}
	return c.Render(http.StatusOK, "base.html", task)
}

// API 处理器
func createTaskAPI(c echo.Context) error {
	var input CreateTaskInput
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "无效的输入",
		})
	}

	if input.Priority == 0 {
		input.Priority = 3
	}

	uuid := uuid.New().String()
	tagsStr := strings.Join(input.Tags, ",")

	task := Task{
		UUID:        uuid,
		Title:       input.Title,
		Description: input.Description,
		Status:      "pending",
		Priority:    input.Priority,
		Tags:        &tagsStr,
		Assignee:    &input.Assignee,
		AgentType:   &input.AgentType,
		AgentModel:  &input.AgentModel,
		CreatedAt:   time.Now().Format(time.RFC3339),
		UpdatedAt:   time.Now().Format(time.RFC3339),
	}

	if err := database.CreateTask(task); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "创建失败",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"id":      task.ID,
		"uuid":    task.UUID,
		"title":   task.Title,
		"status":  task.Status,
		"message": "任务已创建",
	})
}

func queryTasksAPI(c echo.Context) error {
	status := c.QueryParam("status")
	priority := 0
	assignee := c.QueryParam("assignee")
	keyword := c.QueryParam("keyword")
	limit := 20

	tasks, err := database.QueryTasks(status, priority, assignee, keyword, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "查询失败",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"total": len(tasks),
		"tasks": tasks,
	})
}

func updateTaskStatusAPI(c echo.Context) error {
	uuid := c.Param("uuid")
	var input struct {
		NewStatus string `json:"new_status"`
	}

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "无效的输入",
		})
	}

	if err := database.UpdateTaskStatus(uuid, input.NewStatus); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "更新失败",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"uuid":       uuid,
		"new_status": input.NewStatus,
		"message":    "状态已更新",
	})
}

func getStatsAPI(c echo.Context) error {
	stats, err := database.GetDashboardStats()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "获取统计失败",
		})
	}

	return c.JSON(http.StatusOK, stats)
}

// 辅助函数
func contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

func jsonEscape(str string) string {
	data, _ := json.Marshal(str)
	return string(data)[1 : len(data)-1]
}
