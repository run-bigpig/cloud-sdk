package utils

// 计算翻页 CalcOffsetAndLimit
func CalcOffsetAndLimit(page, pageSize int64) (offset, limit int64) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	offset = (page - 1) * pageSize
	limit = pageSize
	return
}
