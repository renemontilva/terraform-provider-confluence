package confluence

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

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
	Id   uint   `json:"id,omitempty"`
	Key  string `json:"key,omitempty"`
	Name string `json:"name,omitempty"`
	Type string `json:"type,omitempty"`
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

func (a *API) GetContents() ([]Content, error) {
	return []Content{}, nil
}

func (a *API) GetContentById(id string) (*Content, error) {
	var content Content
	resp, err := a.requestAPI("GET", fmt.Sprintf("/content/%v", id), []byte(`{}`))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Content ID: %s, not found", id)
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(b, &content)
	if err != nil {
		return nil, err
	}
	return &content, nil
}

func (a *API) CreateContent(c *Content) error {
	body, err := json.Marshal(c)

	if err != nil {
		return err
	}
	resp, err := a.requestAPI("POST", "/content", body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Create Content API failed, message: %s", string(b))
	}
	json.Unmarshal(b, c)

	return nil
}

func (a *API) UpdateContent(c *Content) error {
	content, err := a.GetContentById(c.Id)
	if err != nil {
		return err
	}
	c.Version = &Version{
		Number: content.Version.Number + 1,
	}
	body, err := json.Marshal(c)
	if err != nil {
		return err
	}
	resp, err := a.requestAPI("PUT", fmt.Sprintf("/content/%s", c.Id), body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {

		return fmt.Errorf("Update Content API failed, message:: %v", string(b))
	}

	return nil
}

func (a *API) DeleteContent(id string) error {
	_, err := a.GetContentById(id)
	if err != nil {
		return err
	}
	resp, err := a.requestAPI("DELETE", fmt.Sprintf("content/%s", id), []byte(``))
	if err != nil {
		return fmt.Errorf("Delete Content error: %v", resp.StatusCode)
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("Delete Content API failed, message:: %v", string(b))
	}

	return nil
}
