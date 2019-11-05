package mydramalist

type Drama struct {
	Title         string
	Country       string
	NumEpisodes   uint
	Duration      uint
	Rating        float32
	Summary       string
	NativeTitle   string
	OtherTitles   []string
	Screenwriters []string
	Directors     []string
	Genres        []string
	Tags          []string
}
