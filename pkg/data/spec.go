// Code generated by schema-generate. DO NOT EDIT.

package data


// Root
type Root interface{}

// ClientCheck 
type ClientCheck struct {
  Count float64 `json:"count,omitempty"`
  Description string `json:"description"`
  Group string `json:"group,omitempty"`
  Name string `json:"name"`
  Parent string `json:"parent,omitempty"`
  Result int `json:"result"`
  Time float64 `json:"time,omitempty"`
  Value string `json:"value"`
  Visual string `json:"visual,omitempty"`
}

// ClientMeta 
type ClientMeta struct {
  Host string `json:"host"`
  Notifications *ClientNotifications `json:"notifications,omitempty"`
  Result int `json:"result"`
  Tags []string `json:"tags,omitempty"`
  Time float64 `json:"time,omitempty"`
  Ttl int `json:"ttl,omitempty"`
  Website string `json:"website"`
}

// ClientNotifications 
type ClientNotifications struct {
  Email []interface{} `json:"email,omitempty"`
  Slack map[string]string `json:"slack,omitempty"`
}

// ClientResponse 
type ClientResponse struct {
  Checks []*ClientCheck `json:"checks"`
  Meta *ClientMeta `json:"meta"`
}
