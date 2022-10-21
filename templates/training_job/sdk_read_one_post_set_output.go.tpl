rm.customSetOutput(&resource{ko})
wp_err := rm.customSetWarmPoolOutput(&resource{ko})
if wp_err != nil{
	return &resource{ko}, wp_err
}