// TODO: to another application (geo app)?
package event

type Point struct {
	Lt float64
	Ln float64
}

type Place struct {
	Point Point
	Title string
}
