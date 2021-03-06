package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	pdfcontent "github.com/unidoc/unidoc/pdf/contentstream"
	pdf "github.com/unidoc/unidoc/pdf/model"
)

func main() {
	filename := flag.String("f", "", "Provide the name|path of the file")
	flag.Parse()
	fileextension := filepath.Ext(*filename)
	// fmt.Println(fileextension)
	if fileextension == ".csv" {
		CsvFileReader(*filename)
	} else if fileextension == ".txt" {
		TxtFileReader(*filename)
	} else if fileextension == ".xml" {
		XmlFileReader(*filename)
	} else if fileextension == ".json" {
		JsonFilereader(*filename)
	} else if fileextension == ".pdf" {
		PdfFileReader(*filename)
	} else {
		fmt.Println("This Tool doesn't support " + fileextension + " file")
	}

	// fmt.Println("Extension of the file:", fileextension)
}

//CsvFileReader  reader the All data of the csv file
func CsvFileReader(filename string) {
	// Open the file
	csvfile, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer csvfile.Close()
	// str2 := strings.Join(csvfile, " ")
	// str2 = string(str2)
	// t := strings.Split(str2, "\n")
	// fmt.Println(string(t[1]))
	// fmt.Println(" Date JobTitle City Trending ")
	// for _, line := range csvfile {
	// 	// str := string(line)
	// 	// fmt.Println(string(t))
	// 	// fmt.Println(line)
	// 	// for _, single := range t {
	// 	// 	fmt.Println(string(single))
	// 	// }
	// }
	r := csv.NewReader(csvfile)
	// Iterate through the records
	html := "<center>"
	html += "<table border='1' style='border-collapse: collapse;text-align:center'>"
	html += "<thead>"
	html += "<tr>"
	html += "<th>Date</th>"
	html += "<th>JobTitle</th>"
	html += "<th>City</th>"
	html += "<th>Trending</th>"
	html += "</tr>"
	html += "</thead>"
	for {
		// Read each record from csv
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		// fmt.Println(record[0])
		html += "<tbody>"
		html += "<tr>"
		html += "<td>" + record[0] + "</td>"
		html += "<td>" + record[1] + "</td>"
		html += "<td>" + record[2] + "</td>"
		html += "<td>" + record[3] + "</td>"
		html += "</tr>"
		html += "</tbody>"
	}
	html += "</table>"
	html += "</center>"
	fmt.Println(html)
}

//TxtFileReader  reader the All data of the txt file
func TxtFileReader(filename string) {
	//Open the file
	txtfile, err := scanLines(filename)
	if err != nil {
		panic(err)
	}
	for _, line := range txtfile {
		fmt.Println(line)
	}

}

//XmlFileReader reader the All data of the xml file
func XmlFileReader(filename string) {
	//Open the file
	xmlfile, err := scanLines(filename)
	if err != nil {
		panic(err)
	}
	for _, line := range xmlfile {
		fmt.Println(line)
	}
}

//JsonFileReader  reader the All data of the json file in array
func JsonFilereader(filename string) {
	//Open the file
	jsonFile, err := scanLines(filename)
	if err != nil {
		panic(err)
	}
	for _, line := range jsonFile {
		fmt.Println(line)
	}
}

//PdfFileReader function Read the All content of the Pdf file
func PdfFileReader(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	pdfReader, err := pdf.NewPdfReader(file)
	if err != nil {
		panic(err)
	}

	isEncrypted, err := pdfReader.IsEncrypted()
	if err != nil {
		panic(err)
	}

	if isEncrypted {
		_, err = pdfReader.Decrypt([]byte(""))
		if err != nil {
			panic(err)
		}
	}

	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		panic(err)
	}
	fmt.Println(numPages)
	fmt.Printf("--------------------\n")
	fmt.Printf("PDF to text extraction:\n")
	fmt.Printf("--------------------\n")
	for i := 0; i < numPages; i++ {
		pageNum := i + 1

		page, err := pdfReader.GetPage(pageNum)
		if err != nil {
			panic(err)
		}
		// fmt.Println(pageNum)

		contentStreams, err := page.GetContentStreams()
		if err != nil {
			panic(err)
		}
		// fmt.Println(contentStreams)

		// If the value is an array, the effect shall be as if all of the streams in the array were concatenated,
		// in order, to form a single stream.
		pageContentStr := ""
		for _, cstream := range contentStreams {
			pageContentStr += cstream
		}

		// fmt.Printf("Page %d - content streams %d:\n", pageNum, len(contentStreams))
		cstreamParser := pdfcontent.NewContentStreamParser(pageContentStr)
		txt, err := cstreamParser.ExtractText()
		if err != nil {
			panic(err)
		}
		fmt.Printf("\"%s\"\n", txt)
	}
}

//scanLines line by line in all format files
func scanLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, nil
}
