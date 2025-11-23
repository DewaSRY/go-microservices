package integrations

import (
	"DewaSRY/go-microservices/shared/types"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func GetRoute(ctx context.Context, pickup *types.Coordinate, destination *types.Coordinate) (*types.OsrmApiResponse, error) {
	url := fmt.Sprintf(
		"http://router.project-osrm.org/route/v1/driving/%f,%f;%f,%f?overview=full&geometries=geojson", pickup.Longitude, pickup.Latitude, destination.Longitude, destination.Latitude)

	res, err := http.Get(url)
	if err != nil {
		log.Print(err)
		return nil, fmt.Errorf("failed_to_parse:%v", err)
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Print(err)
		return nil, fmt.Errorf("failed_to_read:%v", err)
	}

	var routeResponse types.OsrmApiResponse
	if err := json.Unmarshal(body, &routeResponse); err != nil {
		log.Print(err)
		return nil, fmt.Errorf("failed_to_unmarshal:%v", err)
	}

	return &routeResponse, nil
}
