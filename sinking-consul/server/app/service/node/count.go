package node

// CountByOnlineStatus 统计online_status数量
func (s *service) CountByOnlineStatus(onlineStatus int) (total int64, err error) {
	return s.repository.CountByOnlineStatus(onlineStatus)
}

// CountAll 统计status数量
func (s *service) CountAll() (total int64, err error) {
	return s.repository.CountAll()
}
