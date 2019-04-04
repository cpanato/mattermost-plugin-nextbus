package main

import (
	"fmt"
	"strings"

	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/plugin"
)

func getCommand() *model.Command {
	return &model.Command{
		Trigger:          "nextbus",
		DisplayName:      "Next Bus",
		Description:      "Next Bus Bot",
		AutoComplete:     true,
		AutoCompleteDesc: "Available commands: agencies, routes, stops, schedule, prediction, help",
		AutoCompleteHint: "[command]",
	}
}

func (p *Plugin) ExecuteCommand(c *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	split := strings.Fields(args.Command)
	command := split[0]
	action := ""
	if len(split) > 1 {
		action = strings.TrimSpace(split[1])
	}

	if command != "/nextbus" {
		return &model.CommandResponse{}, nil
	}

	if action == "" {
		return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, "Missing command, please run `/nextbus help` to check all commands available."), nil
	}

	helpMsg := `run:
	/nextbus agencies - to list all agencies
	/nextbus prediction <AgencyID> <RouteID> <StopID> - to get the prediction for the nexxt bus for the specific route
	/nextbus routes <AgencyID> - to get the routes for the specific agency
	/nextbus stops <AgencyID> <RouteID> - to list all stops for a specific route
	/nextbus schedule - Not Implemented yet
	`

	if action == "help" {
		msg := "run:\n/nextbus xoxoxooxoxoxoox"
		return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, msg), nil
	}

	switch action {
	case "agencies":
		resp, err := p.handleAgencyList(args)
		return resp, err
	case "prediction":
		resp, err := p.handlePrediction(args)
		return resp, err
	case "routes":
		resp, err := p.handleRoute(args)
		return resp, err
	case "stops":
		resp, err := p.handleStops(args)
		return resp, err
	case "schedules":
		return getCommandResponse(model.COMMAND_RESPONSE_TYPE_IN_CHANNEL, "Not implemented yet"), nil
	case "help":
		return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, helpMsg), nil
	default:
		return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, helpMsg), nil
	}
}

func getCommandResponse(responseType, text string) *model.CommandResponse {
	return &model.CommandResponse{
		ResponseType: responseType,
		Text:         text,
		Type:         model.POST_DEFAULT,
	}
}

func (p *Plugin) sendEphemeralMessage(msg, channelId, userId string) {
	ephemeralPost := &model.Post{
		Message:   msg,
		ChannelId: channelId,
		UserId:    p.botUserID,
		Props: model.StringInterface{
			"from_webhook": "true",
		},
	}

	p.API.LogDebug("Will send an ephemeralPost", "msg", msg)

	p.API.SendEphemeralPost(userId, ephemeralPost)
}

func (p *Plugin) handleAgencyList(args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	agencies, err := p.nextBusClient.AgencyList()
	if err != nil {
		msg := fmt.Sprintf("failed to list nextBus agencies... %v", err)
		return getCommandResponse(model.COMMAND_RESPONSE_TYPE_IN_CHANNEL, msg), nil
	}

	if agencies == nil {
		return getCommandResponse(model.COMMAND_RESPONSE_TYPE_IN_CHANNEL, "No agencies found. Maybe an issue in the NextBus API"), nil
	}

	var agenciesSTR []string
	for _, agency := range agencies.Agencies {
		agencySTR := fmt.Sprintf("%s - **Region:** %s **Tag:** %s", agency.Title, agency.RegionTitle, agency.Tag)
		agenciesSTR = append(agenciesSTR, agencySTR)
	}

	post := &model.Post{
		Message: strings.Join(agenciesSTR, "\n"),
		Type:    "custom_nextbus_agencies",
	}

	if _, appErr := p.CreateBotDMPost(args.UserId, post); appErr != nil {
		return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, "Error creating the NextBus Agencies post"), nil
	}

	return &model.CommandResponse{}, nil
}

func (p *Plugin) handleRoute(args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	split := strings.Fields(args.Command)

	parameters := []string{}
	if len(split) > 2 {
		parameters = split[2:]
	}

	if len(parameters) != 1 {
		return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, "Missing agency ID"), nil
	}

	routes, err := p.nextBusClient.RouteList(parameters[0])
	if err != nil {
		msg := fmt.Sprintf("failed to list nextBus routes... %v", err)
		return getCommandResponse(model.COMMAND_RESPONSE_TYPE_IN_CHANNEL, msg), nil
	}

	if routes == nil {
		return getCommandResponse(model.COMMAND_RESPONSE_TYPE_IN_CHANNEL, "No route found. Maybe an issue in the NextBus API"), nil
	}

	routeTitle := fmt.Sprintf("Routes for Agency: %s", parameters[0])
	attachment := &model.SlackAttachment{
		Title: routeTitle,
	}
	var fields []*model.SlackAttachmentField
	for _, route := range routes.Routes {
		fields = addFields(fields, "Route", route.Title, true)
		fields = addFields(fields, "Tag", route.Tag, true)
	}

	attachment.Fields = fields
	post := &model.Post{
		Type: "custom_nextbus_routes",
	}

	model.ParseSlackAttachment(post, []*model.SlackAttachment{attachment})
	if _, appErr := p.CreateBotDMPost(args.UserId, post); appErr != nil {
		return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, "Error creating the NextBus routes post"), nil
	}

	return &model.CommandResponse{}, nil
}

func (p *Plugin) handleStops(args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	split := strings.Fields(args.Command)

	parameters := []string{}
	if len(split) > 2 {
		parameters = split[2:]
	}

	if len(parameters) < 2 {
		return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, "Missing agency and/or route"), nil
	}

	routeCfg, err := p.nextBusClient.RouteConfig(parameters[0], parameters[1])
	if err != nil {
		msg := fmt.Sprintf("failed to list nextBus routes config... %v", err)
		return getCommandResponse(model.COMMAND_RESPONSE_TYPE_IN_CHANNEL, msg), nil
	}

	if routeCfg == nil {
		return getCommandResponse(model.COMMAND_RESPONSE_TYPE_IN_CHANNEL, "No route config found. Maybe an issue in the NextBus API"), nil
	}

	routeTitle := fmt.Sprintf("Routes for Agency: %s", parameters[0])
	attachment := &model.SlackAttachment{
		Title: routeTitle,
	}
	var fields []*model.SlackAttachmentField

	var stops []string
	for _, stop := range routeCfg.Route.Stops {
		msg := fmt.Sprintf("%s - %s", stop.Title, stop.Tag)
		stops = append(stops, msg)

	}
	fields = addFields(fields, "Stops", strings.Join(stops, "\n"), false)
	attachment.Title = fmt.Sprintf("Stops for route %s", routeCfg.Route.Title)
	attachment.Fields = fields
	post := &model.Post{
		Type: "custom_nextbus_routes",
	}

	model.ParseSlackAttachment(post, []*model.SlackAttachment{attachment})
	if _, appErr := p.CreateBotDMPost(args.UserId, post); appErr != nil {
		return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, "Error creating the NextBus route config post"), nil
	}

	return &model.CommandResponse{}, nil
}

func (p *Plugin) handlePrediction(args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	split := strings.Fields(args.Command)

	parameters := []string{}
	if len(split) > 2 {
		parameters = split[2:]
	}

	if len(parameters) < 3 {
		return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, "Missing agency and/or route and/or stop"), nil
	}

	prediction, err := p.nextBusClient.PredictionsForStopTag(parameters[0], parameters[1], parameters[2])
	if err != nil {
		msg := fmt.Sprintf("failed to list nextBus agencies... %v", err)
		return getCommandResponse(model.COMMAND_RESPONSE_TYPE_IN_CHANNEL, msg), nil
	}

	if prediction == nil {
		return getCommandResponse(model.COMMAND_RESPONSE_TYPE_IN_CHANNEL, "No predicion found for now."), nil
	}

	// attachments := make([]*model.SlackAttachment, 0)
	var fields []*model.SlackAttachmentField
	var gmaps string
	fields = addFields(fields, "Route", prediction.Prediction.RouteTitle, true)
	fields = addFields(fields, "Stop", prediction.Prediction.StopTitle, true)
	for _, direction := range prediction.Prediction.Directions {
		// var fieldsPred []*model.SlackAttachmentField
		fields = addFields(fields, "Direction", direction.Title, false)
		for i, directionPrediction := range direction.Predictions {
			departure := fmt.Sprintf("%d minutes %d seconds", directionPrediction.Minutes, directionPrediction.Seconds)
			fields = addFields(fields, "Next departure", departure, false)
			if i == 2 {
				break
			}
		}
	}

	gmaps = fmt.Sprintf("[NextBus Google Maps](https://www.nextbus.com/googleMap/?a=%s&r=%s&s=%s)", parameters[0], parameters[1], parameters[2])
	fields = addFields(fields, "", gmaps, false)
	attachment := &model.SlackAttachment{
		Title:  prediction.Prediction.AgencyTitle,
		Fields: fields,
	}

	post := &model.Post{
		Type: "custom_nextbus_prediction",
	}

	model.ParseSlackAttachment(post, []*model.SlackAttachment{attachment})
	if _, appErr := p.CreateBotDMPost(args.UserId, post); appErr != nil {
		return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, "Error creating the NextBus prediction post"), nil
	}

	return &model.CommandResponse{}, nil
}

func addFields(fields []*model.SlackAttachmentField, title, msg string, short bool) []*model.SlackAttachmentField {
	return append(fields, &model.SlackAttachmentField{
		Title: title,
		Value: msg,
		Short: model.SlackCompatibleBool(short),
	})
}

func (p *Plugin) CreateBotDMPost(userID string, post *model.Post) (*model.Post, *model.AppError) {
	channel, err := p.API.GetDirectChannel(userID, p.botUserID)
	if err != nil {
		p.API.LogError("Couldn't get bot's DM channel", "user_id", userID, "err", err)
		return nil, err
	}

	post.UserId = p.botUserID
	post.ChannelId = channel.Id

	created, err := p.API.CreatePost(post)
	if err != nil {
		p.API.LogError("Couldn't send bot DM", "user_id", userID, "err", err)
		return nil, err
	}

	return created, nil
}
