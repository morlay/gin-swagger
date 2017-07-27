package enum

type EnumOption struct {
	// 选项
	Option string `json:"option"`
	// 说明
	Label string `json:"label"`
}

var EnumMap = map[string][]EnumOption{}

func RegistryEnum(enumType string, optionValue string, label string) {
	if _, ok := EnumMap[enumType]; !ok {
		EnumMap[enumType] = []EnumOption{}
	}
	EnumMap[enumType] = append(EnumMap[enumType], EnumOption{Option: optionValue, Label: label})
}

func GetEnumValueList(enumType string) (enumList []EnumOption, found bool) {
	enumList, found = EnumMap[enumType]
	return
}
