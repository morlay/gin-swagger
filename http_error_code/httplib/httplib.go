package httplib

type ErrorFieldModel struct {
	// 输入中出错的字段信息,这个信息为一个json字符串,方便客户端进行定位错误原因
	// 例如输入中{"name":{"alias" : "test"}}中的alias出错,则返回"name.alias"
	// 如果 alias 是数组, 且第2个元素的a字段错误,则返回"name.alias[2].a"
	Field string `json:"field"`
	// 错误信息
	Msg string `json:"msg"`
	// 错误字段位置, body, query, header, path, formData
	In string `json:"in"`
}

type GeneralError struct {
	// 详细描述
	Code int32 `json:"code"`
	// 错误信息
	Msg string `json:"msg"`
	// 错误代码
	Desc string `json:"desc"`
	// 是否能作为错误话术
	CanBeErrorTalk bool `json:"canBeTalkError"`
	// 出错字段
	ErrorFields []ErrorFieldModel `json:"errorFields"`
	// 错误溯源
	Source []string `json:"source"`
	// Request Id
	Id string `json:"id"`
}

func (g GeneralError) AddErrorField(field string, in string, msg string) GeneralError {
	g.ErrorFields = append(g.ErrorFields, ErrorFieldModel{
		Field: field,
		In:    in,
		Msg:   msg,
	})
	return g
}

func (g GeneralError) AppendSource(s string) GeneralError {
	g.Source = append(g.Source, s)
	return g
}
