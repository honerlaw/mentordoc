package server

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"testing"
	"time"
)

var integration = flag.Bool("it", false, "run integration tests")

func TestMain(m *testing.M) {
	flag.Parse()

	var server *http.Server
	if *integration {
		cmd := exec.Command("bash", "-c", "docker kill mentordoc-mysql; docker rm mentordoc-mysql; docker run --name mentordoc-mysql -p 33060:3306 -e MYSQL_USER=userlocal -e MYSQL_ROOT_PASSWORD=password -e MYSQL_PASSWORD=password -e MYSQL_DATABASE=mentor_doc --tmpfs /var/lib/mysql -d mysql:5.7")
		_, err := cmd.CombinedOutput()
		if err != nil {
			log.Fatal(err)
		}

		// @todo actually check when the docker container is up and running
		time.Sleep(10 * time.Second)

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

		server = StartServer(nil)
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
