package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/google/go-github/v30/github"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/pelletier/go-toml"
	"golang.org/x/oauth2"
)

type Repository struct {
	Name  string
	Owner string
}
type Config struct {
	Repository Repository
}

type StudentEnvironmentRequest struct {
	Id     string   `json:"id"`
	Number int      `json:"number"`
	Labels []string `json:"labels"`
	Email  string   `json:"email"`
	Course string   `json:"course"`
}

func main() {
	// TODO:
	// 1. webserver to listen to student environment requests
	// 2. get location for toml from ENV variable
	// 3. read toml from disk with configuration for repo to write to
	// 4. get location for secret from ENV variable
	// 5. read secret from disk for API Token
	configLocation := os.Getenv("CONFIG_LOCATION")
	if configLocation == "" {
		configLocation = "config.toml"
	}

	tokenLocation := os.Getenv("TOKEN_LOCATION")
	if tokenLocation == "" {
		tokenLocation = "apitoken"
	}

	config := Config{}
	tomlTree, err := toml.LoadFile(configLocation)
	if err != nil {
		fmt.Printf("We could not read the config toml from %v\n", configLocation)
	}
	err = tomlTree.Unmarshal(&config)
	if err != nil {
		fmt.Printf("We could not parse the config toml: %v\n", err)
	} else {
		fmt.Printf("Successfully parsed the config toml: %v\n", config)
	}

	tokenData, err := ioutil.ReadFile(tokenLocation)
	if err != nil {
		fmt.Printf("We could not read the api token from %v\n", tokenLocation)
	}
	apiToken := string(tokenData)
	if apiToken != "" {
		fmt.Println("We retrieved the API Token")
	} else {
		fmt.Println("Warning, we did not retrieve the API Token")
	}

	e := echo.New()
	e.Use(func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &CustomContext{
				Context:         c,
				RepositoryOwner: config.Repository.Owner,
				Repository:      config.Repository.Name,
				APIToken:        apiToken,
			}
			return h(cc)
		}
	})
	e.GET("/health", healthCheck)
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.POST("/ser", studentEnvironmentRequest)
	port := os.Getenv("PORT")
	if port == "" {
		port = "1323"
	}
	e.Logger.Fatal(e.Start(":" + port))
}

func healthCheck(c echo.Context) error {
	cc := c.(*CustomContext)
	if cc.RepositoryOwner != "" && cc.Repository != "" && cc.APIToken != "" {
		return c.String(http.StatusOK, "OK")
	}
	return c.String(http.StatusInternalServerError, "No Configuration found")
}

func studentEnvironmentRequest(c echo.Context) error {
	ser := new(StudentEnvironmentRequest)
	if err := c.Bind(ser); err != nil {
		return err
	}

	ser.Id = uuid.Must(uuid.NewRandom()).String()
	cc := c.(*CustomContext)

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: cc.APIToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)
	filePath := "test.json"
	fileContent, _, _, err := client.Repositories.GetContents(ctx, cc.RepositoryOwner, cc.Repository, filePath, &github.RepositoryContentGetOptions{})
	var path, sha string
	if err != nil {
		fmt.Printf("We could not read from the repo: %v\n", err)
	} else if fileContent == nil {
		fmt.Printf("We could not read the file from repo: %v\n", filePath)
	} else {
		if fileContent.Path != nil {
			path = *fileContent.Path
		}
		if fileContent.SHA != nil {
			sha = *fileContent.SHA
		}
		fmt.Printf("Found file %v, sha: %v\n", path, sha)
	}
	updateLabel := fmt.Sprintf("updateId=%s", uuid.Must(uuid.NewRandom()).String())
	ser.Labels = []string{updateLabel}
	if fileContent.SHA != nil {
		author := "joostvdg"
		branch := "master"
		message := "CLI Auto Update"
		email := "joostvdg@gmail.com"

		commitAuthor := github.CommitAuthor{
			Name:  &author,
			Email: &email,
		}

		content, _ := json.Marshal(ser)

		fileContentOption := github.RepositoryContentFileOptions{
			Author:  &commitAuthor,
			Branch:  &branch,
			Message: &message,
			SHA:     &sha,
			Content: content,
		}
		response, _, err := client.Repositories.UpdateFile(ctx, cc.RepositoryOwner, cc.Repository, filePath, &fileContentOption)
		if err != nil {
			fmt.Printf("We could not update the file %v in repo: %v, because: %v\n", filePath, cc.Repository, err)
		} else {
			sha := ""
			if response.SHA != nil {
				sha = *response.SHA
			}
			fmt.Printf("Successfully updated the file %v, new SHA: %v\n", filePath, sha)
		}
	}
	return c.JSON(http.StatusCreated, ser)
}
