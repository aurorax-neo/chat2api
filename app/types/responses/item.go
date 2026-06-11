package responses

type OutputItem struct {
	ID            string        `json:"id"`
	Type          string        `json:"type"`
	Status        string        `json:"status"`
	Role          string        `json:"role,omitempty"`
	Content       []ContentPart `json:"content,omitempty"`
	CallID        string        `json:"call_id,omitempty"`
	Name          string        `json:"name,omitempty"`
	Arguments     string        `json:"arguments,omitempty"`
	Result        string        `json:"result,omitempty"`
	RevisedPrompt string        `json:"revised_prompt,omitempty"`
}

type ContentPart struct {
	Type        string        `json:"type"`
	Text        string        `json:"text"`
	Annotations []interface{} `json:"annotations"`
}

func ImageOutputItem(id string, result string, prompt string) OutputItem {
	return OutputItem{
		ID:            id,
		Type:          "image_generation_call",
		Status:        "completed",
		Result:        result,
		RevisedPrompt: prompt,
	}
}

func TextOutputItem(id string, text string, status string) OutputItem {
	return OutputItem{
		ID:     id,
		Type:   "message",
		Status: status,
		Role:   "assistant",
		Content: []ContentPart{{
			Type:        "output_text",
			Text:        text,
			Annotations: []interface{}{},
		}},
	}
}

func TextOutputItemStarted(id string) OutputItem {
	return OutputItem{
		ID:     id,
		Type:   "message",
		Status: "in_progress",
		Role:   "assistant",
	}
}

func FunctionCallOutputItem(id string, callID string, name string, arguments string, status string) OutputItem {
	return OutputItem{
		ID:        id,
		Type:      "function_call",
		Status:    status,
		CallID:    callID,
		Name:      name,
		Arguments: arguments,
	}
}
