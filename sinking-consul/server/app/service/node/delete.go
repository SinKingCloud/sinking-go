package node

// DeleteByAddress 通过address删除
func (s *service) DeleteByAddress(addresses []string) (err error) {
	list2, err2 := s.repository.SelectInAddress(addresses)
	err = s.repository.DeleteByAddress(addresses)
	if list2 != nil && err2 == nil {
		for _, key := range list2 {
			if key.Group != "" && key.Name != "" {
				s.Delete(key.Group, key.Address)
			}
		}
	}
	return
}
