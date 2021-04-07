package main

var main_template = `
package main
import(
	"log"
	"os"
	"fmt"
	
	"github.com/T-c33/smicro/server"
	{{if not .Prefix}}
	"router"
	{{else}}
		"{{.Prefix}}/router"
	{{end}}
	"{{.ImportPath}}"
)

var routerServer = &router.RouterServer{}
var (
	BUILD_TIME string
	GO_VERSION string
	GIT_COMMIT string
)
	
func main() {

	if len(os.Args) >= 2 && (os.Args[1] == "--version"|| os.Args[1] == "-v") {
		fmt.Printf("build time:%s\n", BUILD_TIME)
		fmt.Printf("go version:%s\n", GO_VERSION)
		fmt.Printf("git commit:%s\n", GIT_COMMIT)
		return
	}

	err := server.Init("{{.ServiceName}}")
	if err != nil {
		log.Fatal("init service failed, err:%v", err)
		return
	}

	{{.PackageName}}.Register{{.Service.Name}}Server(server.GRPCServer(), routerServer)
	server.Run()
}
`
