package main

import (
	"time"

	"github.com/dadosjusbr/alba/storage"
	"github.com/urfave/cli"
)

// AddCollector is a function to add a collector in the database
func AddCollector(c *cli.Context) error {
	newCollector, err := parseCollectorFromContext(c)
	if err != nil {
		return err
	}
	if err := storage.InsertCollector(newCollector); err != nil {
		return err
	}
	return nil
}

func parseCollectorFromContext(c *cli.Context) (storage.Collector, error) {
	id := c.String("id")
	entity := c.String("entity")
	city := c.String("city")
	fu := c.String("fu")
	path := c.String("path")
	frequency := c.Int("frequency")
	startDay := c.Int("startDay")
	limitMonthBackward := c.Int("limitMonthBackward")
	limitYearBackward := c.Int("limitYearBackward")
	updateDate := time.Now()

	if id == "" {
		return storage.Collector{}, cli.Exit("--id were not provided completely. Please provide all parameters to continue", 1)
	}
	if entity == "" {
		return storage.Collector{}, cli.Exit("--entity were not provided completely. Please provide all parameters to continue", 1)

	}
	if city == "" {
		return storage.Collector{}, cli.Exit("--city were not provided completely. Please provide all parameters to continue", 1)

	}
	if fu == "" {
		return storage.Collector{}, cli.Exit("--fu were not provided completely. Please provide all parameters to continue", 1)

	}
	if path == "" {
		return storage.Collector{}, cli.Exit("--path were not provided completely. Please provide all parameters to continue", 1)

	}
	if frequency == 0 {
		return storage.Collector{}, cli.Exit("--frequency were not provided completely. Please provide all parameters to continue", 1)

	}
	if startDay == 0 {
		return storage.Collector{}, cli.Exit("--startDay were not provided completely. Please provide all parameters to continue", 1)

	}
	if limitMonthBackward == 0 {
		return storage.Collector{}, cli.Exit("--limitMonthBackward were not provided completely. Please provide all parameters to continue", 1)

	}
	if limitYearBackward == 0 {
		return storage.Collector{}, cli.Exit("--limitYearBackward were not provided completely. Please provide all parameters to continue", 1)
	}

	newCollector := storage.Collector{
		ID:                 id,
		Entity:             entity,
		City:               city,
		FU:                 fu,
		UpdateDate:         updateDate,
		Path:               path,
		Frequency:          frequency,
		StartDay:           startDay,
		LimitMonthBackward: limitMonthBackward,
		LimitYearBackward:  limitYearBackward}

	return newCollector, nil
}
