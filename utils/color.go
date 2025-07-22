package utils

import (
	"github.com/fatih/color"
)

// 由专业的库来处理多平台的颜色输出效果
var (
	Red     = color.New(color.FgRed)                // 红色 31
	Green   = color.New(color.FgGreen)              // 绿色 32
	Yellow  = color.New(color.FgYellow)             // 黄色 33
	Blue    = color.New(color.FgBlue, color.Bold)   // 蓝色 34
	Magenta = color.New(color.FgMagenta)            // 紫红色 35
	Cyan    = color.New(color.FgHiCyan, color.Bold) // 青色 36
	White   = color.New(color.FgWhite)              // 白色 37
)
