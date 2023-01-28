package confluence

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (a *API) GetSpace(key string) (*Space, error) {
	var space Space
	resp, err := a.requestAPI(http.MethodGet, fmt.Sprintf("/space/%v", key), []byte(`{}`))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		var msg string
		switch resp.StatusCode {
		case http.StatusUnauthorized:
			msg = "Authentication credentials are incorrect or missing from the request"
		case http.StatusNotFound:
			msg = `Not Found, could be either there is no space with the given key or calling user does not have permission to view the space`
		default:
			msg = fmt.Sprintf("Invalid Status Code: %v", resp.StatusCode)
		}
		return nil, fmt.Errorf("GetSpace gets error: %v, message: %v", msg, string(body))
	}

	err = json.Unmarshal(body, &space)
	if err != nil {
		return nil, fmt.Errorf("Error Marshalling Space struct: %v", err)
	}
	return &space, nil
}

func (a *API) CreateSpace(space *Space) error {
	body, err := json.Marshal(space)
	if err != nil {
		return fmt.Errorf("CreateSpace calls json.Marshal and returns an error: %w", err)
	}

	resp, err := a.requestAPI(http.MethodPost, "/space", body)
	if err != nil {
		return fmt.Errorf("CreateSpace calls a.requestAPI and returns an error: %w", err)
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("CreateSpace calls io.ReadAll and returns an error %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var msg string
		switch resp.StatusCode {
		case http.StatusBadRequest:
			msg = "Bad request, could be either an invalid request or the space already exists"
		case http.StatusUnauthorized:
			msg = "Authentication credentials are incorrect or missing from the request"
		case http.StatusForbidden:
			msg = "The calling user does not have permission to create a space"
		default:
			msg = fmt.Sprintf("Invalid Status Code: %v", resp.StatusCode)
		}
		return fmt.Errorf("CreateSpace gets error: %v, message: %s", msg, string(b))
	}
	err = json.Unmarshal(b, space)
	if err != nil {
		return fmt.Errorf("CreateSpace calls json.Unmarshal and returns an error %w", err)
	}

	return nil
}

func (a *API) UpdateSpace(s *Space) error {
	space, err := a.GetSpace(s.Key)
	if err != nil {
		return fmt.Errorf("UpdateSpace calls a.GetSpace and returns an error %w", err)
	}
	// Update space struct
	s.Id = space.Id
	body, err := json.Marshal(s)
	if err != nil {
		return fmt.Errorf("UpdateSpace calls json.Marshal and returns an error %w", err)
	}

	resp, err := a.requestAPI(http.MethodPut, fmt.Sprintf("/space/%s", s.Key), body)
	if err != nil {
		return fmt.Errorf("UpdateSpace calls a.requestAPI and returns an error %w", err)
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("UpdateSpace calls io.ReadAll and returns an error %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var msg string
		switch resp.StatusCode {
		case http.StatusUnauthorized:
			msg = "Authentication credentials are incorrect or missing from the request"
		case http.StatusNotFound:
			msg = "Not found, could be either there is no space with the given key or the calling user does not have permission to update the space"
		default:
			msg = fmt.Sprintf("Invalid Status Code %v", resp.StatusCode)
		}
		return fmt.Errorf("UpdateSpace gets error: %v, message: %v", msg, string(b))
	}

	return nil
}

func (a *API) DeleteSpace(key string) error {
	space, err := a.GetSpace(key)
	if err != nil {
		return fmt.Errorf("DeleteSpace calls a.GetSpace and returns an error %w", err)
	}
	if space == nil {
		return fmt.Errorf("DeleteSpace calls a.GetSpace and returns an space empty object")
	}

	resp, err := a.requestAPI(http.MethodDelete, fmt.Sprintf("/space/%s", key), []byte{})
	if err != nil {
		return fmt.Errorf("DeleteSpace calls a.requestAPI and returns an error %w", err)
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("DeleteSpace calls io.ReadAll and returns an error %w", err)
	}

	if resp.StatusCode != http.StatusAccepted {
		var msg string
		switch resp.StatusCode {
		case http.StatusUnauthorized:
			msg = "Authentication credentials are incorrect or missing from the request"
		case http.StatusNotFound:
			msg = "Not found, could be either there is no space with the given key or the calling user does not have permission to delete the space"
		default:
			msg = fmt.Sprintf("Invalid Status Code %v", resp.StatusCode)
		}
		return fmt.Errorf("DeleteSpace gets error: %v, message: %v", msg, string(b))
	}

	return nil
}
