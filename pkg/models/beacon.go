package models


const (
	HTTPS uint32 = 1<<iota
	HTTP
	SMB
	DNS
	DNSHTTPS
	TCP
	NBNS
)


type Beacon struct{

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


type ActiveBeacon struct{

	Beacon

	//points to the parent beacon
	PID int64

	//C2-mode defines which protocol the beacon should use for communication to the server (HTTPS,HTTP,SMB...)
	C2m uint32 

	//Peer mode defines which protocol the beacon should use for communication to a peer (HTTPS,HTTP,SMB...)
	Pm uint32 

	// duration of time in which a beacon should ping the server for new commands
	Ping  float64

	//Counter how many pings the beacon missed to send
	MissedPings	int
}

