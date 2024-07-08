package common

type PreferenceChange struct {
    UserId string               `dynamodbav:"userId" json:"userId"`
    Actors map[string]string    `dynamodbav:"actors" json:"actors"`
    Directors map[string]string `dynamodbav:"directors" json:"directors"`
    Genres map[string]string    `dynamodbav:"genres" json:"genres"`

    // adding to the counter which trigger update
    UpdateWeight int            `dynamodbav:"updateWeight" json:"updateWeight"`

    // how much preference for each of the item changes
    ChangeWeight int            `dynamodbav:"changeWeight" json:"changeWeight"`
}
