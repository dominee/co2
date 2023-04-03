package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"

	color "github.com/fatih/color"
	gojsonq "github.com/thedevsaddam/gojsonq/v2"
)

func main() {
	// Parse command line arguments
	showIDs := flag.Bool("i", false, "Display lines containing 'IDs:'")
	flag.Parse()

	scanner := bufio.NewScanner(os.Stdin)
	ipRegex := regexp.MustCompile(`\b\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}\b`)
	tokenRegex := regexp.MustCompile(`\b[0-9a-z]{33}`)
	targetRegex := regexp.MustCompile(`\[\S+\]`)

	for scanner.Scan() {
		line := scanner.Text()

		// Check if the -i flag is set and the line contains "IDs:"
		if *showIDs && !strings.Contains(line, "IDs:") {
			continue
		}

		// Find all IP addresses in the line
		ips := ipRegex.FindAllString(line, -1)

		// Find the query target - filename/hostname
		var target string
		targets := targetRegex.FindAllString(line, -1)
		if len(targets) > 0 {
			// use the last one (2nd)
			target = targets[len(targets)-1]
			// without trailing brackets
			target = target[1 : len(target)-1]
		}

		// Find the token at the end of the line (if any)
		var token string
		tokens := tokenRegex.FindAllString(line, -1)
		if len(tokens) > 0 {
			// adding leading space to match only the one at the end of line
			token = " " + tokens[len(tokens)-1]
		}

		// Colorize the output
		line = colorize(line, ips, target, token)

		fmt.Println(line)
	}
}

// Retrieves the name of the owner of an IP address using the Whois API
func getOwner(ip string) string {

	// json.objects.object[0]["resource-holder"].key = "ORG-Ds65-RIPE";
	// json.objects.object[0]["resource-holder"].name = "Digmia s.r.o.";

	url := fmt.Sprintf("http://rest.db.ripe.net/search.json?query-string=%s&resource-holder=true&type-filter=inetnum", ip)

	// TODO: timeout
	resp, err := http.Get(url)
	if err != nil {
		return "Unknown"
	}
	defer resp.Body.Close()

	result, err := gojsonq.New().Reader(resp.Body).From("objects").FindR("object.[0].resource-holder.name")
	if err != nil {
		// TODO: try to get netname as fallback
		return "Unknown"
	}
	name, _ := result.String()

	return name
}

// Colorizes the output
func colorize(line string, ips []string, target string, token string) string {
	// Define color attributes
	color.NoColor = false
	red := color.New(color.FgRed).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	blue := color.New(color.FgBlue).SprintFunc()

	// Colorize the date and time
	parts := strings.Split(line, " ")
	if len(parts) < 3 {
		return line
	}
	date := parts[0]
	time := parts[1]
	line = strings.Replace(line, date, blue(date), 1)
	line = strings.Replace(line, time, blue(time), 1)

	// Colorize the IP addresses
	for _, ip := range ips {
		// TODO CACHE
		owner := getOwner(ip)
		coloredIP := green(ip) + yellow(fmt.Sprintf(" (%s)", owner))
		line = strings.Replace(line, ip, coloredIP, 1)
	}

	// Colorize the target - file/hostname
	if target != "" {
		line = strings.Replace(line, target, green(target), 2)
	}

	// Colorize the token (if any)
	if token != "" {
		line = strings.Replace(line, token, red(token), 2)
	}

	// Colorize the Interaction type
	if len(parts) >= 4 {
		hostname := parts[4]
		line = strings.Replace(line, hostname, green(hostname), 1)

	}

	return line
}
