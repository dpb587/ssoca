package service

func GetAuthServices(m Manager) []AuthService {
	var services []AuthService

	for _, name := range m.Services() {
		srv, _ := m.Get(name)
		authSrv, ok := srv.(AuthService)
		if ok {
			services = append(services, authSrv)
		}
	}

	return services
}
