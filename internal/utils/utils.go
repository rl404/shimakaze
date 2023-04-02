package utils

import (
	"errors"
	"html"
	"math"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

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

// PtrToBool to convert pointer to bool.
func PtrToBool(b *bool) bool {
	if b != nil {
		return *b
	}
	return false
}

// StrToPtrBool to convert string to bool pointer.
func StrToPtrBool(str string) *bool {
	if str == "" {
		return nil
	}
	b, _ := strconv.ParseBool(str)
	return &b
}

// StrToStrSlice to convert string to slice of string.
func StrToStrSlice(str string) []string {
	if str == "" {
		return nil
	}
	return strings.Split(str, ",")
}

// GetLastPathFromURL to get the last path from url.
func GetLastPathFromURL(str string) string {
	url, err := url.Parse(str)
	if err != nil {
		return ""
	}
	if url.Path[len(url.Path)-1] == '/' {
		url.Path = url.Path[:len(url.Path)-1]
	}
	splitPath := strings.Split(url.Path, "/")
	return splitPath[len(splitPath)-1]
}

const (
	parsingPeriod = iota
	parsingTime
)

// ParseDuration to parse ISO 8601 duration.
func ParseDuration(durationStr string, startFromTime ...bool) (time.Duration, error) {
	state := parsingPeriod
	if len(startFromTime) > 0 && startFromTime[0] {
		state = parsingTime
	}

	duration := 0 * time.Second
	num := ""

	err := errors.New("invalid duration")

	for _, c := range durationStr {
		switch c {
		case 'P', 'p':
			state = parsingPeriod
		case 'T', 't':
			state = parsingTime
		case 'Y', 'y':
			if state != parsingPeriod {
				return 0, err
			}

			y, err := strconv.ParseFloat(num, 64)
			if err != nil {
				return 0, err
			}

			duration += time.Duration(math.Round(y)) * 365 * 24 * time.Hour

			num = ""
		case 'M', 'm':
			if state == parsingPeriod {
				m, err := strconv.ParseFloat(num, 64)
				if err != nil {
					return 0, err
				}

				duration += time.Duration(math.Round(m)) * 30 * 24 * time.Hour

				num = ""
			} else if state == parsingTime {
				m, err := strconv.ParseFloat(num, 64)
				if err != nil {
					return 0, err
				}

				duration += time.Duration(math.Round(m)) * time.Minute

				num = ""
			}
		case 'W', 'w':
			if state != parsingPeriod {
				return 0, err
			}

			w, err := strconv.ParseFloat(num, 64)
			if err != nil {
				return 0, err
			}

			duration += time.Duration(math.Round(w)) * 7 * 24 * time.Hour

			num = ""
		case 'D', 'd':
			if state != parsingPeriod {
				return 0, err
			}

			d, err := strconv.ParseFloat(num, 64)
			if err != nil {
				return 0, err
			}

			duration += time.Duration(math.Round(d)) * 24 * time.Hour

			num = ""
		case 'H', 'h':
			if state != parsingTime {
				return 0, err
			}

			h, err := strconv.ParseFloat(num, 64)
			if err != nil {
				return 0, err
			}

			duration += time.Duration(math.Round(h)) * time.Hour

			num = ""
		case 'S', 's':
			if state != parsingTime {
				return 0, err
			}

			s, err := strconv.ParseFloat(num, 64)
			if err != nil {
				return 0, err
			}

			duration += time.Duration(math.Round(s)) * time.Second

			num = ""
		default:
			if unicode.IsNumber(c) || c == '.' {
				num += string(c)
				continue
			}

			return 0, err
		}
	}

	return duration, nil
}
