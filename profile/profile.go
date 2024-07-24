package profile

type (
	Profile struct {
		Name  Name
		Image Image
	}

	Name  string
	Image string
)

func NewProfile(name Name, image Image) Profile {
	return Profile{
		Name:  name,
		Image: image,
	}
}
