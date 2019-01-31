package capabilities

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"google.golang.org/api/compute/v1"
)

// NewGKENetworksHandler creates a new GKENetworksHandler
func NewGKENetworksHandler() *GKENetworksHandler {
	return &GKENetworksHandler{}
}

// GKENetworksHandler lists available networks in GKE
type GKENetworksHandler struct {
}

func (g *GKENetworksHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	writer.Header().Set("Content-Type", "application/json")

	var body zoneCapabilitiesRequestBody
	err := extractRequestBody(writer, req, &body)

	if err != nil {
		handleErr(writer, err)
		return
	}

	err = validateCapabilitiesRequestBody(writer, &body.capabilitiesRequestBody)

	if err != nil {
		handleErr(writer, err)
		return
	}

	if body.Credentials == "" {
		writer.WriteHeader(http.StatusBadRequest)
		handleErr(writer, fmt.Errorf("invalid credentials"))
		return
	}

	client, err := g.getServiceClient(context.Background(), body.Credentials)

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		handleErr(writer, err)
		return
	}

	result, err := client.Networks.List(body.ProjectID).Do()

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		handleErr(writer, err)
		return
	}

	serialized, err := json.Marshal(result)

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		handleErr(writer, err)
		return
	}

	writer.Write(serialized)
}

func (g *GKENetworksHandler) getServiceClient(ctx context.Context, credentialContent string) (*compute.Service, error) {
	client, err := getOAuthClient(ctx, credentialContent)

	if err != nil {
		return nil, err
	}

	service, err := compute.New(client)

	if err != nil {
		return nil, err
	}
	return service, nil
}
