package utils

func Pagination(oldPage, oldSize int32) (page, size int32) {
	page, size = oldPage, oldSize
	if oldSize == -1 {
		size = 999
		page = 1
		return page, size
	}
	if oldSize == 0 {
		size = 10
	}
	if oldPage == 0 {
		page = 1
	}
	return page, size
}
