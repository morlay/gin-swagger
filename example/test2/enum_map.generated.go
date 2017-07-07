package test2

type EnumOption struct {
	// 选项
	Option string `json:"option"`
	// 说明
	Label string `json:"label"`
}

var EnumsMap = map[string][]EnumOption{}

func addEnumMap(enum string, option string, label string) {
	if _, ok := EnumsMap[enum]; !ok {
		EnumsMap[enum] = []EnumOption{}
	}
	EnumsMap[enum] = append(EnumsMap[enum], EnumOption{Option: option, Label: label})
}

func GetEnumValueList(enum string) (enumList []EnumOption, found bool) {
	enumList, found = EnumsMap[enum]
	return
}

func init() {
	addEnumMap("State", "ONE", "one")
	addEnumMap("State", "TWO", "two")
	addEnumMap("State", "THREE", "three")
	addEnumMap("State", "FOUR", "four")
}
