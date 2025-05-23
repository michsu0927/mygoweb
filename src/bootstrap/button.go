package bootstrap

import "fmt"

/*Button return bootstrap button map element
*	@eleid string html elemnt id
*	@description Alert main content
*	@ele string type
*	@argc map options
 */
func Button(eleid string, text string, ele string, argc ...map[string]interface{}) map[string]interface{} {

	retsult := map[string]interface{}{}
	retsult["Id"] = eleid
	retsult["Text"] = text
	switch ele {
	case "a":
		retsult["Element"] = 1
		break
	case "button":
		retsult["Element"] = 2
		break
	case "input":
		retsult["Element"] = 3
		break
	default:
		retsult["Element"] = 2
	}

	options := map[string]interface{}{}
	if len(argc) > 0 {
		options = argc[0]
	}

	//options
	retsult["ButtonClass"] = "btn-primary"
	if buttonClass, _ := options["buttonClass"]; buttonClass != nil {
		retsult["ButtonClass"] = buttonClass
	}

	if active, _ := options["active"]; active != nil {
		retsult["Active"] = "active"
	}

	if attr, _ := options["attr"]; attr != nil {
		retsult["Attr"] = attr
	}

	if strType, _ := options["type"]; strType != nil {
		retsult["Type"] = strType
	}

	if link, _ := options["link"]; link != nil {
		retsult["Link"] = link
	}

	if target, _ := options["target"]; target != nil {
		retsult["Target"] = target
	}

	if size, _ := options["size"]; size != nil {
		retsult["SizeClass"] = "btn-" + fmt.Sprintf("%s", size)
	}

	if block, _ := options["block"]; block != nil {
		retsult["BlockClass"] = "btn-block"
	}

	return retsult
}

/*ButtonGroup  returns bootstrap ButtonGroup map element
*
 */
func ButtonGroup(eleid string, argc ...map[string]interface{}) map[string]interface{} {
	retsult := map[string]interface{}{}
	retsult["Id"] = eleid
	group := map[string]interface{}{}
	if len(argc) > 0 {
		group = argc[0]
	}
	retsult["GroupButton"] = group

	options := map[string]interface{}{}
	if len(argc) > 1 {
		options = argc[1]
	}

	retsult["ButtonClass"] = "btn-primary"
	if buttonClass, _ := options["buttonClass"]; buttonClass != nil {
		retsult["ButtonClass"] = buttonClass
	}

	if ele, _ := options["radio"]; ele != nil {
		retsult["Radio"] = 1
	}

	if ele, _ := options["group"]; ele != nil {
		retsult["Group"] = 1
	}

	retsult["GroupClass"] = "btn-group"
	if groupclass := options["vertical"]; groupclass != nil {
		retsult["GroupClass"] = "btn-group-vertical"
	}

	if size, _ := options["size"]; size != nil {
		retsult["SizeClass"] = "btn-" + fmt.Sprintf("%s", size)
	}

	return retsult
}
