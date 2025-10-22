package models

import (
	"database/sql"
	"time"
)

type Point struct {
	ID   int64
	Name string
	Lat  float64
	Lon  float64
}

type RoutePoint struct {
	PointID int64
	Seq     int
	Point   Point
}

type Route struct {
	ID       int64
	DepartAt time.Time
	MaxW     float64
	MaxV     float64
	Status   string
	Points   []RoutePoint
}

type Request struct {
	ID       int64
	OriginID sql.NullInt64
	DestID   sql.NullInt64
	Lat      float64
	Lon      float64
	Weight   float64
	Volume   float64
	ReadyAt  time.Time
	Deadline sql.NullTime
	Status   string
}

type Assignment struct {
	RequestID int64
	Seq       int
	ETA       time.Time
}
