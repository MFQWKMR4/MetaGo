// Code generated by ggs; DO NOT EDIT.

package example

import (
)


func GetFirstname(v *Person) string {
	return v.firstName
}

func SetFirstname(v *Person, value string) {
	v.firstName = value
}

func GetLastname(v *Person) string {
	return v.lastName
}

func SetLastname(v *Person, value string) {
	v.lastName = value
}

func GetAge(v *Person) int {
	return v.Age
}

func SetAge(v *Person, value int) {
	v.Age = value
}

func GetDescription(v *Person) *string {
	return v.Description
}

func SetDescription(v *Person, value *string) {
	v.Description = value
}

func GetTags(v *Person) []int {
	return v.Tags
}

func SetTags(v *Person, value []int) {
	v.Tags = value
}

func GetGeo(v *Person) *Geo {
	return v.Geo
}

func SetGeo(v *Person, value *Geo) {
	v.Geo = value
}

func GetLat(v *Geo) float64 {
	return v.Lat
}

func SetLat(v *Geo, value float64) {
	v.Lat = value
}

func GetLng(v *Geo) float64 {
	return v.Lng
}

func SetLng(v *Geo, value float64) {
	v.Lng = value
}

func GetAddress(v *Geo) string {
	return v.Address
}

func SetAddress(v *Geo, value string) {
	v.Address = value
}

func GetIsrent(v *Geo) bool {
	return v.IsRent
}

func SetIsrent(v *Geo, value bool) {
	v.IsRent = value
}
