package v1alpha1

type V1alpha1 struct {
	User UserController
}

func NewV1alpha1Controllers(user UserController) V1alpha1 {
	return V1alpha1{
		User: user,
	}
}
