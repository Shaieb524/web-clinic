package helpers

func RoleValidator(toBeCheckedRole string, desiredRole string) string {
	res := "allowed"
	if toBeCheckedRole != desiredRole {
		return "not allowed"
	}
	return res
}
