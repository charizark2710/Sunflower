package urlconst

var URLRoles = map[string][]string{
	"GET" + GetAllDevices:   {"Admin", "User"},
	"POST" + PostDevice:     {"Admin"},
	"GET" + GetDetailDevice: {"Admin", "User"},
	"PUT" + PutDetailDevice: {"Admin"},
	"DELETE" + DeleteDevice: {"Admin"},
}
