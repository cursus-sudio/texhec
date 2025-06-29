package media

// TODO
// make this make sense after researching go-raylib
// because currently this is special orderer with some random "maybe to use structs"
// remember that camera is not included here by default but it has to be appliable on objects
// (at least that's the way i imagine this)

// 2D

type Position2D struct{ X, Y int } // (0, 0) means top left corner

// moves clockwise
// 0 deg points up ( means no rotation )
// 90 deg points right
type Rotation2D struct{ Deg float64 }

// images are scalled
type Size2D struct{ Width, Height int }

// options are passed in ctor
//
// e.g. options
//
//	type Options struct {
//		Position Position2D
//		Rotation Rotation2D
//		Size     Size2D
//	}
type Drawable interface {
	Draw()
}

//

// by z-index asset calling order is determined
type ZIndex float64 // small z-index is drawn first and big z-index is drawn last

type DrawApi interface {
	// is executed on flush

	Clear() error
	Draw(Drawable, ZIndex) error

	// nothing is done until flush
	Flush() error
}
