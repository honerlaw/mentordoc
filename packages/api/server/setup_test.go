package server

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/honerlaw/mentordoc/server/util"
	"log"
	"net/http"
	"os"
	"os/exec"
	"testing"
)

var integration = flag.Bool("it", false, "run integration tests")
var itTestDatabaseConnection *sql.DB

func TestMain(m *testing.M) {
	flag.Parse()

	var server *http.Server
	if *integration {
		cmd := exec.Command("bash", "-c", "docker kill mentordoc-mysql; docker rm mentordoc-mysql; docker run --name mentordoc-mysql -p 33060:3306 -e MYSQL_USER=userlocal -e MYSQL_ROOT_PASSWORD=password -e MYSQL_PASSWORD=password -e MYSQL_DATABASE=mentor_doc --tmpfs /var/lib/mysql -d mysql:5.7")
		_, err := cmd.CombinedOutput()
		if err != nil {
			log.Fatal(err)
		}

		// set envs required for everything to work
		// @todo these should be loaded from a dotenv or something similarish
		os.Setenv("HOST", "localhost")
		os.Setenv("PORT", "9000")
		os.Setenv("DATABASE_NAME", "mentor_doc")
		os.Setenv("DATABASE_USERNAME", "userlocal")
		os.Setenv("DATABASE_PASSWORD", "password")
		os.Setenv("DATABASE_HOST", "0.0.0.0")
		os.Setenv("DATABASE_PORT", "33060")
		os.Setenv("MIGRATION_DIR", "../migrations")
		os.Setenv("JWT_SIGNING_KEY", "it-test-key")

		server = StartServer(nil)

		itTestDatabaseConnection = util.NewDb()
	}

	result := m.Run()

	if *integration {
		StopServer(server)
	}

	os.Exit(result)
}

func GetTestServerAddress(path string) string {
	return fmt.Sprintf("http://%s:%s/v1%s", os.Getenv("HOST"), os.Getenv("PORT"), path)
}
