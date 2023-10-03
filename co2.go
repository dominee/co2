package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"

	color "github.com/fatih/color"
)

// ipCache is map to store the runtime cache for ip owners to save API calls and identify the first occurance of IP in log
type ipcache map[string]string

var ipCache = make(ipcache)
var apiKey string = ""

// isCached checks if the IP is already cached
func isCached(ipCache ipcache, ip string) bool {
	if _, ok := ipCache[ip]; ok {
		return true
	} else {
		return false
	}
}

func cacheIt(ipCache ipcache, ip string, owner string) {
	ipCache[ip] = owner
}

func getCached(ipCache ipcache, ip string) string {
	return ipCache[ip]
}

func cacheSave(cache ipcache, filename string) {
	// Marshal it to JSON, indented with 4 spaces
	cachej, error := json.MarshalIndent(cache, "", "   ")
	if error != nil {
		fmt.Println(error)
	}

	// Write it to file
	error = os.WriteFile(filename, cachej, 0644)
	if error != nil {
		fmt.Println(error)
	}
}

func cacheLoad(cache ipcache, filename string) ipcache {
	cached, err := os.ReadFile(filename)
	if err != nil {
		fmt.Print(err)
	}
	err = json.Unmarshal(cached, &cache)
	if err != nil {
		fmt.Println("error:", err)
	}
	return cache
}

func main() {
	// Parse command line arguments
	showIDs := flag.Bool("i", false, "Display lines containing interacions only.")
	cacheFileArg := flag.String("c", "cache.json", "JSON cache file to use.")
	flag.Parse()

	cacheFile := *cacheFileArg

	// Restore cache from file if it exists
	if _, err := os.Stat(cacheFile); err == nil {
		_ = cacheLoad(ipCache, cacheFile)
	}

	apiKey = os.Getenv("IP2LOCATION_API_KEY")
	if apiKey == "" {
		fmt.Println("IP2LOCATION_API_KEY environment variable is not set.")
		os.Exit(1)
	}

	// Regexp for line parts
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

	// Save the cache before we finish
	cacheSave(ipCache, cacheFile)
}

func getOwner(ip string) (string, bool) {
	var owner string
	var cacheHit bool
	if isCached(ipCache, ip) {
		owner = getCached(ipCache, ip)
		cacheHit = true
	} else {
		owner = getNewOwner(ip)
		cacheIt(ipCache, ip, owner)
		cacheHit = false
	}
	return owner, cacheHit
}

// Retrieves the name of the owner of an IP address using the Whois API
func getNewOwner(ip string) string {

	type Location struct {
		IP          string  `json:"ip"`
		CountryCode string  `json:"country_code"`
		CountryName string  `json:"country_name"`
		RegionName  string  `json:"region_name"`
		CityName    string  `json:"city_name"`
		Latitude    float64 `json:"latitude"`
		Longitude   float64 `json:"longitude"`
		ZipCode     string  `json:"zip_code"`
		TimeZone    string  `json:"time_zone"`
		ASN         string  `json:"asn"`
		AS          string  `json:"as"`
		IsProxy     bool    `json:"is_proxy"`
	}

	url := fmt.Sprintf("https://api.ip2location.io/?key=%s&ip=%s", apiKey, ip)

	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var loc Location
	err = json.Unmarshal(body, &loc)
	if err != nil {
		panic(err)
	}

	name := fmt.Sprintf("%s / %s", loc.CountryCode, loc.AS)

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
		// get IP owner
		owner, isCached := getOwner(ip)
		var coloredIP string
		if isCached {
			// cahe HIT to indicate an already known IP
			coloredIP = green(ip) + yellow(fmt.Sprintf(" (%s)", owner))
		} else {
			// cache MISS to indicate a new IP
			coloredIP = red(ip) + yellow(fmt.Sprintf(" (%s)", owner))
		}

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
