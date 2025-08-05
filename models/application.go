package models

type Application struct {
	Name     string
	Exec     string
	Icon     string // crudo: path,index o nombre
	IconPath string // expandido
	IconIdx  int
}
