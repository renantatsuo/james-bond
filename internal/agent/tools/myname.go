package tools

var MyName = Tool{
	Name:        "getMyName",
	Description: "Get the user's name.",
	Args:        nil,
	Fn:          MyNameFn,
}

func MyNameFn(_input []byte) (string, error) {
	return "Renan Tatsuo", nil
}
