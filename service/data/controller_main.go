package data

// ControllerOutput controllerの出力のルート
type ControllerOutput struct {
	DBMS     string
	Version  string
	Database string
	Tables   []*Table
}
