package bootstrap

import (
	"fmt"
	"strconv"
)

/*Card return bootstrap Card map element
*	@eleid string html elemnt id
*	@description Alert main content
*	@ele string type
*	@argc map options
 */
func Card(eleid string, title string, text string, argc ...map[string]interface{}) map[string]interface{} {

	retsult := map[string]interface{}{}
	retsult["Id"] = eleid
	retsult["Title"] = title
	retsult["Text"] = text

	options := map[string]interface{}{}
	if len(argc) > 0 {
		options = argc[0]
	}

	//options
	retsult["LinkClass"] = "btn btn-primary"
	if linkClass, _ := options["linkClass"]; linkClass != nil {
		retsult["LinkClass"] = linkClass
	}

	retsult["LinkText"] = ""
	if linkText, _ := options["linkText"]; linkText != nil {
		retsult["LinkText"] = linkText
	}

	if linkTarget, _ := options["linkTarget"]; linkTarget != nil {
		retsult["LinkTarget"] = linkTarget
	}

	if link, _ := options["link"]; link != nil {
		retsult["Link"] = link
	}

	//width 0~100
	retsult["Width"] = "25"
	if width, _ := options["width"]; width != nil {
		w, _ := strconv.Atoi(fmt.Sprintf("%v", width))
		if (w > 100) || (w < 0) {
			retsult["Width"] = "25"
		} else {
			retsult["Width"] = width
		}
	}

	if image, _ := options["image"]; image != nil {
		retsult["Image"] = image
	}

	if imageTitle, _ := options["imageTitle"]; imageTitle != nil {
		retsult["ImageTitle"] = imageTitle
	}

	return retsult
}
