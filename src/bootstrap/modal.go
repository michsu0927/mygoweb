package bootstrap

import (
	"fmt"
	"reflect"
)

/*Modal return bootstrap modal map element
*	@modalid string html elemnt id
*	@title string Modal header
*	@description Modal main content
*	@argc map options
 */
func Modal(modalid string, title string, description string, argc ...map[string]interface{}) map[string]interface{} {
	retsult := map[string]interface{}{}
	retsult["Id"] = modalid
	retsult["Title"] = title
	retsult["Description"] = description
	options := map[string]interface{}{}
	if len(argc) > 0 {
		options = argc[0]
	}
	//fmt.Println(options)
	//options
	retsult["CloseText"] = "Close"
	if closeText, _ := options["closeText"]; closeText != nil {
		retsult["CloseText"] = closeText
	}

	retsult["EnterText"] = "Enter"
	if enterText, _ := options["enterText"]; enterText != nil {
		retsult["EnterText"] = enterText
	}

	if scrollable, _ := options["scrollable"]; scrollable == "enabled" {
		retsult["scrollableClass"] = "modal-dialog-scrollable"
	}

	if verticalcentered, _ := options["verticalcentered"]; verticalcentered == "enabled" {
		retsult["verticalcenterClass"] = "modal-dialog-centered"
	}

	size, _ := options["verticalcentered"]
	if size == "sm" {
		retsult["sizeClass"] = "modal-sm"
	}
	if size == "lg" {
		retsult["sizeClass"] = "modal-lg"
	}
	if size == "xl" {
		retsult["sizeClass"] = "modal-xl"
	}

	if nobutton, _ := options["noButton"]; nobutton == "enabled" {
		retsult["noButton"] = "Y"
	}

	if blocks, _ := options["blocks"]; fmt.Sprintf("%s", reflect.TypeOf(blocks)) == "[]string" {
		retsult["Blocks"] = blocks
	}

	//modal-dialog-centered
	//modal-dialog-scrollable
	return retsult
}
