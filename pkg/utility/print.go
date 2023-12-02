package utility

import (
	"fmt"
	"github.com/gookit/color"
	"github.com/mukezhz/geng/pkg/constant"
	"github.com/mukezhz/geng/pkg/model"
)

func PrintColorizeProjectDetail(data model.ModuleData) {
	color.Cyanln(`
	    GENG: GENERATE GOLANG PROJECT
	
	 ██████╗ ███████╗███╗   ██╗       ██████╗ 
	██╔════╝ ██╔════╝████╗  ██║      ██╔════╝ 
	██║  ███╗█████╗  ██╔██╗ ██║█████╗██║  ███╗
	██║   ██║██╔══╝  ██║╚██╗██║╚════╝██║   ██║
	╚██████╔╝███████╗██║ ╚████║      ╚██████╔╝
	 ╚═════╝ ╚══════╝╚═╝  ╚═══╝       ╚═════╝ 
											  
	
	`)
	color.Greenln("\tThe information you have provided:\n")
	color.Cyanf("\t%-20s💻: %-15s\n", constant.ProjectName, data.ProjectName)
	color.Cyanf("\t%-20s📂: %-15s\n", constant.ProjectModuleName, data.ProjectModuleName)
	color.Cyanf("\t%-20s📚: %-15s\n", constant.ProjectDescription, data.ProjectDescription)
	color.Cyanf("\t%-20s🆚: %-15s\n", constant.GoVersion, data.GoVersion)
	color.Cyanf("\t%-20s🤓: %-15s\n", constant.Author, data.Author)
	PrintFinalStepAfterProjectInitialization(data)
	color.Redln("\n\tThank You For using 🙏🇳🇵🙏:\n")

}

func PrintColorizeModuleDetail(data model.ModuleData) {
	color.Cyanln(`
	    GENG: GENERATE GOLANG MODULE
	
	 ██████╗ ███████╗███╗   ██╗       ██████╗ 
	██╔════╝ ██╔════╝████╗  ██║      ██╔════╝ 
	██║  ███╗█████╗  ██╔██╗ ██║█████╗██║  ███╗
	██║   ██║██╔══╝  ██║╚██╗██║╚════╝██║   ██║
	╚██████╔╝███████╗██║ ╚████║      ╚██████╔╝
	 ╚═════╝ ╚══════╝╚═╝  ╚═══╝       ╚═════╝ 
											  
	
	`)
	color.Greenln("\tThe information you have provided:\n")
	color.Cyanf("\t%-20s💻: %-15s\n", constant.ProjectName, data.ProjectName)
	color.Cyanf("\t%-20s📂: %-15s\n", constant.ProjectModuleName, data.ProjectModuleName)
	color.Cyanf("\t%-20s🆚: %-15s\n", constant.GoVersion, data.GoVersion)
	PrintFinalStepAfterModuleInitialization(data)
	color.Redln("\n\tThank You For using 🙏🇳🇵🙏:\n")
}

func PrintFinalStepAfterModuleInitialization(data model.ModuleData) {
	output := fmt.Sprintf(`
	🎉 Successfully created module %v

	↪️ Restart the server to see the changes:

	🌐 Navigate to the following path:
	    %v

`, data.ModuleName, "/api/"+data.PackageName)
	color.Yellowf(output)
}

func PrintFinalStepAfterProjectInitialization(data model.ModuleData) {
	output := fmt.Sprintf(`
	💻 Change directory to project:
	    cd %v
	
	💾 Initalize git repository:
	    git init

	📚 Sync dependencies:
	    go mod tidy
	
	🕵 Copy .env.example to .env:
	    cp .env.example .env
	
	🏃 Start Project 🏃:
	    go run main.go app:serve
`, data.Directory)
	color.Yellowf(output)
}
