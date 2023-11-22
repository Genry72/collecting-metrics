package server

// PingDataBase проверка подключения к базе данных
func (uc *Server) PingDataBase() error {
	return uc.database.Ping()
}
