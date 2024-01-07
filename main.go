package main

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"reflect"
	"runtime"
	"strings"

	"github.com/gin-gonic/gin"
)

//go:embed assets/* templates/*
var f embed.FS

func runCommand(cmd string) error {
	args := strings.Fields(cmd)
	out, err := os.exec.Command(args[0], args[1:]...).Output()
	log.Printf("%s", out)
	if err != nil {
		return err
	}
	return nil
}

type configFile struct {
	filename string
	content  string
}

func getFieldString(c *configFile, field string) string {
	r := reflect.ValueOf(c)
	f := reflect.Indirect(r).FieldByName(field)
	return f.String()
}

func writeConfigFile(c configFile) error {
	filename := getFieldString(&c, "filename")
	content := getFieldString(&c, "content")

	if err := os.WriteFile(filename, []byte(content), 0644); err != nil {
		return fmt.Errorf("%v for file %v", err, filename)
	}

	return nil
}

func bootstrap(hostname string, hostgroup string, git string, vault string) error {
	if runtime.GOOS != "linux" {
		return fmt.Errorf("OS is not supported")
	}

	commands := []string{
		"apt-get update",
		"apt-get --yes install git python3 python3-pip",
		"pip3 install --disable-pip-version-check ansible",
		"mkdir --verbose --parent /etc/ansible",
		"mkdir --verbose --parent /etc/control-repository",
	}

	for _, c := range commands {
		if err := runCommand(c); err != nil {
			return err
		}
	}

	ansibleConfig := configFile{
		filename: "/etc/ansible/ansible.cfg",
		content:  "[defaults]\ninterpreter_python = auto_silent",
	}

	if err := writeConfigFile(ansibleConfig); err != nil {
		return err
	}

	ansibleInventory := configFile{
		filename: "/etc/ansible/hosts",
		content:  fmt.Sprintf("[%s]\n%s ansible_host=127.0.0.1", hostgroup, hostname),
	}

	if err := writeConfigFile(ansibleInventory); err != nil {
		return err
	}

	ansibleVaultPassword := configFile{
		filename: "/etc/ansible/.vault_pass.txt",
		content:  vault,
	}

	if err := writeConfigFile(ansibleVaultPassword); err != nil {
		return err
	}

	return nil
}

func getHomepage(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

func getBootstrap(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Bootstrap",
	})
}

func postBootstrap(c *gin.Context) {
	inputHostname := c.PostForm("inputHostname")
	inputHostgroup := c.PostForm("inputHostgroup")
	inputGitRepo := c.PostForm("inputGitRepo")
	inputVaultPassword := c.PostForm("inputVaultPassword")

	log.Printf("hostname: %s; hostgroup: %s; gitrepo: %s; vaultpwd: %s", inputHostname, inputHostgroup, inputGitRepo, inputVaultPassword)

	err := bootstrap(inputHostname, inputHostgroup, inputGitRepo, inputVaultPassword)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("%s", err),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"message": "Bootstrap successful",
		})
	}
}

func main() {
	router := gin.Default()

	templ := template.Must(template.New("").ParseFS(f, "templates/*.html"))
	router.SetHTMLTemplate(templ)
	router.StaticFS("/public", http.FS(f))

	router.GET("/", getHomepage)
	api := router.Group("/api")
	{
		api.GET("/bootstrap", getBootstrap)
		api.POST("/bootstrap", postBootstrap)
	}

	router.Run(":5000")
}
