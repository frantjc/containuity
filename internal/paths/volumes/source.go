package volumes

import (
	"fmt"

	"github.com/frantjc/sequence/github/actions"
)

func GetActionSource(action *actions.Reference) string {
	return fmt.Sprintf("sequence-actions-%s-%s", action.Owner, action.Repository) // TODO , action.Version())
}
