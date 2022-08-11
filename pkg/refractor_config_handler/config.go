package refractor_config_handler

import (
	"fmt"
	"strings"
)

const (
	bf2GameDirName      = "Battlefield 2"
	profilesDirName     = "Profiles"
	globalConFileName   = "Global.con"
	profileConFileName  = "Profile.con"
	quoteChar           = "\""
	multiValueSeparator = ";"
)

type ErrNoSuchKey struct {
	key string
}

func (e *ErrNoSuchKey) Error() string {
	return fmt.Sprintf("no such key in config: %s", e.key)
}

type Config struct {
	content map[string]Value
}

func ConfigFromBytes(data []byte) *Config {
	// Split on \n in order to make parsing work with either \r\n or just \n line breaks
	lines := strings.Split(string(data), "\n")

	parsed := map[string]Value{}
	for _, line := range lines {
		// Trim any \r from line and split on first space
		elements := strings.SplitN(strings.Trim(line, "\r"), " ", 2)

		// TODO do something other than ignoring any invalid lines here?
		if len(elements) == 2 {
			// Add key, value or append to value
			key, content := elements[0], elements[1]
			current, exists := parsed[key]
			if exists {
				content = strings.Join([]string{current.content, content}, multiValueSeparator)
			}
			parsed[key] = Value{content: content}
		}
	}

	return &Config{
		content: parsed,
	}
}

func (c *Config) GetValue(key string) (Value, error) {
	value, ok := c.content[key]
	if !ok {
		return Value{}, &ErrNoSuchKey{key: key}
	}
	return value, nil
}

func (c *Config) SetValue(key string, value Value) {
	c.content[key] = value
}

func (c *Config) Delete(key string) {
	delete(c.content, key)
}

type Value struct {
	content string
}

func NewValue(content string) *Value {
	return &Value{
		content: content,
	}
}

func (v *Value) String() string {
	if isQuotedValue(v.content) {
		return strings.Trim(v.content, quoteChar)
	}
	return v.content
}

func (v *Value) Slice() []string {
	values := strings.Split(v.content, multiValueSeparator)
	for i, item := range values {
		if isQuotedValue(item) {
			values[i] = strings.Trim(item, quoteChar)
		}
	}
	return values
}

// isQuotedValue Checks whether a config value is a quoted string (starts and ends with a quote character, with no other quote characters in between)
func isQuotedValue(value string) bool {
	return strings.HasPrefix(value, quoteChar) && strings.HasSuffix(value, quoteChar) && strings.Count(value, quoteChar) == 2
}
