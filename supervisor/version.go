package main

import (
	"fmt"
)

// 版本定义
const (
	// VERSION 程序版本号（用于命令行显示）
	VERSION = "v0.8.0"

	// APIVersion API 兼容版本号（用于 XML-RPC 接口）
	APIVersion = "3.0"
)

// VersionCommand implement the flags.Commander interface
type VersionCommand struct {
}

var versionCommand VersionCommand

// Execute implement Execute() method defined in flags.Commander interface
func (v VersionCommand) Execute(args []string) error {
	fmt.Println(VERSION)
	return nil
}

func init() {
	parser.AddCommand("version",
		"show the version of supervisor",
		"display the supervisor version",
		&versionCommand)
}
