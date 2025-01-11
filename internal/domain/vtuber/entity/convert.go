package entity

import (
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/rl404/shimakaze/internal/domain/wikia/entity"
	"github.com/rl404/shimakaze/internal/utils"
	emoji "github.com/tmdvs/Go-Emoji-Utils"
)

// WikiaPageToVtuber to convert wikia page to vtuber.
func WikiaPageToVtuber(page entity.Page) Vtuber {
	data := page.Content

	// Clean up content.
	data = utils.CleanWikiaTag(data, "ref", true)
	data = utils.CleanWikiaComment(data)
	data = utils.NormalizeWikiaInternalLink(data)
	data = utils.NormalizeNewLine(data)

	// Parse data.
	var vtuber Vtuber
	vtuber.ID, vtuber.Name = page.ID, page.Title
	vtuber.OriginalNames, _ = parseOriginalNames(data)
	vtuber.Nicknames, _ = parseNickNames(data)
	vtuber.Caption, _ = parseCaption(data)
	vtuber.DebutDate, _ = parseDate("debut_date", data)
	vtuber.RetirementDate, _ = parseDate("retirement_date", data)
	vtuber.Affiliations, _ = parseAffiliation(data)
	vtuber.Channels, _ = parseChannels(data)
	vtuber.SocialMedias, _ = parseSocialMedias(data)
	vtuber.OfficialWebsites, _ = parseOfficialWebsites(data)
	vtuber.Gender, _ = parseGender(data)
	vtuber.Age, _ = parseDecimal("age", data)
	vtuber.Birthday, _ = parseDate("birthday", data)
	vtuber.Height, _ = parseDecimal("height", data)
	vtuber.Weight, _ = parseDecimal("weight", data)
	vtuber.BloodType, _ = parseBloodType(data)
	vtuber.ZodiacSign, _ = parseZodiacSign(data)
	vtuber.Emoji, _ = parseEmoji(data)
	return vtuber
}

func parseData(key, data string) (string, string) {
	dataRegex := regexp.MustCompile(`\|` + key + `\s*=\s*(.*?)(<br>)?(\||}})`)

	raw := dataRegex.FindString(data)
	if raw == "" {
		return "", ""
	}

	submatch := dataRegex.FindStringSubmatch(data)

	if len(submatch) < 2 {
		return "", raw
	}

	return strings.TrimSpace(submatch[1]), raw
}

func parseOriginalNames(data string) ([]string, string) {
	value, raw := parseData("original_name", data)

	var names []string
	for _, n := range strings.Split(value, "<br>") {
		n = strings.TrimSpace(n)
		n = utils.RemoveAllHTMLTag(n)
		if n != "" {
			names = append(names, n)
		}
	}

	return names, raw
}

func parseNickNames(data string) ([]string, string) {
	value, raw := parseData("nick_name", data)

	var names []string
	for _, n := range strings.Split(value, "<br>") {
		n = strings.TrimSpace(n)
		n = utils.WikiaInternalLinkToStr(n)
		n = utils.WikiaExternalLinkToStr(n)
		n = utils.RemoveAllHTMLTag(n)
		if n != "" {
			names = append(names, n)
		}
	}

	return names, raw
}

func parseCaption(data string) (string, string) {
	value, raw := parseData("caption1", data)

	caption := utils.WikiaExternalLinkToStr(value)
	caption = utils.RemoveAllHTMLTag(caption)

	return caption, raw
}

func parseDate(key, data string) (*time.Time, string) {
	value, raw := parseData(key, data)

	date1Str := regexp.MustCompile(`\d{4}\/\d{1,2}\/\d{1,2}`).FindString(value)
	if date1Str != "" {
		if date, err := time.Parse("2006/01/_2", date1Str); err == nil {
			return &date, raw
		}

		if date, err := time.Parse("2006/1/_2", date1Str); err == nil {
			return &date, raw
		}

		if date, err := time.Parse("2006/_2/01", date1Str); err == nil {
			return &date, raw
		}
	}

	date2Str := regexp.MustCompile(`\d{1,2}.*\s\w+\s\d{4}`).FindString(value)
	if date2Str != "" {
		date2Split := strings.Split(date2Str, " ")
		date2Split[0] = regexp.MustCompile(`[^\d]+`).ReplaceAllString(date2Split[0], "")
		date2Str = strings.Join(date2Split, " ")

		if date, err := time.Parse("_2 January 2006", date2Str); err == nil {
			return &date, raw
		}
	}

	date10Str := regexp.MustCompile(`\d{1,2}.*\s\w+`).FindString(value)
	if date10Str != "" {
		date10Split := strings.Split(date10Str, " ")
		date10Split[0] = regexp.MustCompile(`[^\d]+`).ReplaceAllString(date10Split[0], "")
		date10Str = strings.Join(date10Split, " ")

		if date, err := time.Parse("_2 January", date10Str); err == nil {
			return &date, raw
		}
	}

	date4Str := regexp.MustCompile(`\d{1,2}\/\d{1,2}\/\d{4}`).FindString(value)
	if date4Str != "" {
		if date, err := time.Parse("_2/01/2006", date4Str); err == nil {
			return &date, raw
		}

		if date, err := time.Parse("_2/1/2006", date4Str); err == nil {
			return &date, raw
		}

		if date, err := time.Parse("01/_2/2006", date4Str); err == nil {
			return &date, raw
		}

		if date, err := time.Parse("1/_2/2006", date4Str); err == nil {
			return &date, raw
		}
	}

	date6Str := regexp.MustCompile(`\d{4}\/\d{2}`).FindString(value)
	if date6Str != "" {
		if date, err := time.Parse("2006/01", date6Str); err == nil {
			return &date, raw
		}
	}

	date8Str := regexp.MustCompile(`\d{2}\/\w{3}\/\d{4}`).FindString(value)
	if date8Str != "" {
		if date, err := time.Parse("02/Jan/2006", date8Str); err == nil {
			return &date, raw
		}
	}

	date5Str := regexp.MustCompile(`[^=\s]+\s\d{1,2}.+\d{4}`).FindString(value)
	if date5Str != "" {
		date5Split := strings.Split(date5Str, " ")
		date5Split[1] = regexp.MustCompile(`[^\d]+`).ReplaceAllString(date5Split[1], "")
		date5Str = strings.Join(date5Split, " ")

		if date, err := time.Parse("January _2 2006", date5Str); err == nil {
			return &date, raw
		}
	}

	date9Str := regexp.MustCompile(`[^=\s]+\s\d{4}`).FindString(value)
	if date9Str != "" {
		if date, err := time.Parse("January 2006", date9Str); err == nil {
			return &date, raw
		}
	}

	date11Str := regexp.MustCompile(`[^=\s]+\s\d{1,2}`).FindString(value)
	if date11Str != "" {
		if date, err := time.Parse("January _2", date11Str); err == nil {
			return &date, raw
		}
	}

	date7Str := regexp.MustCompile(`\d{4}`).FindString(value)
	if date7Str != "" {
		if date, err := time.Parse("2006", date7Str); err == nil {
			return &date, raw
		}
	}

	return nil, raw
}

func parseAffiliation(data string) ([]string, string) {
	value, raw := parseData("affiliation", data)

	var names []string
	for _, n := range strings.Split(value, "<br>") {
		n = strings.TrimSpace(n)
		n = utils.WikiaInternalLinkToStr(n)
		n = utils.WikiaExternalLinkToStr(n)
		n = utils.RemoveAllHTMLTag(n)
		if n != "" {
			names = append(names, n)
		}
	}

	return names, raw
}

func parseChannels(data string) ([]Channel, string) {
	value, raw := parseData("channel", data)

	var channels []Channel
	for _, n := range strings.Split(value, "<br>") {
		link := utils.GetWikiaExternalLink(n)
		if link == "" {
			continue
		}

		channels = append(channels, Channel{
			Type: ParseChannelType(link),
			URL:  link,
		})
	}

	return channels, raw
}

// ParseChannelType to parse url to channel type.
func ParseChannelType(link string) ChannelType {
	u, err := url.Parse(link)
	if err != nil {
		return ChannelOther
	}

	parts := strings.Split(u.Hostname(), ".")
	if len(parts) < 2 {
		return ChannelOther
	}

	switch strings.ToLower(parts[len(parts)-2]) {
	case "youtube":
		return ChannelYoutube
	case "twitch":
		return ChannelTwitch
	case "bilibili":
		return ChannelBilibili
	case "nicovideo":
		if strings.ToLower(u.Hostname()) == "www.nicovideo.jp" {
			return ChannelNiconico
		}
		return ChannelOther
	default:
		return ChannelOther
	}
}

func parseSocialMedias(data string) ([]string, string) {
	value, raw := parseData("social_media", data)

	var links []string
	for _, n := range strings.Split(value, "<br>") {
		link := utils.GetWikiaExternalLink(n)
		if link == "" {
			continue
		}

		links = append(links, link)
	}

	return links, raw
}

func parseOfficialWebsites(data string) ([]string, string) {
	value, raw := parseData("official_website", data)

	var links []string
	for _, n := range strings.Split(value, "<br>") {
		link := utils.GetWikiaExternalLink(n)
		if link == "" {
			continue
		}

		links = append(links, link)
	}

	return links, raw
}

func parseGender(data string) (string, string) {
	value, raw := parseData("gender", data)

	value = strings.ReplaceAll(value, "<br>", " ")
	value = regexp.MustCompile(`\s+`).ReplaceAllString(value, " ")
	value = utils.WikiaInternalLinkToStr(value)
	value = utils.RemoveAllHTMLTag(value)

	return value, raw
}

var uncountableNumber float64 = -1
var invalidNumber float64 = -2

func parseDecimal(key, data string) (*float64, string) {
	value, raw := parseData(key, data)

	if value == "" {
		return nil, raw
	}

	numStr := regexp.MustCompile(`(\d+\,)*\d+(\.\d+)?`).FindString(value)
	if numStr == "" {
		return &uncountableNumber, raw
	}

	numStr = strings.ReplaceAll(numStr, ",", "")

	num, err := strconv.ParseFloat(numStr, 64)
	if err != nil {
		return &invalidNumber, raw
	}

	return &num, raw
}

func parseBloodType(data string) (string, string) {
	value, raw := parseData("blood_type", data)

	value = strings.ReplaceAll(value, "<br>", ", ")

	return value, raw
}

func parseZodiacSign(data string) (string, string) {
	zodiacRegex := regexp.MustCompile(`{{Zodiac\|(.+)}}`)
	data = zodiacRegex.ReplaceAllString(data, "$1")

	value, raw := parseData("zodiac_sign", data)

	value = strings.ReplaceAll(value, "<br>", ", ")

	if zodiac := toZodiac(value); zodiac != "" {
		return zodiac, raw
	}

	return value, raw
}

func toZodiac(value string) string {
	layouts := []string{
		"2 January",
		"January 2",
	}

	var date time.Time
	var err error

	for _, layout := range layouts {
		date, err = time.Parse(layout, value)
		if err == nil {
			break
		}
	}

	if err != nil {
		return ""
	}

	month := date.Month()
	day := date.Day()

	switch {
	case (month == time.March && day >= 21) || (month == time.April && day <= 19):
		return "Aries"
	case (month == time.April && day >= 20) || (month == time.May && day <= 20):
		return "Taurus"
	case (month == time.May && day >= 21) || (month == time.June && day <= 20):
		return "Gemini"
	case (month == time.June && day >= 21) || (month == time.July && day <= 22):
		return "Cancer"
	case (month == time.July && day >= 23) || (month == time.August && day <= 22):
		return "Leo"
	case (month == time.August && day >= 23) || (month == time.September && day <= 22):
		return "Virgo"
	case (month == time.September && day >= 23) || (month == time.October && day <= 22):
		return "Libra"
	case (month == time.October && day >= 23) || (month == time.November && day <= 21):
		return "Scorpio"
	case (month == time.November && day >= 22) || (month == time.December && day <= 21):
		return "Sagittarius"
	case (month == time.December && day >= 22) || (month == time.January && day <= 19):
		return "Capricorn"
	case (month == time.January && day >= 20) || (month == time.February && day <= 18):
		return "Aquarius"
	case (month == time.February && day >= 19) || (month == time.March && day <= 20):
		return "Pisces"
	default:
		return ""
	}
}

func parseEmoji(data string) (string, string) {
	value, raw := parseData("emoji", data)

	emojis := emoji.FindAll(value)
	value = ""
	for _, e := range emojis {
		value += e.Match.Value
	}

	s := strings.Split(value, "")
	sort.Strings(s)

	return strings.Join(s, ""), raw
}

// StrsToChannelTypes to convert slice of string to slice of ChannelType.
func StrsToChannelTypes(strs []string) []ChannelType {
	ct := make([]ChannelType, len(strs))
	for i, str := range strs {
		ct[i] = ChannelType(str)
	}
	return ct
}
