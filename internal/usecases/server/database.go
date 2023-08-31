package server

func (uc *Server) PingDataBase() error {
	return uc.database.Ping()
}
