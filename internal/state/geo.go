package state

import (
	"math"

	"open-transit-rt/internal/gtfs"
)

const earthMetersPerDegree = 111320.0

type projection struct {
	DistanceMeters float64
	ShapeDistance  float64
	Bearing        float64
	Found          bool
}

func projectToShape(lat float64, lon float64, points []gtfs.ShapePoint) projection {
	if len(points) < 2 {
		return projection{}
	}

	refLat := lat * math.Pi / 180
	best := projection{DistanceMeters: math.MaxFloat64}
	for i := 0; i < len(points)-1; i++ {
		a := points[i]
		b := points[i+1]
		ax, ay := lonLatToXY(a.Lat, a.Lon, refLat)
		bx, by := lonLatToXY(b.Lat, b.Lon, refLat)
		px, py := lonLatToXY(lat, lon, refLat)
		dx := bx - ax
		dy := by - ay
		lengthSquared := dx*dx + dy*dy
		if lengthSquared == 0 {
			continue
		}
		t := ((px-ax)*dx + (py-ay)*dy) / lengthSquared
		if t < 0 {
			t = 0
		}
		if t > 1 {
			t = 1
		}
		qx := ax + t*dx
		qy := ay + t*dy
		dist := math.Hypot(px-qx, py-qy)
		if dist < best.DistanceMeters {
			shapeDistance := 0.0
			if a.HasDistance && b.HasDistance {
				shapeDistance = a.DistTraveled + t*(b.DistTraveled-a.DistTraveled)
			}
			best = projection{
				DistanceMeters: dist,
				ShapeDistance:  shapeDistance,
				Bearing:        bearing(a.Lat, a.Lon, b.Lat, b.Lon),
				Found:          true,
			}
		}
	}
	return best
}

func lonLatToXY(lat float64, lon float64, refLatRadians float64) (float64, float64) {
	x := lon * math.Cos(refLatRadians) * earthMetersPerDegree
	y := lat * earthMetersPerDegree
	return x, y
}

func bearing(lat1 float64, lon1 float64, lat2 float64, lon2 float64) float64 {
	phi1 := lat1 * math.Pi / 180
	phi2 := lat2 * math.Pi / 180
	delta := (lon2 - lon1) * math.Pi / 180
	y := math.Sin(delta) * math.Cos(phi2)
	x := math.Cos(phi1)*math.Sin(phi2) - math.Sin(phi1)*math.Cos(phi2)*math.Cos(delta)
	degrees := math.Atan2(y, x) * 180 / math.Pi
	return math.Mod(degrees+360, 360)
}

func bearingDelta(a float64, b float64) float64 {
	diff := math.Abs(a - b)
	if diff > 180 {
		diff = 360 - diff
	}
	return diff
}
