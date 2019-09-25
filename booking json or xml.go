type BookingSummary struct {
	XMLName       xml.Name `xml:"BookingSummary"`
	Text          string   `xml:",chardata"`
	BookedByName  string   `xml:"bookedByName"`
	BookedForName string   `xml:"bookedForName"`
	EndDateTime   string   `xml:"endDateTime"`
	Icon          string   `xml:"icon"`
	ID            string   `xml:"id"`
	ResourceName  string   `xml:"resourceName"`
	StartDateTime string   `xml:"startDateTime"`
} 
