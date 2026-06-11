package responses

import "encoding/json"

type Event struct {
	EventID      string       `json:"event_id,omitempty"`
	Type         string       `json:"type"`
	ResponseID   string       `json:"response_id,omitempty"`
	Response     *Response    `json:"response,omitempty"`
	OutputIndex  int          `json:"output_index,omitempty"`
	ContentIndex int          `json:"content_index,omitempty"`
	ItemID       string       `json:"item_id,omitempty"`
	Item         *OutputItem  `json:"item,omitempty"`
	Part         *ContentPart `json:"part,omitempty"`
	Delta        string       `json:"delta,omitempty"`
	Text         string       `json:"text,omitempty"`
	Name         string       `json:"name,omitempty"`
	CallID       string       `json:"call_id,omitempty"`
	Arguments    string       `json:"arguments,omitempty"`
}

func CreatedEvent(responseID string, model string, created int64) Event {
	return Event{EventID: EventID(), Type: "response.created", Response: &Response{
		ID:                responseID,
		Object:            "response",
		CreatedAt:         created,
		Status:            "in_progress",
		Error:             nil,
		IncompleteDetails: nil,
		Model:             model,
		Output:            []OutputItem{},
		ParallelToolCalls: false,
	}}
}

func InProgressEvent(responseID string, model string, created int64) Event {
	return Event{EventID: EventID(), Type: "response.in_progress", Response: &Response{
		ID:                responseID,
		Object:            "response",
		CreatedAt:         created,
		Status:            "in_progress",
		Error:             nil,
		IncompleteDetails: nil,
		Model:             model,
		Output:            []OutputItem{},
		ParallelToolCalls: false,
	}}
}

func CompletedEvent(responseID string, model string, created int64, output []OutputItem) Event {
	return Event{EventID: EventID(), Type: "response.completed", Response: &Response{
		ID:                responseID,
		Object:            "response",
		CreatedAt:         created,
		Status:            "completed",
		Error:             nil,
		IncompleteDetails: nil,
		Model:             model,
		Output:            output,
		ParallelToolCalls: false,
	}}
}

func OutputItemAddedEvent(responseID string, outputIndex int, item *OutputItem) Event {
	return Event{EventID: EventID(), Type: "response.output_item.added", ResponseID: responseID, OutputIndex: outputIndex, Item: item}
}

func ContentPartAddedEvent(responseID string, itemID string, outputIndex int, contentIndex int, part ContentPart) Event {
	return Event{EventID: EventID(), Type: "response.content_part.added", ResponseID: responseID, ItemID: itemID, OutputIndex: outputIndex, ContentIndex: contentIndex, Part: &part}
}

func ContentPartDoneEvent(responseID string, itemID string, outputIndex int, contentIndex int, part ContentPart) Event {
	return Event{EventID: EventID(), Type: "response.content_part.done", ResponseID: responseID, ItemID: itemID, OutputIndex: outputIndex, ContentIndex: contentIndex, Part: &part}
}

func SSE(event Event) string {
	data, _ := json.Marshal(event)
	if event.Type == "" {
		return "data: " + string(data) + "\n\n"
	}
	return "event: " + event.Type + "\ndata: " + string(data) + "\n\n"
}
