package packet

import (
	"math"
)

type Quaternion [4]float64
type Vector3 [3]float64

func quatConj(q Quaternion) Quaternion {
	return Quaternion{-q[0], -q[1], -q[2], q[3]}
}

func quatVec(q Quaternion) Vector3 {
	return Vector3{q[0], q[1], q[2]}
}

func quatScalar(q Quaternion) float64 {
	return q[3]
}

func cross(a, b Vector3) Vector3 {
	return Vector3{
		a[1]*b[2] - a[2]*b[1],
		a[2]*b[0] - a[0]*b[2],
		a[0]*b[1] - a[1]*b[0],
	}
}

func add(a, b Vector3) Vector3 {
	return Vector3{a[0] + b[0], a[1] + b[1], a[2] + b[2]}
}

func sub(a, b Vector3) Vector3 {
	return Vector3{a[0] - b[0], a[1] - b[1], a[2] - b[2]}
}

func scale(a Vector3, s float64) Vector3 {
	return Vector3{a[0] * s, a[1] * s, a[2] * s}
}

// QuatNorm normalizes a quaternion
func QuatNorm(q Quaternion) Quaternion {
	// Calculate the magnitude (norm) of the quaternion
	magnitude := math.Sqrt(q[0]*q[0] + q[1]*q[1] + q[2]*q[2] + q[3]*q[3])

	// If the magnitude is zero or very close to zero, return the original quaternion
	// to avoid division by zero
	if magnitude < 1e-10 {
		return q
	}

	// Divide each component by the magnitude
	return Quaternion{
		q[0] / magnitude,
		q[1] / magnitude,
		q[2] / magnitude,
		q[3] / magnitude,
	}
}

func quatRot(v Vector3, q Quaternion) Vector3 {
	qv := quatVec(q)
	u := cross(qv, v)
	w := quatScalar(q)
	p := add(u, scale(v, w))
	t := scale(qv, 2)
	rr := cross(t, p)
	r := add(v, rr)
	return r
}

func rollPitchYaw(q Quaternion) (float64, float64, float64) {
	// see http://www.euclideanspace.com/maths/geometry/rotations/conversions/quaternionToEuler/Quaternions.pdf
	// permute quaternion coefficients to determine
	// order of roll/pitch/yaw rotation axes, scalar comes first always
	p := [4]float64{q[3], q[2], q[0], q[1]}

	pNorm := QuatNorm(p)

	// e is for handedness of axes
	e := -1.0
	xP := 2 * (pNorm[0]*pNorm[2] + float64(e)*pNorm[1]*pNorm[3])
	pitch := math.Asin(xP)

	var roll, yaw float64

	if math.Abs(pitch) > math.Pi/2-1e-6 {
		// handle singularity when pitch is +-90°
		yaw = 0
		roll = math.Atan2(pNorm[1], pNorm[0])
	} else {
		yR := 2 * (pNorm[0]*pNorm[1] - float64(e)*pNorm[2]*pNorm[3])
		xR := 1 - 2*(pNorm[1]*pNorm[1]+pNorm[2]*pNorm[2])
		roll = math.Atan2(yR, xR)

		yY := 2 * (pNorm[0]*pNorm[3] - float64(e)*pNorm[1]*pNorm[2])
		xY := 1 - 2*(pNorm[2]*pNorm[2]+pNorm[3]*pNorm[3])
		yaw = math.Atan2(yY, xY)
	}

	return -roll, -pitch, -yaw
}
