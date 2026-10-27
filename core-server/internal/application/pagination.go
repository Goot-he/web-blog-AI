package application

const (
	defaultPageSize uint32 = 10  // 不传分页大小时使用
	maxPageSize     uint32 = 100 // 单页最多返回的数量
)

// pagination 把从 1 开始的页码换算成数据库的 offset/limit
func pagination(page, size uint32) (offset, limit int) {
	if page == 0 {
		page = 1
	}
	if size == 0 {
		size = defaultPageSize
	}
	if size > maxPageSize {
		size = maxPageSize
	}
	return int(uint64(page-1) * uint64(size)), int(size)
}
