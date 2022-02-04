package groupie

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
)

/*This var is a pointer towards template.Template that is a
pointer to help process the html.*/
var tpl *template.Template

/*This init function, once it's initialised, makes it so that each html file
in the templates folder is parsed i.e. they all get looked through once and
then stored in the memory ready to go when needed*/
func init() {
	tpl = template.Must(template.ParseGlob("templates/*html"))
}

var (
	ArtistID              []int
	ArtistImage           []string
	ArtistName            []string
	ArtistMembers         [][]string
	ArtistCreationDate    []int
	ArtistFirstAlbum      []string
	ArtistLocations       [][]string
	ArtistConcertDates    [][]string
	ArtistsDatesLocations map[string][]string
)

type Artists []struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Locations    string   `json:"locations"`
	ConcertDates string   `json:"concertDates"`
	Relations    string   `json:"relations"`
}

type Dates struct {
	Dates []dates `json:"index"`
}

type dates struct {
	ID    int      `json:"id"`
	Dates []string `json:"dates"`
}
type Locations struct {
	Locations []locations `json:"index"`
}

type locations struct {
	ID        int      `json:"id"`
	Locations []string `json:"locations"`
	Dates     string   `json:"dates"`
}

type Relations struct {
	Relations []relations `json:"index"`
}

type relations struct {
	ID             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

func main() {

	UnmarshalArtistData()

	fmt.Println(ArtistsDatesLocations["saitama-japan"])

}

func UnmarshalArtistData() {

	responseArtists, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		panic("Couldn't get Artists info from API")
	}
	defer responseArtists.Body.Close()

	responseArtistsData, err := ioutil.ReadAll(responseArtists.Body)
	if err != nil {
		panic("Couldn't read data for Artists!")
	}

	var responseObjectArtists Artists

	json.Unmarshal(responseArtistsData, &responseObjectArtists)

	for i := 0; i < len(responseObjectArtists); i++ {
		ArtistFirstAlbum = append(ArtistFirstAlbum, responseObjectArtists[i].FirstAlbum)
	}

	for i := 0; i < len(responseObjectArtists); i++ {
		ArtistID = append(ArtistID, responseObjectArtists[i].ID)
	}

	for i := 0; i < len(responseObjectArtists); i++ {
		ArtistImage = append(ArtistImage, responseObjectArtists[i].Image)
	}

	for i := 0; i < len(responseObjectArtists); i++ {
		ArtistMembers = append(ArtistMembers, responseObjectArtists[i].Members)
	}

	for i := 0; i < len(responseObjectArtists); i++ {
		ArtistCreationDate = append(ArtistCreationDate, responseObjectArtists[i].CreationDate)
	}

	for i := 0; i < len(responseObjectArtists); i++ {
		ArtistName = append(ArtistName, responseObjectArtists[i].Name)
	}

	responseRelations, err := http.Get("https://groupietrackers.herokuapp.com/api/relation")
	if err != nil {
		panic("Couldn't get the relations data!")
	}

	responseData, err := ioutil.ReadAll(responseRelations.Body)
	if err != nil {
		panic("Couldn't read data for the Artists")
	}

	var responseObjectRelations Relations

	json.Unmarshal(responseData, &responseObjectRelations)

	//x := responseObjectRelations.Relations[0].DatesLocations

	// for k, v := range x {

	// 	fmt.Println(k, v)

	// }
	for i := 0; i < len(Artists{}); i++ {

		ArtistsDatesLocations = responseObjectRelations.Relations[i].DatesLocations

	}

	responseDates, err := http.Get("https://groupietrackers.herokuapp.com/api/dates")
	if err != nil {
		panic("Couldn't get Dates info from the API!")
	}
	defer responseDates.Body.Close()

	responseDatesData, err := ioutil.ReadAll(responseDates.Body)
	if err != nil {
		panic("Couldn't read data for Dates")
	}

	var responseObjectDates Dates
	json.Unmarshal(responseDatesData, &responseObjectDates)

	for i := 0; i < len(responseObjectDates.Dates); i++ {
		ArtistConcertDates = append(ArtistConcertDates, responseObjectDates.Dates[i].Dates)
	}

	responseLocations, err := http.Get("https://groupietrackers.herokuapp.com/api/locations")
	if err != nil {
		panic("Couldn't get Location info from API")
	}
	defer responseLocations.Body.Close()

	responseLocationsData, err := ioutil.ReadAll(responseLocations.Body)
	if err != nil {
		panic("Couldn't read data for Locations!")
	}

	var responseObjectLocations Locations
	json.Unmarshal(responseLocationsData, &responseObjectLocations)

	//fmt.Println(responseObjectLocations.Locations[0].Locations)

	for i := 0; i < len(responseObjectLocations.Locations); i++ {
		ArtistLocations = append(ArtistLocations, responseObjectLocations.Locations[i].Locations)
	}

}

func Requests() {

	http.HandleFunc("/", index)
	http.HandleFunc("/info", artistInfo)
	http.ListenAndServe(":8080", nil)
	log.Println("Server started on: http://localhost:8080")
}

func index(w http.ResponseWriter, r *http.Request) {

	//-------------Create a struct to hold unmarshalled data-----------

	var TotalInfo struct {
		ArtistID           []int
		ArtistImage        string
		ArtistName         []string
		ArtistMembers      [][]string
		ArtistCreationDate []int
		ArtistFirstAlbum   []string
		ArtistLocations    [][]string
		ArtistConcertDates [][]string
	}

	if r.URL.Path != "/" {
		http.Error(w, "404 address not found: wrong address entered!", http.StatusNotFound)
	} else {

		tpl.ExecuteTemplate(w, "index.html", TotalInfo)
	}
}

func artistInfo(w http.ResponseWriter, r *http.Request) {

	response, err := http.Get("https://groupietrackers.herokuapp.com/api/relation")
	if err != nil {
		panic("Couldn't get the relations data!")
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic("Couldn't read data for the Artists")
	}

	var responseObject Relations

	json.Unmarshal(responseData, &responseObject)

	if r.URL.Path != "/info" {
		http.Error(w, "404 address not found: wrong address entered!", http.StatusNotFound)
	} else {

		tpl.ExecuteTemplate(w, "info.html", ArtistsDatesLocations)
	}

}
