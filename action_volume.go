package sequence

import (
	"fmt"

	"github.com/frantjc/sequence/github/actions"
)

func GetActionVolumeName(action actions.Reference) string {
	return fmt.Sprintf("sequence-actions-%s-%s", action.Owner(), action.Repository()) // TODO , action.Version())
}
