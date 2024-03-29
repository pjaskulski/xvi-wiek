package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path"
	"regexp"
	"strings"
	"time"
)

// funkcja zwraca nazwę pliku bez rozszerzenia
func filenameWithoutExtension(fn string) string {
	return strings.TrimSuffix(fn, path.Ext(fn))
}

func prepareGeoHTML(geo string) string {
	html := `<span><a href="[geo]" target="_blank" rel="noopener" class="no-tufte-underline">
	<img src="/static/img/world.png" class="small-icon" alt="Położenie geograficzne na mapie"/></a>
	</span>`
	pos := strings.Split(geo, ",")
	if len(pos) == 2 {
		url := strings.Replace(`https://www.openstreetmap.org/?mlat=[lat]&mlon=[lon]&zoom=9`, "[lat]", pos[0], 1)
		url = strings.Replace(url, "[lon]", pos[1], 1)
		html = strings.Replace(html, "[geo]", url, 1)
	}
	return html
}

func prepareBookURLHTML(url string, urlname string) string {
	html := `<a href="[url]" target="_blank" rel="noopener">[urlname]</a>`
	html = strings.Replace(html, "[url]", url, -1)
	html = strings.Replace(html, "[urlname]", urlname, -1)
	return html
}

func prepareBookHTML(image string) string {
	html := `
	<p><label for="[image]" class="margin-toggle">&#8853;</label>
	<input type="checkbox" id="[image]" class="margin-toggle"/>
	<span class="marginnote"><img src="/static/books/[image]" width="160" alt=""/>
	</p>`
	html = strings.Replace(html, "[image]", image, -1)
	return html
}

func prepareImageHTML(image string, imageInfo string) string {
	html := `
	<figure>
	  <label for="fact-image" class="margin-toggle">&#8853;</label>
	  <input type="checkbox" id="fact-image" class="margin-toggle"/>
	  <span class="marginnote">[imageInfo]</span>
	  <img src="/static/gallery/[image]" alt="[imageInfo]" />
	</figure>`

	html = strings.Replace(html, "[image]", image, -1)
	html = strings.Replace(html, "[imageInfo]", imageInfo, -1)

	return html
}

func prepareTextStyle(content string, clear bool) string {
	// kapitaliki
	var rgx = regexp.MustCompile(`\{\{\{(.*?)\}\}\}`)
	var pre, post string
	if !clear {
		pre = "<span class=\"newthought\">"
		post = "</span>"
	}

	textToSmallCaps := rgx.FindAllString(content, -1)

	for _, item := range textToSmallCaps {
		textHTML := pre + item[3:len(item)-3] + post
		content = strings.Replace(content, item, textHTML, -1)
	}

	// pogrubienie
	var rgxb = regexp.MustCompile(`\{\{(.*?)\}\}`)
	if !clear {
		pre = "<strong>"
		post = "</strong>"
	}

	textBold := rgxb.FindAllString(content, -1)

	for _, item := range textBold {
		textHTML := pre + item[2:len(item)-2] + post
		content = strings.Replace(content, item, textHTML, -1)
	}

	// italiki
	var rgxi = regexp.MustCompile(`\{(.*?)\}`)
	if !clear {
		pre = "<em>"
		post = "</em>"
	}

	textItalic := rgxi.FindAllString(content, -1)

	for _, item := range textItalic {
		textHTML := pre + item[1:len(item)-1] + post
		content = strings.Replace(content, item, textHTML, -1)
	}

	// złamanie wiersza
	content = strings.Replace(content, `\\`, `</p><p>`, -1)

	return content
}

func prepareFactHTML(content string, id string, sources []Source) string {

	content = prepareTextStyle(content, false)

	pre := `<label for="%s" class="margin-toggle sidenote-number"></label>
<input type="checkbox" id="%s" class="margin-toggle"/>
<span class="sidenote">`
	post := `</span>`

	var tmpValue string

	for _, item := range sources {
		idQuote := fmt.Sprintf("%s-%s", id, item.ID)
		preItem := fmt.Sprintf(pre, idQuote, idQuote)

		if item.Type == "reference" {
			newValue, found := ReferenceMap[item.Value]
			if found {
				item.Value = newValue
			}
		}

		if item.Page != "" {
			tmpValue = item.Value + ", " + item.Page
		} else {
			tmpValue = item.Value
		}
		value := prepareTextStyle(tmpValue, false)
		if item.URL != "" {
			var nameURL string
			if item.URLName != "" {
				nameURL = item.URLName
			} else {
				nameURL = item.URL
			}
			value += fmt.Sprintf(" <a href=\"%s\" target=\"_blank\" rel=\"noopener\">%s</a> ", item.URL, nameURL)
		}
		content = strings.Replace(content, "["+item.ID+"]", preItem+value+post, -1)
	}

	return content
}

// randomInt - funkcja zwraca losową liczbę całkowitą z podanego zakresu
func randomInt(min int, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min+1) + min
}

// isRunByRun - funkcja sprawdza czy uruchomiono program przez go run
// czy też program skompilowany, funkcja dla systemu Linux
func isRunByRun() bool {
	return strings.Contains(os.Args[0], "/tmp/go-build")
}

// getPrevNextHTML func
func getPrevNextHTML(month int, day int) string {
	var prevnext string

	daysInMonth := []int{0, 31, 29, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}

	tmp := `<a href="/dzien/%d/%d" class="no-tufte-underline">&#8678;</a> <a href="/dzien/%d/%d">Poprzedni</a> : Dzień : <a href="/dzien/%d/%d">Następny</a> <a href="/dzien/%d/%d" class="no-tufte-underline">&#8680;</a>`
	tmpf := `<span style="color: gray;">Poprzedni</span> : Dzień : <a href="/dzien/%d/%d">Następny</a> <a href="/dzien/%d/%d" class="no-tufte-underline">&#8680;</a>`
	tmpl := `<a href="/dzien/%d/%d" class="no-tufte-underline">&#8678;</a> <a href="/dzien/%d/%d">Poprzedni</a> : Dzień : <span style="color: gray;">Następny</span>`

	if day == 1 && month == 1 {
		prevnext = fmt.Sprintf(tmpf, 1, 2, 1, 2)
	} else if day == 31 && month == 12 {
		prevnext = fmt.Sprintf(tmpl, 12, 30, 12, 30)
	} else if day == 1 {
		prevnext = fmt.Sprintf(tmp, month-1, daysInMonth[month-1], month-1, daysInMonth[month-1], month, day+1, month, day+1)
	} else if day == daysInMonth[month] {
		prevnext = fmt.Sprintf(tmp, month, day-1, month, day-1, month+1, 1, month+1, 1)
	} else {
		prevnext = fmt.Sprintf(tmp, month, day-1, month, day-1, month, day+1, month, day+1)
	}

	return prevnext
}

// prepareFactLinkHTML func
func prepareFactLinkHTML(month int, day int, id string) string {
	html := `<a href="/dzien/%d/%d#%s" class="no-tufte-underline"> &#8680; </a>`
	return fmt.Sprintf(html, month, day, id)
}

// IsValidXML - func from: https://stackoverflow.com/questions/53476012/how-to-validate-a-xml
func IsValidXML(input string) bool {
	decoder := xml.NewDecoder(strings.NewReader(input))
	for {
		err := decoder.Decode(new(interface{}))
		if err != nil {
			return err == io.EOF
		}
	}
}

// inSlice - function checks if the specified string is present in the specified slice
func inSlice(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

// inSliceInt - function checks if the specified integer is present in the specified slice
func inSliceInt(slice []int, val int) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

// inSliceInt - function checks if the specified KeywordFact struct (Title field must be unique)
// is present in the specified slice
func inSliceKeywordFact(slice []KeywordFact, val KeywordFact) bool {
	for _, item := range slice {
		if item.Title == val.Title {
			return true
		}
	}
	return false
}

// prepareNavigationPeopleIndexHTML func
func prepareNavigationIndexHTML(indexLetter []string) string {

	var html string
	var counter int

	for _, item := range indexLetter {
		counter += 1
		if counter < len(indexLetter) {
			html += fmt.Sprintf(`<a href="#%s">%s</a>&#10625;`, item, strings.ToUpper(item)) + "\n"
		} else {
			html += fmt.Sprintf(`<a href="#%s">%s</a>`, item, strings.ToUpper(item)) + "\n"
		}
	}

	return html
}
