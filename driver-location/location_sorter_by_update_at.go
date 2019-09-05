package driver_location

import "github.com/heetch/jose-odg-technical-test/driver-location/internal"

type ByUpdatedAt []domain.Location

func (a ByUpdatedAt) Len() int      { return len(a) }
func (a ByUpdatedAt) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByUpdatedAt) Less(i, j int) bool {
	return a[i].UpdatedAt().Date().Before(a[j].UpdatedAt().Date())
}
