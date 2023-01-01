package entity

// Page is entity for page.
type Page struct {
	ID      int64
	Title   string
	Content string
}

// PageImage is entity for page image.
type PageImage struct {
	ID    int64
	Title string
	Image string
}

// CategoryMember is entity for category member.
type CategoryMember struct {
	ID    int64
	Title string
}

// PageCategory is entity for page category.
type PageCategory struct {
	Title string
}
