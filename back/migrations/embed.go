// Package migrations 提供东东电竞按业务分类的 PostgreSQL 版本化迁移文件。
package migrations

import "embed"

// Files 仅包含主迁移链的 up/down SQL，不含 legacy 历史补丁。
// 由 data 层迁移工具选择方向、管理版本和执行顺序，本包不主动执行 SQL。
//
//go:embed *.up.sql *.down.sql
var Files embed.FS
