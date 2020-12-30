package data

type Skateboard struct {
  Brand string `json:"brand"`
  Type string `json:"type"`
  New bool `json:"new"`
  Id int `json:"id"`
}

var Skateboards = []Skateboard{
  {Brand: "Element", Type: "skateboard", New: true, Id: 1},
  {Brand: "Santa Cruz", Type: "longboard", New: false, Id: 2},
  {Brand: "Magneto", Type: "mini cruiser", New: true, Id: 3},
}
