package beamlines

import (
	"strings"

	srvConfig "github.com/CHESSComputing/golib/config"
	utils "github.com/CHESSComputing/golib/utils"
)

// SchemaFileName obtains schema file name from schema name
func SchemaFileName(sname string) string {
	var fname string
	for _, f := range srvConfig.Config.CHESSMetaData.SchemaFiles {
		if strings.Contains(f, sname) {
			fname = f
			break
		}
	}
	return utils.FullPath(fname)
}

// SchemaName extracts schema name from schema file name
func SchemaName(fname string) string {
	arr := strings.Split(fname, "/")
	return strings.Split(arr[len(arr)-1], ".")[0]
}
