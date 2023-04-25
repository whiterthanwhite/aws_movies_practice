package movies

type Movie struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var Movies = []Movie{
	{
		ID:   1,
		Name: "Avengers",
	},
	{
		ID:   2,
		Name: "Ant-Man",
	},
	{
		ID:   3,
		Name: "Thor",
	},
	{
		ID:   4,
		Name: "Hulk",
	},
	{
		ID:   5,
		Name: "Doctor Strange",
	},
}
