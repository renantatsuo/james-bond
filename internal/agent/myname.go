package agent

var ToolMyName = Tool{
	Name:        "getMyName",
	Description: "Get the user's name.",
	Args:        nil,
	Fn:          MyName,
}

func MyName(_input []byte) (string, error) {
	return "Renan Tatsuo", nil
}
