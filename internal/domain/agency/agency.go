package agency

import "errors"

var (
	ErrAgencyNotFound = errors.New("agency not found")
	ErrInvalidAgency  = errors.New("invalid agency data")
)

type Agency struct {
	id 			string
	name 		string
	url 		string
	timezone	string
	lang		string
	phone		string
	fareURL		string
}

// NewAgency garante que uma nova agência seja criada com dados válidos
func NewAgency(id, name, url, timezone, lang, phone, fareURL string) (*Agency, error) {
	if id == "" || name == "" {
		return nil, ErrInvalidAgency
	}

	return &Agency{
		id:       id,
		name:     name,
		url:      url,
		timezone: timezone,
		lang:     lang,
		phone:    phone,
		fareURL:  fareURL,
	}, nil
}

func RestoreAgency(id, name, url, timezone, lang, phone, fareURL string) *Agency {
	return &Agency{
		id:       id,
		name:     name,
		url:      url,
		timezone: timezone,
		lang:     lang,
		phone:    phone,
		fareURL:  fareURL,
	}
}

// Métodos de acesso (getters) mantêm o encapsulamento
func (a *Agency) ID() string { return a.id }
func (a *Agency) Name() string { return a.name }
func (a *Agency) Url() string { return a.url }
func (a *Agency) Timezone() string { return a.timezone }
func (a *Agency) Lang() string { return a.lang }
func (a *Agency) Phone() string { return a.phone }
func (a *Agency) FareURL() string { return a.fareURL }