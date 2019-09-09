package test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	http2 "github.com/honerlaw/mentordoc/server/http"
	"github.com/honerlaw/mentordoc/server/lib/shared"
	"github.com/honerlaw/mentordoc/server/lib/util"
	"github.com/joho/godotenv"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"testing"
)

type GlobalTestData struct {
	ItTestDatabaseConnection *sql.DB
	Integration              *bool
	TestServer               *http2.Server
}

func InitTestData(envPath string, migrationDir string) *GlobalTestData {
	data := &GlobalTestData{
		Integration: flag.Bool("it", false, "run integration tests"),
	}

	flag.Parse()

	err := godotenv.Load(envPath)
	if err != nil {
		log.Fatal(err)
	}

	if *data.Integration {
		cmd := exec.Command("bash", "-c", fmt.Sprintf("docker kill mentordoc-mysql; docker rm mentordoc-mysql; docker run --name mentordoc-mysql -p %s:3306 -e MYSQL_USER=%s -e MYSQL_ROOT_PASSWORD=%s -e MYSQL_PASSWORD=%s -e MYSQL_DATABASE=%s --tmpfs /var/lib/mysql -d mysql:5.7", os.Getenv("DATABASE_PORT"), os.Getenv("DATABASE_USERNAME"), os.Getenv("DATABASE_PASSWORD"), os.Getenv("DATABASE_PASSWORD"), os.Getenv("DATABASE_NAME")))
		_, err := cmd.CombinedOutput()
		if err != nil {
			log.Fatal(err)
		}

		_ = os.Setenv("MIGRATION_DIR", migrationDir)

		data.TestServer = http2.StartServer(nil)
		data.ItTestDatabaseConnection = util.NewDb()
	}

	return data
}

func RunTests(m *testing.M, data *GlobalTestData) {
	result := m.Run()

	if *data.Integration {
		http2.StopServer(data.TestServer)

		exec.Command("bash", "-c", "docker kill mentordoc-mysql; docker rm mentordoc-mysql")
	}

	os.Exit(result)
}

type AuthData struct {
	User         *shared.User
	AccessToken  string
	RefreshToken string
	Organization *shared.Organization
}

func SetupAuthentication(t *testing.T, data *GlobalTestData) *AuthData {
	user := &shared.User{}
	user.Id = uuid.NewV4().String()
	user.Email = fmt.Sprintf("%s@example.com", user.Id)
	_, err := data.ItTestDatabaseConnection.Exec("insert into user (id, email, password, created_at, updated_at) values (?, ?, 'hash', 0, 0)", user.Id, user.Email)
	assert.Nil(t, err)

	// setup the org or the user
	org, err := data.TestServer.OrganizationService.Create("test")
	assert.Nil(t, err)
	err = data.TestServer.AclService.LinkUserToRole(user, "organization:owner", org.Id)
	assert.Nil(t, err)

	tokenService := util.NewTokenService()
	accessToken, err := tokenService.GenerateToken(user.Id, util.TokenAccess)
	assert.Nil(t, err)
	refreshToken, err := tokenService.GenerateToken(user.Id, util.TokenRefresh)
	assert.Nil(t, err)
	// generate the tokens we need

	return &AuthData{
		User:         user,
		Organization: org,
		AccessToken:  *accessToken,
		RefreshToken: *refreshToken,
	}
}

type RequestOptions struct {
	Method        string
	Path          string
	Headers       map[string]string
	Body          interface{}
	ResponseModel interface{}
}

/**
Utility method to make an http request and get the response body
*/
func Request(options *RequestOptions) (int, interface{}, error) {

	// attempt to convert the body to byte array
	var data []byte
	if options.Body != nil {
		marshalled, err := json.Marshal(options.Body)
		if err != nil {
			return -1, nil, err
		}
		data = marshalled
	}

	// build the url
	url := fmt.Sprintf("http://%s:%s/v1%s", os.Getenv("API_HOST"), os.Getenv("API_PORT"), options.Path)

	// create the request
	req, err := http.NewRequest(options.Method, url, bytes.NewBuffer(data))
	if err != nil {
		log.Print(err)
		return -1, nil, err
	}

	// set headers
	req.Header.Set("Content-Type", "application/json")
	for key, value := range options.Headers {
		req.Header.Set(key, value)
	}

	// do the request
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Print(err)
		return -1, nil, err
	}

	// read the content in
	data, err = ioutil.ReadAll(response.Body)
	if err != nil {
		log.Print(err)
		return -1, nil, err
	}

	if len(data) == 0 || options.ResponseModel == true {
		return response.StatusCode, data, nil
	}

	err = json.Unmarshal(data, options.ResponseModel)
	if err != nil {
		log.Print(err)
		return -1, nil, err
	}

	return response.StatusCode, options.ResponseModel, nil
}

func ConvertModel(source interface{}, target interface{}) interface{} {
	data, _ := json.Marshal(source)
	_ = json.Unmarshal(data, target)
	return target
}
