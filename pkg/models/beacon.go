package models

type Beacon struct {

	//unique id
	Id int64

	//path to beacon template
	Path string

	//holds information about the os (windows 6.1,ubunutu 14.1 ..)
	Os string

	//holds information about architecture x86 x64?
	Arch string

	//which language the beacon is build on assembly, python,powershell
	Lang string
}
