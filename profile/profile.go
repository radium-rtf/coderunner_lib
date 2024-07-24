package profile

type Profile struct {
	Name  string
	Image string
}

func NewProfile(name, image string) Profile {
	return Profile{
		Name:  name,
		Image: image,
	}
}
