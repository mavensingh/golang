package main

import (
	"fmt"
	"os/exec"
	"strings"
)

func FetchData() {
	var Domain_List [14]string
	Domain_List[0] = "https://driverus.us"
	Domain_List[1] = "https://chaufferus.us"
	Domain_List[2] = "https://chaufferjob.us"
	Domain_List[3] = "https://driveus.us"
	Domain_List[4] = "https://ownboss.us"
	Domain_List[5] = "https://driverjob.us"
	Domain_List[6] = "https://doordash.givearide.us"
	Domain_List[7] = "https://tutree.com"
	Domain_List[8] = "https://givearide.us"
	Domain_List[9] = "https://51talkus.com"
	Domain_List[10] = "https://sprout4future.us"
	Domain_List[11] = "https://rouchi.us"
	Domain_List[12] = "https://qkids.co"
	Domain_List[13] = "https://delivery.givearide.us"

	for domains := 0; domains < len(Domain_List); domains++ {
		data_store1 := Domain_List[domains]         //store all domains inside of this variable
		fmt.Println("<h1>" + data_store1 + "</h1>") //All domains name will print
		Collect_url := [4]string{"/fbapply/", "/fbapply-i/", "/fbapply-dd/", "/fbapply-dd-a/"}
		for url_segments := 0; url_segments < len(Collect_url); url_segments++ {
			data_store2 := Collect_url[url_segments]
			// fmt.Println("***********************************************************************************************")
			fmt.Println("<h2>" + data_store2 + "</h2>")
			// fmt.Println("<xmp>")
			// Execute Command to collect Data
			out, err := exec.Command("curl", data_store1+data_store2).Output()
			if err != nil {
				fmt.Println(err)
			}
			// output := string(out[0:])
			output := string(out[23:56]) // Select range
			// filter_fields := strings.Split(output, ";")
			filter_values := strings.Trim(output, ";</scri")
			// filter_values := strings.Trim(output, ";url=")
			// filter := filter_data(output)
			// filter_values := filter_data(output)
			// fmt.Println(filter_data(output))
			fmt.Println("<ul>" + "<li>" + filter_values + "</li>" + "</ul>")
			// output := string(out[0:])
			// sp_values := strings.Split(output, "=")
			// fmt.Println(output)
			// fmt.Println(filter_values)
			// fmt.Println(data_store1 + data_store2)

			// fmt.Println("<xmp>")

		}
		// fmt.Println("***********************************************************************************************")
		fmt.Println("<h2>" + "Ping" + "</h2>")
		out, err := exec.Command("curl", data_store1+"/open-positions/ping").Output()
		if err != nil {
			fmt.Println(err)
		}
		output := string(out[20:36]) // Select range
		fmt.Println("<ul>" + "<li>" + output + "</li>" + "</ul>")
		// fmt.Println(data_store1)
		fmt.Println("<hr/>")
		// fmt.Println("***********************************************************************************************")
		// fmt.Println("")

	}

}

func main() {
	FetchData()
}

