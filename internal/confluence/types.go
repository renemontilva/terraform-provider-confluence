package confluence

type Content struct {
	Id        string    `json:"id,omitempty"`
	Type      string    `json:"type,omitempty"`
	Title     string    `json:"title,omitempty"`
	Space     *Space    `json:"space,omitempty"`
	Status    string    `json:"status,omitempty"`
	Ancestors []Content `json:"ancestors,omitempty"`
	Body      Body      `json:"body,omitempty"`
	Version   *Version  `json:"version,omitempty"`
}

type Space struct {
	Id          uint              `json:"id,omitempty"`
	Key         string            `json:"key,omitempty"`
	Name        string            `json:"name,omitempty"`
	Type        string            `json:"type,omitempty"`
	Status      string            `json:"status,omitempty"`
	Description *SpaceDescription `json:"description,omitempty"`
}

type Storage struct {
	Value          string `json:"value,omitempty"`
	Representation string `json:"representation,omitempty"`
}

type Body struct {
	Storage Storage `json:"storage,omitempty"`
}

type Version struct {
	Number int `json:"number,omitempty"`
}

type SpaceDescription struct {
	Plain Plain `json:"plain,omitempty"`
}

type Plain struct {
	Value          string `json:"value,omitempty"`
	Representation string `json:"representation,omitempty"`
}
