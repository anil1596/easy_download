package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"unicode"

	"golang.org/x/net/publicsuffix"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter Channel name: ")
	hint, _ := reader.ReadString('\n')
	response, err := http.Get("http://google.com/search?q=youtube" + hint)
	check_error(err)

	var channel_links = make(map[int]string)

	defer response.Body.Close()
	doc, err := goquery.NewDocumentFromReader(io.Reader(response.Body))
	check_error(err)

	doc.Find("h3.r a").Each(func(i int, s *goquery.Selection) {
		str, exists := s.Attr("href")
		if exists {
			u, err := url.Parse(str)
			check_error(err)
			m, _ := url.ParseQuery(u.RawQuery)
			fmt.Println(i+1, ") \033[1;35m"+s.Text()+"\033[0m", m["q"][0])
			channel_links[i+1] = m["q"][0]
		} else {
			fmt.Println(s.Text())
		}
	})
	number := 0
	fmt.Println("\n\nEnter the channel number from where you would like to DOWNLOAD:")
	_, err = fmt.Scanf("%d", &number)

	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}

	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
	}

	options := cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	}
	jar, err := cookiejar.New(&options)
	check_error(err)

	client := http.Client{Transport: transport, Jar: jar}

	res, _ := client.Get(channel_links[number])
	defer res.Body.Close()
	main_link := channel_links[number]

	doc, err = goquery.NewDocumentFromResponse(res)
	check_error(err)

	var videos_links = make(map[int]string)

	doc.Find(".yt-lockup-content").Each(func(i int, s *goquery.Selection) {
		title := s.Text()
		t := stringMinifier(title)
		result := strings.Split(t, "-")
		if result[0] == "" {
			result = strings.Split(t, "Duration")
		}
		link, _ := s.Find("a").Attr("href")
		if i < 12 {
			videos_links[i+1] = main_link + link
			fmt.Printf("\n%d) %s", i+1, result[0])
			fmt.Printf("\n    %s\n", result[1])
		}
	})

	to_get := 1

	fmt.Printf("\n\n       Enter index of video  to download:")
	_, err = fmt.Scanf("%d", &to_get)
	check_error(err)

	fmt.Printf("\n\n   DOWNLOADING YOUR VIDEO:(wait for some time)")
	cmd := exec.Command("youtube-dl", videos_links[to_get])
	_, err = cmd.Output()
	check_error(err)
	print("\n\n DOWNLOAD  SUCCESSFUL\n")

}

func check_error(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

//to repair retreived string
func stringMinifier(in string) (out string) {
	white := false
	for _, c := range in {
		if unicode.IsSpace(c) {
			if !white {
				out = out + " "
			}
			white = true
		} else {
			out = out + string(c)
			white = false
		}
	}
	return
}
