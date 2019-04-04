package nextbus

import (
	"encoding/xml"
	"errors"
	"strings"
)

var errNotImplemented error

func init() {
	errNotImplemented = errors.New("not implemented")
}

type response interface {
	responseError() error
}

type Response struct {
	XMLName xml.Name        `xml:"body"`
	Error   *Response_Error `xml:"Error"`
}

type Response_Error struct {
	ShouldRetry bool   `xml:"shouldRetry,attr"`
	Message     string `xml:",chardata"`
}

func (r *Response) responseError() error {
	if r.Error != nil {
		msg := strings.Trim(r.Error.Message, "\n ")
		return errors.New(msg)
	}

	return nil
}

type AgencyListResponse struct {
	Response
	Agencies []AgencyListResponse_Agency `xml:"agency"`
}

type AgencyListResponse_Agency struct {
	Tag         string `xml:"tag,attr"`
	Title       string `xml:"title,attr"`
	RegionTitle string `xml:"regionTitle,attr"`
}

type RouteListResponse struct {
	Response
	Routes []RouteListResponse_Route `xml:"route"`
}

type RouteListResponse_Route struct {
	Tag   string `xml:"tag,attr"`
	Title string `xml:"title,attr"`
}

type RouteConfigResponse struct {
	Response
	Route RouteConfigResponse_Route `xml:"route"`
}

type RouteConfigResponse_Route struct {
	Tag           string                          `xml:"tag,attr"`
	Title         string                          `xml:"title,attr"`
	Color         string                          `xml:"color,attr"`
	OppositeColor string                          `xml:"oppositeColor,attr"`
	LatMin        float32                         `xml:"latMin,attr"`
	LatMax        float32                         `xml:"latMax,attr"`
	LonMin        float32                         `xml:"lonMin,attr"`
	LonMax        float32                         `xml:"lonMax,attr"`
	Stops         []RouteConfigResponse_Stop      `xml:"stop"`
	Directions    []RouteConfigResponse_Direction `xml:"direction"`
	Paths         []RouteConfigResponse_Path      `xml:"path"`
}

type RouteConfigResponse_Direction struct {
	Tag      string                        `xml:"tag,attr"`
	Title    string                        `xml:"title,attr"`
	Name     string                        `xml:"name,attr"`
	UseForUI bool                          `xml:"useForUI,attr"`
	Stops    []RouteConfigResponse_StopTag `xml:"stop"`
}

type RouteConfigResponse_Stop struct {
	Tag    string  `xml:"tag,attr"`
	Title  string  `xml:"title,attr"`
	StopID string  `xml:"stopId,attr"`
	Lat    float32 `xml:"lat,attr"`
	Lon    float32 `xml:"lon,attr"`
}

type RouteConfigResponse_Path struct {
	Points []RouteConfigResponse_Point `xml:"point"`
}

type RouteConfigResponse_Point struct {
	Lat float32 `xml:"lat,attr"`
	Lon float32 `xml:"lon,attr"`
}

type RouteConfigResponse_StopTag struct {
	Tag string `xml:"tag,attr"`
}

type PredictionsResponse struct {
	Response
	Prediction PredictionsResponse_Predictions `xml:"predictions"`
}

type PredictionsResponse_Predictions struct {
	AgencyTitle string                          `xml:"agencyTitle,attr"`
	RouteTitle  string                          `xml:"routeTitle,attr"`
	RouteTag    string                          `xml:"routeTag,attr"`
	StopTitle   string                          `xml:"stopTitle,attr"`
	StopTag     string                          `xml:"stopTag,attr"`
	Directions  []PredictionsResponse_Direction `xml:"direction"`
	Message     PredictionsResponse_Message     `xml:"message"`
}

type PredictionsResponse_Message struct {
	Text     string `xml:"text,attr"`
	Priority string `xml:"priority,attr"`
}

type PredictionsResponse_Direction struct {
	Title       string                           `xml:"title,attr"`
	Predictions []PredictionsResponse_Prediction `xml:"prediction"`
}

type PredictionsResponse_Prediction struct {
	EpochTime         uint64 `xml:"epochTime,attr"`
	Seconds           uint32 `xml:"seconds,attr"`
	Minutes           uint32 `xml:"minutes,attr"`
	IsDeparture       bool   `xml:"isDeparture,attr"`
	DirTag            string `xml:"dirTag,attr"`
	Vehicle           string `xml:"vehicle,attr"`
	VehiclesInConsist uint32 `xml:"vehiclesInConsist,attr"`
	Block             string `xml:"block,attr"`
	TripTag           string `xml:"tripTag,attr"`
}

type ScheduleResponse struct {
	// TODO: not implemented
	Response
}

type MessagesResponse struct {
	// TODO: not implemented
	Response
}

type VehicleLocationsResponse struct {
	// TODO: not implemented
	Response
}
