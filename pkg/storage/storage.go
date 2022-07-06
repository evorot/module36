package storage

// Публикация, получаемая из RSS.
type Post struct {
	ID      int    // номер записи
	Title   string // заголовок публикации
	Content string // содержание публикации
	PubTime int64  // время публикации
	Link    string // ссылка на источник
}

// Interface задаёт контракт на работу с БД.
type Interface interface {
	Posts(quantity int) ([]Post, error) // получение n-ого кол-ва публикаций
	AddPost(Post) error                 // создание новой публикации
	PostsMany([]Post) error             // создание n-ого кол-ва публикаций
	UpdatePost(Post) error              // обновление публикации
	DeletePost(Post) error              // удаление публикации по ID
}
