package book

import (
  "fmt"
  "regexp"
  "strings"
)

type ISBN struct {
  value string
}

func NewISBN(value string) (*ISBN, error) {
  // Step 1: Normalize input - remove spaces and hyphens
  // Accepts: "978-3-16-148410-0" or "9783161484100"
  cleanValue := strings.ReplaceAll(strings.ReplaceAll(value, " ", ""), "-", "")

  // Step 2: Validate format - must be exactly 13 digits
  matched, _ := regexp.MatchString(`^\d{13}$`, cleanValue)
  if !matched {
    return nil, fmt.Errorf("ISBN must be 13 digits. Got: %s", value)
  }

  // Step 3: Store clean value (no hyphens)
  return &ISBN{value: cleanValue}, nil
}

func (i *ISBN) GetValue() string {
  return i.value
}

// Step 4: Provide formatted output when needed
// Returns: "978-3-16-148410-0"
func (i *ISBN) GetFormatted() string {
  return fmt.Sprintf("%s-%s-%s-%s-%s",
    i.value[0:3],
    i.value[3:4],
    i.value[4:6],
    i.value[6:12],
    i.value[12:13])
}

func (i *ISBN) Equals(other *ISBN) bool {
  return i.value == other.value
}

func (i *ISBN) String() string {
  return i.value
}