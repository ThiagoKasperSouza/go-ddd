package dto

type CreateAgencyDTO struct {
	AgencyID       string  `json:"agency_id"`
	AgencyName     string  `json:"agency_name"`
	AgencyURL      string  `json:"agency_url"`
	AgencyTimezone string  `json:"agency_timezone"`
	AgencyLang     *string `json:"agency_lang,omitempty"`
	AgencyPhone    *string `json:"agency_phone,omitempty"`
	AgencyFareURL  *string `json:"agency_fare_url,omitempty"`
}

type UpdateAgencyDTO struct {
	AgencyName     string  `json:"agency_name"`
	AgencyURL      string  `json:"agency_url"`
	AgencyTimezone string  `json:"agency_timezone"`
	AgencyLang     *string `json:"agency_lang"`
	AgencyPhone    *string `json:"agency_phone"`
	AgencyFareURL  *string `json:"agency_fare_url"`
}

type AgencyDTO struct {
	AgencyID       string  `json:"agency_id"`
	AgencyName     string  `json:"agency_name"`
	AgencyURL      string  `json:"agency_url"`
	AgencyTimezone string  `json:"agency_timezone"`
	AgencyLang     *string `json:"agency_lang,omitempty"`
	AgencyPhone    *string `json:"agency_phone,omitempty"`
	AgencyFareURL  *string `json:"agency_fare_url,omitempty"`
}

type ListAgenciesOutputDTO struct {
	Agencies []AgencyDTO `json:"agencies"`
}