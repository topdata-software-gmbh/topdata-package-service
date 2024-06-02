package commands

import (
	"flag"
	"fmt"
	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/topdata-software-gmbh/topdata-package-util/config"
	"github.com/topdata-software-gmbh/topdata-package-util/controllers"
	"github.com/topdata-software-gmbh/topdata-package-util/gin_middleware"
	"github.com/topdata-software-gmbh/topdata-package-util/globals"
	"github.com/topdata-software-gmbh/topdata-package-util/model"
	"log"
	"net/http"
)

var webserverConfig model.WebserverConfig

var (
	portFromCliOption string
)

var webserverConfigFile string

var webserverCommand = &cobra.Command{
	Use:   "webserver",
	Short: "Start the webserver",
	Run: func(cmd *cobra.Command, args []string) {
		flag.Parse()
		router := gin.Default()

		var err error

		// ---- webserver config .. FIXME.. we use now just one global config for all functionalities (except the repo portfolio)
		fmt.Printf("---- Reading webserver config file: %s\n", webserverConfigFile)
		webserverConfig, err = config.LoadWebserverConfig(webserverConfigFile)
		if err != nil {
			log.Fatalf("Failed to load webserverConfig: %s", err)
		}

		if webserverConfig.Username != nil && webserverConfig.Password != nil {
			router.Use(gin.BasicAuth(gin.Accounts{
				*webserverConfig.Username: *webserverConfig.Password,
			}))
		}

		// pkgConfigList := config.LoadPackagePortfolioFile(PackagePortfolioFile)

		// ---- register loaded configs in middlewares
		router.Use(gin_middleware.WebserverConfigMiddleware(webserverConfig))
		router.Use(gin_middleware.PkgConfigListMiddleware(globals.PkgConfigList))

		// ---- define routes
		router.GET("/", welcomeHandler)
		router.GET("/ping", pingHandler)
		router.GET("/repositories", controllers.GetRepositoriesHandler)
		router.GET("/repository-details/:name", controllers.GetRepositoryDetailsHandler)

		// ----
		color.Cyan("Loaded %d repository configs\n", len(globals.PkgConfigList.Items))

		// ---- get port (TODO: remove, use spf13/viper)
		finalPort := portFromCliOption
		if finalPort == "" {
			if webserverConfig.Port != 0 {
				finalPort = fmt.Sprint(webserverConfig.Port)
			} else {
				finalPort = "8080"
			}
		}

		// ---- start the server
		fmt.Println("Starting server at http://localhost:" + finalPort)
		err = router.Run(":" + finalPort)
		if err != nil {
			log.Fatalf("Failed to start server: %s", err)
		}
	},
}

func welcomeHandler(c *gin.Context) {
	c.String(http.StatusOK, "Welcome to the TopData Package Service!")
}

func pingHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func init() {
	webserverCommand.Flags().StringVarP(&webserverConfigFile, "webserver-config-file", "w", "webserver-config.yaml", "Path to config file with settings for the webserver")
	appRootCommand.AddCommand(webserverCommand)
}
