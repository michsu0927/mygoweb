package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"web/src/bootstrap"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4" //use v4
)

// Demo demo page
func Demo(c echo.Context) error {
	//return c.String(http.StatusOK, "Hello, Demo!")
	data := map[string]interface{}{}
	getparam := map[string]interface{}{}
	var err error
	page, err := strconv.Atoi(c.Param("page"))
	if err != nil {
		page = 1
	}
	//c.Request().Body.Read
	defer c.Request().Body.Close()
	//取得session id
	sess, _ := session.Get("session", c)
	fmt.Printf("%s", sess.Values["SessionID"])

	options := map[string]interface{}{}
	options["size"] = "sm"                        //lg large sm small
	options["justify"] = "justify-content-center" //justify-content-center or justify-content-end empty is justify-content-start尸
	options["prevText"] = "«"                     //empty is Prev
	options["nextText"] = "»"                     //empty is Next
	//options["prevnext"] = "disabled" //do not show prev and next

	res := bootstrap.Pagination("/hello", page, 5, getparam, options)
	//fmt.Println(res)
	data["PageNavi"] = res

	//HTML := ""
	//HTML, _ = bootstrap.ExecHTML(res, "pagination") // return pagination HTML
	//fmt.Println(HTML)

	options = map[string]interface{}{}
	options["size"] = "sm"
	options["verticalcentered"] = "enabled"
	options["scrollable"] = "enabled"
	options["closeText"] = "X"
	options["enterText"] = "O"

	// blocks := []string{} //extra content blocks
	// blocks = append(blocks, "description Blocks...")
	// blocks = append(blocks, "description Blocks...")
	// options["blocks"] = blocks
	//options[""]
	//options["noButton"] = "enabled" // no display button

	//if you want to put html code inside shout add htmlSafe function https://curtisvermeeren.github.io/2017/09/14/Golang-Templates-Cheatsheet#functions
	res = bootstrap.Modal("ModalDialogExample", "This is example", "Some description ...Some description ...", options)
	//fmt.Println(res)
	data["ModalDialog"] = res
	//HTML, _ = bootstrap.ExecHTML(res, "modalDialog") // return pagination HTML
	//fmt.Println(HTML)                                //print

	options = map[string]interface{}{}
	options["strongText"] = "Strong Text"
	//options["h4Text"] = "H4 Text" //Header Title
	options["alertClass"] = "alert-warning"

	// blocks := []string{} //extra content blocks
	// blocks = append(blocks, "description Blocks...")
	// blocks = append(blocks, "description Blocks...")
	// options["blocks"] = blocks

	res = bootstrap.Alert("AlertExample", "Alert Message!", options)
	data["Alert"] = res
	//HTML, _ = bootstrap.ExecHTML(res, "alert") // return pagination HTML
	//fmt.Println(HTML)

	options = map[string]interface{}{}
	options["number"] = 3

	res = bootstrap.Badge("BadgeExample", "info btn", "button", options)
	data["BadgeButton"] = res

	options["badgeClass"] = "badge-warning"
	res = bootstrap.Badge("BadgeExample", "info link", "link", options)
	data["BadgeLink"] = res

	options["badgeClass"] = "badge-dark"
	options["badgePill"] = true
	res = bootstrap.Badge("BadgeExample", "info span", "span", options)
	data["BadgeSpan"] = res

	res = map[string]interface{}{}
	arrayBread := []map[string]string{}
	arrayBread = append(arrayBread, map[string]string{"Link": "#1", "Text": "Link1"})
	arrayBread = append(arrayBread, map[string]string{"Link": "#2", "Text": "Link2", "Target": "_blank"})
	arrayBread = append(arrayBread, map[string]string{"Link": "#3", "Text": "Link3", "Active": "active"})
	res["Breadcrumb"] = arrayBread
	res["BackgroudColor"] = "antiquewhite"
	data["TplBreadcrumb"] = res

	options = map[string]interface{}{}
	options["link"] = "http://google.com"
	options["target"] = "__blank"
	options["buttonClass"] = "badge-dark"
	res = bootstrap.Button("ButtonLinkExample", "button link example", "a", options)
	data["ButtonLinkExample"] = res

	options = map[string]interface{}{}
	options["type"] = "input"
	res = bootstrap.Button("ButtonInputExample", "button input example", "button", options)
	data["ButtonInputExample"] = res

	options = map[string]interface{}{}
	options["type"] = "button"
	options["buttonClass"] = "btn-outline-danger"
	res = bootstrap.Button("ButtonButtonExample", "button button example", "button", options)
	data["ButtonButtonExample"] = res

	options = map[string]interface{}{}
	options["radio"] = 1
	groupbutton := map[string]interface{}{}
	groupbutton["0"] = map[string]string{"Text": "buttonText", "Value": "buttonValue", "Name": "buttonName", "Active": "active", "Id": "id0"}
	groupbutton["1"] = map[string]string{"Text": "buttonText1", "Value": "buttonValue1", "Name": "buttonName", "Id": "id1"}
	groupbutton["2"] = map[string]string{"Text": "buttonText2", "Value": "buttonValue2", "Name": "buttonName", "Id": "id2"}
	res = bootstrap.ButtonGroup("ButtonGroupExample", groupbutton, options)
	data["ButtonGroup"] = res

	options = map[string]interface{}{}
	options["group"] = 1
	groupbutton = map[string]interface{}{}
	groupbutton["0"] = map[string]string{"Text": "Link1", "Id": "dropid0", "Link": "#1", "Active": "active", "Target": "_blank"}

	dropdown := []map[string]string{}
	dropdown = append(dropdown, map[string]string{"Text": "Down1", "Link": "#d1", "Id": "dropid2", "Target": "__blank"})
	dropdown = append(dropdown, map[string]string{"Text": "Down2", "Link": "#d2", "Id": "dropid3", "Target": "__blank"})
	groupbutton["1"] = map[string]interface{}{"Text": "TopDropDown", "Id": "dropid1", "Sub": dropdown}

	res = bootstrap.ButtonGroup("ButtonGroupExample", groupbutton, options)
	data["ButtonGroup2"] = res

	//html will be encoded as html specialchars
	data["HTMLTest"] = "Test <div><h1>This is Test</h1></div>"

	options = map[string]interface{}{}
	options["image"] = "https://via.placeholder.com/350x150"
	options["link"] = "#Card"
	options["linkText"] = "This is a link"
	res = bootstrap.Card("CardExample", "Card Title", "Card Text Card Text Card Text Card Text Card Text", options)
	data["Card"] = res
	//return c.String(http.StatusOK, "Hello, World demo 153!")
	return c.Render(http.StatusOK, "hello", data)
}
