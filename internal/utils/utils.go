package utils

import (
	"html"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// CleanWikiaComment to remove wikia comment.
func CleanWikiaComment(str string) string {
	return regexp.MustCompile(`<!--.+?-->`).ReplaceAllString(str, "")
}

// CleanWikiaTag to remove wikia tag.
func CleanWikiaTag(str, tag string, deleteContent bool) string {
	str = regexp.MustCompile(`<`+tag+` .+/>`).ReplaceAllString(str, "")
	if deleteContent {
		doc, _ := goquery.NewDocumentFromReader(strings.NewReader(str))
		doc.Find(tag).Remove()
		str, _ = doc.Html()
		return html.UnescapeString(str)
	}
	return regexp.MustCompile(`</*`+tag+`>`).ReplaceAllString(str, "")
}

// NormalizeNewLine to normalize new line to <br>.
func NormalizeNewLine(str string) string {
	return regexp.MustCompile(`(?i)((\n|\\n|<br\s*\/*>)+)`).ReplaceAllString(str, "<br>")
}

// WikiaExternalLinkToStr to convert wikia external link to string.
func WikiaExternalLinkToStr(str string) string {
	linkRegex := regexp.MustCompile(`\[.+?\s+(.+?)\]`)
	if linkRegex.FindString(str) == "" {
		return str
	}
	return linkRegex.ReplaceAllString(str, "$1")
}

// WikiaInternalLinkToStr to convert wikia internal link to string.
func WikiaInternalLinkToStr(str string) string {
	linkRegex := regexp.MustCompile(`\[\[(.+?)\]\]`)
	if linkRegex.FindString(str) == "" {
		return str
	}
	return linkRegex.ReplaceAllString(str, "$1")
}

// NormalizeWikiaInternalLink to normalize wikia internal link.
// [[NIJISANJI (main branch)|NIJISANJI]] => [[NIJISANJI]]
func NormalizeWikiaInternalLink(str string) string {
	linkRegex := regexp.MustCompile(`\[\[([^\|\]]+\|)*(.+?)(\]\])?\]\]`)
	return linkRegex.ReplaceAllString(str, "[[$2]]")
}

// GetWikiaExternalLink to get wikia external link.
func GetWikiaExternalLink(str string) string {
	linkRegex := regexp.MustCompile(`\[([^\]\s]+).*?\]?\]`)
	if linkRegex.FindString(str) == "" {
		return ""
	}

	submatch := linkRegex.FindStringSubmatch(str)

	if len(submatch) < 2 {
		return ""
	}

	return strings.TrimSpace(submatch[1])
}

// RemoveAllHTMLTag to remove all html tag.
func RemoveAllHTMLTag(str string) string {
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(str))
	return html.UnescapeString(doc.Text())
}
