package main

import (
	"fmt"
	"path"
	"regexp"
	"strings"
)

func filenameWithoutExtension(fn string) string {
	return strings.TrimSuffix(fn, path.Ext(fn))
}

func prepareBookURLHTML(url string, urlname string) string {
	html := `<a href="[url]">[urlname]</a>`
	html = strings.Replace(html, "[url]", url, -1)
	html = strings.Replace(html, "[urlname]", urlname, -1)
	return html
}

func prepareBookHTML(image string) string {
	html := `
	<p><label for="[image]" class="margin-toggle">&#8853;</label>
	<input type="checkbox" id="[image]" class="margin-toggle"/>
	<span class="marginnote"><img src="/static/books/[image]" width="140" alt=""/>
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

	return content
}

func prepareFactHTML(content string, id string, sources []Source) string {

	content = prepareTextStyle(content, false)

	pre := ` <label for="%s" class="margin-toggle sidenote-number"></label>
<input type="checkbox" id="%s" class="margin-toggle"/>
<span class="sidenote"> `
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
			value += fmt.Sprintf(" <a href=\"%s\">%s</a> ", item.URL, nameURL)
		}
		content = strings.Replace(content, "["+item.ID+"]", preItem+value+post, -1)
	}

	return content
}
