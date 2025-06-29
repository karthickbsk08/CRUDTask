package connectdb

import "log"

var GDB DBInstance

func Build_Global_Connections() error {
	var lErr error
	GDB.GPostgres, GDB.GRMPostgres, lErr = LocalDBConnect(POSTGRES)
	if lErr != nil {
		log.Println("Error while building postgres connection : ", lErr.Error())
		return lErr
	}
	return nil
}
