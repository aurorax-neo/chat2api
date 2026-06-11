package service

import (
	"chat2api/app/types/completions"
	"chat2api/app/types/responses"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func streamResponsesTextEvents(c *gin.Context, model string, resp *http.Response) (*chatResult, error) {
	c.Header("Content-Type", "text/event-stream")
	responseID := responses.ResponseID()
	itemID := responses.MessageID()
	created := time.Now().Unix()
	if _, err := c.Writer.WriteString(responses.SSE(responses.CreatedEvent(responseID, model, created))); err != nil {
		return nil, err
	}
	if _, err := c.Writer.WriteString(responses.SSE(responses.InProgressEvent(responseID, model, created))); err != nil {
		return nil, err
	}
	item := responses.TextOutputItemStarted(itemID)
	if _, err := c.Writer.WriteString(responses.SSE(responses.OutputItemAddedEvent(responseID, 0, &item))); err != nil {
		return nil, err
	}
	part := responses.ContentPart{Type: "output_text", Text: "", Annotations: []interface{}{}}
	if _, err := c.Writer.WriteString(responses.SSE(responses.ContentPartAddedEvent(responseID, itemID, 0, 0, part))); err != nil {
		return nil, err
	}
	c.Writer.Flush()
	result, err := handleChatStream(resp, func(event chatStreamEvent) error {
		if event.Delta == "" {
			return nil
		}
		if _, err := c.Writer.WriteString(responses.SSE(responses.Event{EventID: responses.EventID(), Type: "response.output_text.delta", ResponseID: responseID, ItemID: itemID, OutputIndex: 0, ContentIndex: 0, Delta: event.Delta})); err != nil {
			return err
		}
		c.Writer.Flush()
		return nil
	})
	if err != nil {
		return nil, err
	}
	if _, err := c.Writer.WriteString(responses.SSE(responses.Event{EventID: responses.EventID(), Type: "response.output_text.done", ResponseID: responseID, ItemID: itemID, OutputIndex: 0, ContentIndex: 0, Text: result.Content})); err != nil {
		return nil, err
	}
	donePart := responses.ContentPart{Type: "output_text", Text: result.Content, Annotations: []interface{}{}}
	if _, err := c.Writer.WriteString(responses.SSE(responses.ContentPartDoneEvent(responseID, itemID, 0, 0, donePart))); err != nil {
		return nil, err
	}
	completedItem := responses.TextOutputItem(itemID, result.Content, "completed")
	if _, err := c.Writer.WriteString(responses.SSE(responses.Event{EventID: responses.EventID(), Type: "response.output_item.done", ResponseID: responseID, OutputIndex: 0, Item: &completedItem})); err != nil {
		return nil, err
	}
	if _, err := c.Writer.WriteString(responses.SSE(responses.CompletedEvent(responseID, model, created, []responses.OutputItem{completedItem}))); err != nil {
		return nil, err
	}
	_, _ = c.Writer.WriteString("data: [DONE]\n\n")
	c.Writer.Flush()
	return result, nil
}

func streamResponsesFunctionCallingEvents(c *gin.Context, apiReq *completions.ApiReq, resp *http.Response) (*chatResult, error) {
	c.Header("Content-Type", "text/event-stream")
	responseID := responses.ResponseID()
	itemID := responses.MessageID()
	created := time.Now().Unix()
	if _, err := c.Writer.WriteString(responses.SSE(responses.CreatedEvent(responseID, apiReq.Model, created))); err != nil {
		return nil, err
	}
	if _, err := c.Writer.WriteString(responses.SSE(responses.InProgressEvent(responseID, apiReq.Model, created))); err != nil {
		return nil, err
	}
	c.Writer.Flush()

	detector := completions.NewStreamToolDetector(completions.ToolifyTriggerSignal)
	textItemStarted := false
	toolStreamFinished := false

	result, err := handleChatStream(resp, func(event chatStreamEvent) error {
		if detector.State() == "tool_parsing" {
			detector.AppendParsing(event.Delta)
			if !detector.HasCompleteToolBlock() {
				return nil
			}
			calls := detector.Finalize()
			if len(calls) == 0 || completions.ValidateParsedToolCalls(calls, apiReq.Tools) != nil {
				toolStreamFinished = true
				return errToolCallsStreamFinished
			}
			event.Result.ToolCalls = completions.ToolCallsFromParsed(calls, false)
			event.Result.ToolContent = completions.ToolCallPrefixText(event.Text)
			event.Result.FinishReason = "tool_calls"
			prefixItems := []responses.OutputItem(nil)
			startIndex := 0
			if textItemStarted {
				textItem, err := finishResponsesTextItem(c, responseID, itemID, 0, event.Result.ToolContent)
				if err != nil {
					return err
				}
				prefixItems = append(prefixItems, textItem)
				startIndex = 1
			}
			if err := writeResponsesToolCallEvents(c, responseID, apiReq.Model, created, event.Result.ToolCalls, prefixItems, startIndex); err != nil {
				return err
			}
			toolStreamFinished = true
			return errToolCallsStreamFinished
		}

		detected, content := detector.ProcessChunk(event.Delta)
		if content != "" {
			if !textItemStarted {
				item := responses.TextOutputItemStarted(itemID)
				if _, err := c.Writer.WriteString(responses.SSE(responses.OutputItemAddedEvent(responseID, 0, &item))); err != nil {
					return err
				}
				textItemStarted = true
			}
			if _, err := c.Writer.WriteString(responses.SSE(responses.Event{EventID: responses.EventID(), Type: "response.output_text.delta", ResponseID: responseID, ItemID: itemID, OutputIndex: 0, ContentIndex: 0, Delta: content})); err != nil {
				return err
			}
			c.Writer.Flush()
		}
		if detected {
			return nil
		}
		return nil
	})
	if err != nil && err != errToolCallsStreamFinished {
		return nil, err
	}
	if result == nil {
		result = &chatResult{}
	}
	if toolStreamFinished {
		_, _ = c.Writer.WriteString("data: [DONE]\n\n")
		c.Writer.Flush()
		return result, nil
	}
	if detector.State() == "tool_parsing" {
		if calls := detector.Finalize(); len(calls) > 0 && completions.ValidateParsedToolCalls(calls, apiReq.Tools) == nil {
			result.ToolCalls = completions.ToolCallsFromParsed(calls, false)
			result.ToolContent = completions.ToolCallPrefixText(result.Content)
			result.FinishReason = "tool_calls"
			prefixItems := []responses.OutputItem(nil)
			startIndex := 0
			if textItemStarted {
				textItem, err := finishResponsesTextItem(c, responseID, itemID, 0, result.ToolContent)
				if err != nil {
					return nil, err
				}
				prefixItems = append(prefixItems, textItem)
				startIndex = 1
			}
			if err := writeResponsesToolCallEvents(c, responseID, apiReq.Model, created, result.ToolCalls, prefixItems, startIndex); err != nil {
				return nil, err
			}
			_, _ = c.Writer.WriteString("data: [DONE]\n\n")
			c.Writer.Flush()
			return result, nil
		}
	} else if text := detector.FlushText(); text != "" {
		if !textItemStarted {
			item := responses.TextOutputItemStarted(itemID)
			if _, err := c.Writer.WriteString(responses.SSE(responses.OutputItemAddedEvent(responseID, 0, &item))); err != nil {
				return nil, err
			}
			textItemStarted = true
		}
		if _, err := c.Writer.WriteString(responses.SSE(responses.Event{EventID: responses.EventID(), Type: "response.output_text.delta", ResponseID: responseID, ItemID: itemID, OutputIndex: 0, ContentIndex: 0, Delta: text})); err != nil {
			return nil, err
		}
	}
	if textItemStarted {
		completedItem, err := finishResponsesTextItem(c, responseID, itemID, 0, result.Content)
		if err != nil {
			return nil, err
		}
		if _, err := c.Writer.WriteString(responses.SSE(responses.CompletedEvent(responseID, apiReq.Model, created, []responses.OutputItem{completedItem}))); err != nil {
			return nil, err
		}
	} else {
		if _, err := c.Writer.WriteString(responses.SSE(responses.CompletedEvent(responseID, apiReq.Model, created, []responses.OutputItem{}))); err != nil {
			return nil, err
		}
	}
	_, _ = c.Writer.WriteString("data: [DONE]\n\n")
	c.Writer.Flush()
	return result, nil
}

func finishResponsesTextItem(c *gin.Context, responseID string, itemID string, outputIndex int, text string) (responses.OutputItem, error) {
	donePart := responses.ContentPart{Type: "output_text", Text: text, Annotations: []interface{}{}}
	if _, err := c.Writer.WriteString(responses.SSE(responses.Event{EventID: responses.EventID(), Type: "response.output_text.done", ResponseID: responseID, ItemID: itemID, OutputIndex: outputIndex, ContentIndex: 0, Text: text})); err != nil {
		return responses.OutputItem{}, err
	}
	if _, err := c.Writer.WriteString(responses.SSE(responses.ContentPartDoneEvent(responseID, itemID, outputIndex, 0, donePart))); err != nil {
		return responses.OutputItem{}, err
	}
	completedItem := responses.TextOutputItem(itemID, text, "completed")
	if _, err := c.Writer.WriteString(responses.SSE(responses.Event{EventID: responses.EventID(), Type: "response.output_item.done", ResponseID: responseID, OutputIndex: outputIndex, Item: &completedItem})); err != nil {
		return responses.OutputItem{}, err
	}
	return completedItem, nil
}

func writeResponsesToolCallEvents(c *gin.Context, responseID string, model string, created int64, toolCalls []completions.ToolCall, prefixItems []responses.OutputItem, startIndex int) error {
	output := make([]responses.OutputItem, 0, len(prefixItems)+len(toolCalls))
	output = append(output, prefixItems...)
	for i, toolCall := range toolCalls {
		outputIndex := startIndex + i
		itemID := responses.MessageID()
		item := responses.FunctionCallOutputItem(itemID, toolCall.ID, toolCall.Function.Name, "", "in_progress")
		if _, err := c.Writer.WriteString(responses.SSE(responses.OutputItemAddedEvent(responseID, outputIndex, &item))); err != nil {
			return err
		}
		if _, err := c.Writer.WriteString(responses.SSE(responses.Event{EventID: responses.EventID(), Type: "response.function_call_arguments.delta", ResponseID: responseID, ItemID: itemID, OutputIndex: outputIndex, Delta: toolCall.Function.Arguments})); err != nil {
			return err
		}
		if _, err := c.Writer.WriteString(responses.SSE(responses.Event{EventID: responses.EventID(), Type: "response.function_call_arguments.done", ResponseID: responseID, ItemID: itemID, OutputIndex: outputIndex, Arguments: toolCall.Function.Arguments})); err != nil {
			return err
		}
		completed := responses.FunctionCallOutputItem(itemID, toolCall.ID, toolCall.Function.Name, toolCall.Function.Arguments, "completed")
		if _, err := c.Writer.WriteString(responses.SSE(responses.Event{EventID: responses.EventID(), Type: "response.output_item.done", ResponseID: responseID, OutputIndex: outputIndex, Item: &completed})); err != nil {
			return err
		}
		output = append(output, completed)
	}
	if _, err := c.Writer.WriteString(responses.SSE(responses.CompletedEvent(responseID, model, created, output))); err != nil {
		return err
	}
	c.Writer.Flush()
	return nil
}
