package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
)

type SitemapIndex struct {
	Locations []Location `xml:"sitemap"`
}

type Location struct {
	Loc string `xml:"loc"`
}

func main() {
	domain := flag.String("d", "domain.com", "provide domain")
	protocol := flag.String("p", "https", "Protocol")
	flag.Parse()
	response, err := http.Get(*protocol + "://" + *domain + "/robots.txt")
	Error(err)
	defer response.Body.Close()
	responsebody, err := ioutil.ReadAll(response.Body)
	Error(err)
	str := string(responsebody)
	xmlfinder(str)

	//xml parsing
	response, err = http.Get(*protocol + "://" + *domain + "/robots.txt")
	Error(err)
	defer response.Body.Close()
	responsebody, err = ioutil.ReadAll(response.Body)
	Error(err)
	str = string(responsebody)
	// fmt.Println(str)
}

// func XMldatatype(xmldata []uint8) string {
// 	return string(xmldata)
// }
func (l Location) String() string {
	return fmt.Sprintf(l.Loc)
}

func xmldatacollect(xmlfile string) {
	response, err := http.Get(xmlfile)
	Error(err)
	defer response.Body.Close()
	responsebody, err := ioutil.ReadAll(response.Body)
	Error(err)
	// fmt.Println(string(responsebody))
	// parseXmldata(responsebody)
	bodyxml(responsebody)
	// return string(responsebody)
}

func bodyxml(data []uint8) {
	var s SitemapIndex
	xml.Unmarshal(data, &s)
	// fmt.Println(s.Locations)
	// lastindex := len(s.Locations) - 1
	for _, Location := range s.Locations {
		// fmt.Printf("%s\n", s.Locations[lastindex])
		fmt.Printf("%s\n", Location)
	}
}

func xmlfinder(domain string) {
	sitemap := regexp.MustCompile(`sitemap: (.*)`)
	xmlpath := sitemap.FindAllStringSubmatch(domain, -1)
	// fmt.Println(xmlpath)
	for _, element := range xmlpath {
		// fmt.Println(element[1])
		xmldatacollect(element[1])
	}
}

func Error(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
