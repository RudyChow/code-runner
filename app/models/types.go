package models

type ContainerOption struct {
	Image          string
	Cmd            []string
	SourceFilePath string
	TargetFilePath string
}
