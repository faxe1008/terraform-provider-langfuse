package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Client is a Langfuse API client using the Admin API key.
type Client struct {
	baseURL    string
	adminKey   string
	httpClient *http.Client
}

// NewClient creates a new Langfuse Client with baseURL and adminKey.
func NewClient(baseURL, adminKey string) *Client {
	return &Client{
		baseURL:    baseURL,
		adminKey:   adminKey,
		httpClient: &http.Client{},
	}
}

// Organization represents a Langfuse organization.
type Organization struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Project represents a Langfuse project.
type Project struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	OrganizationID string `json:"organizationId"`
	PublicKey      string `json:"publicKey"`
	SecretKey      string `json:"secretKey"`
}

// CreateOrganization calls POST /api/admin/organizations.
func (c *Client) CreateOrganization(ctx context.Context, name string) (*Organization, error) {
	url := fmt.Sprintf("%s/api/admin/organizations", c.baseURL)
	body := map[string]string{"name": name}
	data, _ := json.Marshal(body)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.adminKey)
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		b, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("create organization failed: %s", string(b))
	}
	var org Organization
	if err := json.NewDecoder(resp.Body).Decode(&org); err != nil {
		return nil, err
	}
	return &org, nil
}

// GetOrganization calls GET /api/admin/organizations/{orgId}.
func (c *Client) GetOrganization(ctx context.Context, orgID string) (*Organization, error) {
	url := fmt.Sprintf("%s/api/admin/organizations/%s", c.baseURL, orgID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.adminKey)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("organization %s not found", orgID)
	}
	if resp.StatusCode >= 300 {
		b, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("get organization failed: %s", string(b))
	}
	var org Organization
	if err := json.NewDecoder(resp.Body).Decode(&org); err != nil {
		return nil, err
	}
	return &org, nil
}

// UpdateOrganization calls PUT /api/admin/organizations/{orgId}.
func (c *Client) UpdateOrganization(ctx context.Context, orgID, name string) (*Organization, error) {
	url := fmt.Sprintf("%s/api/admin/organizations/%s", c.baseURL, orgID)
	body := map[string]string{"name": name}
	data, _ := json.Marshal(body)
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.adminKey)
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		b, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("update organization failed: %s", string(b))
	}
	var org Organization
	if err := json.NewDecoder(resp.Body).Decode(&org); err != nil {
		return nil, err
	}
	return &org, nil
}

// DeleteOrganization calls DELETE /api/admin/organizations/{orgId}.
func (c *Client) DeleteOrganization(ctx context.Context, orgID string) error {
	url := fmt.Sprintf("%s/api/admin/organizations/%s", c.baseURL, orgID)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+c.adminKey)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		b, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("delete organization failed: %s", string(b))
	}
	return nil
}

// CreateProject calls POST /api/admin/organizations/{orgId}/projects.
func (c *Client) CreateProject(ctx context.Context, orgID, name string) (*Project, error) {
	url := fmt.Sprintf("%s/api/admin/organizations/%s/projects", c.baseURL, orgID)
	body := map[string]string{"name": name}
	data, _ := json.Marshal(body)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.adminKey)
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		b, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("create project failed: %s", string(b))
	}
	var proj Project
	if err := json.NewDecoder(resp.Body).Decode(&proj); err != nil {
		return nil, err
	}
	return &proj, nil
}

// GetProject calls GET /api/admin/organizations/{orgId}/projects/{projId}.
func (c *Client) GetProject(ctx context.Context, orgID, projID string) (*Project, error) {
	url := fmt.Sprintf("%s/api/admin/organizations/%s/projects/%s", c.baseURL, orgID, projID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.adminKey)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("project %s not found", projID)
	}
	if resp.StatusCode >= 300 {
		b, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("get project failed: %s", string(b))
	}
	var proj Project
	if err := json.NewDecoder(resp.Body).Decode(&proj); err != nil {
		return nil, err
	}
	return &proj, nil
}

// UpdateProject calls PUT /api/admin/organizations/{orgId}/projects/{projId}.
func (c *Client) UpdateProject(ctx context.Context, orgID, projID, name string) (*Project, error) {
	url := fmt.Sprintf("%s/api/admin/organizations/%s/projects/%s", c.baseURL, orgID, projID)
	body := map[string]string{"name": name}
	data, _ := json.Marshal(body)
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.adminKey)
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		b, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("update project failed: %s", string(b))
	}
	var proj Project
	if err := json.NewDecoder(resp.Body).Decode(&proj); err != nil {
		return nil, err
	}
	return &proj, nil
}

// DeleteProject calls DELETE /api/admin/organizations/{orgId}/projects/{projId}.
func (c *Client) DeleteProject(ctx context.Context, orgID, projID string) error {
	url := fmt.Sprintf("%s/api/admin/organizations/%s/projects/%s", c.baseURL, orgID, projID)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+c.adminKey)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		b, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("delete project failed: %s", string(b))
	}
	return nil
}
