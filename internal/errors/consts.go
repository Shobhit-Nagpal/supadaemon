package errors

const (
	FileErr             = "FileError"
	ServiceErr          = "ServiceError"
	SystemErr           = "SystemError"
	FileNotFound        = "FileNotFound"
	ReadFile            = "ReadFile"
	WriteFile           = "WriteFile"
	DeleteFile          = "DeleteFile"
	RootAccess          = "RootAccess"
	ReloadService       = "ReloadService"
	StopService         = "StopService"
	DisableService      = "DisableService"
	EnableService       = "EnableService"
	StartService        = "StartService"
	UserLookup          = "UserLookup"
	ExecuteableLookup   = "ExecuteableLookup"
	ReloadDaemon        = "ReloadDaemon"
	DataSerialization   = "DataSerialization"
	DataDeserialization = "DataDeserialization"
)
