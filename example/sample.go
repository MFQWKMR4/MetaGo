package example

//go:generate go run ../cmd/ggs/ggs.go -input $GOFILE

type Person struct {
	firstName, lastName string  //+ggs
	Age                 int     `json:"age"`         //+ggs
	Description         *string `json:"description"` //+ggs
	Tags                []int   `json:"tags"`        //+ggs
	Geo                 *Geo    //+ggs
}

type Geo struct {
	Lat, Lng float64 //+ggs
	Address  string  //+ggs
	IsRent   bool    //+ggs
}
