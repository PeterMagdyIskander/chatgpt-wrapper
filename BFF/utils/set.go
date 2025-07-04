package utils

// Set implementation using map for O(1) operations
type Set map[string]struct{}

// NewSet creates a new empty set
func NewSet() Set {
	return make(Set)
}

// Add adds an item to the set
func (s Set) Add(item string) {
	s[item] = struct{}{}
}

// Contains checks if an item exists in the set
func (s Set) Contains(item string) bool {
	_, ok := s[item]
	return ok
}

// Remove removes an item from the set
func (s Set) Remove(item string) {
	delete(s, item)
}

// Size returns the number of items in the set
func (s Set) Size() int {
	return len(s)
}

// Values returns all items in the set as a slice
func (s Set) Values() []string {
	result := make([]string, 0, len(s))
	for k := range s {
		result = append(result, k)
	}
	return result
}
