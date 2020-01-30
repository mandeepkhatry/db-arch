package def

//ApplicationSpecificType maps language specific data type to application specific data type
var ApplicationSpecificType = map[string]string{
	"string":    "string",
	"float32":   "double",
	"float64":   "double",
	"int":       "int",
	"bool":      "bool",
	"time.Time": "datetime",
}
