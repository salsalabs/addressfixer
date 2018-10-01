package addressfixer

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

const (
	//CA is the regex that matches postal codes in Canada.
	CA string = `(^[A-Za-z]\d[A-Za-z])([ -]?\d[A-Za-z]\d)*$`
	//GB is the regex that matches postal codes in Great Britain.  The official
	//regex appears to not work well.  This shorter one appears to work well enough.
	GB string = `^[A-Z]{1,2}[0-9][0-9A-Z]?\s?[0-9][A-Z]{2}`
	//IE is the regex for Ireland (IE).  see https://www.grcdi.nl/gsb/ireland.htm
	IE string = `^A([A|C-F|H|K|N|P|R|T-Z][0-9][0-9|W]( )[A|C-F|H|K|N|P|R|T-Z][A|C-F|H|K|N|P|R|T-Z|0-9][A|C-F|H|K|N|P|R|T-Z|0-9][A|C-F|H|K|N|P|R|T-Z|0-9])$ `
	//NL the the regex that mathes postal codes for the Netherlands.  Also very long...
	NL string = `(?:NL-)?(?:[1-9]\d{3} ?(?:[A-EGHJ-NPRTVWXZ][A-EGHJ-NPRSTVWXZ]|S[BCEGHJ-NPRTVWXZ]))$`
)

//ZPlace is returned by Zippopotamus for a match with the zipcode.
type ZPlace struct {
	Name  string `json:"place name"`
	Long  string `json:"longitude"`
	State string
	Abbr  string `json:"state abbreviation"`
	Lat   string `json:"latitude"`
}

//ZResult is return by Zippotamus for all places that match a
//postal code.
type ZResult struct {
	PostCode    string `json:"post code"`
	Country     string `json:"country"`
	CountryCode string `json:"country abbreviation"`
	Places      []ZPlace
}

//PMap maps a regex pattern to a list of matching country codes.
type PMap map[string][]string

//MatchPostal searches a list of regexes for a zipcode.  Returns
//a matched indicator and the country code.  Sources are StackTrace,
//Wikipedia (https://en.wikipedia.org/wiki/List_of_postal_codes),
//various postal services and https://rgxdb.com/
func MatchPostal(s *Supporter) (bool, string) {
	m := PMap{
		CA:        []string{"CA"},
		GB:        []string{"GB"},
		IE:        []string{"IE"},
		NL:        []string{"NL"},
		`^\d{6}$`: []string{"BY", "CN", "NN", "EC", "KZ", "KG", "NG", "RO", "RU", "SG", "TJ", "TT", "TM", "UZ", "VN"},
		`^\d{5}$`: []string{"AX", "AX", "BA", "BR", "BT", "CC", "CP", "CP", "CR", "DE", "DO", "DZ", "EE", "EG", "ES",
			"FR", "GT", "HR", "ID", "IQ", "IT", "KH", "KR", "KW", "LA", "LB", "LK", "MA", "ME", "MM",
			"MN", "MU", "MV", "MX", "MY", "NI", "NP", "PE", "PK", "PO", "PO", "PO", "RS", "RS", "SD",
			"TH", "TR", "TZ", "UA", "UY", "XK", "ZM"},
		`^\d{5}-?\d{3}$`: []string{"BR"},
		`^\d{4}$`: []string{"AF", "AL", "AR", "AM", "AU", "AT", "BD", "BE", "BG", "CV",
			"CX", "CC", "CY", "DK", "SV", "ET", "GE", "DE", "GL", "GW",
			"HT", "HU", "LR", "LI", "LU", "MK", "MZ", "NZ", "NE", "NF",
			"NO", "PA", "PY", "PH", "PT", "SG", "ZA", "CH", "SJ", "TN"},
		`^\d{3}$`:                     []string{"FO", "GN", "IS", "LS", "NG", "OM", "PS", "PG"},
		`^00120$`:                     []string{"VA"},
		`^00[6-9](?:[-\s]\d{4})`:      []string{"PR"},
		`^008[0-5]\d`:                 []string{"VI"},
		`^4789\d$`:                    []string{"SM"},
		`^96799(?:[-\s]\d{4})$`:       []string{"AS"},
		`^9691\d{2}(?:[-\s]\d{4})?$`:  []string{"GU"},
		`^9695[0-2](?:[-\s]\d{4})?$`:  []string{"MP"},
		`^96960$`:                     []string{"PW"},
		`^969[6-7]\d(?:[-\s]\d{4})?$`: []string{"MH"},
		`^9694[1-4](?:[-\s]\d{4})?$`:  []string{"FM"},
		`^971\d{2}$`:                  []string{"GP"},
		`^97133$`:                     []string{"BL"},
		`^97150$`:                     []string{"MF"},
		`^972\d{2}`:                   []string{"MQ"},
		`^973\d{2}$`:                  []string{"GF"},
		`^974\d{2}$`:                  []string{"RE"},
		`^975\d{2}$`:                  []string{"PM"},
		`^976\d{2}$`:                  []string{"YT"},
		`^980\d{2}$`:                  []string{"MC"},
		`^986\d{2}$`:                  []string{"WF"},
		`^987\d{2}$`:                  []string{"PF"},
		`^988\d{2}$`:                  []string{"NC"},
		`^LC`:                         []string{"LC"},
		`^PCRN`:                       []string{"PN"},
		`^SIQQ`:                       []string{"GS"},
		`^TKCA`:                       []string{"TC"},
	}

	if len(s.Zip) == 0 {
		return false, ""
	}
	for p, c := range m {
		if regexp.MustCompile(p).MatchString(s.Zip) {
			if len(c) == 1 {
				return true, c[0]
			}
			for _, x := range c {
				e := strings.ToUpper(s.Email)
				if strings.HasSuffix(e, "."+x) {
					if x == "US" && len(s.Country) == 0 {
						return true, s.Country
					}
					return true, x
				}
			}
		}
	}
	// Default to the US.  Open for discussion.
	return false, "US"
}

//City checks to see if the supporter's state is correct.  If not, then
//the record is changed and a Mod is added to the list of modifications.
func City(s *Supporter, p ZPlace) (modified bool) {
	modified = false
	s.City = strings.TrimSpace(s.City)
	name := p.Name
	// Zippopotamus shows neighboring towns in a neighborhood
	// in parens.  Not really a good city name.
	if strings.Contains(name, "(") {
		name = strings.Split(name, "(")[0]
		name = strings.TrimSpace(name)
	}
	if len(s.City) == 0 {
		s.City = name
		modified = true
	}
	return modified
}

//State checks to see if the supporter's state is correct.  If not, then
//the record is changed and a Mod is added to the list of modifications.
func State(s *Supporter, p ZPlace) (modified bool) {
	modified = false
	if len(p.Abbr) != 0 && s.State != p.Abbr {
		if strings.Contains(p.Abbr, "Whistler") {
			s.State = "BC"
		} else {
			s.State = p.Abbr
			modified = true
		}
	}
	modified = caKludge(s) || modified
	return modified
}

//Country checks to see if the supporter's state is correct.  If not, then
//the record is changed and a Mod is added to the list of modifications.
func Country(s *Supporter, t ZResult) {
	if s.Country != t.CountryCode {
		s.Country = t.CountryCode
	}
}

//Zippopatamus has a hard time assigning Canadian provinces.
//See https://en.wikipedia.org/wiki/Postal_codes_in_Canada for a chart
//of the first letter of the postal code to province.
func caKludge(s *Supporter) bool {
	modified := false
	provinceMap := map[string]string{
		"A": "NL",
		"B": "NS",
		"C": "PE",
		"E": "NE",
		"G": "QC",
		"H": "QC",
		"I": "QC",
		"K": "ON",
		"L": "ON",
		"M": "ON",
		"N": "ON",
		"P": "ON",
		"R": "MB",
		"S": "SK",
		"T": "AB",
		"V": "BC",
		"X": "NU",
		"Y": "YT",
	}
	if s.Country == "CA" && len(s.State) == 0 && len(s.Zip) > 0 {
		k := s.Zip[0:1]
		k = strings.ToUpper(k)
		v, ok := provinceMap[k]
		if ok {
			s.State = v
			modified = true
		}
	}
	return modified
}

//Fetch retrieves information for a zip code.
func Fetch(s *Supporter) (ZResult, error) {
	// Make adjustment to the postal code submitted to Zippopotamus.
	p := s.Zip
	c := s.Country
	switch c {
	case "CA":
		if len(p) > 2 {
			p = p[0:3]
		}
	case "GB":
		if len(p) > 2 {
			p = p[0:3]
		}
	case "":
		c = "US"
	}
	if c == "US" {
		if strings.Contains(s.Zip, "-") {
			p = strings.Split(s.Zip, "-")[0]
		}
	}
	u := fmt.Sprintf("http://api.zippopotam.us/%v/%v", c, p)
	var body []byte
	var zr ZResult
	resp, err := http.Get(u)
	if resp == nil {
		err = fmt.Errorf("Null response object")
		return zr, err
	}
	if resp.StatusCode != 200 {
		err = fmt.Errorf("Zippo: %v", resp.Status)
	}
	if err != nil {
		return zr, err
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return zr, err
	}
	err = json.Unmarshal(body, &zr)
	if err != nil {
		return zr, err
	}
	if len(zr.Places) == 0 {
		err = fmt.Errorf("Place not found")
		return zr, err
	}
	return zr, err
}

//FixShortZips adds a leading zero to a Zip code if the country is "US",
//the postal code has four digits, and the state is one of the US states
//that has a leading zero.
func FixShortZips(s *Supporter) (modified bool) {
	modified = false
	re := regexp.MustCompile(`^\d{4}$`)
	if s.Country == "US" && re.MatchString(s.Zip) {
		zeroStates := strings.Split("CT,MA,MN,NH,NJ,PR,RI,VT,VI", ",")
		for _, x := range zeroStates {
			if s.State == x {
				z := "0" + s.Zip
				s.Zip = z
				modified = false
			}
		}
	}
	return modified
}

//Zippo does a lookup using the free service from http://www.zippopotam.us/.
//Note that ambiguous results from Zippopotamus are not applied.
func Zippo(s *Supporter) (modified bool, err error) {
	modified = false
	s.Country = strings.TrimSpace(s.Country)
	s.Zip = strings.TrimSpace(s.Zip)
	s.Zip = strings.ToUpper(s.Zip)
	if len(s.Zip) == 0 {
		return modified, nil
	}
	FixShortZips(s)
	m, c := MatchPostal(s)
	if m {
		if c != s.Country {
			s.Country = c
			modified = true
		}
	}
	zr, err := Fetch(s)
	if err != nil {
		return modified, err
	}
	Country(s, zr)
	if len(zr.Places) > 0 {
		zp := zr.Places[0]
		modified = modified || City(s, zp)
		modified = modified || State(s, zp)
	}
	return modified, nil
}
