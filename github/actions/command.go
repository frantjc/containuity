package actions

import "fmt"

type Command struct {
	Command    string
	Parameters map[string]string
	Value      string
}

func (c *Command) String() string {
	s := fmt.Sprintf("::%s", c.Command)

	paramSpl := " "
	numParams := len(c.Parameters)
	paramsAdded := 0
	for k, v := range c.Parameters {
		s = fmt.Sprintf("%s%s%s=%s", s, paramSpl, k, v)
		paramSpl = ","
		paramsAdded++
		if paramsAdded >= numParams {
			paramSpl = ""
		}
	}

	return fmt.Sprintf("%s::%s", s, c.Value)
}
