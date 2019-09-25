package main

import (
	"encoding/xml"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"fmt"
)

// GetBookings will return the current bookings from Andy's booking system
func GetBookings() []byte {

	//creating the proxyURL
	proxyStr := "http://empweb2.ey.net:8080/"
	proxyURL, err := url.Parse(proxyStr)
	if err != nil {
		log.Println(err)
	}

	//creating the URL to be loaded through the proxy
	urlStr := "https://booking.minervabathrc.org.uk/api/booking/all"

	url, err := url.Parse(urlStr)
	if err != nil {
		log.Println(err)
	}

	//adding the proxy settings to the Transport object
	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}

	//adding the Transport object to the http Client
	client := &http.Client{
		Transport: transport,
	}

	//generating the HTTP GET request
	request, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		log.Println(err)
	}

	//calling the URL
	response, err := client.Do(request)
	if err != nil {
		log.Println(err)
	}

	//getting the response
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
	}
	//printing the response
	return data
}

// ArrayOfBookingSummary don't care about this - just an XML holder for the data
type ArrayOfBookingSummary struct {
	XMLName        xml.Name         `xml:"ArrayOfBookingSummary"`
	Text           string           `xml:",chardata"`
	I              string           `xml:"i,attr"`
	Xmlns          string           `xml:"xmlns,attr"`
	BookingSummary []BookingSummary `xml:"BookingSummary"`
}

// BookingSummary the booking elements
type BookingSummary struct {
	Text          string `xml:",chardata"`
	BookedByName  string `xml:"bookedByName"`
	BookedForName string `xml:"bookedForName"`
	EndDateTime   string `xml:"endDateTime"`
	Icon          string `xml:"icon"`
	ID            string `xml:"id"` // UID Int
	ResourceName  string `xml:"resourceName"`
	StartDateTime string `xml:"startDateTime"`
}

// ParseBookings will return a slice of booking structs
func ParseBookings() []BookingSummary {

	// Bookings get parsed here

	b := GetBookings()

	var a ArrayOfBookingSummary

	xml.Unmarshal(b, &a)

	aslice := a.BookingSummary

	return aslice
}

func storeBookings(b []BookingSummary) {
	// internal function to send bookings to a database - this should run at least once daily
	// to capture midnight booking status
	// will require some uniquification...?
}

// func main() {
// 	b := ParseBookings()

// 	fmt.Println(b)
// }
