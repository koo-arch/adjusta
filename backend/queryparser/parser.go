package queryparser

import (
	"fmt"
	"time"
	"strconv"

	"github.com/gin-gonic/gin"
)

type QueryParser struct {
	c *gin.Context
}

func NewQueryParser(c *gin.Context) *QueryParser {
	return &QueryParser{c: c}
}

func (qp *QueryParser) ParseTime(key string) (*time.Time, error) {
	timeStr := qp.c.Query(key)
	if timeStr == "" {
		return nil, nil
	}

	t, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse time: %w", err)
	}

	return &t, nil
}

func(qp *QueryParser) ParseDefaultTime(key string, defaultValue time.Time) (*time.Time, error) {
	timeStr := qp.c.Query(key)
	if timeStr == "" {
		return &defaultValue, nil
	}

	t, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse time: %w", err)
	}

	return &t, nil
}

func (qp *QueryParser) ParseString(key string) (*string, error) {
	str := qp.c.Query(key)
	if str == "" {
		return nil, nil
	}
	
	return &str, nil
}

func(qp *QueryParser) ParseDefaultString(key string, defaultValue string) (*string, error) {
	str := qp.c.Query(key)
	if str == "" {
		return &defaultValue, nil
	}

	return &str, nil
}

func (qp *QueryParser) ParseInt(key string) (*int, error) {
	str := qp.c.Query(key)
	if str == "" {
		return nil, nil
	}
	
	i, err := strconv.Atoi(str)
	if err != nil {
		return nil, fmt.Errorf("failed to parse int: %w", err)
	}
	
	return &i, nil
}

func(qp *QueryParser) ParseDefaultInt(key string, defaultValue int) (*int, error) {
	str := qp.c.Query(key)
	if str == "" {
		return &defaultValue, nil
	}

	i, err := strconv.Atoi(str)
	if err != nil {
		return nil, fmt.Errorf("failed to parse int: %w", err)
	}

	return &i, nil
}

func(qp *QueryParser) ParseBool(key string) (*bool, error) {
	str := qp.c.Query(key)
	if str == "" {
		return nil, nil
	}

	b, err := strconv.ParseBool(str)
	if err != nil {
		return nil, fmt.Errorf("failed to parse bool: %w", err)
	}

	return &b, nil
}

func(qp *QueryParser) ParseDefaultBool(key string, defaultValue bool) (*bool, error) {
	str := qp.c.Query(key)
	if str == "" {
		return &defaultValue, nil
	}

	b, err := strconv.ParseBool(str)
	if err != nil {
		return nil, fmt.Errorf("failed to parse bool: %w", err)
	}

	return &b, nil
}