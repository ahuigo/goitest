package interpolate

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/tidwall/gjson"
)

const (
	startAnchor = "${"
	endAnchor = "}"
)

var (
	referencePattern = regexp.MustCompile(`\${[-a-zA-Z_0-9.\[\]'*:?@()!=><~,]+}`)
	fullReferencePattern = regexp.MustCompile(`^\${[-a-zA-Z_0-9.\[\]'*:?@()!=><~,]+}$`)
)

//1、将slice形式的result[0]--->result.0.
//2、将循环取用的..--->.#.
// TODO: complete me! Full json path convert to gson path implementation.
func jsonPath2GsonPath(str string) string {
	strArray := strings.Split(str, ".")
	var newStringArray []string
	for _, strItem := range strArray {
		var newStrItem string
		if strItem == "" {
			newStrItem = "#"
		} else if strings.Contains(strItem, "[") {
			leftBracket := strings.IndexAny(strItem, "[")
			rightBracket := strings.IndexAny(strItem, "]")
			newStrItem = fmt.Sprintf("%s.%s", strItem[:leftBracket], strItem[leftBracket+1:rightBracket])
		} else {
			newStrItem = strItem
		}
		newStringArray = append(newStringArray, newStrItem)
	}
	return strings.Join(newStringArray, ".")
}

func Interpolation(expr string, envBytes []byte) (interface{}, error) {
	if !referencePattern.MatchString(expr) {
		// not interpolation string, return itself
		return expr, nil
	}
	// envBytes, err := json.Marshal(environment())
	// if err != nil {
	// 	return nil, err
	// }
	if fullReferencePattern.MatchString(expr) {
		// full variable reference pattern. i.e. "${a.b}"
		jsonPathExpr := expr[len(startAnchor):len(expr) - len(endAnchor)]
		gsonPathExpr := jsonPath2GsonPath(jsonPathExpr)
		value := gjson.GetBytes(envBytes, gsonPathExpr)
		return value.Value(), nil
	}
	return referencePattern.ReplaceAllStringFunc(expr, func(anchor string) string {
		jsonPathExpr := anchor[len(startAnchor):len(anchor) - len(endAnchor)]
		gsonPathExpr := jsonPath2GsonPath(jsonPathExpr)
		value := gjson.GetBytes(envBytes, gsonPathExpr)
		return value.String()
	}), nil
}