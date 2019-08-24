package server_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/honerlaw/mentordoc/server"
	"github.com/honerlaw/mentordoc/server/model"
	"github.com/honerlaw/mentordoc/server/util"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"testing"
)

var integration = flag.Bool("it", false, "run integration tests")
var itTestDatabaseConnection *sql.DB
var testServer *server.Server

func TestMain(m *testing.M) {
	flag.Parse()

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

		testServer = server.StartServer(nil)

		itTestDatabaseConnection = util.NewDb()
	}

	result := m.Run()

	if *integration {
		server.StopServer(testServer)
	}

	os.Exit(result)
}

type AuthData struct {
	user         *model.User
	accessToken  string
	refreshToken string
}

func SetupAuthentication(t *testing.T) *AuthData {
	user := &model.User{}
	user.Id = uuid.NewV4().String()
	user.Email = fmt.Sprintf("%s@example.com", user.Id)
	_, err := itTestDatabaseConnection.Exec("insert into user (id, email, password, created_at, updated_at) values (?, ?, 'hash', 0, 0)", user.Id, user.Email)
	assert.Nil(t, err)

	authService := server.NewAuthenticationService()
	accessToken, err := authService.GenerateToken(user.Id, server.TokenAccess)
	assert.Nil(t, err)
	refreshToken, err := authService.GenerateToken(user.Id, server.TokenRefresh)
	assert.Nil(t, err)
	// generate the tokens we need

	return &AuthData{
		user:         user,
		accessToken:  *accessToken,
		refreshToken: *refreshToken,
	}
}

type PostOptions struct {
	Path    string
	Headers map[string]string
}

func PostItTest(options *PostOptions, body interface{}, resp interface{}) (int, interface{}, error) {
	data, err := json.Marshal(body)
	if err != nil {
		return -1, nil, err
	}

	url := fmt.Sprintf("http://%s:%s/v1%s", os.Getenv("HOST"), os.Getenv("PORT"), options.Path)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return -1, nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	for key, value := range options.Headers {
		req.Header.Set(key, value)
	}

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return -1, nil, err
	}

	data, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return -1, nil, err
	}

	if resp == true {
		return response.StatusCode, data, nil
	}

	err = json.Unmarshal(data, resp)
	if err != nil {
		return -1, nil, err

	}

	return response.StatusCode, resp, nil
}
