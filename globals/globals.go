package globals

import (
	"github.com/panjf2000/ants/v2"
)

// 全局协程池（首字母大写，允许其他包访问）
var TaskPool *ants.Pool
