package tglib

func (s *Storage) Get(key string) any {
	return s.Data[key]
}

func (s *Storage) Set(key string, value any) {
	s.Data[key] = value
}
