package confluence

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (a *API) GetContents() ([]Content, error) {
	return []Content{}, nil
}

func (a *API) GetContentById(id string) (*Content, error) {
	var content Content
	resp, err := a.requestAPI("GET", fmt.Sprintf("/content/%v", id), []byte(`{}`))
	if err != nil {
		return nil, fmt.Errorf("GetContentById calls to a.requestAPI and returns an error: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		var msg string
		switch resp.StatusCode {
		case http.StatusBadRequest:
			msg = "Bad request, sub-expansions limit exceeds"
		case http.StatusUnauthorized:
			msg = "Authentication credentials are incorrect or missing from the request"
		case http.StatusForbidden:
			msg = "User does not have correct permission to read this content"
		case http.StatusNotFound:
			msg = "User does not have permission to view the requested content"
		default:
			msg = fmt.Sprintf("Invalid Status Code: %v", resp.StatusCode)
		}
		return nil, fmt.Errorf("GetContentById gets error: %v", msg)
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("GetContentById calls io.ReadAll returns an error: %w", err)
	}
	err = json.Unmarshal(b, &content)
	if err != nil {
		return nil, fmt.Errorf("GetContentById calls json.Unmarshal returns an error: %w", err)
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
		var msg string
		switch resp.StatusCode {
		case http.StatusBadRequest:
			msg = "Bad request could be either content type invalid or sub-expansions limit exceeds"
		case http.StatusUnauthorized:
			msg = "Authentication credentials are incorrect or missing from the request"
		case http.StatusForbidden:
			msg = `Forbidden could be either space that the content is being created in does not exist or
			 requesting user does not have permission to create content in it`
		case http.StatusNotFound:
			msg = "The parent content is invalid -- user does not have permission or it does not exist"
		case 413:
			msg = "Request is too large. Requests for this resource can be at most 5242880 bytes"
		default:
			msg = fmt.Sprintf("Invalid Status Code: %v", resp.StatusCode)
		}
		return fmt.Errorf("CreateContent gets error:, %v, message: %s", msg, string(b))
	}
	json.Unmarshal(b, c)

	return nil
}

func (a *API) UpdateContent(c *Content) error {
	content, err := a.GetContentById(c.Id)
	if err != nil {
		return fmt.Errorf("UpdateContent calls a.GetContentById and returns an error: %w", err)
	}
	c.Version = &Version{
		Number: content.Version.Number + 1,
	}
	body, err := json.Marshal(c)
	if err != nil {
		return fmt.Errorf("UpdateContent calls json.Marshal and returns an error: %w", err)
	}
	resp, err := a.requestAPI("PUT", fmt.Sprintf("/content/%s", c.Id), body)
	if err != nil {
		return fmt.Errorf("UpdateContent calls a.requestAPI and returns an error: %w", err)
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("UpdateContent calls io.ReadAll and returns an error: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		var msg string
		switch resp.StatusCode {
		case http.StatusBadRequest:
			msg = "Bad request could be either request body is missing required parameters (version, type, title) or  type property has been set incorrectly"
		case http.StatusUnauthorized:
			msg = "Authentication credentials are incorrect or missing from the request"
		case http.StatusForbidden:
			msg = "The calling user can not update the content with specified id"
		case http.StatusNotFound:
			msg = "A draft with current content cannot be found when setting the status parameter to draft and the content status is current"
		case http.StatusConflict:
			msg = "Conflict could be either the page is a draft (draft pages cannot be updated) or the version property has not been set correctly for the content type"
		case http.StatusNotImplemented:
			msg = `Indicates that the server does not support the functionality needed to serve the request. For example,
			trying to publish edit draft and passing content status as historical
			`
		default:
			msg = fmt.Sprintf("UpdateContent invalid status code: %v", resp.StatusCode)
		}
		return fmt.Errorf("UpdateContent gets error:, %v, message: %s", msg, string(b))
	}

	return nil
}

func (a *API) DeleteContent(id string) error {
	_, err := a.GetContentById(id)
	if err != nil {
		return fmt.Errorf("DeleteContent calls a.GetContentById and returns an error: %w", err)
	}
	resp, err := a.requestAPI("DELETE", fmt.Sprintf("content/%s", id), []byte(``))
	if err != nil {
		return fmt.Errorf("DeleteContent calls a.requestAPI and returns an error: %w", err)
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("DeleteContent calls io.ReadAll and returns an error:%w", err)
	}
	if resp.StatusCode != http.StatusNoContent {
		var msg string
		switch resp.StatusCode {
		case http.StatusBadRequest:
			msg = "The content id is invalid"
		case http.StatusUnauthorized:
			msg = "Authentication credentials are incorrect or missing from the request"
		case http.StatusForbidden:
			msg = "The calling user can not delete the content with specified id"
		case http.StatusNotFound:
			msg = "Not Found could be either there is no content with the given ID or the requesting user does not have permission to trash or purge the content"
		default:
			msg = fmt.Sprintf("DeleteContent invalid status code: %v", resp.StatusCode)
		}
		return fmt.Errorf("DeleteContent gets error:, %v, message: %s", msg, string(b))
	}
	return nil
}
