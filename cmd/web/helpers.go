package main

import (
	"fmt"
	"math/rand"
	"os"
	"path"
	"regexp"
	"strings"
	"time"
)

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

	if textToSmallCaps != nil {
		for _, item := range textToSmallCaps {
			textHTML := pre + item[3:len(item)-3] + post
			content = strings.Replace(content, item, textHTML, -1)
		}
	}

	// pogrubienie
	var rgxb = regexp.MustCompile(`\{\{(.*?)\}\}`)
	if !clear {
		pre = "<strong>"
		post = "</strong>"
	}

	textBold := rgxb.FindAllString(content, -1)

	if textBold != nil {
		for _, item := range textBold {
			textHTML := pre + item[2:len(item)-2] + post
			content = strings.Replace(content, item, textHTML, -1)
		}
	}

	// italiki
	var rgxi = regexp.MustCompile(`\{(.*?)\}`)
	if !clear {
		pre = "<em>"
		post = "</em>"
	}

	textItalic := rgxi.FindAllString(content, -1)

	if textItalic != nil {
		for _, item := range textItalic {
			textHTML := pre + item[1:len(item)-1] + post
			content = strings.Replace(content, item, textHTML, -1)
		}
	}

	// złamanie wiersza
	content = strings.Replace(content, "\\", "<br>", -1)

	return content
}

func prepareFactHTML(content string, id string, sources []Source) string {

	content = prepareTextStyle(content, false)

	pre := `<label for="%s" class="margin-toggle sidenote-number"></label>
<input type="checkbox" id="%s" class="margin-toggle"/>
<span class="sidenote">`
	post := `</span>`

	for _, item := range sources {
		idQuote := fmt.Sprintf("%s-%s", id, item.ID)
		preItem := fmt.Sprintf(pre, idQuote, idQuote)
		value := prepareTextStyle(item.Value, false)
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
	if strings.Index(os.Args[0], "/tmp/go-build") != -1 {
		return true
	}
	return false
}

// getPrevNextHTML func
func getPrevNextHTML(month int, day int) string {
	var prevnext string

	daysInMonth := []int{0, 31, 29, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}

	tmp := `<a href="/dzien/%d/%d">Poprzedni</a> dzień : Dzień <a href="/dzien/%d/%d">następny</a>`
	tmpf := `<span style="color: gray;">Poprzedni dzień</span> : Dzień <a href="/dzien/%d/%d">następny</a>`
	tmpl := `<a href="/dzien/%d/%d">Poprzedni</a> dzień : <span style="color: gray;">Dzień następny</span>`

	if day == 1 && month == 1 {
		prevnext = fmt.Sprintf(tmpf, 1, 2)
	} else if day == 31 && month == 12 {
		prevnext = fmt.Sprintf(tmpl, 12, 30)
	} else if day == 1 {
		prevnext = fmt.Sprintf(tmp, month-1, daysInMonth[month-1], month, day+1)
	} else if day == daysInMonth[month] {
		prevnext = fmt.Sprintf(tmp, month, day-1, month+1, 1)
	} else {
		prevnext = fmt.Sprintf(tmp, month, day-1, month, day+1)
	}

	return prevnext
}

// prepareFactLinkHTML func
func prepareFactLinkHTML(month int, day int, id string) string {
	html := `<a href="/dzien/%d/%d#%s" class="no-tufte-underline"> &#8680; </a>`
	return fmt.Sprintf(html, month, day, id)
}
