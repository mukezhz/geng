package model

type ModuleData struct {
	ModuleName         string
	PackageName        string
	ProjectModuleName  string
	ProjectName        string
	GoVersion          string
	ProjectDescription string
	Author             string
	Directory          string
	IsModuleGenerated  bool
}

type GoMod struct {
	Module    string
	GoVersion string
}
