package main

import (
	"bytes"
	"compress/gzip"
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
<<<<<<< HEAD

=======
	"sync"
>>>>>>> bcb8952a409d7fd8293a94e9e471c35ba79bc93e
)

type SitemapIndex struct {
	Locations []Location `xml:"sitemap"`
}

type URLSets struct {
	Urls []SitemapURL `xml:"sitemap"`
}

type URLSet struct {
	Urls []SitemapURL `xml:"url"`
}

type SitemapURL struct {
	Location         string `xml:"loc"`
	LastModifiedDate string `xml:"lastmod"`
	ChangeFrequency  string `xml:"changefreq"`
	Priority         string `xml:"priority"`
}

type Location struct {
	Loc string `xml:"loc"`
}

func (l Location) String() string {
	return fmt.Sprintf(l.Loc)
}

var waitGroup sync.WaitGroup

//Unzip files from the webpage
func gunzipWrite(w io.Writer, data []byte) error {
	// Write gzipped data to the client
	gr, err := gzip.NewReader(bytes.NewBuffer(data))
	defer gr.Close()
	data, err = ioutil.ReadAll(gr)
	if err != nil {
		return err
	}
	w.Write(data)
	return nil
}

func DetectType(urlSitemap, pageType string) {
	defer waitGroup.Done()
	re := regexp.MustCompile(`/(xml|x-gzip|html)`)
	ContentType := re.FindAllStringSubmatch(string(pageType), -1)
	for _, contents := range ContentType {
		log.Println(contents[1])
		if contents[1] == "xml" {
			go ParseXml(urlSitemap)
			fmt.Println("********************************************************************")
			return
		} else if contents[1] == "x-gzip" {
			zipbytes := FetchUrl(urlSitemap)
			if string(zipbytes) != "" {
				var buf bytes.Buffer
				err := gunzipWrite(&buf, zipbytes)
				if err != nil {
					log.Fatal(err)
					return
				}
				// fmt.Println("decompressed:\t", buf.String())
				var s URLSet
				xml.Unmarshal(buf.Bytes(), &s)
				// for _, Location := range s.Urls {
				// fmt.Printf("%s\n", Location)
				// }
				urlCount := 0
				for i := range s.Urls {
					url := s.Urls[i]
					log.Println(url)
					// createfile(string(filename.Host), string(url))
					urlCount++
					//log.Printf(">>%s\n", url)
				}
				var N URLSets
				xml.Unmarshal(buf.Bytes(), &s)
				// for _, Location := range s.Urls {
				// fmt.Printf("%s\n", Location)
				// }
				for i := range N.Urls {
					url := N.Urls[i]
					// log.Println(url)
					// createfile(string(filename.Host), string(url))
					// urlCount++
					//log.Printf(">>%s\n", url)
					zipbytes := FetchUrl(url.Location)
					if string(zipbytes) != "" {
						var buf bytes.Buffer
						err := gunzipWrite(&buf, zipbytes)
						if err != nil {
							log.Fatal(err)
							return
						}
						// fmt.Println("decompressed:\t", buf.String())
						var s URLSet
						xml.Unmarshal(buf.Bytes(), &s)
						// for _, Location := range s.Urls {
						// fmt.Printf("%s\n", Location)
						// }
						for i := range s.Urls {
							url := s.Urls[i]
							log.Println(url)
							// createfile(string(filename.Host), string(url))
							urlCount++
							//log.Printf(">>%s\n", url)
						}
					}
				}
				log.Println("Numbers of url found", urlCount)
			} else {
				log.Println("Page response in Nil")
			}

			fmt.Println("********************************************************************")
			return

		} else {
			// xmlData := FetchUrl(urlSitemap)
			// if string(xmlData) == "" {
			// 	log.Println("Xml file is Empty")
			// } else {
			// 	var s URLSet
			// 	xml.Unmarshal(xmlData, &s)
			// 	urlCount := 0
			// 	for i := range s.Urls {
			// 		url := s.Urls[i]
			// 		log.Println(url)
			// 		urlCount++
			// 	}
			// 	var N URLSets
			// 	xml.Unmarshal(xmlData, &N)
			// 	// for _, Location := range s.Urls {
			// 	// fmt.Printf("%s\n", Location)
			// 	// }
			// 	for i := range N.Urls {
			// 		url := N.Urls[i]
			// 		log.Println(url)
			// 		urlCount++

			// 	}
			// 	log.Println("Numbers of url found", urlCount)
			// 	return
			// }
			ParseXml(urlSitemap)
			// log.Println("Program not able to read the content of the page")
			fmt.Println("********************************************************************")
			return
		}
	}
}

//Detect which type of url we getting from Robots.txt file
func DetectTypeOfFiles(urlSitemap string) {
	fmt.Println(urlSitemap)
	headerInfo, err := http.Get(urlSitemap)
	if err != nil {
		log.Println("response not found", err)
	}
	defer headerInfo.Body.Close()
	// log.Println(headerInfo.Header.Get("Content-type"))
	pageType1 := headerInfo.Header.Get("Content-type")
	pageType2 := headerInfo.Header.Get("content-type")
	waitGroup.Add(1)
	if headerInfo.Header.Get("Content-type") == pageType1 {
		go DetectType(urlSitemap, pageType1)
	} else if headerInfo.Header.Get("content-type") == pageType2 {
		go DetectType(urlSitemap, pageType2)
	}
	waitGroup.Wait()
}

//fetch all domain+urls
func FetchUrl(url string) []byte {
	response, err := http.Get(url)
	if err != nil {
		log.Println("having problem to find url", err)
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println("Nothing found in the page", err)
	}
	// log.Println(string(body))
	return body

}

//fetch all domain+urls
func VisitRobotsTxt(domain string) string {
	response, err := http.Get(domain + "robots.txt")
	if err != nil {
		log.Println("having problem to find url", err)
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println("Nothing found in the page", err)
	}
	// fmt.Println(string(body))
	return string(body)

}

//collect All sitemaps from the robots.txt
func getSitemapFromRobotsTxt(domain string) {
	// re := regexp.MustCompile(`(.*)-robots.txt`)
	robotsFile := VisitRobotsTxt(domain)
	filelower := strings.ToLower(robotsFile)
	re := regexp.MustCompile(`sitemap:(.*)`)
	FileToDomain := re.FindAllStringSubmatch(string(filelower), -1)
	if FileToDomain == nil {
		log.Println("No sitemap url found", domain)
		return
	}
	waitGroup.Add(len(FileToDomain))
	for _, Domain := range FileToDomain {
		response := strings.Replace(Domain[1], " ", "", -1)
		// log.Println(Domain[1])
		go DetectTypeOfFiles(response)
		// CreateDomainsFile(string(Domain[1]))
	}
	waitGroup.Wait()
	waitGroup.Done()
	return
}

//Parse All url from Xml file
func ParseXml(url string) {
	xmlData := FetchUrl(url)
	if string(xmlData) == "" {
		log.Println("Xml file is Empty")
	} else {
		// filename, err := url.Parse(urlSitemap)
		// if err != nil {
		// log.Println("Parsing failed", err)
		// }
		var s URLSet
		xml.Unmarshal(xmlData, &s)
		// for _, Location := range s.Urls {
		// fmt.Printf("%s\n", Location)
		// }
		urlCount := 0
		for i := range s.Urls {
			url := s.Urls[i]
			log.Println(url)
			// createfile(string(filename.Host), string(url))
			urlCount++
			//log.Printf(">>%s\n", url)
		}
		var N URLSets
		xml.Unmarshal(xmlData, &N)
		// for _, Location := range s.Urls {
		// fmt.Printf("%s\n", Location)
		// }
		for i := range N.Urls {
			url := N.Urls[i]
			// log.Println(url)
			// createfile(string(filename.Host), string(url))
			// urlCount++
			//log.Printf(">>%s\n", url)
			zipbytes := FetchUrl(url.Location)
			// log.Println(string(zipbytes))
			// if string(zipbytes) != "" {
			// 	var buf bytes.Buffer
			// 	err := gunzipWrite(&buf, zipbytes)
			// 	if err != nil {
			// 		log.Fatal(err)
			// 		// return
			// 	}
			// 	// fmt.Println("decompressed:\t", buf.String())
			// 	var s URLSet
			// 	xml.Unmarshal(buf.Bytes(), &s)
			// 	// for _, Location := range s.Urls {
			// 	// fmt.Printf("%s\n", Location)
			// 	// }
			// 	for i := range s.Urls {
			// 		url := s.Urls[i]
			// 		log.Println(url)
			// 		// createfile(string(filename.Host), string(url))
			// 		urlCount++
			// 		//log.Printf(">>%s\n", url)
			// 	}

			// }
			var SS URLSet
			xml.Unmarshal(zipbytes, &SS)
			// for _, Location := range s.Urls {
			// fmt.Printf("%s\n", Location)
			// }
			for i := range SS.Urls {
				url := SS.Urls[i]
				log.Println(url)
				// createfile(string(filename.Host), string(url))
				urlCount++
				//log.Printf(">>%s\n", url)
			}
		}
		log.Println("Numbers of url found", urlCount)
		return
	}
	waitGroup.Done()

}

func IfBodyDataIsSimple(pageContent string) {
	split := strings.Split(pageContent, " ")
	log.Println(split)
	return
	// createfile(filename, string(split))
}

func createfile(filename, output string) {
	out, err := os.Create("./output/" + filename + ".xml")
	if err != nil {
		log.Println(err)
		return

	}
	defer out.Close()
	//Write the body into file
	_, err = io.WriteString(out, output)
	if err != nil {
		log.Println(err)
		return
	}
}

func ReadFile(filepath string) {
	readfile, err := ioutil.ReadFile(filepath)
	if err != nil {
		log.Println("Error to reading the file", err)
		os.Exit(1)
	}
	split := strings.Split(string(readfile), "\n")
	// EachDomains := []string{}
	waitGroup.Add(len(split))
	for _, line := range split {
		// timeout := time.Duration(1 * time.Second)
		// conn, err := net.DialTimeout("tcp","mysyte:myport", timeout)
		// if err != nil {
		// log.Println("Site unreachable, error: ", err)
		// }
		// fetch(line)
		go getSitemapFromRobotsTxt(line)
		// EachDomains = append(EachDomains, line)

	}
	// fmt.Println(EachDomains)
	waitGroup.Wait()
	waitGroup.Done()
}

func main() {
	domain := flag.String("f", "", "Provide the path of the file")
	flag.Parse()
	waitGroup.Add(1)
	// go ReadFile(*domain)
	go getSitemapFromRobotsTxt(*domain)
	// fetch("https://chaufferus.us/https-sitemap_index.xml")
	waitGroup.Wait()
	waitGroup.Done()

}
