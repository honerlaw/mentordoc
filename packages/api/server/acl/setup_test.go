package acl_test

import (
	"database/sql"
	"flag"
	"github.com/honerlaw/mentordoc/server/acl"
	"github.com/honerlaw/mentordoc/server/util"
	"github.com/joho/godotenv"
	"log"
	"os"
	"os/exec"
	"testing"
)

var integration = flag.Bool("it", false, "run integration tests")
var itTestDatabaseConnection *sql.DB

func TestMain(m *testing.M) {
	flag.Parse()

	_ = godotenv.Load(".env.test")

	if *integration {
		cmd := exec.Command("bash", "-c", "docker kill mentordoc-mysql-acl; docker rm mentordoc-mysql-acl; docker run --name mentordoc-mysql-acl -p 33066:3306 -e MYSQL_USER=userlocal -e MYSQL_ROOT_PASSWORD=password -e MYSQL_PASSWORD=password -e MYSQL_DATABASE=mentor_doc --tmpfs /var/lib/mysql -d mysql:5.7")
		_, err := cmd.CombinedOutput()
		if err != nil {
			log.Fatal(err)
		}

		_ = os.Setenv("MIGRATION_DIR", "../../migrations")

		itTestDatabaseConnection = util.NewDb()

		service := acl.NewAclService(util.NewTransactionManager(itTestDatabaseConnection, nil), itTestDatabaseConnection, nil)
		err = service.Init()
		if err != nil {
			panic(err)
		}
	}

	result := m.Run()

	os.Exit(result)
}