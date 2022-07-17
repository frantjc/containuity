package sequence

import (
	"encoding/json"
	"fmt"

	"github.com/frantjc/sequence/pkg/github/actions"
)

const (
	// ActionMetadataKey is the key in a Step's
	// Step_Out.Metadata map that holds the json
	// encoding of the action that the step cloned.
	ActionMetadataKey = "__sqnc_action_metadata"
)

func (o *Step_Out) GetActionMetadata() (*actions.Metadata, error) {
	if actionMetadataJSON, ok := o.Metadata[ActionMetadataKey]; ok {
		actionMetadata := &actions.Metadata{}
		if err := json.Unmarshal([]byte(actionMetadataJSON), actionMetadata); err != nil {
			return nil, err
		}

		return actionMetadata, nil
	}

	return nil, fmt.Errorf("action metadata not found")
}
