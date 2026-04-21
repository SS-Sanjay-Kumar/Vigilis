package services

type Album struct {
    ID     uint16 `json:"id"`
    Title  string `json:"title"`
    Artist string `json:"artist"`
    Price  string `json:"price"`
}

var Albums = []Album{
    {ID: 1, Title: "Bully", Artist: "ye", Price: "12 USD"},
    {ID: 2, Title: "Damn", Artist: "Kendrick Lamar", Price: "26 USD"},
    {ID: 3, Title: "Beku", Artist: "Bekandi", Price: "20 USD"},
}