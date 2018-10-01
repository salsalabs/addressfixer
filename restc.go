package addressfixer

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

//RResp is what's returned by RestCountries for
//an input.
type RResp struct {
	Input      string `json:"name"`
	Alpha2Code string `json:"alpha2Code"`
}

//RFetch retrieves information for a zip code.
func RFetch(s Supporter) (RResp, error) {
	u := fmt.Sprintf("https://restcountries.eu/rest/v2/name/%v", s.Country)
	var body []byte
	var rr RResp
	var target []RResp
	resp, err := http.Get(u)
	if resp == nil {
		err = fmt.Errorf("RestC:  %-8s null response on %v", s.Key, u)
	} else {
		if resp.StatusCode != 200 {
			//fmt.Printf("RFetch: %v  %v\n", resp.Status, u)
			err = fmt.Errorf("RestC: %v", resp.Status)
		}
	}
	if err != nil {
		return rr, err
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return rr, err
	}
	err = json.Unmarshal(body, &target)
	if len(target) == 0 {
		err = fmt.Errorf("country %v returned no restcountry results", s.Country)
	}
	if err != nil {
		return rr, err
	}
	return target[0], err
}

//RestCountries accepts a supporter record.  If the country
//is empty or two digits, then nothing happens.  If the country
//is three digits or more, then we'll use that as input for
//a lookup on restcountries.eu.  Matches modify the supporter
//record and add a Mod to the list of modifications.
func RestCountries(s *Supporter) (bool, error) {
	modified := false
	if len(s.Country) > 2 {
		rr, err := RFetch(*s)
		if err != nil {
			return modified, err
		}
		if rr.Alpha2Code == "UM" {
			rr.Alpha2Code = "US"
		}
		modified = s.Country == rr.Alpha2Code
		s.Country = rr.Alpha2Code
	}
	return modified, nil
}
