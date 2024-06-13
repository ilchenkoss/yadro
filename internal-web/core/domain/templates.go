package domain

type HomeTemplateData struct {
	Logged bool
	Login  string
}

type LoginTemplateData struct {
	LoginErr string
	Logged   bool
}

type ComicsTemplateData struct {
	Logged    bool
	Comics    []Comics
	SearchErr string
}
