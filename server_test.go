package gproxy_test

//func newServer(t *testing.T) *gproxy.Server {
//	t.Helper()
//	gotenv.Must(gotenv.Load)
//	cfg := config.New()
//	db, err := postgres.New(cfg)
//	assert.NoError(t, err)
//	server := gproxy.NewServer(db)
//	return server
//}
//
//func TestNewServer(t *testing.T) {
//	var err error
//	cfg := config.New()
//	db, err := postgres.New(cfg)
//	assert.NoError(t, err)
//	server := gproxy.NewServer(db)
//	go func() {
//		err = server.Serve()
//	}()
//	time.Sleep(time.Millisecond * 100)
//	assert.NoError(t, err)
//	server.GracefulStop()
//}
