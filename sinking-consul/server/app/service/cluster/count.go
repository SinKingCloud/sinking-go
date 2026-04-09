package cluster

// CountByStatus 统计status数量
func (s *service) CountByStatus(status int) (total int64, err error) {
	return s.repository.CountByStatus(status)
}

// CountAll 统计status数量
func (s *service) CountAll() (total int64, err error) {
	return s.repository.CountAll()
}
