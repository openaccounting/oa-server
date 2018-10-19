package types

import (
	"net/url"
	"strconv"
)

type QueryOptions struct {
	Limit                 int    `json:"limit"`
	Skip                  int    `json:"skip"`
	SinceInserted         int    `json:"sinceInserted"`
	SinceUpdated          int    `json:"sinceUpdated"`
	BeforeInserted        int    `json:"beforeInserted"`
	BeforeUpdated         int    `json:"beforeUpdated"`
	StartDate             int    `json:"startDate"`
	EndDate               int    `json:"endDate"`
	DescriptionStartsWith string `json:"descriptionStartsWith"`
	IncludeDeleted        bool   `json:"includeDeleted"`
	Sort                  string `json:"string"`
}

func QueryOptionsFromURLQuery(urlQuery url.Values) (*QueryOptions, error) {
	qo := &QueryOptions{}

	var err error

	if urlQuery.Get("limit") != "" {
		qo.Limit, err = strconv.Atoi(urlQuery.Get("limit"))

		if err != nil {
			return nil, err
		}
	}

	if urlQuery.Get("skip") != "" {
		qo.Skip, err = strconv.Atoi(urlQuery.Get("skip"))

		if err != nil {
			return nil, err
		}
	}

	if urlQuery.Get("sinceInserted") != "" {
		qo.SinceInserted, err = strconv.Atoi(urlQuery.Get("sinceInserted"))

		if err != nil {
			return nil, err
		}
	}

	if urlQuery.Get("sinceUpdated") != "" {
		qo.SinceUpdated, err = strconv.Atoi(urlQuery.Get("sinceUpdated"))

		if err != nil {
			return nil, err
		}
	}

	if urlQuery.Get("beforeInserted") != "" {
		qo.BeforeInserted, err = strconv.Atoi(urlQuery.Get("beforeInserted"))

		if err != nil {
			return nil, err
		}
	}

	if urlQuery.Get("beforeUpdated") != "" {
		qo.BeforeUpdated, err = strconv.Atoi(urlQuery.Get("beforeUpdated"))

		if err != nil {
			return nil, err
		}
	}

	if urlQuery.Get("startDate") != "" {
		qo.StartDate, err = strconv.Atoi(urlQuery.Get("startDate"))

		if err != nil {
			return nil, err
		}
	}

	if urlQuery.Get("endDate") != "" {
		qo.EndDate, err = strconv.Atoi(urlQuery.Get("endDate"))

		if err != nil {
			return nil, err
		}
	}

	if urlQuery.Get("descriptionStartsWith") != "" {
		qo.DescriptionStartsWith = urlQuery.Get("descriptionStartsWith")
	}

	if urlQuery.Get("includeDeleted") == "true" {
		qo.IncludeDeleted = true
	}

	if urlQuery.Get("sort") != "" {
		qo.Sort = urlQuery.Get("sort")
	}

	return qo, nil
}
