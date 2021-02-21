package short

// Stats -> struct that holds the stats of a ShortURL
type Stats struct {
	// Calls -> Global Calls
	Calls uint64

	// Calls60 -> Calls in 60 minutes
	Calls60 uint64
}
