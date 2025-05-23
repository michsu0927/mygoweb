package bootstrap

import (
	"fmt"
	"reflect"
)

/*Alert return bootstrap alert map element
*	@alertid string html elemnt id
*	@description Alert main content
*	@argc map options
 */
func Alert(alertid string, description string, argc ...map[string]interface{}) map[string]interface{} {
	retsult := map[string]interface{}{}
	retsult["Id"] = alertid
	retsult["Description"] = description
	options := map[string]interface{}{}
	if len(argc) > 0 {
		options = argc[0]
	}
	//fmt.Println(options)
	//options
	retsult["AlertClass"] = "alert-primary"
	if alertClass, _ := options["alertClass"]; alertClass != nil {
		retsult["AlertClass"] = alertClass
	}

	if strongText, _ := options["strongText"]; strongText != nil {
		retsult["Strong"] = strongText
	}

	if h4Text, _ := options["h4Text"]; h4Text != "" {
		retsult["H4"] = h4Text
	}

	if blocks, _ := options["blocks"]; fmt.Sprintf("%s", reflect.TypeOf(blocks)) == "[]string" {
		retsult["Blocks"] = blocks
	}

	return retsult
}
