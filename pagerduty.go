package pagerduty

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

// inspired by https://github.com/mdcollins05/phoneduty

// PagerDuty the client
type PagerDuty struct {
	SubDomain string
	Token     string
}

// GetScheduleUsersInput is stuff
type GetScheduleUsersInput struct {
	ScheduleID string
	Since      *time.Time
	Until      *time.Time
}

// User ...
type User struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Email          string `json:"email"`
	TimeZone       string `json:"time_zone"`
	Color          string `json:"color"`
	Role           string `json:"role"`
	AvatarURL      string `json:"avatar_url"`
	UserURL        string `json:"user_url"`
	InvitationSent bool   `json:"invitation_sent"`
	// "billed":true,
	// "description":"",
	// "user_url":"/users/PZX9XQ3",
	// "marketing_opt_out":false,
	ContactMethods []ContactMethod `json:"contact_methods"`
}

// GetScheduleUsersOutput ...
type GetScheduleUsersOutput struct {
	Users []User `json:"users"`
}

// ContactMethod ...
type ContactMethod struct {
	ID          string `json:"id"`
	Label       string `json:"label"`
	Address     string `json:"address"`
	Type        string `json:"type"`
	CountryCode int    `json:"country_code"`
	PhoneNumber string `json:"phone_number"`
}

// GetUserOutput ...
type GetUserOutput struct {
	User User `json:"user"`
}

// GetScheduleUsers ...
func (p *PagerDuty) GetScheduleUsers(in *GetScheduleUsersInput) (*GetScheduleUsersOutput, error) {
	v := url.Values{}
	if in.Since != nil {
		v.Add("since", in.Since.Format(time.RFC3339))
	}
	if in.Until != nil {
		v.Add("until", in.Until.Format(time.RFC3339))
	}
	urlString := fmt.Sprintf("https://%s/api/v1/schedules/%s/users?%s", p.SubDomain, in.ScheduleID, v.Encode())
	req, err := p.makeRequest("GET", urlString)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var out GetScheduleUsersOutput
	err = json.Unmarshal(buf, &out)
	if err != nil {
		return nil, err
	}

	return &out, nil
}

// GetUser ...
func (p *PagerDuty) GetUser(id string) (*GetUserOutput, error) {
	v := url.Values{}
	v.Add("include[]", "contact_methods")

	urlString := fmt.Sprintf("https://%s/api/v1/users/%s?%s", p.SubDomain, id, v.Encode())
	req, err := p.makeRequest("GET", urlString)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var out GetUserOutput
	err = json.Unmarshal(buf, &out)
	if err != nil {
		return nil, err
	}

	return &out, nil
}

func (p *PagerDuty) makeRequest(method, urlString string) (*http.Request, error) {
	req, err := http.NewRequest(method, urlString, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Token token=%s", p.Token))
	return req, nil
}

// TimezoneMap is because PagerDuty timezones only make sense to them
var TimezoneMap = map[string]string{
	"International Date Line West": "Pacific/Midway",
	"Midway Island":                "Pacific/Midway",
	"American Samoa":               "Pacific/Pago_Pago",
	"Hawaii":                       "Pacific/Honolulu",
	"Alaska":                       "America/Juneau",
	"Pacific Time (US & Canada)":   "America/Los_Angeles",
	"Tijuana":                      "America/Tijuana",
	"Mountain Time (US & Canada)":  "America/Denver",
	"Arizona":                      "America/Phoenix",
	"Chihuahua":                    "America/Chihuahua",
	"Mazatlan":                     "America/Mazatlan",
	"Central Time (US & Canada)":   "America/Chicago",
	"Saskatchewan":                 "America/Regina",
	"Guadalajara":                  "America/Mexico_City",
	"Mexico City":                  "America/Mexico_City",
	"Monterrey":                    "America/Monterrey",
	"Central America":              "America/Guatemala",
	"Eastern Time (US & Canada)":   "America/New_York",
	"Indiana (East)":               "America/Indiana/Indianapolis",
	"Bogota":                       "America/Bogota",
	"Lima":                         "America/Lima",
	"Quito":                        "America/Lima",
	"Atlantic Time (Canada)":       "America/Halifax",
	"Caracas":                      "America/Caracas",
	"La Paz":                       "America/La_Paz",
	"Santiago":                     "America/Santiago",
	"Newfoundland":                 "America/St_Johns",
	"Brasilia":                     "America/Sao_Paulo",
	"Buenos Aires":                 "America/Argentina/Buenos_Aires",
	"Montevideo":                   "America/Montevideo",
	"Georgetown":                   "America/Guyana",
	"Greenland":                    "America/Godthab",
	"Mid-Atlantic":                 "Atlantic/South_Georgia",
	"Azores":                       "Atlantic/Azores",
	"Cape Verde Is.":               "Atlantic/Cape_Verde",
	"Dublin":                       "Europe/Dublin",
	"Edinburgh":                    "Europe/London",
	"Lisbon":                       "Europe/Lisbon",
	"London":                       "Europe/London",
	"Casablanca":                   "Africa/Casablanca",
	"Monrovia":                     "Africa/Monrovia",
	"UTC":                          "Etc/UTC",
	"Belgrade":                     "Europe/Belgrade",
	"Bratislava":                   "Europe/Bratislava",
	"Budapest":                     "Europe/Budapest",
	"Ljubljana":                    "Europe/Ljubljana",
	"Prague":                       "Europe/Prague",
	"Sarajevo":                     "Europe/Sarajevo",
	"Skopje":                       "Europe/Skopje",
	"Warsaw":                       "Europe/Warsaw",
	"Zagreb":                       "Europe/Zagreb",
	"Brussels":                     "Europe/Brussels",
	"Copenhagen":                   "Europe/Copenhagen",
	"Madrid":                       "Europe/Madrid",
	"Paris":                        "Europe/Paris",
	"Amsterdam":                    "Europe/Amsterdam",
	"Berlin":                       "Europe/Berlin",
	"Bern":                         "Europe/Berlin",
	"Rome":                         "Europe/Rome",
	"Stockholm":                    "Europe/Stockholm",
	"Vienna":                       "Europe/Vienna",
	"West Central Africa":          "Africa/Algiers",
	"Bucharest":                    "Europe/Bucharest",
	"Cairo":                        "Africa/Cairo",
	"Helsinki":                     "Europe/Helsinki",
	"Kyiv":                         "Europe/Kiev",
	"Riga":                         "Europe/Riga",
	"Sofia":                        "Europe/Sofia",
	"Tallinn":                      "Europe/Tallinn",
	"Vilnius":                      "Europe/Vilnius",
	"Athens":                       "Europe/Athens",
	"Istanbul":                     "Europe/Istanbul",
	"Minsk":                        "Europe/Minsk",
	"Jerusalem":                    "Asia/Jerusalem",
	"Harare":                       "Africa/Harare",
	"Pretoria":                     "Africa/Johannesburg",
	"Kaliningrad":                  "Europe/Kaliningrad",
	"Moscow":                       "Europe/Moscow",
	"St. Petersburg":               "Europe/Moscow",
	"Volgograd":                    "Europe/Volgograd",
	"Samara":                       "Europe/Samara",
	"Kuwait":                       "Asia/Kuwait",
	"Riyadh":                       "Asia/Riyadh",
	"Nairobi":                      "Africa/Nairobi",
	"Baghdad":                      "Asia/Baghdad",
	"Tehran":                       "Asia/Tehran",
	"Abu Dhabi":                    "Asia/Muscat",
	"Muscat":                       "Asia/Muscat",
	"Baku":                         "Asia/Baku",
	"Tbilisi":                      "Asia/Tbilisi",
	"Yerevan":                      "Asia/Yerevan",
	"Kabul":                        "Asia/Kabul",
	"Ekaterinburg":                 "Asia/Yekaterinburg",
	"Islamabad":                    "Asia/Karachi",
	"Karachi":                      "Asia/Karachi",
	"Tashkent":                     "Asia/Tashkent",
	"Chennai":                      "Asia/Kolkata",
	"Kolkata":                      "Asia/Kolkata",
	"Mumbai":                       "Asia/Kolkata",
	"New Delhi":                    "Asia/Kolkata",
	"Kathmandu":                    "Asia/Kathmandu",
	"Astana":                       "Asia/Dhaka",
	"Dhaka":                        "Asia/Dhaka",
	"Sri Jayawardenepura":          "Asia/Colombo",
	"Almaty":                       "Asia/Almaty",
	"Novosibirsk":                  "Asia/Novosibirsk",
	"Rangoon":                      "Asia/Rangoon",
	"Bangkok":                      "Asia/Bangkok",
	"Hanoi":                        "Asia/Bangkok",
	"Jakarta":                      "Asia/Jakarta",
	"Krasnoyarsk":                  "Asia/Krasnoyarsk",
	"Beijing":                      "Asia/Shanghai",
	"Chongqing":                    "Asia/Chongqing",
	"Hong Kong":                    "Asia/Hong_Kong",
	"Urumqi":                       "Asia/Urumqi",
	"Kuala Lumpur":                 "Asia/Kuala_Lumpur",
	"Singapore":                    "Asia/Singapore",
	"Taipei":                       "Asia/Taipei",
	"Perth":                        "Australia/Perth",
	"Irkutsk":                      "Asia/Irkutsk",
	"Ulaanbaatar":                  "Asia/Ulaanbaatar",
	"Seoul":                        "Asia/Seoul",
	"Osaka":                        "Asia/Tokyo",
	"Sapporo":                      "Asia/Tokyo",
	"Tokyo":                        "Asia/Tokyo",
	"Yakutsk":                      "Asia/Yakutsk",
	"Darwin":                       "Australia/Darwin",
	"Adelaide":                     "Australia/Adelaide",
	"Canberra":                     "Australia/Melbourne",
	"Melbourne":                    "Australia/Melbourne",
	"Sydney":                       "Australia/Sydney",
	"Brisbane":                     "Australia/Brisbane",
	"Hobart":                       "Australia/Hobart",
	"Vladivostok":                  "Asia/Vladivostok",
	"Guam":                         "Pacific/Guam",
	"Port Moresby":                 "Pacific/Port_Moresby",
	"Magadan":                      "Asia/Magadan",
	"Srednekolymsk":                "Asia/Srednekolymsk",
	"Solomon Is.":                  "Pacific/Guadalcanal",
	"New Caledonia":                "Pacific/Noumea",
	"Fiji":                         "Pacific/Fiji",
	"Kamchatka":                    "Asia/Kamchatka",
	"Marshall Is.":                 "Pacific/Majuro",
	"Auckland":                     "Pacific/Auckland",
	"Wellington":                   "Pacific/Auckland",
	"Nuku'alofa":                   "Pacific/Tongatapu",
	"Tokelau Is.":                  "Pacific/Fakaofo",
	"Chatham Is.":                  "Pacific/Chatham",
	"Samoa":                        "Pacific/Apia",
}
