//go:build dbtest
// +build dbtest

package url

var db *postgres.DB

func init() {
	time.Local = time.UTC

	var err error
	db, err = postgres.New(
		context.Background(),
		postgres.WithDBName("orchestrator_test"),
		postgres.WithHost("localhost"),
		postgres.WithPort(5432),
		postgres.WithUser("postgres"),
	)
	if err != nil {
		panic(err)
	}
}
