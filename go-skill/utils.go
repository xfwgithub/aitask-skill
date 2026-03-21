package main

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// generateUUID 生成 UUID
func generateUUID() string {
	return uuid.New().String()
}

// getCurrentTime 获取当前时间字符串
func getCurrentTime() string {
	return time.Now().Format(time.RFC3339)
}

// formatPriority 格式化优先级显示
func formatPriority(priority int) string {
	switch priority {
	case 1:
		return "Critical"
	case 2:
		return "High"
	case 3:
		return "Medium"
	case 4:
		return "Low"
	default:
		return fmt.Sprintf("P%d", priority)
	}
}
