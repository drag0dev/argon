package common

type PreferenceChange struct {
    UserId string      `dynamodbav:"userId" json:"userId"`
    Actors []string    `dynamodbav:"actors" json:"actors"`
    Directors []string `dynamodbav:"directors" json:"directors"`
    Genres []string    `dynamodbav:"genres" json:"genres"`

    // adding to the counter which trigger update
    UpdateWeight float64   `dynamodbav:"updateWeight" json:"updateWeight"`

    // how much preference for each of the item changes
    ChangeWeight float64   `dynamodbav:"changeWeight" json:"changeWeight"`
}

type UserPreference struct {
    Actors map[string]float64
    Directors map[string]float64
    Genres map[string]float64
    UpdateCounter float64
}
