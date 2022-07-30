package requests

import (
	"fmt"
	"log"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

const ROBOTS_FILENAME = "robots.txt"

type Rule struct {
	pattern    string
	allow      bool
	lineNumber int
}

type Robot struct {
	Rules         map[string][]Rule
	Sitemaps      []string
	PreferredHost string
	Url           string
	CrawlDelay    map[string]int
}

func (r *Robot) addRule(userAgents []string, pattern string, allow bool, lineNumber int) {
	var rules = r.Rules

	for _, userAgent := range userAgents {
		if _, ok := rules[userAgent]; !ok {
			rules[userAgent] = []Rule{}
		}

		if pattern == "" {
			return
		}

		rules[userAgent] = append(rules[userAgent], Rule{
			pattern:    pattern,
			allow:      allow,
			lineNumber: lineNumber,
		})
	}
}

func (r *Robot) setCrawlDelay(userAgents []string, delayStr string) {
	var rules = r.CrawlDelay
	var delay, err = strconv.Atoi(delayStr)
	if err != nil {
		return
	}
	for _, userAgent := range userAgents {
		rules[userAgent] = delay
	}
}

func (r *Robot) addSitemap(url string) {
	r.Sitemaps = append(r.Sitemaps, url)
}

// Normalises the user-agent string by converting it to
// lower case and removing any version numbers.
func formatUserAgent(userAgent string) string {
	var formattedUserAgent string = strings.ToLower(userAgent)

	// Strip the version number from robot/1.0 user agents

	var idx int = strings.Index(formattedUserAgent, "/")
	if idx > -1 {
		formattedUserAgent = formattedUserAgent[0:idx]
	}

	return strings.Trim(formattedUserAgent, " \n")

}

// Remove comments from lines
func removeComments(lines []string) []string {
	var result = []string{}
	for _, line := range lines {
		var commentStartIndex int = strings.Index(line, "#")
		if commentStartIndex > -1 {
			result = append(result, line[0:commentStartIndex])
		} else {
			result = append(result, line)
		}
	}
	return result
}

// Splits a line at the first occurrence of :
func splitLines(lines []string) [][]string {
	fmt.Printf("Number of lines %v\n", len(lines))
	var result [][]string = [][]string{}
	for _, line := range lines {
		var idx int = strings.Index(line, ":")

		if idx < 0 {
			result = append(result, []string{})
		} else {
			result = append(result, []string{strings.Trim(line[0:idx], " "), strings.Trim(line[idx+1:], " ")})
		}
	}
	return result
}

func ParseRobots(contents string, urlStr string) Robot {
	var newlineRegex = "\n"

	var robots = Robot{
		Rules:         make(map[string][]Rule),
		Sitemaps:      make([]string, 0),
		PreferredHost: "",
		Url:           urlStr,
		CrawlDelay:    make(map[string]int),
	}
	var lines []string = strings.Split(contents, newlineRegex)
	lines = removeComments(lines)
	var currentUserAgents []string = []string{}
	var isNoneUserAgentState = true

	for i, line := range splitLines(lines) {
		if len(line) == 0 {
			continue
		}

		switch strings.ToLower(line[0]) {
		case "user-agent":
			if isNoneUserAgentState {
				log.Printf("resetting current user agents\n")
				currentUserAgents = []string{}
			}

			if len(line) > 1 {
				var agent string = formatUserAgent(line[1])
				log.Printf("user agent added to current user agents %s\n", agent)
				currentUserAgents = append(currentUserAgents, agent)
			}
		case "disallow":
			robots.addRule(currentUserAgents, line[1], false, i+1)
		case "allow":
			robots.addRule(currentUserAgents, line[1], true, i+1)
		case "crawl-delay":
			robots.setCrawlDelay(currentUserAgents, line[1])
		case "sitemap":
			if len(line) > 1 {
				robots.addSitemap(line[1])
			}
		case "host":
			if len(line) > 1 {
				robots.PreferredHost = strings.ToLower(line[1])
			}
		}
		isNoneUserAgentState = strings.ToLower(line[0]) != "user-agent"
	}
	return robots
}

// Matches a pattern with the specified path
// Uses same algorithm to match patterns as the Google implementation in
// google/robotstxt so it should be consistent with the spec.
func matches(pattern string, path string) bool {
	// I've added extra comments to try make this easier to understand

	// Stores the lengths of all the current matching substrings.
	// Maximum number of possible matching lengths is every length in path plus
	// 1 to handle 0 length too (if pattern starts with * which is zero or more)
	var matchingLengths []int = make([]int, len(path)+1)
	var numMatchingLengths int = 1

	// Initially longest match is 0
	matchingLengths[0] = 0

	for p := range pattern {
		// If $ is at the end of pattern then we must match the whole path.
		// Which is true if the longest matching length matches path length
		if pattern[p] == '$' && p+1 == len(pattern) {
			return matchingLengths[numMatchingLengths-1] == len(path)
		}

		// Handle wildcards
		if pattern[p] == '*' {
			// Wildcard so all substrings minus the current smallest matching
			// length are matches
			numMatchingLengths = len(path) - matchingLengths[0] + 1

			// Update matching lengths to include the smallest all the way up
			// to numMatchingLengths
			// Don't update smallest possible match as * matches zero or more
			// so the smallest current match is also valid
			for i := 1; i < numMatchingLengths; i++ {
				matchingLengths[i] = matchingLengths[i-1] + 1
			}
		} else {
			// Check the char at the matching length matches the pattern, if it
			// does increment it and add it as a valid length, ignore if not.
			var numMatches int = 0
			for i := 0; i < numMatchingLengths; i++ {
				if matchingLengths[i] < len(path) && path[matchingLengths[i]] == pattern[p] {
					matchingLengths[numMatches] = matchingLengths[i] + 1
					numMatches++
				}
			}

			// No paths matched the current pattern char so not a match
			if numMatches == 0 {
				return false
			}

			numMatchingLengths = numMatches
		}
	}

	return true
}

// Returns if a pattern is allowed by the specified rules.
func findRule(path string, rules []Rule) *Rule {
	var matchedRule *Rule = nil
	fmt.Println(len(rules))
	fmt.Println(path)
	for i := 0; i < len(rules); i++ {
		var rule = rules[i]
		fmt.Println(rule.pattern)
		match, _ := regexp.MatchString(rule.pattern, path)
		fmt.Println(match)
		if !match {
			continue
		}

		// The longest matching rule takes precedence
		// If rules are the same length then allow takes precedence
		if matchedRule == nil || len(rule.pattern) > len(matchedRule.pattern) {
			matchedRule = &rule
		} else if len(rule.pattern) == len(matchedRule.pattern) && rule.allow && !matchedRule.allow {
			matchedRule = &rule
		}
	}
	fmt.Println(matchedRule)
	return matchedRule
}

// Get rule that would apply to the given url and user agent
func (r Robot) GetRule(urlStr string, ua string) (*Rule, error) {
	parsedUrl, err := url.Parse(urlStr)
	if ua == "" {
		ua = "*"
	}
	if err != nil {
		return nil, err
	}
	var userAgent string = formatUserAgent(ua)
	robotUrl, err2 := url.Parse(r.Url)

	if err2 != nil {
		return nil, err2
	}

	// The base URL must match otherwise this robots.txt is not valid for it.
	if parsedUrl.Scheme != robotUrl.Scheme || parsedUrl.Host != robotUrl.Host {
		return nil, fmt.Errorf("schemes don't match %s %s", parsedUrl.Scheme, robotUrl.Scheme)
	}

	var rules, ok = r.Rules[userAgent]
	if !ok {
		rules, ok = r.Rules["*"]
		if !ok {
			rules = []Rule{}
		}
	}

	var path = parsedUrl.Path
	if parsedUrl.ForceQuery {
		path = path + "?"
	}

	path = path + parsedUrl.RawQuery
	var rule = findRule(strings.ToLower(path), rules)

	return rule, nil
}

func (r Robot) IsAllowed(url string, ua string) (bool, error) {
	var rule, err = r.GetRule(url, ua)

	if rule == nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return (*rule).allow, nil
}

func (r Robot) IsDisallowed(url string, ua string) (bool, error) {
	var rule, err = r.GetRule(url, ua)
	if rule == nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return !rule.allow, nil
}
