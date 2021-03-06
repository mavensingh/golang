package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/gookit/color"
	"golang.org/x/net/html"
)

func Domains() {
	//fetch data from given url
	// domain := os.Args
	domain := flag.String("d", "sample.com", "Name of the Domain")
	protocol := flag.String("p", "https|http", "Protocol of the domain")
	flag.Parse()
	response, err := http.Get(*protocol + "://" + *domain + "/robots.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	// Read data from url
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	storeDomains := string(body)
	re := regexp.MustCompile(`sitemap: (.*)`)
	// re := regexp.MustCompile(`http://(.*)`)
	submatchall := re.FindAllStringSubmatch(storeDomains, -1)
	for _, element := range submatchall {
		// fmt.Println(element[1])
		response, err := http.Get(element[1])
		if err != nil {
			log.Fatal(err)
		}
		defer response.Body.Close()
		// Read data from url
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}
		storeDomains := string(body)
		re := regexp.MustCompile(`<loc>(.*)</loc>`)
		// re := regexp.MustCompile(`http://(.*)`)
		submatchall := re.FindAllStringSubmatch(storeDomains, -1)

		//find last index of an string
		lastindex := len(submatchall) - 1
		urls := (submatchall[lastindex][1])
		// fmt.Println(urls)
		resp, err := http.Get(urls)
		if err != nil {
			fmt.Println(err)
		}
		defer resp.Body.Close()
		bdy, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
		}
		storeurls := string(bdy)
		regx := regexp.MustCompile(`<loc>(.*)</loc>`)
		submatch := regx.FindAllStringSubmatch(storeurls, -1)
		// lastindex = len(submatch) - 1
		// urls = (submatch[lastindex][1])
		// fmt.Println(storeurls)
		for _, elements := range submatch {
			// fmt.Println(elements[1])
			if resp, err := http.Get(elements[1]); err == nil {
				defer resp.Body.Close()

				// log.Println("Load page complete")

				if resp != nil {
					// log.Println("Page response is NOT nil")
					// --------------
					data, _ := ioutil.ReadAll(resp.Body)
					hdata := strings.Replace(string(data), "<noscript>", "", -1)
					hdata = strings.Replace(hdata, "</noscript>", "", -1)
					// --------------

					if document, err := html.Parse(strings.NewReader(hdata)); err == nil {
						yellow := color.FgYellow.Render
						var parser func(*html.Node)
						parser = func(n *html.Node) {
							if n.Type == html.ElementNode && n.Data == "img" {

								var imgSrcUrl string

								for _, element := range n.Attr {
									if element.Key == "src" {
										imgSrcUrl = element.Val
									}
									// if element.Key == "data-original" {
									// 	imgDataOriginal = element.Val
									// }
								}
								// fmt.Println(*domain + " ---------> " + imgSrcUrl)
								re := strings.Replace(imgSrcUrl, *protocol+"://"+*domain, "", -1)
								re = strings.Replace(imgSrcUrl, "http://"+*domain, "", -1)
								// log.Println("imageUrl " + re)

								resp, err := http.Get(*protocol + "://" + *domain + re)
								red := color.FgRed.Render

								if err != nil {
									fmt.Println(err)
								}
								if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
									fmt.Println(re)
									fmt.Print()
								} else {
									fmt.Println(red(re))
								}
								// resp.Body.Close()
								// bodyElement,err := ioutil.ReadAll(resp)
							}

							for c := n.FirstChild; c != nil; c = c.NextSibling {
								parser(c)
							}

						}
						fmt.Println(yellow("***************************************************************************************************"))
						parser(document)
					} else {
						log.Panicln("Parse html error", err)
					}

				} else {
					log.Println("Page response IS nil")
				}
			}
		}
	}

}

func main() {
	Domains()
	os.Exit(0)
	// http.HandleFunc("/v1/domains", Domains)
	// fmt.Println("Development server started localhost:8000/v1/domains")
	// log.Fatal(http.ListenAndServe(":8000", nil))
}
