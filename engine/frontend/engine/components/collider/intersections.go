package collider

import "shared/utils/httperrors"

func collides(shape1, shape2 any) (Intersection, error) {
	switch s1 := shape1.(type) {
	case *box:
		switch s2 := shape2.(type) {
		case *box:
			return intersectsBoxBox(s1, s2)
		case *ray:
			return intersectsBoxRay(s1, s2)
		case *sphere:
			return intersectsBoxSphere(s1, s2)
		}
		break
	case *sphere:
		switch s2 := shape2.(type) {
		case *box:
			return intersectsSphereBox(s1, s2)
		case *ray:
			return intersectsSphereRay(s1, s2)
		case *sphere:
			return intersectsSphereSphere(s1, s2)
		}
		break
	case *ray:
		switch s2 := shape2.(type) {
		case *box:
			return intersectsRayBox(s1, s2)
		case *ray:
			return intersectsRayRay(s1, s2)
		case *sphere:
			return intersectsRaySphere(s1, s2)
		}
		break
	}
	return nil, httperrors.Err501
}

func intersectsBoxBox(shape1 *box, shape2 *box) (Intersection, error) {
	return nil, nil
}
func intersectsBoxRay(shape1 *box, shape2 *ray) (Intersection, error) {
	return nil, nil
}
func intersectsBoxSphere(shape1 *box, shape2 *sphere) (Intersection, error) {
	return nil, nil
}

func intersectsRayBox(shape1 *ray, shape2 *box) (Intersection, error) {
	return nil, nil
}
func intersectsRayRay(shape1 *ray, shape2 *ray) (Intersection, error) {
	return nil, nil
}
func intersectsRaySphere(shape1 *ray, shape2 *sphere) (Intersection, error) {
	return nil, nil
}

func intersectsSphereBox(shape1 *sphere, shape2 *box) (Intersection, error) {
	return nil, nil
}
func intersectsSphereRay(shape1 *sphere, shape2 *ray) (Intersection, error) {
	return nil, nil
}
func intersectsSphereSphere(shape1 *sphere, shape2 *sphere) (Intersection, error) {
	return nil, nil
}
