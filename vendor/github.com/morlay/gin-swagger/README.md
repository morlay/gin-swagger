# Gin Swagger

Pick Swagger from code which wrote by [gin](https://github.com/gin-gonic/gin)

## Rules

### Path & Method

- [x] pick path and method from `*gin.Engine` or `*gin.RouterGroup`
- [x] path of `*gin.RouterGroup` will be prefix of path of `*gin.Engine`
- [x] only support `GET` `POST` `PUT` `PATCH` `HEAD` `DELETE` `OPTIONS`
- [x] gin-style router will be convert to swagger-style, `:id` => `{id}`, double check with parameter definitions, undefined path parameter will be use `0` instead.

### Operation

- [x] pick operation in scope of gin-handler (not support anonymous func, we need func name as operationId)
- [x] name of gin-handler will be `operationId`
- [x] only support single gin-handler.
 
#### Parameter

- [x] struct type of variable `req` or `request` in scope of gin-handler will be used for picking parameters.
- [x] tag `in` of struct field must be defined, expect body parameter, but need to use fieldName `Body`. 
- [x] tag `json` will be used as `name` 
- [x] struct type of anonymous struct field will be picked too.
- [x] others will be same as Schema

#### Response

- [x] `status` will be picked by gin-context render method.
- [x] `c.JSON` will set schema by type of return value and with produce `application/json`
- [x] `c.HTML` will with produce `application/html`
- [x] `c.Rediect` `c.Data` and `c.Render` will be no responce

#### Schema

- [x] only support `json`
- [x] basic type will be translated, but `json:"key,string"` will force converting to `string`
- [x] tag `default` will be set the default value, if it not exists, we will set field `required`
- [x] tag `validate` will be set common validations, for example, `validate:"@int[0,100)"` will be `{ "minimum": 0, "exclusiveMinimum": true, "maximum": 100 }`
- [x] anonymous struct field will be used with `allOf`

##### Enums

- [x] pick `enum` from commented `swagger:enum` type

- [x] string `enum` from const

```go
// swagger:enum State
type State int

const (
	STATE_UNKNOWN = iota
	STATE__ONE    // one
	STATE__TWO    // two
	STATE__THREE  // three
)
``` 
will be 
```json
{
  "enum": [
    "ONE",
    "TWO",
    "THREE"
  ],
  "x-enum-labels": [
    "one",
    "two",
    "three"
  ],
  "x-enum-type": "State"
}
```

- [x] `validate:"@string{ONE,TWO}"` or `validate:"@int{1,2}"` will be used for partial pick enum values;  

##### String format

- [x] pick format from commented `swagger:strfmt <format-name>` type

#### Definitions

- [x] only collect the named complex type, like struct type, slice type, map type

