package file

type (
	File struct {
		Name    string
		Content Content
	}
)

func NewFile(name string, content Content) File {
	return File{
		Name:    name,
		Content: content,
	}
}
