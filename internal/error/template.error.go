package errorTemplate

import "fmt"

var ErrInvalidInputs = fmt.Errorf("invalid inputs")
var ErrEmailPasswordIncorrect = fmt.Errorf("email or password are incorrect")
var ErrDoubleSubmittedTestimony = fmt.Errorf("you have already submitted a testimony")
var ErrEmailExist = fmt.Errorf("Email is Exist")