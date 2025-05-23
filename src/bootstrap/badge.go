package bootstrap

/*Badge return bootstrap alert map element
*	@badgeid string html elemnt id
*	@description Modal main content
*	@bagetype string bage type
*	@argc map options
 */
func Badge(badgeid string, text string, bagetype string, argc ...map[string]interface{}) map[string]interface{} {
	retsult := map[string]interface{}{}
	retsult["Id"] = badgeid
	retsult["Text"] = text

	switch bagetype {
	case "span":
		retsult["Type"] = 2
		break
	case "link":
		retsult["Type"] = 3
		break
	case "button":
		retsult["Type"] = 1
		break
	default:
		retsult["Type"] = 2
	}

	options := map[string]interface{}{}
	if len(argc) > 0 {
		options = argc[0]
	}
	//fmt.Println(options)
	//options
	retsult["BadgeClass"] = "badge-primary"
	if badgeClass, _ := options["badgeClass"]; badgeClass != nil {
		retsult["BadgeClass"] = badgeClass
	}

	retsult["BadgePill"] = ""
	if pillClass, _ := options["badgePill"]; pillClass != nil {
		retsult["BadgePill"] = "badge-pill"
	}

	retsult["BadgeNum"] = 0
	if number, _ := options["number"]; number != nil {
		retsult["BadgeNum"] = number
	}

	return retsult
}
