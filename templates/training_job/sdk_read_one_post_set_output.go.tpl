wp_err := rm.customSetOutput(&resource{ko})
if wp_err != nil{
	return &resource{ko}, wp_err
}