package wikicollector

import (
	"fmt"
	"net/url"
	"strings"
)

func (ra *RecentAction) String() string {
	res := ""
	res += fmt.Sprintf("<a href='%s'><b>%s</b></a> edited by <a href='%s'><b>%s</b></a> (%s%d) \n\n",
		ra.getURLTitle(),
		ra.Title,
		ra.getURLUser(),
		ra.User,
		addPlusIfNonNegative(ra.NewLen-ra.OldLen),
		ra.NewLen-ra.OldLen,
	)
	if len(ra.Comment) > 0 {
		res += fmt.Sprintf("<code>💬  %s</code> \n\n", escapeCharacters(ra.Comment))
	}
	res += fmt.Sprintf("See the changes <a href='%s'><b>here</b></a>",
		ra.getURLDiff(),
	)
	return res
}

func (ra *RecentAction) getURLDiff() string {
	params := url.Values{}
	params.Add("title", replaceSpaceAndEncode(ra.Title))
	params.Add("curid", fmt.Sprint(ra.PageID))
	params.Add("diff", fmt.Sprint(ra.RevID))
	params.Add("oldid", fmt.Sprint(ra.OldRevID))
	return fmt.Sprintf("%s/index.php?", ra.WikiURL) + params.Encode()
}

func (ra *RecentAction) getURLTitle() string {
	return fmt.Sprintf("%s/index.php/%s", ra.WikiURL, replaceSpaceAndEncode(ra.Title))
}

func (ra *RecentAction) getURLUser() string {
	return fmt.Sprintf("%s/index.php/User:%s", ra.WikiURL, replaceSpaceAndEncode(ra.User))
}

func replaceSpaceAndEncode(str string) string {
	return url.QueryEscape(strings.ReplaceAll(str, " ", "_"))
}

func addPlusIfNonNegative(value int) string {
	if value >= 0 {
		return "+"
	}
	return ""
}

func escapeCharacters(str string) string {
	str = strings.ReplaceAll(str, "<", "&lt;")
	str = strings.ReplaceAll(str, ">", "&gt;")
	return str
}
