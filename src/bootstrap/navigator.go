package bootstrap

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"os"
)

/*Pagination return bootstrap pagination map element
*	@uri string web link url
*	@page int page active
*	@argc[0] getParam : map HTTP GET parameters
* 	@argc[1] options : map Pagination options
 */
func Pagination(uri string, page int, totalpage int, argc ...map[string]interface{}) map[string]interface{} {
	//var retsult map[string]interface{} 只有定義 type 還沒 new variable
	retsult := map[string]interface{}{}
	retsult["prev"] = map[string]interface{}{}
	retsult["next"] = map[string]interface{}{}
	//retsult["naviLoop"] = []interface{}{}
	naviLoop := []interface{}{}
	link := uri + "/"
	strget := ""
	if len(argc) > 0 {
		getParam := argc[0]
		if len(getParam) > 0 {
			strget = "?"
			for key, value := range getParam {
				strget = strget + fmt.Sprintf("%s=%s&", key, value)
			}
		}
	}
	options := map[string]interface{}{}
	if len(argc) > 1 {
		options = argc[1]
	}
	//fmt.Println(options)

	prevText, _ := options["prevText"]
	if prevText == nil {
		prevText = "Prev"
	}
	//fmt.Println("prevText:" + prevText)

	nextText, _ := options["nextText"]
	if nextText == nil {
		nextText = "Next"
	}
	//fmt.Println("nextText:" + nextText)
	//options

	if prevnext, _ := options["prevnext"]; prevnext != "disabled" {
		if page == 1 {
			if prevText == nil {
				retsult["prev"] = map[string]interface{}{"link": "#", "disabled": true}
			} else {
				retsult["prev"] = map[string]interface{}{"link": "#", "disabled": true, "text": prevText}
			}
		}
		if page > 1 {
			link = uri + "/"
			link = link + fmt.Sprintf("%d", page-1) + "/" + strget
			if prevText == nil {
				retsult["prev"] = map[string]interface{}{"link": link, "disabled": false}
			} else {
				retsult["prev"] = map[string]interface{}{"link": link, "disabled": false, "text": prevText}
			}
		}
		if page == totalpage {
			if nextText == nil {
				retsult["next"] = map[string]interface{}{"link": "#", "disabled": true}
			} else {
				retsult["next"] = map[string]interface{}{"link": "#", "disabled": true, "text": nextText}
			}
		}
		if page < totalpage && page >= 1 {
			link = uri + "/"
			link = link + fmt.Sprintf("%d", page+1) + "/" + strget
			if nextText == nil {
				retsult["next"] = map[string]interface{}{"link": link, "disabled": false}
			} else {
				retsult["next"] = map[string]interface{}{"link": link, "disabled": false, "text": nextText}
			}
		}
	}

	if justify, _ := options["justify"]; justify != nil {
		retsult["justify"] = justify
	}

	// target := ""
	// target, _ = options["target"]
	// if target != "" {
	// 	retsult["target"] = target
	// }
	// fmt.Println("target:" + target)

	if size, _ := options["size"]; size != nil {
		if size == "lg" {
			retsult["sizeClass"] = "pagination-lg"
		}
		if size == "sm" {
			retsult["sizeClass"] = "pagination-sm"
		}
	}

	for i := 1; i <= totalpage; i++ {
		navi := map[string]interface{}{}
		if totalpage < 10 {
			if i == page {
				navi = map[string]interface{}{"link": "#", "disabled": false, "active": true, "num": i}
			} else {
				link = uri + "/"
				link = link + fmt.Sprintf("%d", i) + "/" + strget
				navi = map[string]interface{}{"link": link, "disabled": false, "num": i}
			}
			naviLoop = append(naviLoop, navi)
		} else {
			if i < 3 {
				if i == page {
					navi = map[string]interface{}{"link": "#", "disabled": false, "active": true, "num": i}
				} else {
					link = uri + "/"
					link = link + fmt.Sprintf("%d", i) + "/" + strget
					navi = map[string]interface{}{"link": link, "disabled": false, "num": i}
				}
				if page == 1 && i == 2 {
					navi["dotNext"] = true
				}
				naviLoop = append(naviLoop, navi)
			} else if i > (totalpage - 2) {
				if i == page {
					navi = map[string]interface{}{"link": "#", "disabled": false, "active": true, "num": i}
				} else {
					link = uri + "/"
					link = link + fmt.Sprintf("%d", i) + "/" + strget
					navi = map[string]interface{}{"link": link, "disabled": false, "num": i}
				}
				if page == totalpage && i == (totalpage-1) {
					navi["dotPrev"] = true
				}
				naviLoop = append(naviLoop, navi)
			} else if i >= (page-1) && i <= (page+1) {
				if i == page {
					navi = map[string]interface{}{"link": "#", "disabled": false, "active": true, "num": i}
				} else {
					link = uri + "/"
					link = link + fmt.Sprintf("%d", i) + "/" + strget
					navi = map[string]interface{}{"link": link, "disabled": false, "num": i}
				}
				if i == (page - 1) {
					navi["dotPrev"] = true
				}
				if i == (page + 1) {
					navi["dotNext"] = true
				}
				naviLoop = append(naviLoop, navi)
			}
		}
		//fmt.Println(navi)
	}
	//fmt.Println(naviLoop)
	retsult["naviLoop"] = naviLoop
	return retsult
}

/*ExecHTML return HTML
*	@m map data for template parse data
*	@name element template name
 */
func ExecHTML(m map[string]interface{}, name string) (string, error) {
	var err error
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}

	html := ""
	//t := template.New("") single template execution no need this line
	//path + "/resource/views/" + name + ".html" is the template path
	t, err := template.ParseFiles(path + "/resource/views/" + name + ".html")
	if err != nil {
		return html, err
	}
	var tpl bytes.Buffer
	if err := t.ExecuteTemplate(&tpl, "bootstrap-"+name, m); err != nil {
		return html, err
	}
	html = tpl.String()
	return html, err
}
